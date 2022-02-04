// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: orchestrator.proto

package orchestrator

import (
	assessment "clouditor.io/clouditor/api/assessment"
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OrchestratorClient is the client API for Orchestrator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrchestratorClient interface {
	// Registers the passed assessment tool
	RegisterAssessmentTool(ctx context.Context, in *RegisterAssessmentToolRequest, opts ...grpc.CallOption) (*AssessmentTool, error)
	// Lists all assessment tools assessing evidences for the metric given by the
	// passed metric id
	ListAssessmentTools(ctx context.Context, in *ListAssessmentToolsRequest, opts ...grpc.CallOption) (*ListAssessmentToolsResponse, error)
	// Returns assessment tool given by the passed tool id
	GetAssessmentTool(ctx context.Context, in *GetAssessmentToolRequest, opts ...grpc.CallOption) (*AssessmentTool, error)
	// Updates the assessment tool given by the passed id
	UpdateAssessmentTool(ctx context.Context, in *UpdateAssessmentToolRequest, opts ...grpc.CallOption) (*AssessmentTool, error)
	// Remove assessment tool with passed id from the list of active assessment
	// tools
	DeregisterAssessmentTool(ctx context.Context, in *DeregisterAssessmentToolRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Stores the assessment result provided by an assessment tool
	StoreAssessmentResult(ctx context.Context, in *StoreAssessmentResultRequest, opts ...grpc.CallOption) (*StoreAssessmentResultResponse, error)
	// Stores stream of assessment results provided by an assessment tool
	StoreAssessmentResults(ctx context.Context, opts ...grpc.CallOption) (Orchestrator_StoreAssessmentResultsClient, error)
	// List all assessment results. Part of the public API, also exposed as REST.
	ListAssessmentResults(ctx context.Context, in *assessment.ListAssessmentResultsRequest, opts ...grpc.CallOption) (*assessment.ListAssessmentResultsResponse, error)
	// Creates a new metric
	CreateMetric(ctx context.Context, in *CreateMetricRequest, opts ...grpc.CallOption) (*assessment.Metric, error)
	// Updates an existing metric
	UpdateMetric(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*assessment.Metric, error)
	// Returns the metric with the passed metric id
	GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*assessment.Metric, error)
	// List all metrics provided by the metric catalog
	ListMetrics(ctx context.Context, in *ListMetricsRequest, opts ...grpc.CallOption) (*ListMetricsResponse, error)
	// Registers a new target cloud service
	RegisterCloudService(ctx context.Context, in *RegisterCloudServiceRequest, opts ...grpc.CallOption) (*CloudService, error)
	// Registers a new target cloud service
	UpdateCloudService(ctx context.Context, in *UpdateCloudServiceRequest, opts ...grpc.CallOption) (*CloudService, error)
	// Retrieves a target cloud service
	GetCloudService(ctx context.Context, in *GetCloudServiceRequest, opts ...grpc.CallOption) (*CloudService, error)
	// Lists all target cloud services
	ListCloudServices(ctx context.Context, in *ListCloudServicesRequest, opts ...grpc.CallOption) (*ListCloudServicesResponse, error)
	// Removes a target cloud service
	RemoveCloudService(ctx context.Context, in *RemoveCloudServiceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Retrieves a metric configuration for a specific service and metric ID
	GetMetricConfiguration(ctx context.Context, in *GetMetricConfigurationRequest, opts ...grpc.CallOption) (*assessment.MetricConfiguration, error)
	// Lists all a metric configurations for a specific service and metric ID
	ListMetricConfigurations(ctx context.Context, in *ListMetricConfigurationRequest, opts ...grpc.CallOption) (*ListMetricConfigurationResponse, error)
}

type orchestratorClient struct {
	cc grpc.ClientConnInterface
}

func NewOrchestratorClient(cc grpc.ClientConnInterface) OrchestratorClient {
	return &orchestratorClient{cc}
}

func (c *orchestratorClient) RegisterAssessmentTool(ctx context.Context, in *RegisterAssessmentToolRequest, opts ...grpc.CallOption) (*AssessmentTool, error) {
	out := new(AssessmentTool)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/RegisterAssessmentTool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) ListAssessmentTools(ctx context.Context, in *ListAssessmentToolsRequest, opts ...grpc.CallOption) (*ListAssessmentToolsResponse, error) {
	out := new(ListAssessmentToolsResponse)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/ListAssessmentTools", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) GetAssessmentTool(ctx context.Context, in *GetAssessmentToolRequest, opts ...grpc.CallOption) (*AssessmentTool, error) {
	out := new(AssessmentTool)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/GetAssessmentTool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) UpdateAssessmentTool(ctx context.Context, in *UpdateAssessmentToolRequest, opts ...grpc.CallOption) (*AssessmentTool, error) {
	out := new(AssessmentTool)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/UpdateAssessmentTool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) DeregisterAssessmentTool(ctx context.Context, in *DeregisterAssessmentToolRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/DeregisterAssessmentTool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) StoreAssessmentResult(ctx context.Context, in *StoreAssessmentResultRequest, opts ...grpc.CallOption) (*StoreAssessmentResultResponse, error) {
	out := new(StoreAssessmentResultResponse)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/StoreAssessmentResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) StoreAssessmentResults(ctx context.Context, opts ...grpc.CallOption) (Orchestrator_StoreAssessmentResultsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Orchestrator_ServiceDesc.Streams[0], "/clouditor.Orchestrator/StoreAssessmentResults", opts...)
	if err != nil {
		return nil, err
	}
	x := &orchestratorStoreAssessmentResultsClient{stream}
	return x, nil
}

type Orchestrator_StoreAssessmentResultsClient interface {
	Send(*assessment.AssessmentResult) error
	CloseAndRecv() (*emptypb.Empty, error)
	grpc.ClientStream
}

type orchestratorStoreAssessmentResultsClient struct {
	grpc.ClientStream
}

func (x *orchestratorStoreAssessmentResultsClient) Send(m *assessment.AssessmentResult) error {
	return x.ClientStream.SendMsg(m)
}

func (x *orchestratorStoreAssessmentResultsClient) CloseAndRecv() (*emptypb.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(emptypb.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *orchestratorClient) ListAssessmentResults(ctx context.Context, in *assessment.ListAssessmentResultsRequest, opts ...grpc.CallOption) (*assessment.ListAssessmentResultsResponse, error) {
	out := new(assessment.ListAssessmentResultsResponse)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/ListAssessmentResults", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) CreateMetric(ctx context.Context, in *CreateMetricRequest, opts ...grpc.CallOption) (*assessment.Metric, error) {
	out := new(assessment.Metric)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/CreateMetric", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) UpdateMetric(ctx context.Context, in *UpdateMetricRequest, opts ...grpc.CallOption) (*assessment.Metric, error) {
	out := new(assessment.Metric)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/UpdateMetric", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) GetMetric(ctx context.Context, in *GetMetricRequest, opts ...grpc.CallOption) (*assessment.Metric, error) {
	out := new(assessment.Metric)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/GetMetric", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) ListMetrics(ctx context.Context, in *ListMetricsRequest, opts ...grpc.CallOption) (*ListMetricsResponse, error) {
	out := new(ListMetricsResponse)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/ListMetrics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) RegisterCloudService(ctx context.Context, in *RegisterCloudServiceRequest, opts ...grpc.CallOption) (*CloudService, error) {
	out := new(CloudService)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/RegisterCloudService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) UpdateCloudService(ctx context.Context, in *UpdateCloudServiceRequest, opts ...grpc.CallOption) (*CloudService, error) {
	out := new(CloudService)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/UpdateCloudService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) GetCloudService(ctx context.Context, in *GetCloudServiceRequest, opts ...grpc.CallOption) (*CloudService, error) {
	out := new(CloudService)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/GetCloudService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) ListCloudServices(ctx context.Context, in *ListCloudServicesRequest, opts ...grpc.CallOption) (*ListCloudServicesResponse, error) {
	out := new(ListCloudServicesResponse)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/ListCloudServices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) RemoveCloudService(ctx context.Context, in *RemoveCloudServiceRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/RemoveCloudService", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) GetMetricConfiguration(ctx context.Context, in *GetMetricConfigurationRequest, opts ...grpc.CallOption) (*assessment.MetricConfiguration, error) {
	out := new(assessment.MetricConfiguration)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/GetMetricConfiguration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orchestratorClient) ListMetricConfigurations(ctx context.Context, in *ListMetricConfigurationRequest, opts ...grpc.CallOption) (*ListMetricConfigurationResponse, error) {
	out := new(ListMetricConfigurationResponse)
	err := c.cc.Invoke(ctx, "/clouditor.Orchestrator/ListMetricConfigurations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrchestratorServer is the server API for Orchestrator service.
// All implementations must embed UnimplementedOrchestratorServer
// for forward compatibility
type OrchestratorServer interface {
	// Registers the passed assessment tool
	RegisterAssessmentTool(context.Context, *RegisterAssessmentToolRequest) (*AssessmentTool, error)
	// Lists all assessment tools assessing evidences for the metric given by the
	// passed metric id
	ListAssessmentTools(context.Context, *ListAssessmentToolsRequest) (*ListAssessmentToolsResponse, error)
	// Returns assessment tool given by the passed tool id
	GetAssessmentTool(context.Context, *GetAssessmentToolRequest) (*AssessmentTool, error)
	// Updates the assessment tool given by the passed id
	UpdateAssessmentTool(context.Context, *UpdateAssessmentToolRequest) (*AssessmentTool, error)
	// Remove assessment tool with passed id from the list of active assessment
	// tools
	DeregisterAssessmentTool(context.Context, *DeregisterAssessmentToolRequest) (*emptypb.Empty, error)
	// Stores the assessment result provided by an assessment tool
	StoreAssessmentResult(context.Context, *StoreAssessmentResultRequest) (*StoreAssessmentResultResponse, error)
	// Stores stream of assessment results provided by an assessment tool
	StoreAssessmentResults(Orchestrator_StoreAssessmentResultsServer) error
	// List all assessment results. Part of the public API, also exposed as REST.
	ListAssessmentResults(context.Context, *assessment.ListAssessmentResultsRequest) (*assessment.ListAssessmentResultsResponse, error)
	// Creates a new metric
	CreateMetric(context.Context, *CreateMetricRequest) (*assessment.Metric, error)
	// Updates an existing metric
	UpdateMetric(context.Context, *UpdateMetricRequest) (*assessment.Metric, error)
	// Returns the metric with the passed metric id
	GetMetric(context.Context, *GetMetricRequest) (*assessment.Metric, error)
	// List all metrics provided by the metric catalog
	ListMetrics(context.Context, *ListMetricsRequest) (*ListMetricsResponse, error)
	// Registers a new target cloud service
	RegisterCloudService(context.Context, *RegisterCloudServiceRequest) (*CloudService, error)
	// Registers a new target cloud service
	UpdateCloudService(context.Context, *UpdateCloudServiceRequest) (*CloudService, error)
	// Retrieves a target cloud service
	GetCloudService(context.Context, *GetCloudServiceRequest) (*CloudService, error)
	// Lists all target cloud services
	ListCloudServices(context.Context, *ListCloudServicesRequest) (*ListCloudServicesResponse, error)
	// Removes a target cloud service
	RemoveCloudService(context.Context, *RemoveCloudServiceRequest) (*emptypb.Empty, error)
	// Retrieves a metric configuration for a specific service and metric ID
	GetMetricConfiguration(context.Context, *GetMetricConfigurationRequest) (*assessment.MetricConfiguration, error)
	// Lists all a metric configurations for a specific service and metric ID
	ListMetricConfigurations(context.Context, *ListMetricConfigurationRequest) (*ListMetricConfigurationResponse, error)
	mustEmbedUnimplementedOrchestratorServer()
}

// UnimplementedOrchestratorServer must be embedded to have forward compatible implementations.
type UnimplementedOrchestratorServer struct {
}

func (UnimplementedOrchestratorServer) RegisterAssessmentTool(context.Context, *RegisterAssessmentToolRequest) (*AssessmentTool, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterAssessmentTool not implemented")
}
func (UnimplementedOrchestratorServer) ListAssessmentTools(context.Context, *ListAssessmentToolsRequest) (*ListAssessmentToolsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAssessmentTools not implemented")
}
func (UnimplementedOrchestratorServer) GetAssessmentTool(context.Context, *GetAssessmentToolRequest) (*AssessmentTool, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAssessmentTool not implemented")
}
func (UnimplementedOrchestratorServer) UpdateAssessmentTool(context.Context, *UpdateAssessmentToolRequest) (*AssessmentTool, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAssessmentTool not implemented")
}
func (UnimplementedOrchestratorServer) DeregisterAssessmentTool(context.Context, *DeregisterAssessmentToolRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeregisterAssessmentTool not implemented")
}
func (UnimplementedOrchestratorServer) StoreAssessmentResult(context.Context, *StoreAssessmentResultRequest) (*StoreAssessmentResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreAssessmentResult not implemented")
}
func (UnimplementedOrchestratorServer) StoreAssessmentResults(Orchestrator_StoreAssessmentResultsServer) error {
	return status.Errorf(codes.Unimplemented, "method StoreAssessmentResults not implemented")
}
func (UnimplementedOrchestratorServer) ListAssessmentResults(context.Context, *assessment.ListAssessmentResultsRequest) (*assessment.ListAssessmentResultsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAssessmentResults not implemented")
}
func (UnimplementedOrchestratorServer) CreateMetric(context.Context, *CreateMetricRequest) (*assessment.Metric, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMetric not implemented")
}
func (UnimplementedOrchestratorServer) UpdateMetric(context.Context, *UpdateMetricRequest) (*assessment.Metric, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMetric not implemented")
}
func (UnimplementedOrchestratorServer) GetMetric(context.Context, *GetMetricRequest) (*assessment.Metric, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetric not implemented")
}
func (UnimplementedOrchestratorServer) ListMetrics(context.Context, *ListMetricsRequest) (*ListMetricsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMetrics not implemented")
}
func (UnimplementedOrchestratorServer) RegisterCloudService(context.Context, *RegisterCloudServiceRequest) (*CloudService, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterCloudService not implemented")
}
func (UnimplementedOrchestratorServer) UpdateCloudService(context.Context, *UpdateCloudServiceRequest) (*CloudService, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCloudService not implemented")
}
func (UnimplementedOrchestratorServer) GetCloudService(context.Context, *GetCloudServiceRequest) (*CloudService, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCloudService not implemented")
}
func (UnimplementedOrchestratorServer) ListCloudServices(context.Context, *ListCloudServicesRequest) (*ListCloudServicesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCloudServices not implemented")
}
func (UnimplementedOrchestratorServer) RemoveCloudService(context.Context, *RemoveCloudServiceRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveCloudService not implemented")
}
func (UnimplementedOrchestratorServer) GetMetricConfiguration(context.Context, *GetMetricConfigurationRequest) (*assessment.MetricConfiguration, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMetricConfiguration not implemented")
}
func (UnimplementedOrchestratorServer) ListMetricConfigurations(context.Context, *ListMetricConfigurationRequest) (*ListMetricConfigurationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMetricConfigurations not implemented")
}
func (UnimplementedOrchestratorServer) mustEmbedUnimplementedOrchestratorServer() {}

// UnsafeOrchestratorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrchestratorServer will
// result in compilation errors.
type UnsafeOrchestratorServer interface {
	mustEmbedUnimplementedOrchestratorServer()
}

func RegisterOrchestratorServer(s grpc.ServiceRegistrar, srv OrchestratorServer) {
	s.RegisterService(&Orchestrator_ServiceDesc, srv)
}

func _Orchestrator_RegisterAssessmentTool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterAssessmentToolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).RegisterAssessmentTool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/RegisterAssessmentTool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).RegisterAssessmentTool(ctx, req.(*RegisterAssessmentToolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_ListAssessmentTools_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAssessmentToolsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).ListAssessmentTools(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/ListAssessmentTools",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).ListAssessmentTools(ctx, req.(*ListAssessmentToolsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_GetAssessmentTool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAssessmentToolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).GetAssessmentTool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/GetAssessmentTool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).GetAssessmentTool(ctx, req.(*GetAssessmentToolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_UpdateAssessmentTool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAssessmentToolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).UpdateAssessmentTool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/UpdateAssessmentTool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).UpdateAssessmentTool(ctx, req.(*UpdateAssessmentToolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_DeregisterAssessmentTool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeregisterAssessmentToolRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).DeregisterAssessmentTool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/DeregisterAssessmentTool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).DeregisterAssessmentTool(ctx, req.(*DeregisterAssessmentToolRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_StoreAssessmentResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreAssessmentResultRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).StoreAssessmentResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/StoreAssessmentResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).StoreAssessmentResult(ctx, req.(*StoreAssessmentResultRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_StoreAssessmentResults_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OrchestratorServer).StoreAssessmentResults(&orchestratorStoreAssessmentResultsServer{stream})
}

type Orchestrator_StoreAssessmentResultsServer interface {
	SendAndClose(*emptypb.Empty) error
	Recv() (*assessment.AssessmentResult, error)
	grpc.ServerStream
}

type orchestratorStoreAssessmentResultsServer struct {
	grpc.ServerStream
}

func (x *orchestratorStoreAssessmentResultsServer) SendAndClose(m *emptypb.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *orchestratorStoreAssessmentResultsServer) Recv() (*assessment.AssessmentResult, error) {
	m := new(assessment.AssessmentResult)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Orchestrator_ListAssessmentResults_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(assessment.ListAssessmentResultsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).ListAssessmentResults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/ListAssessmentResults",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).ListAssessmentResults(ctx, req.(*assessment.ListAssessmentResultsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_CreateMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).CreateMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/CreateMetric",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).CreateMetric(ctx, req.(*CreateMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_UpdateMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).UpdateMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/UpdateMetric",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).UpdateMetric(ctx, req.(*UpdateMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_GetMetric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).GetMetric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/GetMetric",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).GetMetric(ctx, req.(*GetMetricRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_ListMetrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMetricsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).ListMetrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/ListMetrics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).ListMetrics(ctx, req.(*ListMetricsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_RegisterCloudService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterCloudServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).RegisterCloudService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/RegisterCloudService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).RegisterCloudService(ctx, req.(*RegisterCloudServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_UpdateCloudService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCloudServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).UpdateCloudService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/UpdateCloudService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).UpdateCloudService(ctx, req.(*UpdateCloudServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_GetCloudService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCloudServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).GetCloudService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/GetCloudService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).GetCloudService(ctx, req.(*GetCloudServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_ListCloudServices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCloudServicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).ListCloudServices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/ListCloudServices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).ListCloudServices(ctx, req.(*ListCloudServicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_RemoveCloudService_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveCloudServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).RemoveCloudService(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/RemoveCloudService",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).RemoveCloudService(ctx, req.(*RemoveCloudServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_GetMetricConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMetricConfigurationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).GetMetricConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/GetMetricConfiguration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).GetMetricConfiguration(ctx, req.(*GetMetricConfigurationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orchestrator_ListMetricConfigurations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMetricConfigurationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrchestratorServer).ListMetricConfigurations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/clouditor.Orchestrator/ListMetricConfigurations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrchestratorServer).ListMetricConfigurations(ctx, req.(*ListMetricConfigurationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Orchestrator_ServiceDesc is the grpc.ServiceDesc for Orchestrator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Orchestrator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "clouditor.Orchestrator",
	HandlerType: (*OrchestratorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterAssessmentTool",
			Handler:    _Orchestrator_RegisterAssessmentTool_Handler,
		},
		{
			MethodName: "ListAssessmentTools",
			Handler:    _Orchestrator_ListAssessmentTools_Handler,
		},
		{
			MethodName: "GetAssessmentTool",
			Handler:    _Orchestrator_GetAssessmentTool_Handler,
		},
		{
			MethodName: "UpdateAssessmentTool",
			Handler:    _Orchestrator_UpdateAssessmentTool_Handler,
		},
		{
			MethodName: "DeregisterAssessmentTool",
			Handler:    _Orchestrator_DeregisterAssessmentTool_Handler,
		},
		{
			MethodName: "StoreAssessmentResult",
			Handler:    _Orchestrator_StoreAssessmentResult_Handler,
		},
		{
			MethodName: "ListAssessmentResults",
			Handler:    _Orchestrator_ListAssessmentResults_Handler,
		},
		{
			MethodName: "CreateMetric",
			Handler:    _Orchestrator_CreateMetric_Handler,
		},
		{
			MethodName: "UpdateMetric",
			Handler:    _Orchestrator_UpdateMetric_Handler,
		},
		{
			MethodName: "GetMetric",
			Handler:    _Orchestrator_GetMetric_Handler,
		},
		{
			MethodName: "ListMetrics",
			Handler:    _Orchestrator_ListMetrics_Handler,
		},
		{
			MethodName: "RegisterCloudService",
			Handler:    _Orchestrator_RegisterCloudService_Handler,
		},
		{
			MethodName: "UpdateCloudService",
			Handler:    _Orchestrator_UpdateCloudService_Handler,
		},
		{
			MethodName: "GetCloudService",
			Handler:    _Orchestrator_GetCloudService_Handler,
		},
		{
			MethodName: "ListCloudServices",
			Handler:    _Orchestrator_ListCloudServices_Handler,
		},
		{
			MethodName: "RemoveCloudService",
			Handler:    _Orchestrator_RemoveCloudService_Handler,
		},
		{
			MethodName: "GetMetricConfiguration",
			Handler:    _Orchestrator_GetMetricConfiguration_Handler,
		},
		{
			MethodName: "ListMetricConfigurations",
			Handler:    _Orchestrator_ListMetricConfigurations_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StoreAssessmentResults",
			Handler:       _Orchestrator_StoreAssessmentResults_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "orchestrator.proto",
}
