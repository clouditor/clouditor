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

// DefaultApiKeyPassword is the default password to protect the API key
const DefaultApiKeyPassword = "changeme"

const DefaultApiKeyPath = "~/.clouditor/api.key"

// Service is an implementation of the gRPC Authentication service
type Service struct {
	auth.UnimplementedAuthenticationServer

	// TODO(oxisto): re-use private key / store somewhere secure (pw protected)
	apiKey *ecdsa.PrivateKey
}

// UserClaims extend jwt.StandardClaims with more detailed claims about a user
type UserClaims struct {
	jwt.RegisteredClaims
	FullName string `json:"full_name"`
	EMail    string `json:"email"`
}

func init() {
	log = logrus.WithField("component", "auth")
}

func NewService() *Service {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	return &Service{
		apiKey: key,
	}
}

func (s Service) loadApiKey() {
	//keyFile := DefaultApiKeyPath

	// Check, if we already have a persisted API key
	//os.OpenFile(keyFile, os.O_RDONLY, 0600)

	//x509.ParsePKCS1PrivateKey()
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
