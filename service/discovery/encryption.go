package discovery

type AtRestEncryption struct {
	Enabled    bool
	Algorithm  string
	KeyManager string
}

type TransportEncryption struct {
	Enabled    bool
	Enforced   bool
	TlsVersion string
}
