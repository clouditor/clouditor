FROM golang:1.18-alpine as builder

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN apk update && apk add protobuf gcc libc-dev
RUN apk update && apk add nodejs npm && npm install @bufbuild/buf

RUN go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/google/gnostic/cmd/protoc-gen-openapi \
    github.com/srikrsna/protoc-gen-gotag

ADD . .

RUN go generate ./...
RUN go build -o /build/engine ./cmd/engine/engine.go
RUN go build -o /build/cl cmd/cli/cl.go

FROM alpine

WORKDIR /app

#RUN apk update && apk add gcc libc-dev

COPY --from=builder /build/engine .
COPY --from=builder /build/cl .
COPY --from=builder /build/policies ./policies
COPY --from=builder /build/service/orchestrator/metrics.json .

# Expose port fer rest gateway (For OAuth to work you also should publish this port when running the container image)
EXPOSE 8080
# Expose port for grpc
EXPOSE 9090


# Set program arguments via ENV variables
CMD ["./engine"]
