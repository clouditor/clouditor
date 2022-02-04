package common

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var JwkURL = "http://localhost/.well-known/jwks.json"

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

func RetrieveJWK() (KeySet, error) {
	req, _ := http.Get(JwkURL)
	b, _ := ioutil.ReadAll(req.Body)

	var set KeySet

	err := json.Unmarshal(b, &set)

	return set, err
}

type KeySet struct {
	Keys []JWK `json:keys`
}

type JWK struct {
	Kty string `json:kty`
	Crv string `json:crv`
	X   string `json:x`
	Y   string `json:y`
	Kid string `json:kid`
}

func (k KeySet) ByID(kid string) (key crypto.PublicKey, err error) {
	for _, j := range k.Keys {
		if j.Kid == kid {
			x := new(big.Int)
			x.SetString(j.X, 10)
			y := new(big.Int)
			y.SetString(j.Y, 10)
			key = &ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     x,
				Y:     y,
			}
			return key, nil
		}
	}

	return nil, nil
}

func KeyFuncJwk(t *jwt.Token) (interface{}, error) {
	set, err := RetrieveJWK()

	if err != nil {
		return nil, fmt.Errorf("could not retrieve JWKS: %w", err)
	}

	// get kid
	kid, ok := t.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid is not a string")
	}

	// get key from key set
	key, err := set.ByID(kid)

	return &key, err
}
