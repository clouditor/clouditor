FROM golang

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN apt update && apt install -y protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go 
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc 
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

ADD . .

RUN go generate ./...
RUN go build ./...
