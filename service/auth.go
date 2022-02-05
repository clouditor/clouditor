package service

import (
	"context"
	"fmt"
	"time"

	"clouditor.io/clouditor/rest"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var log *logrus.Entry

type AuthConfig struct {
	jwksUrl string

	Jwks     *keyfunc.JWKS
	AuthFunc grpc_auth.AuthFunc
}

var DefaultJwkUrl = fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", rest.DefaultAPIHTTPPort)

func init() {
	log = logrus.WithField("component", "service-auth")
}

type AuthOption func(*AuthConfig)

func WithJwksUrl(url string) AuthOption {
	return func(ac *AuthConfig) {
		ac.jwksUrl = url
	}
}

type AuthContextKey string

func ConfigureAuth(opts ...AuthOption) *AuthConfig {
	var config = &AuthConfig{
		jwksUrl: DefaultJwkUrl,
	}

	// Apply options
	for _, o := range opts {
		o(config)
	}

	config.AuthFunc = func(ctx context.Context) (newCtx context.Context, err error) {
		// Lazy loading of JWKS
		if config.Jwks == nil {
			config.Jwks, err = keyfunc.Get(config.jwksUrl, keyfunc.Options{
				RefreshInterval: time.Hour,
			})
			if err != nil {
				log.Debugf("Could not retrieve JWKS. API authentication will fail: %v", err)
				return nil, status.Errorf(codes.FailedPrecondition, "could not retrieve JWKS: %v", err)
			}
		}

		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}

		tokenInfo, err := parseToken(token, config)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}

		newCtx = context.WithValue(ctx, AuthContextKey("token"), tokenInfo)

		return newCtx, nil
	}

	return config
}

func parseToken(token string, authConfig *AuthConfig) (jwt.Claims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, authConfig.Jwks.Keyfunc)

	if err != nil {
		return nil, fmt.Errorf("could not validate JWT: %w", err)
	}

	return parsedToken.Claims, nil
}
