FROM golang:1.18 as builder

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN apt update && apt install -y protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/google/gnostic/cmd/protoc-gen-openapi \
    github.com/srikrsna/protoc-gen-gotag

ADD . .

RUN go generate ./...
RUN go build -o /build/engine ./cmd/engine/engine.go
RUN go build -o /build/cl cmd/cli/cl.go

FROM debian:stable-slim

COPY --from=builder /build/engine /engine
COPY --from=builder /build/cl /cl
COPY --from=builder /build/policies /policies
COPY --from=builder /build/service/orchestrator/metrics.json /metrics.json

CMD ["./engine", "--db-in-memory"]
