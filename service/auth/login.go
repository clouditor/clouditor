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

package auth

import (
	"context"
	"errors"
	"time"

	"clouditor.io/clouditor"
	"clouditor.io/clouditor/persistence"
	"github.com/dgrijalva/jwt-go"
	"github.com/oxisto/go-httputil/argon2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

//go:generate protoc -I ../../proto -I ../../third_party auth.proto --go_out=../.. --go-grpc_out=../..

var apiSecret = "changeme"
var apiIssuer = "clouditor"

// Service is an implementation of the gRPC Authentication service
type Service struct {
	clouditor.UnimplementedAuthenticationServer
}

// UserClaims extend jwt.StandardClaims with more detailed claims about a user
type UserClaims struct {
	jwt.StandardClaims
	FullName string `json:"full_name"`
	EMail    string `json:"email"`
}

// Login handles a login request
func (s Service) Login(ctx context.Context, request *clouditor.LoginRequest) (response *clouditor.LoginResponse, err error) {
	var result bool
	var user *clouditor.User

	if result, user, err = verifyLogin(request); err != nil {
		// a returned error means, something has gone wrong
		return nil, grpc.Errorf(codes.Internal, "error during login: %v", err)
	}

	if result == false {
		// authentication error
		return nil, grpc.Errorf(codes.Unauthenticated, "login failed")
	}

	var token string

	if token, err = issueToken(user.Username, user.FullName, user.Email, time.Now().Add(1*3600*24)); err != nil {
		return nil, grpc.Errorf(codes.Internal, "token issue failed: %w", err)
	}

	response = &clouditor.LoginResponse{Token: token}

	return response, nil
}

// verifyLogin compares the credentials supplied by request with the ones stored in the database.
// It will return an error in err only if a database error or something else is occurred, this should be
// returned to the user as an internal server error. For security reasons, if authentication failed, only
// the result will be set to false, but no detailed error will be returned to the user.
func verifyLogin(request *clouditor.LoginRequest) (result bool, user *clouditor.User, err error) {
	db := persistence.GetDatabase()

	user = new(clouditor.User)

	err = db.Where("username = ?", request.Username).First(user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// user not found, set result to false, but hide the error
		return false, nil, nil
	} else if err != nil {
		// a database connection error has occurred, return it
		return false, nil, err
	}

	err = argon2.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))

	if errors.Is(err, argon2.ErrMismatchedHashAndPassword) {
		// do not return the error, just set result to false
		return false, nil, nil
	} else if err != nil {
		// some other error occurred, return it
		return false, nil, err
	}

	return true, user, nil
}

// HashPassword returns a hash of password using argon2id.
func (s Service) HashPassword(password string) ([]byte, error) {
	return argon2.GenerateFromPasswordWithParams([]byte(password), argon2.IDParams{
		SaltLength:  16,
		Memory:      65536,
		KeyLength:   16, /* moved over from Java code. might need to be upgraded to 32 */
		Iterations:  3,
		Parallelism: 6,
	})
}

// issueToken issues a JWT token
func issueToken(subject string, fullName string, email string, expiry time.Time) (token string, err error) {
	key := []byte(apiSecret)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&UserClaims{
			FullName: fullName,
			EMail:    email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expiry.Unix(),
				Issuer:    apiIssuer,
				Subject:   subject,
			}},
	)

	token, err = claims.SignedString(key)
	return
}
