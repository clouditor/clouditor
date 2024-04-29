package openpgp

import (
	"bytes"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
)

type Entity = openpgp.Entity
type EntityList = openpgp.EntityList

var ArmoredDetachSignText = openpgp.ArmoredDetachSignText
var NewEntity = openpgp.NewEntity
var ReadArmoredKeyRing = openpgp.ReadArmoredKeyRing
var CheckArmoredDetachedSignature = openpgp.CheckArmoredDetachedSignature

// WriteArmoredKey serializes a [openpgp.Entity] (more specifically its public
// key) in an armored form. It is basically a (missing) counterpart to
// [ReadArmoredKeyRing].
func WriteArmoredKey(key *openpgp.Entity) (armor string, err error) {
	var b bytes.Buffer
	err = key.Serialize(&b)
	if err != nil {
		return "", err
	}

	return doArmor(b.Bytes(), openpgp.PublicKeyType)
}

// doArmor is an internal helper to help with the armoring.
func doArmor(in []byte, blockType string) (out string, err error) {
	var (
		b bytes.Buffer
		w io.WriteCloser
	)

	w, err = armor.Encode(&b, blockType, nil)
	if err != nil {
		return "", err
	}

	_, err = w.Write(in)
	if err != nil {
		return "", err
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
