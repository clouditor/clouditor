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
	"clouditor.io/clouditor/service"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"io/ioutil"
	"net"
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
	"google.golang.org/grpc"
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

const (
	// DefaultApiKeySaveOnCreate specifies whether a created API key will be saved. This is useful to turn of in unit tests, where
	// we only want a temporary key.
	DefaultApiKeySaveOnCreate = true

	// DefaultApiKeyPassword is the default password to protect the API key
	DefaultApiKeyPassword = "changeme"

	// DefaultApiKeyPath is the default path for the API private key
	DefaultApiKeyPath = "~/.clouditor/api.key"

	// DefaultKeyID specifies the default Key ID used in the JWKS of the authentication service
	DefaultKeyID = "1"
)

// Service is an implementation of the gRPC Authentication service
type Service struct {
	auth.UnimplementedAuthenticationServer

	config struct {
		keySaveOnCreate bool
		keyPath         string
		keyPassword     string
	}

	apiKey *ecdsa.PrivateKey

	db persistence.IsDatabase
}

// UserClaims extend jwt.StandardClaims with more detailed claims about a user
type UserClaims struct {
	jwt.RegisteredClaims
	FullName string `json:"full_name"`
	EMail    string `json:"email"`
}

// ServiceOption is a functional option type to configure the authentication service.
type ServiceOption func(s *Service)

// WithApiKeyPath is an option to configure the path for the API key.
func WithApiKeyPath(path string) ServiceOption {
	return func(s *Service) {
		s.config.keyPath = path
	}
}

// WithApiKeyPassword is an option to configure the password that is used to protect the API key.
func WithApiKeyPassword(password string) ServiceOption {
	return func(s *Service) {
		s.config.keyPassword = password
	}
}

// WithApiKeySaveOnCreate is an option to configure whether a key should be saved when it is created.
func WithApiKeySaveOnCreate(saveOnCreate bool) ServiceOption {
	return func(s *Service) {
		s.config.keySaveOnCreate = saveOnCreate
	}
}

// NewService creates a new Service representing an authentication service.
func NewService(db persistence.IsDatabase, opts ...ServiceOption) *Service {
	var err error

	s := &Service{
		config: struct {
			keySaveOnCreate bool
			keyPath         string
			keyPassword     string
		}{
			keySaveOnCreate: DefaultApiKeySaveOnCreate,
			keyPath:         DefaultApiKeyPath,
			keyPassword:     DefaultApiKeyPassword,
		},
		db: db,
	}

	// Apply options
	for _, o := range opts {
		o(s)
	}

	// Load the key
	s.apiKey, err = loadApiKey(s.config.keyPath, []byte(s.config.keyPassword))

	// Recover from an error, so that we have at least a temporary key
	if err != nil {
		s.recoverFromLoadApiKeyError(err, s.config.keyPath == DefaultApiKeyPath)
	}

	// Lets clear out some sensitive things for slightly more security
	s.config.keyPassword = ""

	return s
}

// loadApiKey loads an ecdsa.PrivateKey from a path. The key must in PEM format and protected by
// a password using PKCS#8 with PBES2.
func loadApiKey(path string, password []byte) (key *ecdsa.PrivateKey, err error) {
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

// saveApiKey saves an ecdsa.PrivateKey to a path. The key will be saved in PEM format and protected by
// a password using PKCS#8 with PBES2.
func saveApiKey(apiKey *ecdsa.PrivateKey, keyPath string, password string) (err error) {
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

	data, err := MarshalECPrivateKeyWithPassword(apiKey, []byte(password))
	if err != nil {
		return fmt.Errorf("error while marshalling private key: %w", err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("error while writing file content: %w", err)
	}

	return nil
}

// recoverFromLoadApiKeyError tries to recover from an error during key loading. We treat different errors diffently.
// For example if the path is the default path and the error is os.ErrNotExist, this could be just the first start of Clouditor.
// So we only treat this as an information that we will create a new key, which we will save, based on the config.
//
// If the user specifies a custom path and this one does not exist, we will report an error
// here.
func (s *Service) recoverFromLoadApiKeyError(err error, defaultPath bool) {
	// In any case, create a new temporary API key
	s.apiKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if defaultPath && errors.Is(err, os.ErrNotExist) {
		log.Infof("API key does not exist at the default location yet. We will create a new one")

		if s.config.keySaveOnCreate {
			// Also save the key in this case, so we can load it next time
			err = saveApiKey(s.apiKey, s.config.keyPath, s.config.keyPassword)

			// Error while error handling, meh
			if err != nil {
				log.Errorf("Error while saving the new API key: %v", err)
			}
		}
	} else if err != nil {
		log.Errorf("Could not load key from file, continuing with a temporary key: %v", err)
	}
}

// Login handles a login request
func (s Service) Login(_ context.Context, request *auth.LoginRequest) (response *auth.LoginResponse, err error) {
	var result bool
	var u *auth.User

	if result, u, err = s.verifyLogin(request); err != nil {
		// a returned error means, something has gone wrong
		return nil, status.Errorf(codes.Internal, "error during login: %v", err)
	}

	if !result {
		// authentication error
		return nil, status.Errorf(codes.Unauthenticated, "login failed")
	}

	var token string

	var expiry = time.Now().Add(time.Hour * 24)

	if token, err = s.issueToken(u.Username, u.FullName, u.Email, expiry); err != nil {
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
				Kid: DefaultKeyID,
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
func (s *Service) verifyLogin(request *auth.LoginRequest) (result bool, user *auth.User, err error) {
	var match bool

	user = new(auth.User)

	err = s.db.Read(user, "username = ?", request.Username)

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
// TODO(all): return created user?
func (s *Service) CreateDefaultUser(username string, password string) error {
	var (
		storedUsers []*auth.User
		err         error
	)

	err = s.db.Read(&storedUsers)
	if err != nil && err != gorm.ErrRecordNotFound {
		return status.Errorf(codes.Internal, "db error: %v", err)
	} else if len(storedUsers) == 0 {
		hash, _ := hashPassword(password)

		u := auth.User{
			Username: username,
			FullName: username,
			Password: string(hash),
		}

		log.Infof("Creating default user %s\n", u.Username)

		err = s.db.Create(&u)
		if err != nil {
			return err
		}
	}
	return nil
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
	claims.Header["kid"] = DefaultKeyID

	token, err = claims.SignedString(s.apiKey)
	return
}

// StartDedicatedAuthServer starts a gRPC server containing just the auth service
func (s *Service) StartDedicatedAuthServer(address string) (sock net.Listener, server *grpc.Server, err error) {
	// create a new socket for gRPC communication
	sock, err = net.Listen("tcp", address)
	if err != nil {
		return nil, nil, fmt.Errorf("could not listen: %w", err)
	}

	err = s.CreateDefaultUser("clouditor", "clouditor")
	if err != nil {
		return nil, nil, fmt.Errorf("could not create default user: %v", err)
	}

	authConfig := service.ConfigureAuth(service.WithPublicKey(s.GetPublicKey()))

	server = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_auth.UnaryServerInterceptor(authConfig.AuthFunc),
		),
	)
	auth.RegisterAuthenticationServer(server, s)

	go func() {
		// serve the gRPC socket
		_ = server.Serve(sock)
	}()

	return sock, server, nil
}

func (s Service) GetPublicKey() *ecdsa.PublicKey {
	return &s.apiKey.PublicKey
}

// AuthFuncOverride implements the ServiceAuthFuncOverride interface to override the AuthFunc for this service.
// The reason is to actually disable authentication checking in the auth service, because functions such as login
// need to be publicly available.
func (Service) AuthFuncOverride(ctx context.Context, _ string) (context.Context, error) {
	// No authentication needed for login functions, otherwise we could not log in
	return ctx, nil
}

// expandPath expands a path that possible contains a tilde (~) character into the home directory
// of the user
func expandPath(path string) (out string, err error) {
	var (
		u *user.User
	)

	// Fetch the current user home directory
	u, err = user.Current()
	if err != nil {
		return path, fmt.Errorf("could not find retrieve current user: %w", err)
	}

	if path == "~" {
		return u.HomeDir, nil
	} else if strings.HasPrefix(path, "~") {
		// We only allow ~ at the beginning of the path
		return filepath.Join(u.HomeDir, path[2:]), nil
	}

	return path, nil
}
