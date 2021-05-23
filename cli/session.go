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
	"os"
	"strings"

	"clouditor.io/clouditor/api/auth"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Session struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func NewSession(url string) *Session {
	return &Session{
		URL: url,
	}
}

func ContinueSession() (session *Session, err error) {
	var (
		home string
		file *os.File
	)

	// try to read from session.json

	// find the home directory
	if home, err = os.UserHomeDir(); err != nil {
		return
	}

	if file, err = os.OpenFile(fmt.Sprintf("%s/.clouditor/session.json", home), os.O_RDONLY, 0600); err != nil {
		return
	}

	defer file.Close()

	session = new(Session)

	if err = json.NewDecoder(file).Decode(&session); err != nil {
		return
	}

	return session, nil
}

func (s *Session) Context() context.Context {
	md := metadata.Pairs("Authorization", fmt.Sprintf("Bearer %s", s.Token))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	return ctx
}

// Save saves the session into the `.clouditor` folder in the home directory
func (s *Session) Save() {
	var (
		err  error
		home string
		file *os.File
	)
	// find the home directory
	if home, err = os.UserHomeDir(); err != nil {
		return
	}

	// create the .clouditor directory
	os.MkdirAll(fmt.Sprintf("%s/.clouditor", home), 0744)

	if file, err = os.OpenFile(fmt.Sprintf("%s/.clouditor/session.json", home), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600); err != nil {
		return
	}

	defer file.Close()

	if err = json.NewEncoder(file).Encode(s); err != nil {
		return
	}
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

	fmt.Printf("%s\n", string(b))

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
