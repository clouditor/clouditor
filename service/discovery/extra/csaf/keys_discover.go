package csaf

import (
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/api/ontology"
	"clouditor.io/clouditor/v2/internal/crypto/openpgp"
	"clouditor.io/clouditor/v2/internal/util"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
)

// discoverKeys discovers the PGP keys and returns the respective keys in the ontology format as well as the keyring
// needed for discovering the Security Advisories (discoverSecurityAdvisories)
func (d *csafDiscovery) discoverKeys(pgpkeys []csaf.PGPKey, parentId string) (keys []ontology.IsResource, keyring openpgp.EntityList) {
	for _, pgpkey := range pgpkeys {
		key, openPGPEntity := d.handleKey(pgpkey, parentId)
		keys = append(keys, key)
		keyring = append(keyring, openPGPEntity)
	}
	return
}

// handleKey handles a [csaf.PGPKey]: First we try to fetch the actual key and provide a [openpgp.Entity]. Then we use
// this information as well as the information provided by [csaf.PGPKey] to create a Key ontology object.
func (d *csafDiscovery) handleKey(pgpkey csaf.PGPKey, parentId string) (key *ontology.Key, openPGPEntity *openpgp.Entity) {
	var (
		err error
		// isAccessible denotes that the key is accessible. We assume it is but set it to false if we could not fetch it
		isAccessible = true
	)

	// 1st: Try to fetch key for creating the OpenGPG entity
	openPGPEntity, err = d.fetchKey(pgpkey)
	if err != nil {
		// If we could not fetch the key we assume that the key exists but is not accessible
		isAccessible = false
		log.Warnf("Could not fetch key '%s': %v", util.Deref(pgpkey.URL), err)
	}

	// 2nd: Create the key in the ontology format
	key = &ontology.Key{
		Algorithm:                  "PGP",
		Id:                         util.Deref(pgpkey.URL),
		Raw:                        discovery.Raw(pgpkey),
		ParentId:                   &parentId,
		InternetAccessibleEndpoint: isAccessible,
	}
	return
}

func (d *csafDiscovery) fetchKey(keyinfo csaf.PGPKey) (key *openpgp.Entity, err error) {
	var (
		res  *http.Response
		keys openpgp.EntityList
	)

	res, err = d.client.Get(util.Deref(keyinfo.URL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	keys, err = openpgp.ReadArmoredKeyRing(res.Body)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, errors.New("no key in key file")
	} else if len(keys) > 1 {
		return nil, errors.New("too many keys in file")
	}

	key = keys[0]

	if !strings.EqualFold(hex.EncodeToString(key.PrimaryKey.Fingerprint), string(keyinfo.Fingerprint)) {
		return nil, errors.New("fingerprints do not match")
	}

	return
}
