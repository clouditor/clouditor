package common

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 0),
		Handler: handle(),
	}

	sock, _ := net.Listen("tcp", srv.Addr)

	port := sock.Addr().(*net.TCPAddr).Port

	go srv.Serve(sock)

	JwkURL = fmt.Sprintf("http://localhost:%d/.well-known/jwks.json", port)

	set, err := RetrieveJWK()

	assert.Nil(t, err)
	assert.NotNil(t, set)

	key, ok := set.ReadOnlyKeys()["Public key used in JWS spec Appendix A.3 example"]

	assert.True(t, ok)
	assert.NotNil(t, key)

	pkey, ok := key.(*ecdsa.PublicKey)

	assert.True(t, ok)
	assert.NotNil(t, pkey)
}

func handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/jwks.json" {
			w.Write([]byte(`{"keys":[{"kty":"EC",
			"crv":"P-256",
			"x":"f83OJ3D2xF1Bg8vub9tLe1gHMzV76e8Tus9uPHvRVEU",
			"y":"x_FEzRu9m36HLN_tue659LNpXW6pCyStikYjKIWI5a0",
			"kid":"Public key used in JWS spec Appendix A.3 example"
		   }]}`))
		}
	})
}
