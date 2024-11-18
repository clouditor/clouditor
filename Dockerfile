FROM golang:1.22-alpine3.18 AS builder

WORKDIR /build

ADD go.mod .
ADD go.sum .
# We need the .git folder for the git tag and commit hash for the Runtime API endpoint
ADD .git .

RUN apk update && apk add protobuf gcc libc-dev git

RUN go install \
    github.com/oxisto/owl2proto/cmd/owl2proto \
    github.com/srikrsna/protoc-gen-gotag

RUN go install github.com/mattn/go-sqlite3
# The latest version does not work with the alpine image. There are problems regarding the glibc vs. musl library
RUN go install github.com/bufbuild/buf/cmd/buf@v1.45.0 

ADD . .

# RUN go generate ./...
RUN go build -ldflags="-X clouditor.io/clouditor/v2/service.version=$(git describe --exact-match --tags --abbrev=0)" -o /build/engine ./cmd/engine
RUN go build -ldflags="-X clouditor.io/clouditor/v2/service.version=$(git describe --exact-match --tags --abbrev=0)" -o /build/cl ./cmd/cli

FROM alpine

WORKDIR /app

COPY --from=builder /build/engine .
COPY --from=builder /build/cl .
COPY --from=builder /build/catalogs ./catalogs
COPY --from=builder /build/policies ./policies
COPY --from=builder /build/service/orchestrator/metrics.json .

# Expose port for rest gateway (For OAuth to work you also should publish this port when running the container image)
EXPOSE 8080
# Expose port for grpc
EXPOSE 9090

ENTRYPOINT ["./engine"]