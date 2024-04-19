package providertest

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"time"

	"clouditor.io/clouditor/v2/internal/util"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"golang.org/x/net/http2"
)

type testProvider struct {
	http.Server
	pmd  *csaf.ProviderMetadata
	sock net.Listener
}

func NewTestProvider() (p *testProvider, err error) {
	mux := http.NewServeMux()

	p = &testProvider{}
	p.Handler = mux

	mux.HandleFunc("/.well-known/csaf/provider-metadata.json", p.handlePMD)

	p.sock, err = net.Listen("tcp4", ":0")
	if err != nil {
		return nil, fmt.Errorf("could not listen: %w", err)
	}

	p.pmd = csaf.NewProviderMetadataDomain(p.Domain(), []csaf.TLPLabel{csaf.TLPLabelWhite})
	p.pmd.Publisher = &csaf.Publisher{
		Name:      util.Ref("Test Vendor"),
		Category:  util.Ref(csaf.CSAFCategoryVendor),
		Namespace: util.Ref("http://localhost"),
	}

	// Create a TLS certificate
	//publickey, privkey, err := ed25519.GenerateKey(rand.Reader)
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		return nil, fmt.Errorf("could not generate ed25519 private key: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: new(big.Int),
		Subject: pkix.Name{
			Organization: []string{"Clouditor"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 1),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:                  true,
		DNSNames:              []string{"localhost", "csaf.data.security.localhost"},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	bytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, fmt.Errorf("could not create x509 certificate: %w", err)
	}
	cert, err := x509.ParseCertificate(bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse x509 certificate: %w", err)
	}

	pool := x509.NewCertPool()
	pool.AddCert(cert)

	p.Server.TLSConfig = &tls.Config{
		RootCAs: pool,
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert.Raw},
				PrivateKey:  priv,
				Leaf:        cert,
			},
		},
	}

	return
}

func (p *testProvider) handlePMD(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(p.pmd)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (p *testProvider) Serve() error {
	return p.Server.ServeTLS(p.sock, "", "")
}

func (p *testProvider) Domain() string {
	port := p.sock.Addr().(*net.TCPAddr).AddrPort()
	return fmt.Sprintf("localhost:%d", port.Port())
}

func (p *testProvider) Client() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: p.TLSConfig,
		},
	}
}
