package common

import (
	"context"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var JwkURL = "http://localhost/.well-known/jwks.json"
var jwks *keyfunc.JWKS

var log *logrus.Entry

func init() {
	var err error

	log = logrus.WithField("component", "service-auth")

	jwks, err = keyfunc.Get(JwkURL, keyfunc.Options{
		RefreshInterval: time.Hour,
	})
	if err != nil {
		log.Errorf("Failed to get the JWKS from the given URL :%v", err)
	}
}

func MyAuth(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	newCtx := context.WithValue(ctx, "something", tokenInfo)

	return newCtx, nil
}

func parseToken(token string) (jwt.Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, KeyFuncJwk)

	if err != nil {
		return nil, fmt.Errorf("could not validate JWT: %w", err)
	}

	return parsedToken.Claims, nil
}

func RetrieveJWK() (*keyfunc.JWKS, error) {

	return jwks, nil
}

func KeyFuncJwk(t *jwt.Token) (interface{}, error) {
	set, err := RetrieveJWK()

	if err != nil {
		return nil, fmt.Errorf("could not retrieve JWKS: %w", err)
	}

	// get key from key set
	return set.Keyfunc(t)
}
