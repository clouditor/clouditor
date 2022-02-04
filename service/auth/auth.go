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
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"time"

	"clouditor.io/clouditor/api/auth"
	"clouditor.io/clouditor/persistence"
	argon2 "github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

const (
	Issuer = "clouditor"
)

var log *logrus.Entry

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

	if token, err = s.issueToken(user.Username, user.FullName, user.Email, time.Now().Add(time.Hour*24)); err != nil {
		return nil, status.Errorf(codes.Internal, "token issue failed: %v", err)
	}

	response = &auth.LoginResponse{Token: token}

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

func (s Service) GetPublicKey() crypto.PublicKey {
	return s.apiKey.PublicKey
}

// StartDedicatedAuthServer starts a gRPC server containing just the auth service
func StartDedicatedAuthServer(address string) (sock net.Listener, server *grpc.Server, authService *Service, err error) {
	// create a new socket for gRPC communication
	sock, err = net.Listen("tcp", address)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not listen: %w", err)
	}

	authService = NewService()
	authService.CreateDefaultUser("clouditor", "clouditor")

	server = grpc.NewServer()
	auth.RegisterAuthenticationServer(server, authService)

	go func() {
		// serve the gRPC socket
		_ = server.Serve(sock)
	}()

	return sock, server, authService, nil
}
