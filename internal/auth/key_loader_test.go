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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_keyLoader_recoverFromLoadApiKeyError(t *testing.T) {
	var tmpFile, _ = ioutil.TempFile("", "api.key")
	// Close it immediately , since we want to write to it
	tmpFile.Close()

	defer func() {
		os.Remove(tmpFile.Name())
	}()

	type fields struct {
		saveOnCreate bool
		path         string
		password     string
	}
	type args struct {
		err         error
		defaultPath bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantKey assert.ValueAssertionFunc
	}{
		{
			name: "Could not load key from custom path",
			fields: fields{
				saveOnCreate: false,
				path:         "doesnotexist",
				password:     "test",
			},
			args: args{
				err:         os.ErrNotExist,
				defaultPath: false,
			},
			wantKey: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				// A temporary key should be created
				return assert.NotNil(tt, i1.(*ecdsa.PrivateKey))
			},
		},
		{
			name: "Could not load key from default path and save it",
			fields: fields{
				saveOnCreate: true,
				path:         tmpFile.Name(),
				password:     "test",
			},
			args: args{
				err:         os.ErrNotExist,
				defaultPath: true,
			},
			wantKey: func(tt assert.TestingT, i1 interface{}, i2 ...interface{}) bool {
				// A temporary key should be created
				if !assert.NotNil(tt, i1.(*ecdsa.PrivateKey)) {
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
			s := &keyLoader{
				path:         tt.fields.path,
				password:     tt.fields.password,
				saveOnCreate: tt.fields.saveOnCreate,
			}
			gotKey := s.recoverFromLoadApiKeyError(tt.args.err, tt.args.defaultPath)

			if tt.wantKey != nil {
				tt.wantKey(t, gotKey, tt.args.err, tt.args.defaultPath)
			}
		})
	}
}

func TestService_loadKeyFromFile(t *testing.T) {
	// Prepare a tmp file that contains a new temporary private key
	var tmpFile, _ = ioutil.TempFile("", "api.key")
	tmpFile.Close()

	// Create a new temporary key
	tmpKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	defer func() {
		os.Remove(tmpFile.Name())
	}()

	// Save a key to it
	err := saveKeyToFile(tmpKey, tmpFile.Name(), "tmp")
	assert.NoError(t, err)

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
			gotKey, err := loadKeyFromFile(tt.args.path, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.loadKeyFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotKey, tt.wantKey) {
				t.Errorf("Service.loadKeyFromFile() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

// mockStorage is a mocked persistence.Storage implementation that returns errors at the specified
// operations.
//
// TODO(lebogg): Extract this struct into our new internal/testutils package
type mockStorage struct {
	createError error
	saveError   error
	updateError error
	getError    error
	listError   error
	countError  error
	deleteError error
}

func (m mockStorage) Create(interface{}) error { return m.createError }

func (m mockStorage) Save(interface{}, ...interface{}) error { return m.saveError }

func (m mockStorage) Update(interface{}, interface{}, ...interface{}) error {
	return m.updateError
}

func (m mockStorage) Get(interface{}, ...interface{}) error { return m.getError }

func (m mockStorage) List(interface{}, ...interface{}) error { return m.listError }

func (m mockStorage) Count(interface{}, ...interface{}) (int64, error) {
	return 0, m.countError
}

func (m mockStorage) Delete(interface{}, ...interface{}) error { return m.deleteError }
