FROM golang:1.18-alpine as builder

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN apk update && apk add protobuf

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/google/gnostic/cmd/protoc-gen-openapi \
    github.com/srikrsna/protoc-gen-gotag

ADD . .

RUN go generate ./...
RUN go build -o /build/engine ./cmd/engine/engine.go
RUN go build -o /build/cl cmd/cli/cl.go

FROM golang:1.18-alpine

WORKDIR /app

COPY --from=builder /build/engine .
COPY --from=builder /build/cl .
COPY --from=builder /build/policies ./policies
COPY --from=builder /build/service/orchestrator/metrics.json .
RUN mkdir "/root/.clouditor"

# TODO(lebogg): Use ENV instead of hardcoded arguments
CMD ["./engine", "--db-in-memory", "--discovery-auto-start", "--discovery-provider=azure", "--dashboard-url=deployment_engine_1:3000"]
