package voc

import "github.com/cayleygraph/quad"

type IsEncryption interface {
	IsEnabled() bool

	isEncryption() bool
}

type Encryption struct {
	rdfType struct{} `quad:"@type > cloud:Encryption"`

	Enabled bool `quad:"cloud:enabled"`
}

func (e *Encryption) IsEnabled() bool {
	return e.Enabled
}

func (e *Encryption) isEncryption() bool {
	return true
}

type AtRestEncryption struct {
	rdfType struct{} `quad:"@type > cloud:AtRestEncryption"`

	Encryption
	Algorithm  string `quad:"cloud:algorithm"`
	KeyManager string `quad:"cloud:keyManager"`
}

func NewAtRestEncryption(enabled bool, algorithm string, keyManager string) *AtRestEncryption {
	return &AtRestEncryption{
		Encryption: Encryption{Enabled: enabled},
		Algorithm:  algorithm,
		KeyManager: keyManager,
	}
}

type TransportEncryption struct {
	rdfType struct{} `quad:"@type > cloud:TransportEncryption"`

	Encryption
	Enforced   bool     `quad:"cloud:enforced,optional"`
	TlsVersion quad.IRI `quad:"cloud:tlsVersion,optional"`
}

func NewTransportEncryption(enabled bool, enforced bool, tlsVersion quad.IRI) *TransportEncryption {
	return &TransportEncryption{
		Encryption: Encryption{Enabled: enabled},
		Enforced:   enforced,
		TlsVersion: tlsVersion,
	}
}
