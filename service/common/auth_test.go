package common

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

var privateKey *ecdsa.PrivateKey

func init() {
	privateKey, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func Test(t *testing.T) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 0),
		Handler: handle(),
	}

	sock, _ := net.Listen("tcp", srv.Addr)

	port := sock.Addr().(*net.TCPAddr).Port
	go srv.Serve(sock)

	// Create a JWT based on our private key
	token := jwt.New(jwt.SigningMethodES256)
	token.Header["kid"] = "1"
	signed, err := token.SignedString(privateKey)

	assert.ErrorIs(t, err, nil)

	config := ConfigureAuth(WithJwksUrl(fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port)))
	ctx := metadata.NewIncomingContext(context.TODO(), metadata.MD{"Authorization": []string{fmt.Sprintf("bearer %s", signed)}})

	newCtx, err := config.AuthFunc(ctx)

	assert.ErrorIs(t, err, nil)
	assert.NotNil(t, newCtx)

	key, ok := config.Jwks.ReadOnlyKeys()["1"]

	assert.True(t, ok)
	assert.NotNil(t, key)

	pkey, ok := key.(*ecdsa.PublicKey)

	assert.True(t, ok)
	assert.NotNil(t, pkey)
}

func handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/jwks.json" {
			var json = fmt.Sprintf(`{"keys":
			[
			  {"kty":"EC",
			   "crv":"P-256",
			   "x":"%s",
			   "y":"%s",
			   "use":"enc",
			   "kid":"1"}
			]
		  }`,
				base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.X.Bytes()),
				base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.Y.Bytes()),
			)

			fmt.Printf(json)
			// Example key from RFC 7517 A.1
			w.Write([]byte(json))
		}
	})
}
