// Copyright 2016-2020 Fraunhofer AISEC
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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "auth")
}

const (
	// DefaultApiKeySaveOnCreate specifies whether a created API key will be saved. This is useful to turn of in unit tests, where
	// we only want a temporary key.
	DefaultApiKeySaveOnCreate = true

	// DefaultApiKeyPassword is the default password to protect the API key
	DefaultApiKeyPassword = "changeme"

	// DefaultApiKeyPath is the default path for the API private key
	DefaultApiKeyPath = DefaultConfigDirectory + "/api.key"

	// DefaultConfigDirectory is the default path for the clouditor configuration, such as keys
	DefaultConfigDirectory = "~/.clouditor"
)

// UserClaims extend jwt.StandardClaims with more detailed claims about a user
type UserClaims struct {
	jwt.RegisteredClaims
	FullName string `json:"full_name"`
	EMail    string `json:"email"`
}

type keyLoader struct {
	path         string
	password     string
	saveOnCreate bool
}

// LoadSigningKeys implements a singing keys func for our internal authorization server
func LoadSigningKeys(path string, password string, saveOnCreate bool) map[int]*ecdsa.PrivateKey {
	// create a key loader with our arguments
	loader := keyLoader{
		path:         path,
		password:     password,
		saveOnCreate: saveOnCreate,
	}

	return map[int]*ecdsa.PrivateKey{
		0: loader.LoadKey(),
	}
}

func (l *keyLoader) LoadKey() (key *ecdsa.PrivateKey) {
	var (
		err error
	)

	// Try to load the key from the given path
	key, err = loadKeyFromFile(l.path, []byte(l.password))
	if err != nil {
		key = l.recoverFromLoadApiKeyError(err, l.path == DefaultApiKeyPath)
	}

	return
}

// recoverFromLoadApiKeyError tries to recover from an error during key loading. We treat different errors differently.
// For example if the path is the default path and the error is os.ErrNotExist, this could be just the first start of Clouditor.
// So we only treat this as an information that we will create a new key, which we will save, based on the config.
//
// If the user specifies a custom path and this one does not exist, we will report an error
// here.
func (l *keyLoader) recoverFromLoadApiKeyError(err error, defaultPath bool) (key *ecdsa.PrivateKey) {
	// In any case, create a new temporary API key
	key, _ = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if defaultPath && errors.Is(err, os.ErrNotExist) {
		log.Infof("API key does not exist at the default location yet. We will create a new one")

		if l.saveOnCreate {
			// Also make sure that default config path exists
			err = ensureConfigFolderExistence()
			// Error while error handling, meh
			if err != nil {
				log.Errorf("Error while saving the new API key: %v", err)
			}

			// Also save the key in this case, so we can load it next time
			err = saveKeyToFile(key, l.path, l.password)

			// Error while error handling, meh
			if err != nil {
				log.Errorf("Error while saving the new API key: %v", err)
			}
		}
	} else if err != nil {
		log.Errorf("Could not load key from file, continuing with a temporary key: %v", err)
	}

	return key
}

// loadKeyFromFile loads an ecdsa.PrivateKey from a path. The key must in PEM format and protected by
// a password using PKCS#8 with PBES2.
func loadKeyFromFile(path string, password []byte) (key *ecdsa.PrivateKey, err error) {
	var (
		keyFile string
	)

	keyFile, err = expandPath(path)
	if err != nil {
		return nil, fmt.Errorf("error while expanding path: %w", err)
	}

	if _, err = os.Stat(keyFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist (yet): %w", err)
	}

	// Check, if we already have a persisted API key
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("error while reading key: %w", err)
	}

	key, err = ParseECPrivateKeyFromPEMWithPassword(data, password)
	if err != nil {
		return nil, fmt.Errorf("error while parsing private key: %w", err)
	}

	return key, nil
}

// saveKeyToFile saves an ecdsa.PrivateKey to a path. The key will be saved in PEM format and protected by
// a password using PKCS#8 with PBES2.
func saveKeyToFile(apiKey *ecdsa.PrivateKey, keyPath string, password string) (err error) {
	keyPath, err = expandPath(keyPath)
	if err != nil {
		return fmt.Errorf("error while expanding path: %w", err)
	}

	// Check, if we already have a persisted API key
	f, err := os.OpenFile(keyPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("error while opening the file: %w", err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Errorf("Error while closing file: %v", err)
		}
	}()

	data, err := MarshalECPrivateKeyWithPassword(apiKey, []byte(password))
	if err != nil {
		return fmt.Errorf("error while marshalling private key: %w", err)
	}

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("error while writing file content: %w", err)
	}

	return nil
}

// expandPath expands a path that possible contains a tilde (~) character into the home directory
// of the user
func expandPath(path string) (out string, err error) {
	var (
		u *user.User
	)

	// Fetch the current user home directory
	u, err = user.Current()
	if err != nil {
		return path, fmt.Errorf("could not find retrieve current user: %w", err)
	}

	if path == "~" {
		return u.HomeDir, nil
	} else if strings.HasPrefix(path, "~") {
		// We only allow ~ at the beginning of the path
		return filepath.Join(u.HomeDir, path[2:]), nil
	}

	return path, nil
}

// ensureConfigesFolderExistence ensures that the config folder exists.
func ensureConfigFolderExistence() (err error) {
	var configPath string

	// Expand the config directory, if it contains any ~ characters.
	configPath, err = expandPath(DefaultConfigDirectory)
	if err != nil {
		// Directly return the error here, no need for additional wrapping
		return err
	}

	// Create the directory, if it not exists
	_, err = os.Stat(configPath)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(configPath, os.ModePerm)
		if err != nil {
			// Directly return the error here, no need for additional wrapping
			return err
		}
	}

	return
}
