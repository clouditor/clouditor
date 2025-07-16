package server

import "crypto/ecdsa"

// DefaultJWKSURL is the default JWKS url pointing to a local authentication server.
const DefaultJWKSURL = "http://localhost:8080/v1/auth/certs"

// WithJWKS is an option to provide a URL that contains a JSON Web Key Set (JWKS). The JWKS will be used
// to validate tokens coming from RPC clients against public keys contains in the JWKS.
func WithJWKS(url string) StartGRPCServerOption {
	return func(c *config) {
		c.ac.JwksURL = url
		c.ac.UseJWKS = true
	}
}

// WithPublicKey is an option to directly provide a ECDSA public key which is used to verify tokens coming from RPC clients.
func WithPublicKey(publicKey *ecdsa.PublicKey) StartGRPCServerOption {
	return func(c *config) {
		c.ac.PublicKey = publicKey
	}
}
