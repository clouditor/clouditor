package assessment

import (
	"context"

	"clouditor.io/clouditor/api/assessment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

var standalone *Service

type standaloneEvidenceStream struct {
	serverChannel chan *assessment.Evidence
	clientChannel chan *emptypb.Empty
	ctx           context.Context
}

type standaloneEvidenceClient struct{}

func (s standaloneEvidenceStream) SendAndClose(*emptypb.Empty) error {
	s.clientChannel <- &emptypb.Empty{}

	return nil
}

func (s standaloneEvidenceStream) CloseAndRecv() (*emptypb.Empty, error) {
	empty := <-s.clientChannel

	return empty, nil
}

func (s standaloneEvidenceStream) CloseSend() error {
	return nil
}

func (s standaloneEvidenceStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (s standaloneEvidenceStream) Trailer() metadata.MD {
	return nil
}

func (s standaloneEvidenceStream) Context() context.Context {
	return nil
}

func (s standaloneEvidenceStream) Send(evidence *assessment.Evidence) error {
	s.serverChannel <- evidence

	return nil
}

func (s standaloneEvidenceStream) SendHeader(metadata.MD) error {
	return nil
}

func (s standaloneEvidenceStream) SetHeader(metadata.MD) error {
	return nil
}

func (s standaloneEvidenceStream) SetTrailer(metadata.MD) {

}

func (s standaloneEvidenceStream) Recv() (*assessment.Evidence, error) {
	evidence := <-s.serverChannel

	return evidence, nil
}

func (s standaloneEvidenceStream) RecvMsg(m interface{}) error {
	return nil
}

func (s standaloneEvidenceStream) SendMsg(m interface{}) error {
	return nil
}

func (s standaloneEvidenceClient) TriggerAssessment(ctx context.Context, in *assessment.TriggerAssessmentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return standalone.TriggerAssessment(ctx, in)
}

func (s standaloneEvidenceClient) StoreEvidence(ctx context.Context, in *assessment.StoreEvidenceRequest, opts ...grpc.CallOption) (*assessment.Evidence, error) {
	return standalone.StoreEvidence(ctx, in)
}

func (s standaloneEvidenceClient) StreamEvidences(ctx context.Context, opts ...grpc.CallOption) (assessment.Assessment_StreamEvidencesClient, error) {
	var stream = &standaloneEvidenceStream{
		serverChannel: make(chan *assessment.Evidence),
		clientChannel: make(chan *emptypb.Empty),
		ctx:           context.Background(),
	}

	go standalone.StreamEvidences(stream)

	return stream, nil
}

func NewInMemoryClient() assessment.AssessmentClient {
	return &standaloneEvidenceClient{}
}

func StandaloneService() *Service {
	standalone = &Service{}

	return standalone
}
