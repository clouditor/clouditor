// Copyright 2016-2020 Fraunhofer AISEC
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

package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/api/orchestrator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var DefaultSessionFolder string
var Output io.Writer = os.Stdout

type Session struct {
	URL   string `json:"url"`
	Token string `json:"token"`

	Folder string `json:"-"`

	*grpc.ClientConn
}

func init() {
	var home string
	var err error

	// find the home directory
	if home, err = os.UserHomeDir(); err != nil {
		return
	}

	DefaultSessionFolder = fmt.Sprintf("%s/.clouditor/", home)
}

func NewSession(url string, opts ...grpc.DialOption) (session *Session, err error) {
	session = &Session{
		URL:    url,
		Folder: viper.GetString("session-directory"),
	}

	if len(opts) == 0 {
		// TODO(oxisto): set flag depending on target url, insecure only for localhost
		opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}

	if session.ClientConn, err = grpc.Dial(session.URL, opts...); err != nil {
		return nil, fmt.Errorf("could not connect: %w", err)
	}

	return session, nil
}

func ContinueSession() (session *Session, err error) {
	var (
		file   *os.File
		folder string
	)

	folder = viper.GetString("session-directory")

	// try to read from session.json
	if file, err = os.OpenFile(fmt.Sprintf("%s/session.json", folder), os.O_RDONLY, 0600); err != nil {
		return
	}

	defer file.Close()

	session = new(Session)
	session.Folder = folder

	if err = json.NewDecoder(file).Decode(&session); err != nil {
		return
	}

	if session.ClientConn, err = grpc.Dial(session.URL, grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
		return nil, fmt.Errorf("could not connect: %w", err)
	}

	return session, nil
}

// Save saves the session into the `.clouditor` folder in the home directory
func (s *Session) Save() (err error) {
	var (
		file *os.File
	)

	// create the session directory
	if err = os.MkdirAll(s.Folder, 0744); err != nil {
		return fmt.Errorf("could not create .clouditor in home directory: %w", err)
	}

	if file, err = os.OpenFile(fmt.Sprintf("%s/session.json", s.Folder), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
		return fmt.Errorf("could not save session.json: %w", err)
	}

	defer file.Close()

	if err = json.NewEncoder(file).Encode(s); err != nil {
		return fmt.Errorf("could not serialize JSON: %w", err)
	}

	return nil
}

// HandleResponse handles the response and error message of an gRPC call
func (s *Session) HandleResponse(msg proto.Message, err error) error {
	if err != nil {
		// check, if it is a gRPC error
		s, ok := status.FromError(err)

		// otherwise, forward the error message
		if !ok {
			return err
		}

		// create a new error with just the message
		return errors.New(s.Message())
	}

	opt := protojson.MarshalOptions{
		Multiline: true,
		Indent:    "  ",
	}

	b, _ := opt.Marshal(msg)

	fmt.Fprintf(Output, "%s\n", string(b))

	return err
}

func PromtForLogin() (loginRequest *auth.LoginRequest, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')

	if err != nil {
		return
	}

	fmt.Print("Enter password: ")
	password, err := reader.ReadString('\n')

	if err != nil {
		return
	}

	loginRequest = &auth.LoginRequest{
		Username: strings.Trim(username, "\n"),
		Password: strings.Trim(password, "\n"),
	}

	return loginRequest, nil
}

// Invoke implements `grpc.ClientConnInterface` and automatically provides an authenticated
// context of this session
func (s *Session) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return s.ClientConn.Invoke(s.AuthenticatedContext(ctx), method, args, reply, opts...)
}

// NewStream implements `grpc.ClientConnInterface` and automatically provides an authenticated
// context of this session
func (s *Session) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return s.ClientConn.NewStream(s.AuthenticatedContext(ctx), desc, method, opts...)
}

func (s *Session) AuthenticatedContext(ctx context.Context) context.Context {
	return metadata.NewOutgoingContext(ctx,
		metadata.Pairs(
			"Authorization",
			fmt.Sprintf("Bearer %s",
				s.Token),
		))
}

func DefaultArgsShellComp(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{}, cobra.ShellCompDirectiveNoFileComp
}

func ValidArgsGetTools(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return getTools(toComplete), cobra.ShellCompDirectiveNoFileComp
}

func ValidArgsGetMetrics(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return getMetrics(toComplete), cobra.ShellCompDirectiveNoFileComp
}

func ValidArgsGetCloudServices(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return getCloudServices(toComplete), cobra.ShellCompDirectiveNoFileComp
}

func getTools(_ string) []string {
	var (
		err     error
		session *Session
		client  orchestrator.OrchestratorClient
		res     *orchestrator.ListAssessmentToolsResponse
	)

	if session, err = ContinueSession(); err != nil {
		fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
		return nil
	}

	client = orchestrator.NewOrchestratorClient(session)

	if res, err = client.ListAssessmentTools(context.Background(), &orchestrator.ListAssessmentToolsRequest{}); err != nil {
		return []string{}
	}

	var tools []string
	for _, v := range res.Tools {
		tools = append(tools, fmt.Sprintf("%s\t%s: %s", v.Id, v.Name, v.Description))
	}

	return tools
}

// TODO(oxisto): This could be an interesting use case for 1.18 Go generics
func getMetrics(_ string) []string {
	var (
		err     error
		session *Session
		client  orchestrator.OrchestratorClient
		res     *orchestrator.ListMetricsResponse
	)

	if session, err = ContinueSession(); err != nil {
		fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
		return nil
	}

	client = orchestrator.NewOrchestratorClient(session)

	if res, err = client.ListMetrics(context.Background(), &orchestrator.ListMetricsRequest{}); err != nil {
		return []string{}
	}

	var metrics []string
	for _, v := range res.Metrics {
		metrics = append(metrics, fmt.Sprintf("%s\t%s: %s", v.Id, v.Name, v.Description))
	}

	return metrics
}

func getCloudServices(_ string) []string {
	var (
		err     error
		session *Session
		client  orchestrator.OrchestratorClient
		res     *orchestrator.ListCloudServicesResponse
	)

	if session, err = ContinueSession(); err != nil {
		fmt.Printf("Error while retrieving the session. Please re-authenticate.\n")
		return nil
	}

	client = orchestrator.NewOrchestratorClient(session)

	if res, err = client.ListCloudServices(context.Background(), &orchestrator.ListCloudServicesRequest{}); err != nil {
		return []string{}
	}

	var metrics []string
	for _, v := range res.Services {
		metrics = append(metrics, fmt.Sprintf("%s\t%s: %s", v.Id, v.Name, v.Description))
	}

	return metrics
}
