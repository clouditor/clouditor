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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"clouditor.io/clouditor/api"
	"clouditor.io/clouditor/api/orchestrator"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var DefaultSessionFolder string

const SessionFolderFlag = "session-directory"

var Output io.Writer = os.Stdout

type Session struct {
	*grpc.ClientConn
	*oauth2.Config

	authorizer api.Authorizer

	// URL is the URL of the gRPC server to connect to
	URL string `json:"url"`

	Folder string `json:"-"`

	// dirty flags that we need to fetch a new token and save the session again
	dirty bool
}

func (s *Session) SetAuthorizer(authorizer api.Authorizer) {
	s.authorizer = authorizer
}

func (s *Session) Authorizer() api.Authorizer {
	return s.authorizer
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

func NewSession(url string, config *oauth2.Config, token *oauth2.Token) (session *Session, err error) {
	session = &Session{
		URL:        url,
		Folder:     viper.GetString(SessionFolderFlag),
		authorizer: api.NewOAuthAuthorizerFromConfig(config, token),
		Config:     config,
	}

	var opts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

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

	folder = viper.GetString(SessionFolderFlag)

	// try to read from session.json
	if file, err = os.OpenFile(fmt.Sprintf("%s/session.json", folder), os.O_RDONLY, 0600); err != nil {
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	session = new(Session)
	session.Folder = folder

	if err = json.NewDecoder(file).Decode(&session); err != nil {
		return nil, fmt.Errorf("could not parse session file: %w", err)
	}

	// If we detect that this session is "dirty", try to save it again
	if session.dirty {
		_ = session.Save()
	}

	if session.ClientConn, err = grpc.Dial(session.URL, api.DefaultGrpcDialOptions(session)...); err != nil {
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

	defer func(file *os.File) {
		err = file.Close()
	}(file)

	if err = json.NewEncoder(file).Encode(&s); err != nil {
		return fmt.Errorf("could not serialize JSON: %w", err)
	}

	s.dirty = false

	return nil
}

// MarshalJSON is custom JSON marshalling implementation that gives us more control over
// the fields we want to serialize. The core problem is that we want to serialize our token, e.g.
// to store a session state for our clients, but we do not want to export the token field in our
// struct. Exporting the field would create problems in multi-threaded environments. Therefore,
// access is only allowed through the Token() function, which keeps the token synchronized using a mutex.
func (s *Session) MarshalJSON() ([]byte, error) {
	token, _ := s.authorizer.Token()

	return json.Marshal(&struct {
		URL    string         `json:"url"`
		Token  *oauth2.Token  `json:"token"`
		Config *oauth2.Config `json:"oauth2"`
	}{
		URL:    s.URL,
		Token:  token,
		Config: s.Config,
	})
}

// UnmarshalJSON is custom JSON marshalling implementation that gives us more control over
// the fields we want to deserialize. See MarshalJSON for a detailed explanation, why this is
// necessary.
func (s *Session) UnmarshalJSON(data []byte) (err error) {
	v := struct {
		URL    string         `json:"url"`
		Token  *oauth2.Token  `json:"token"`
		Config *oauth2.Config `json:"oauth2"`
	}{}

	if err = json.Unmarshal(data, &v); err != nil {
		return
	}

	s.URL = v.URL

	// Mark this session as dirty, if the token is not valid (anymore)
	s.dirty = !v.Token.Valid()
	s.Config = v.Config

	// Check, if oauth config is missing or invalid
	if s.Config == nil {
		return errors.New("missing oauth2 config")
	}

	s.authorizer = api.NewOAuthAuthorizerFromConfig(s.Config, v.Token)
	return
}

// HandleResponse handles the response and error message of an gRPC call
func (*Session) HandleResponse(msg proto.Message, err error) error {
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
		Multiline:       true,
		Indent:          "  ",
		EmitUnpopulated: true,
	}

	b, _ := opt.Marshal(msg)

	_, err = fmt.Fprintf(Output, "%s\n", string(b))

	return err
}

func DefaultArgsShellComp(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
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
