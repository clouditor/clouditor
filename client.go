/*
 * Copyright 2016-2020 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *           $$\                           $$\ $$\   $$\
 *           $$ |                          $$ |\__|  $$ |
 *  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 * $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 * $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
 * $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 * \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package clouditor

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Client struct {
	URL string

	httpClient *http.Client
	token      string
}

func (c *Client) continueSession() (err error) {
	var (
		home   string
		file   *os.File
		result map[string]interface{}
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

	if err = json.NewDecoder(file).Decode(&result); err != nil {
		return
	}

	if token, ok := result["token"].(string); ok {
		// set this client's token
		c.token = token
	}

	return nil
}

func NewClient(url string) *Client {
	var client = &Client{
		URL: url,
	}
	client.httpClient = &http.Client{}
	client.continueSession()

	return client
}

func (c Client) Authenticate() (err error) {
	var home string
	var file *os.File
	var result map[string]interface{}

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

	loginRequest := map[string]string{
		"username": strings.Trim(username, "\n"),
		"password": strings.Trim(password, "\n"),
	}

	b, err := json.Marshal(loginRequest)

	if err != nil {
		return
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("%s/v1/auth/login", c.URL), "application/json", bytes.NewBuffer(b))

	if err != nil {
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	// TODO: actually check if resp is status 200
	fmt.Printf("%+v\n", result["token"])

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

	if err = json.NewEncoder(file).Encode(result); err != nil {
		return
	}

	fmt.Println("Successfully logged in")

	return nil
}

func (c Client) StartDiscovery() (err error) {
	var (
		b      []byte
		result map[string]interface{}
	)

	request := map[string]string{}

	if b, err = json.Marshal(request); err != nil {
		return fmt.Errorf("could not serialize JSON: %w", err)
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("%s/v1/discovery/start", c.URL), "application/json", bytes.NewBuffer(b))

	if err != nil {
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	fmt.Printf("Response: %+v", resp)

	return
}
