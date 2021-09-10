// Copyright 2021 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package standalone

import (
	"context"

	"clouditor.io/clouditor/api/assessment"
	service_assessment "clouditor.io/clouditor/service/assessment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

var assessmentService assessment.AssessmentServer

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
	return assessmentService.TriggerAssessment(ctx, in)
}

func (s standaloneEvidenceClient) AssessEvidence(ctx context.Context, in *assessment.AssessEvidenceRequest, opts ...grpc.CallOption) (*assessment.AssessEvidenceResponse, error) {
	return assessmentService.AssessEvidence(ctx, in)
}

func (s standaloneEvidenceClient) ListAssessmentResults(ctx context.Context, in *assessment.ListAssessmentResultsRequest, opts ...grpc.CallOption) (*assessment.ListAssessmentResultsResponse, error) {
	return assessmentService.ListAssessmentResults(ctx, in)
}

func (s standaloneEvidenceClient) AssessEvidences(ctx context.Context, opts ...grpc.CallOption) (assessment.Assessment_AssessEvidencesClient, error) {
	var stream = &standaloneEvidenceStream{
		serverChannel: make(chan *assessment.Evidence),
		clientChannel: make(chan *emptypb.Empty),
		ctx:           context.Background(),
	}

	go func() {
		_ = assessmentService.AssessEvidences(stream)
	}()

	return stream, nil
}

func NewAssessmentClient() assessment.AssessmentClient {
	return &standaloneEvidenceClient{}
}

func NewAssessmentServer() assessment.AssessmentServer {
	assessmentService = service_assessment.NewService()

	return assessmentService
}
