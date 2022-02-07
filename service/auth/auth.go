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

package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/persistence"
	argon2 "github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

const (
	Issuer = "clouditor"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "auth")
}

// DefaultApiKeyPassword is the default password to protect the API key
const DefaultApiKeyPassword = "changeme"

// DefaultApiKeyPath is the default path for the API private key
const DefaultApiKeyPath = "~/.clouditor/api.key"

// Service is an implementation of the gRPC Authentication service
type Service struct {
	auth.UnimplementedAuthenticationServer

	config struct {
		keyPath     string
		keyPassword string
	}

	apiKey *ecdsa.PrivateKey
}

// UserClaims extend jwt.StandardClaims with more detailed claims about a user
type UserClaims struct {
	jwt.RegisteredClaims
	FullName string `json:"full_name"`
	EMail    string `json:"email"`
}

type ServiceOption func(s *Service)

func WithApiKeyPath(path string) ServiceOption {
	return func(s *Service) {
		s.config.keyPath = path
	}
}

func WithApiKeyPassword(password string) ServiceOption {
	return func(s *Service) {
		s.config.keyPassword = password
	}
}

func NewService(opts ...ServiceOption) *Service {
	s := &Service{
		config: struct {
			keyPath     string
			keyPassword string
		}{
			keyPath:     DefaultApiKeyPath,
			keyPassword: DefaultApiKeyPassword,
		},
	}

	var (
		err error
	)

	// Apply options
	for _, o := range opts {
		o(s)
	}

	s.apiKey, err = s.loadApiKey(s.config.keyPath, []byte(s.config.keyPassword))

	// We treat different errors diffently. For example if the path is the default path
	// and the error is os.ErrNotExist, this could be just the first start of Clouditor.
	// So we only treat this as an information that we will create a new key.
	//
	// If the user specifies a custom path and this one does not exist, we will report an error
	// here.
	if err != nil {
		s.apiKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

		if s.config.keyPath == DefaultApiKeyPath && errors.Is(err, os.ErrNotExist) {
			log.Infof("API key does not exist at the default location yet. We will create a new one")

			// Also save the key in this case, so we can load it next time
			err = s.saveApiKey(s.config.keyPath, s.config.keyPassword)

			// Error while error handling, meh
			if err != nil {
				log.Errorf("Error while saving the new API key: %v", err)
			}
		} else if err != nil {
			log.Errorf("Could not load key from file, continuing with a temporary key: %v", err)
		}
	}

	// Lets clear out some sensitive things for slightly more security
	s.config.keyPassword = ""

	return s
}

func (s Service) loadApiKey(path string, password []byte) (key *ecdsa.PrivateKey, err error) {
	var (
		keyFile string
	)

	keyFile, err = expandPath(path)
	if err != nil {
		return nil, fmt.Errorf("error while expanding path: %w", err)
	}

	if _, err = os.Stat(keyFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist (yet): %w", err)
	}

	// Check, if we already have a persisted API key
	f, err := os.OpenFile(keyFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("error while opening the file: %w", err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error while reading file content: %w", err)
	}

	key, err = ParseECPrivateKeyFromPEMWithPassword(data, password)
	if err != nil {
		return nil, fmt.Errorf("error while parsing private key: %w", err)
	}

	return key, nil
}

func (s *Service) saveApiKey(keyPath string, password string) (err error) {
	keyPath, err = expandPath(keyPath)
	if err != nil {
		return fmt.Errorf("error while expanding path: %w", err)
	}

	// Check, if we already have a persisted API key
	f, err := os.OpenFile(keyPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("error while opening the file: %w", err)
	}
	defer f.Close()

	data, err := MarshalECPrivateKeyWithPassword(s.apiKey, []byte(password))
	if err != nil {
		return fmt.Errorf("error while marshalling private key: %w", err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("error while writing file content: %w", err)
	}

	return nil
}

// Login handles a login request
func (s Service) Login(_ context.Context, request *auth.LoginRequest) (response *auth.LoginResponse, err error) {
	var result bool
	var user *auth.User

	if result, user, err = verifyLogin(request); err != nil {
		// a returned error means, something has gone wrong
		return nil, status.Errorf(codes.Internal, "error during login: %v", err)
	}

	if !result {
		// authentication error
		return nil, status.Errorf(codes.Unauthenticated, "login failed")
	}

	var token string

	var expiry = time.Now().Add(time.Hour * 24)

	if token, err = s.issueToken(user.Username, user.FullName, user.Email, expiry); err != nil {
		return nil, status.Errorf(codes.Internal, "token issue failed: %v", err)
	}

	response = &auth.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		Expiry:      timestamppb.New(expiry),
	}

	return response, nil
}

func (s Service) ListPublicKeys(_ context.Context, _ *auth.ListPublicKeysRequest) (response *auth.ListPublicResponse, err error) {
	response = &auth.ListPublicResponse{
		Keys: []*auth.JsonWebKey{
			{
				Kid: "1",
				Kty: "EC",
				Crv: s.apiKey.Params().Name,
				X:   base64.RawURLEncoding.EncodeToString(s.apiKey.X.Bytes()),
				Y:   base64.RawURLEncoding.EncodeToString(s.apiKey.Y.Bytes()),
			},
		},
	}

	return response, nil
}

// verifyLogin compares the credentials supplied by request with the ones stored in the database.
// It will return an error in err only if a database error or something else is occurred, this should be
// returned to the user as an internal server error. For security reasons, if authentication failed, only
// the result will be set to false, but no detailed error will be returned to the user.
func verifyLogin(request *auth.LoginRequest) (result bool, user *auth.User, err error) {
	var match bool

	db := persistence.GetDatabase()

	user = new(auth.User)

	err = db.Where("username = ?", request.Username).First(user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// user not found, set result to false, but hide the error
		return false, nil, nil
	} else if err != nil {
		// a database connection error has occurred, return it
		return false, nil, err
	}

	match, err = argon2.ComparePasswordAndHash(request.Password, user.Password)

	if err != nil {
		// some other error occurred, return it
		return false, nil, err
	}

	if match {
		return true, user, nil
	} else {
		return false, nil, nil
	}
}

// hashPassword returns a hash of password using argon2id.
func hashPassword(password string) (string, error) {
	return argon2.CreateHash(password, &argon2.Params{
		SaltLength:  16,
		Memory:      65536,
		KeyLength:   32,
		Iterations:  3,
		Parallelism: 6,
	})
}

// CreateDefaultUser creates a default user in the database
func (Service) CreateDefaultUser(username string, password string) {
	db := persistence.GetDatabase()

	var count int64
	db.Model(&auth.User{}).Count(&count)

	if count == 0 {
		hash, _ := hashPassword(password)

		user := auth.User{
			Username: username,
			FullName: username,
			Password: string(hash),
		}

		log.Infof("Creating default user %s\n", user.Username)

		db.Create(&user)
	}
}

// issueToken issues a JWT-based token
func (s Service) issueToken(subject string, fullName string, email string, expiry time.Time) (token string, err error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodES256,
		&UserClaims{
			FullName: fullName,
			EMail:    email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiry),
				Issuer:    Issuer,
				Subject:   subject,
			}},
	)
	claims.Header["kid"] = "1"

	token, err = claims.SignedString(s.apiKey)
	return
}

func (s Service) GetPublicKey() *ecdsa.PublicKey {
	return &s.apiKey.PublicKey
}

// AuthFuncOverride implements the ServiceAuthFuncOverride interface to override the AuthFunc for this service.
// The reason is to actually disable authentication checking in the auth service, because functions such as login
// need to be publically available.
func (Service) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	// No authentication needed for login functions, otherwise we could not login
	return ctx, nil
}

// expandPath expands a path that possible contains a tilde (~) character into the home directory
// of the user
func expandPath(path string) (out string, err error) {
	var (
		usr *user.User
	)

	// Fetch the current user home directory
	usr, err = user.Current()
	if err != nil {
		return path, fmt.Errorf("could not find retrieve current user: %w", err)
	}

	if path == "~" {
		return usr.HomeDir, nil
	} else if strings.HasPrefix(path, "~") {
		// We only allow ~ at the beginning of the path
		return filepath.Join(usr.HomeDir, path[2:]), nil
	}

	return path, nil
}
