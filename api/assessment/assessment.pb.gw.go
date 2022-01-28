// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: assessment.proto

/*
Package assessment is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package assessment

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_Assessment_AssessEvidence_0(ctx context.Context, marshaler runtime.Marshaler, client AssessmentClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq AssessEvidenceRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq.Evidence); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.AssessEvidence(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Assessment_AssessEvidence_0(ctx context.Context, marshaler runtime.Marshaler, server AssessmentServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq AssessEvidenceRequest
	var metadata runtime.ServerMetadata

	newReader, berr := utilities.IOReaderFactory(req.Body)
	if berr != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", berr)
	}
	if err := marshaler.NewDecoder(newReader()).Decode(&protoReq.Evidence); err != nil && err != io.EOF {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.AssessEvidence(ctx, &protoReq)
	return msg, metadata, err

}

func request_Assessment_ListAssessmentResults_0(ctx context.Context, marshaler runtime.Marshaler, client AssessmentClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ListAssessmentResultsRequest
	var metadata runtime.ServerMetadata

	msg, err := client.ListAssessmentResults(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_Assessment_ListAssessmentResults_0(ctx context.Context, marshaler runtime.Marshaler, server AssessmentServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ListAssessmentResultsRequest
	var metadata runtime.ServerMetadata

	msg, err := server.ListAssessmentResults(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterAssessmentHandlerServer registers the http handlers for service Assessment to "mux".
// UnaryRPC     :call AssessmentServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterAssessmentHandlerFromEndpoint instead.
func RegisterAssessmentHandlerServer(ctx context.Context, mux *runtime.ServeMux, server AssessmentServer) error {

	mux.Handle("POST", pattern_Assessment_AssessEvidence_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/clouditor.Assessment/AssessEvidence")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Assessment_AssessEvidence_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Assessment_AssessEvidence_0(ctx, mux, outboundMarshaler, w, req, response_Assessment_AssessEvidence_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Assessment_ListAssessmentResults_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req, "/clouditor.Assessment/ListAssessmentResults")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_Assessment_ListAssessmentResults_0(rctx, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Assessment_ListAssessmentResults_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterAssessmentHandlerFromEndpoint is same as RegisterAssessmentHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterAssessmentHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterAssessmentHandler(ctx, mux, conn)
}

// RegisterAssessmentHandler registers the http handlers for service Assessment to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterAssessmentHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterAssessmentHandlerClient(ctx, mux, NewAssessmentClient(conn))
}

// RegisterAssessmentHandlerClient registers the http handlers for service Assessment
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "AssessmentClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "AssessmentClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "AssessmentClient" to call the correct interceptors.
func RegisterAssessmentHandlerClient(ctx context.Context, mux *runtime.ServeMux, client AssessmentClient) error {

	mux.Handle("POST", pattern_Assessment_AssessEvidence_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/clouditor.Assessment/AssessEvidence")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Assessment_AssessEvidence_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Assessment_AssessEvidence_0(ctx, mux, outboundMarshaler, w, req, response_Assessment_AssessEvidence_0{resp}, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_Assessment_ListAssessmentResults_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req, "/clouditor.Assessment/ListAssessmentResults")
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_Assessment_ListAssessmentResults_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_Assessment_ListAssessmentResults_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

type response_Assessment_AssessEvidence_0 struct {
	proto.Message
}

func (m response_Assessment_AssessEvidence_0) XXX_ResponseBody() interface{} {
	response := m.Message.(*AssessEvidenceResponse)
	return response.Status
}

var (
	pattern_Assessment_AssessEvidence_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"v1", "assessment", "evidences"}, ""))

	pattern_Assessment_ListAssessmentResults_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"v1", "assessment", "results"}, ""))
)

var (
	forward_Assessment_AssessEvidence_0 = runtime.ForwardResponseMessage

	forward_Assessment_ListAssessmentResults_0 = runtime.ForwardResponseMessage
)
