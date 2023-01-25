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

package voc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticityInterface(t *testing.T) {
	tests := []struct {
		isAuthenticity IsAuthenticity
		typ            string
	}{
		{
			isAuthenticity: &NoAuthentication{},
			typ:            "NoAuthentication",
		},
		{
			isAuthenticity: &SingleSignOn{},
			typ:            "SingleSignOn",
		},
		{
			isAuthenticity: &OTPBasedAuthentication{},
			typ:            "OTPBasedAuthentication",
		},
		{
			isAuthenticity: &PasswordBasedAuthentication{},
			typ:            "PasswordBasedAuthentication",
		},
		{
			isAuthenticity: &TokenBasedAuthentication{},
			typ:            "TokenBasedAuthentication",
		},
	}

	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			var isa = tt.isAuthenticity
			assert.Equal(t, tt.typ, isa.Type())
		})
	}
}

func TestAuthorizationInterface(t *testing.T) {
	tests := []struct {
		isAuthorization IsAuthorization
		typ             string
	}{
		{
			isAuthorization: &AccessRestriction{},
			typ:             "AccessRestriction",
		},
		{
			isAuthorization: &ABAC{},
			typ:             "ABAC",
		},
		{
			isAuthorization: &RBAC{},
			typ:             "RBAC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.typ, func(t *testing.T) {
			var isa = tt.isAuthorization
			assert.Equal(t, tt.typ, isa.Type())
		})
	}
}
