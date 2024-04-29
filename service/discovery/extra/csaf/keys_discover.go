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

func (d *csafDiscovery) discoverKeys(pgpkeys []csaf.PGPKey, parentId string) (keys []ontology.IsResource) {
	for _, pgpkey := range pgpkeys {
		keys = append(keys, d.handleKey(pgpkey, parentId))
	}

	return
}

func (d *csafDiscovery) handleKey(pgpkey csaf.PGPKey, parentId string) (key *ontology.Key) {
	return &ontology.Key{
		Algorithm: "PGP",
		Id:        util.Deref(pgpkey.URL),
		Raw:       discovery.Raw(pgpkey),
		ParentId:  &parentId,
		// We can always set it to true because otherwise fetchKey would return an error and stop discovery anyway
		InternetAccessibleEndpoint: true,
	}
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
