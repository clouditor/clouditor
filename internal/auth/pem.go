// Copyright 2022 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

var (
	oidAESCBC         = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 16, 12}
	oidHMACWithSHA256 = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 9}
	oidPBKDF2         = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 12}
	oidPBES2          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 5, 13}
)

var ErrNotECPrivateKey = errors.New("key is not a valid EC private key")

// PBKDF2Params are parameters for PBKDF2. See https://datatracker.ietf.org/doc/html/rfc8018#appendix-A.2.
type PBKDF2Params struct {
	Salt           []byte
	IterationCount int
	PRF            asn1.ObjectIdentifier `asn1:"optional"`
}

// KeyDerivationFunc is part of PBES2 and specify the key derivation function.
// See https://datatracker.ietf.org/doc/html/rfc8018#appendix-A.4.
type KeyDerivationFunc struct {
	Algorithm    asn1.ObjectIdentifier
	PBKDF2Params PBKDF2Params
}

// EncryptionScheme is part of PBES2 and specifies the encryption algorithm.
// See https://datatracker.ietf.org/doc/html/rfc8018#appendix-A.4.
type EncryptionScheme struct {
	EncryptionAlgorithm asn1.ObjectIdentifier
	IV                  []byte
}

// PBES2Params are parameters for PBES2. See https://datatracker.ietf.org/doc/html/rfc8018#appendix-A.4.
type PBES2Params struct {
	KeyDerivationFunc KeyDerivationFunc
	EncryptionScheme  EncryptionScheme
}

// EncryptionAlgorithmIdentifier is the identifier for the encryption algorithm.
// See https://datatracker.ietf.org/doc/html/rfc5958#section-3.
type EncryptionAlgorithmIdentifier struct {
	Algorithm asn1.ObjectIdentifier
	Params    PBES2Params
}

// EncryptedPrivateKeyInfo contains meta-info about the encrypted private key.
// See https://datatracker.ietf.org/doc/html/rfc5958#section-3.
type EncryptedPrivateKeyInfo struct {
	EncryptionAlgorithm EncryptionAlgorithmIdentifier
	EncryptedData       []byte
}

// MarshalECPrivateKeyWithPassword marshals an ECDSA private key protected with a password
// according to PKCS#8 into a byte array
func MarshalECPrivateKeyWithPassword(key *ecdsa.PrivateKey, password []byte) (data []byte, err error) {
	var decrypted []byte
	decrypted, err = x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		// Directly return error here, because we are basically a wrapper around
		// x509.MarshalPKCS8PrivateKey and we want our errors to be similar
		return nil, err
	}

	block, err := EncryptPEMBlock(rand.Reader, decrypted, password)
	if err != nil {
		return nil, fmt.Errorf("could not encrypt PEM block: %w", err)
	}

	return pem.EncodeToMemory(block), nil
}

// ParseECPrivateKeyFromPEMWithPassword ready an ECDSA private key protected with a password according
// to PKCS#8 from a byte array.
func ParseECPrivateKeyFromPEMWithPassword(data []byte, password []byte) (key *ecdsa.PrivateKey, err error) {
	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(data); block == nil {
		return nil, errors.New("could not decode PEM")
	}

	var decrypted []byte
	if decrypted, err = DecryptPEMBlock(block, password); err != nil {
		return nil, fmt.Errorf("could not decrypt PEM block: %w", err)
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(decrypted)
	if err != nil {
		// Directly return error here, because we are basically a wrapper around
		// x509.ParsePKCS8PrivateKey and we want our errors to be similar
		return nil, err
	}

	var ok bool
	if key, ok = parsedKey.(*ecdsa.PrivateKey); !ok {
		return nil, ErrNotECPrivateKey
	}

	return
}

// EncryptPEMBlock encrypts a private key contain in data into a PEM block according to PKCS#8.
func EncryptPEMBlock(rand io.Reader, data, password []byte) (block *pem.Block, err error) {
	var salt = make([]byte, 8)
	if _, err = rand.Read(salt); err != nil {
		return nil, fmt.Errorf("error creating salt: %w", err)
	}

	var iv = make([]byte, 16)
	if _, err = rand.Read(iv); err != nil {
		return nil, fmt.Errorf("error creating IV: %w", err)
	}

	var pad = 16 - len(data)%16

	// Build EncryptedPrivateKeyInfo
	keyInfo := EncryptedPrivateKeyInfo{
		EncryptionAlgorithm: EncryptionAlgorithmIdentifier{
			Algorithm: oidPBES2,
			Params: PBES2Params{
				KeyDerivationFunc: KeyDerivationFunc{
					Algorithm: oidPBKDF2,
					PBKDF2Params: PBKDF2Params{
						IterationCount: 1000,
						Salt:           salt,
						PRF:            oidHMACWithSHA256,
					},
				},
				EncryptionScheme: EncryptionScheme{
					EncryptionAlgorithm: oidAESCBC,
					IV:                  iv,
				},
			},
		},
		EncryptedData: make([]byte, len(data), len(data)+pad), // We will encrypt this later
	}

	// Derive key using PBKDF2
	key := pbkdf2.Key(
		password,
		salt,
		keyInfo.EncryptionAlgorithm.Params.KeyDerivationFunc.PBKDF2Params.IterationCount,
		32,
		sha256.New,
	)

	// Set up symmetric encryption of our block
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create AES cipher mode: %w", err)
	}

	mode := cipher.NewCBCEncrypter(cipherBlock, keyInfo.EncryptionAlgorithm.Params.EncryptionScheme.IV)

	copy(keyInfo.EncryptedData, data)
	for i := 0; i < pad; i++ {
		keyInfo.EncryptedData = append(keyInfo.EncryptedData, byte(pad))
	}

	mode.CryptBlocks(keyInfo.EncryptedData, keyInfo.EncryptedData)

	block = &pem.Block{
		Type:    "ENCRYPTED PRIVATE KEY",
		Headers: make(map[string]string),
	}

	// Marshal key info into ASN1 format, which is the payload of our PEM block
	block.Bytes, err = asn1.Marshal(keyInfo)
	if err != nil {
		return nil, fmt.Errorf("could not marshal ASN1: %w", err)
	}

	return
}

// DecryptPEMBlock is a drop-in replacement for x509.DecryptPEMBlock which only supports
// state-of-the art algorithms such as PBES2.
func DecryptPEMBlock(block *pem.Block, password []byte) ([]byte, error) {
	var (
		keyInfo EncryptedPrivateKeyInfo
		prf     asn1.ObjectIdentifier
		err     error
	)

	if block.Type != "ENCRYPTED PRIVATE KEY" {
		return nil, errors.New("key is not a PKCS#8")
	}

	_, err = asn1.Unmarshal(block.Bytes, &keyInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private key info: %w", err)
	}

	if !keyInfo.EncryptionAlgorithm.Algorithm.Equal(oidPBES2) {
		return nil, errors.New("unsupported encryption algorithm: only PBES2 is supported")
	}

	if !keyInfo.EncryptionAlgorithm.Params.KeyDerivationFunc.Algorithm.Equal(oidPBKDF2) {
		return nil, errors.New("unsupported key derivation algorithm: only PBKDF2 is supported")
	}

	prf = keyInfo.EncryptionAlgorithm.Params.KeyDerivationFunc.PBKDF2Params.PRF

	if prf != nil && !prf.Equal(oidHMACWithSHA256) {
		return nil, errors.New("unsupported pseudo-random function: only HMACWithSHA256 is supported")
	}

	keyParams := keyInfo.EncryptionAlgorithm.Params.KeyDerivationFunc.PBKDF2Params
	keyHash := sha256.New

	symkey := pbkdf2.Key(password, keyParams.Salt, keyParams.IterationCount, 32, keyHash)
	cipherBlock, err := aes.NewCipher(symkey)
	if err != nil {
		return nil, fmt.Errorf("could not create AES cipher mode: %w", err)
	}

	mode := cipher.NewCBCDecrypter(cipherBlock, keyInfo.EncryptionAlgorithm.Params.EncryptionScheme.IV)
	mode.CryptBlocks(keyInfo.EncryptedData, keyInfo.EncryptedData)

	return keyInfo.EncryptedData, nil
}
