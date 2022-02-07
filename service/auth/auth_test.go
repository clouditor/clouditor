// Copyright 2016-2022 Fraunhofer AISEC
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

package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io/ioutil"
	"math/big"
	"os"
	"reflect"
	"testing"

	"clouditor.io/clouditor/api/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestService_ListPublicKeys(t *testing.T) {
	type fields struct {
		apiKey *ecdsa.PrivateKey
	}
	type args struct {
		in0 context.Context
		in1 *auth.ListPublicKeysRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantResponse *auth.ListPublicResponse
		wantErr      bool
	}{
		{
			name: "List single public key",
			fields: fields{
				apiKey: &ecdsa.PrivateKey{
					PublicKey: ecdsa.PublicKey{
						Curve: elliptic.P256(),
						X:     big.NewInt(1),
						Y:     big.NewInt(2),
					},
				},
			},
			args: args{
				in0: context.TODO(),
				in1: &auth.ListPublicKeysRequest{},
			},
			wantResponse: &auth.ListPublicResponse{
				Keys: []*auth.JsonWebKey{
					{
						Kid: "1",
						Kty: "EC",
						Crv: "P-256",
						X:   "AQ",
						Y:   "Ag",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				apiKey: tt.fields.apiKey,
			}
			gotResponse, err := s.ListPublicKeys(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListPublicKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !proto.Equal(gotResponse, tt.wantResponse) {
				t.Errorf("Service.ListPublicKeys() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestService_recoverFromLoadApiKeyError(t *testing.T) {
	var tmpFile, _ = ioutil.TempFile("", "api.key")
	// Close it immediately , since we want to write to it
	tmpFile.Close()

	defer func() {
		os.Remove(tmpFile.Name())
	}()

	type fields struct {
		config struct {
			keySaveOnCreate bool
			keyPath         string
			keyPassword     string
		}
		apiKey *ecdsa.PrivateKey
	}
	type args struct {
		err         error
		defaultPath bool
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantService assert.ValueAssertionFunc
	}{
		{
			name: "Could not load key from custom path",
			fields: fields{
				config: struct {
					keySaveOnCreate bool
					keyPath         string
					keyPassword     string
				}{
					keySaveOnCreate: false,
					keyPath:         "doesnotexist",
					keyPassword:     "test",
				},
			},
			args: args{
				err:         os.ErrNotExist,
				defaultPath: false,
			},
			wantService: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				// A temporary key should be created
				return assert.NotNil(tt, i1.(*Service).apiKey)
			},
		},
		{
			name: "Could not load key from default path and save it",
			fields: fields{
				config: struct {
					keySaveOnCreate bool
					keyPath         string
					keyPassword     string
				}{
					keySaveOnCreate: true,
					keyPath:         tmpFile.Name(),
					keyPassword:     "test",
				},
			},
			args: args{
				err:         os.ErrNotExist,
				defaultPath: true,
			},
			wantService: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				// A temporary key should be created
				if !assert.NotNil(tt, i1.(*Service).apiKey) {
					return false
				}

				f, err := os.OpenFile(tmpFile.Name(), os.O_RDONLY, 0600)
				if !assert.ErrorIs(tt, err, nil) {
					return false
				}

				// Our tmp file should also contain something now
				data, err := ioutil.ReadAll(f)
				if !assert.ErrorIs(tt, err, nil) {
					return false
				}

				return assert.True(tt, len(data) > 0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				config: tt.fields.config,
				apiKey: tt.fields.apiKey,
			}
			s.recoverFromLoadApiKeyError(tt.args.err, tt.args.defaultPath)

			if tt.wantService != nil {
				tt.wantService(t, s, tt.args.err, tt.args.defaultPath)
			}
		})
	}
}

func TestService_loadApiKey(t *testing.T) {
	// Prepare a tmp file that contains a new temporary private key
	var tmpFile, _ = ioutil.TempFile("", "api.key")
	tmpFile.Close()

	// Create a new temporary key
	tmpKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	defer func() {
		os.Remove(tmpFile.Name())
	}()

	// Save a key to it
	err := saveApiKey(tmpKey, tmpFile.Name(), "tmp")
	assert.ErrorIs(t, err, nil)

	type args struct {
		path     string
		password []byte
	}
	tests := []struct {
		name    string
		args    args
		wantKey *ecdsa.PrivateKey
		wantErr bool
	}{
		{
			name: "Load existing key",
			args: args{
				path:     tmpFile.Name(),
				password: []byte("tmp"),
			},
			wantKey: tmpKey,
			wantErr: false,
		},
		{
			name: "Load existing key with wrong password",
			args: args{
				path:     tmpFile.Name(),
				password: []byte("notpassword"),
			},
			wantKey: nil,
			wantErr: true,
		},
		{
			name: "Load not existing key",
			args: args{
				path:     "notexists",
				password: []byte("tmp"),
			},
			wantKey: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, err := loadApiKey(tt.args.path, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.loadApiKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotKey, tt.wantKey) {
				t.Errorf("Service.loadApiKey() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
