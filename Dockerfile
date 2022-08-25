FROM golang:1.18-alpine as builder

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN apk update && apk add protobuf gcc libc-dev

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
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
RUN mkdir "/root/.clouditor" # TODO: Can be removed after https://github.com/clouditor/clouditor/issues/786 is fixed


# Set program arguments via ENV variables
CMD ["./engine"]
