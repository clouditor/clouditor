// Copyright 2024 Fraunhofer AISEC
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

package azure

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"

	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
)

// keyUsage contains for each key an array of resources where it is being used. Currently, only Cosmos DB is supported
var keyUsage = make(map[string][]string)

// secretUsage contains for each secret an array of resources where it is being used. Currently, only web app is supported
var secretUsage = make(map[string][]string)

type azureKeyVaultDiscovery struct {
	*azureDiscovery
	// metricsClient is a client to query Azure Monitor w.r.t. given metrics (e.g. API Hits)
	metricsClient *azquery.MetricsClient
}

func NewKeyVaultDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureKeyVaultDiscovery{
		azureDiscovery: &azureDiscovery{
			discovererComponent: KeyVaultComponent,
			csID:                discovery.DefaultCloudServiceID,
			backupMap:           make(map[string]*backup),
		},
	}
	for _, opt := range opts {
		opt(d.azureDiscovery)
	}
	return d
}

func (*azureKeyVaultDiscovery) Name() string {
	return "Azure Key Vault"
}

func (d *azureKeyVaultDiscovery) List() (list []voc.IsCloudResource, err error) {
	// make sure, we are authorized
	if err = d.authorize(); err != nil {
		err = fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
		return
	}

	// Discover Key Vaults
	log.Info("Discover Azure Key Vaults")
	keyVaults, err := d.discoverKeyVaults()
	if err != nil {
		err = fmt.Errorf("could not discover Key Vaults: %w", err)
		return
	}
	list = append(list, keyVaults...)

	return
}

// discoverKeyVaults discovers all key vaults as well as all belonging keys
func (d *azureKeyVaultDiscovery) discoverKeyVaults() (list []voc.IsCloudResource, err error) {
	// initialize key vault client
	if err = d.initKeyVaultClient(); err != nil {
		return nil, fmt.Errorf("could not initialize key vault client: %v", err)
	}
	// initialize keys client
	if err = d.initKeysClient(); err != nil {
		return nil, fmt.Errorf("could not initialize keys client: %v", err)
	}
	// initialize secrets client
	if err = d.initSecretsClient(); err != nil {
		return nil, fmt.Errorf("could not initialize secrets client: %v", err)
	}

	err = listPager(d.azureDiscovery,
		d.clients.keyVaultClient.NewListPager,
		d.clients.keyVaultClient.NewListByResourceGroupPager,
		func(res armkeyvault.VaultsClientListResponse) (vaults []*armkeyvault.Vault) {
			// TODO(lebogg): Dunno why but here, `res.Value` list is of type Resource instead of Vault (in contrast to the resource group call, see below)
			vaults = make([]*armkeyvault.Vault, len(res.Value))
			for i, r := range res.Value {
				vaults[i] = &armkeyvault.Vault{
					Location: r.Location,
					Tags:     r.Tags,
					ID:       r.ID,
					Name:     r.Name,
					Type:     r.Type,
				}
			}
			return
		},
		func(res armkeyvault.VaultsClientListByResourceGroupResponse) []*armkeyvault.Vault {
			return res.Value
		},
		func(kv *armkeyvault.Vault) error {
			// Handle key vault
			keyVault, err := d.handleKeyVault(kv)
			if err != nil {
				return fmt.Errorf("could not handle key vault: %w", err)
			}

			// Handle Keys and add IDs to the key vault
			keys, err := d.getKeys(kv)
			if err != nil {
				return fmt.Errorf("could not handle keys: %w", err)
			}
			keyVault.Keys = getKeyIDs(keys)

			// Handle secrets and add IDs to the key vault
			secrets, err := d.getSecrets(kv)
			if err != nil {
				return fmt.Errorf("could not handle secrets: %w", err)
			}
			keyVault.Secrets = getSecretIDs(secrets)

			// Handle certificates
			certificates, err := d.getCertificates(kv)
			if err != nil {
				return fmt.Errorf("could not handle secrets: %w", err)
			}
			keyVault.Certificates = getCertificateIDs(certificates)

			// Add all resources (key vaults, keys and secrets) to the list
			log.Infof("Adding key vault '%s'", keyVault.GetName())
			list = append(list, keyVault)
			log.Infof("Adding keys '%s'", gekKeyNames(keys))
			for _, k := range keys {
				list = append(list, k)
			}
			log.Infof("Adding secrets '%s'", getSecretNames(secrets))
			for _, s := range secrets {
				list = append(list, s)
			}
			log.Infof("Adding certificates '%s'", getCertificationNames(certificates))
			for _, c := range certificates {
				list = append(list, c)
			}

			return nil
		})
	if err != nil {
		log.Errorf("Could not run list pager in while discovering key vaults: %v", err)
		list = nil
		return
	}
	return
}

func gekKeyNames(keys []*voc.Key) (keyNames []string) {
	for _, k := range keys {
		keyNames = append(keyNames, k.GetName())
	}
	return keyNames
}

func getSecretNames(secrets []*voc.Secret) (secretNames []string) {
	for _, k := range secrets {
		secretNames = append(secretNames, k.GetName())
	}
	return secretNames
}

func getCertificationNames(certificates []*voc.Certificate) (certificatesNames []string) {
	for _, k := range certificates {
		certificatesNames = append(certificatesNames, k.GetName())
	}
	return certificatesNames

}

func (d *azureKeyVaultDiscovery) initKeyVaultClient() (err error) {
	d.clients.keyVaultClient, err = initClient(d.clients.keyVaultClient, d.azureDiscovery, armkeyvault.NewVaultsClient)
	return
}

func (d *azureKeyVaultDiscovery) initKeysClient() (err error) {
	d.clients.keysClient, err = initClient(d.clients.keysClient, d.azureDiscovery, armkeyvault.NewKeysClient)
	return
}
func (d *azureKeyVaultDiscovery) initSecretsClient() (err error) {
	d.clients.secretsClient, err = initClient(d.clients.secretsClient, d.azureDiscovery, armkeyvault.NewSecretsClient)
	return
}

func (d *azureKeyVaultDiscovery) initCertificatesClient(kv *armkeyvault.Vault) (err error) {
	d.clients.certificationsClient, err = azcertificates.NewClient(util.Deref(kv.Properties.VaultURI), d.cred, nil)
	return
}

func (d *azureKeyVaultDiscovery) initMetricsClient() (err error) {
	d.metricsClient, err = azquery.NewMetricsClient(d.cred, &azquery.MetricsClientOptions{})
	// TODO(all): I cannot use the generic initClient function (see below) since `azquery.NewMetricsClient()`
	// has another structure. But I think that is the reason I cannot successfully use it in tests either...
	// d.metricsClient, err = initClient(d.metricsClient, d.azureDiscovery, azquery.NewMetricsClient)
	return
}

// TODO(lebogg): Test
func (d *azureKeyVaultDiscovery) handleKeyVault(kv *armkeyvault.Vault) (*voc.KeyVault, error) {
	var createdAt *time.Time
	// Find out if key vault is actively used
	isActive, err := d.isActive(kv)
	if err != nil {
		return nil, fmt.Errorf("could not handle key vault: %v", err)
	}

	if kv.SystemData != nil {
		createdAt = kv.SystemData.CreatedAt
	}

	return &voc.KeyVault{
		Resource: discovery.NewResource(d,
			voc.ResourceID(util.Deref(kv.ID)),
			util.Deref(kv.Name),
			createdAt,
			voc.GeoLocation{
				Region: util.Deref(kv.Location),
			},
			labels(kv.Tags),
			resourceGroupID(kv.ID),
			voc.KeyVaultType,
			kv),
		IsActive:     isActive,
		Keys:         []voc.ResourceID{}, // Will be added later when we retrieve the single keys
		PublicAccess: getPublicAccess(kv),
	}, nil
}

// TODO(lebogg): Test in Az-Testbed
func getPublicAccess(kv *armkeyvault.Vault) bool {
	if kv.Properties != nil && util.Deref(kv.Properties.PublicNetworkAccess) == "Enabled" {
		if len(kv.Properties.NetworkACLs.VirtualNetworkRules) == 0 && len(kv.Properties.NetworkACLs.IPRules) == 0 {
			return true
		} else { // There is at least one rule, i.e. we assume that public network access is restricted
			return false
		}
	}
	// Otherwise, public network access is set to "Disabled" or properties is not set -> we assume no access
	return false
}

// getKeyIDs returns the ID values corresponding to the given keys. If slice of keys is empty, return empty slice of
// resourceIDs (not nil slice)
func getKeyIDs(keys []*voc.Key) []voc.ResourceID {
	keyIDs := []voc.ResourceID{}
	for _, k := range keys {
		keyIDs = append(keyIDs, k.GetID())
	}
	return keyIDs
}

// getSecretIDs returns the ID values corresponding to the given secrets. If slice of secrets is empty, return empty
// slice of resourceIDs (not nil slice)
func getSecretIDs(secrets []*voc.Secret) []voc.ResourceID {
	secretIDs := []voc.ResourceID{}
	for _, s := range secrets {
		secretIDs = append(secretIDs, s.GetID())
	}
	return secretIDs
}

// getCertificateIDs returns the ID values corresponding to the given certificates. If slice of certificates is empty,
// return empty slice of resourceIDs (not nil slice)
func getCertificateIDs(certs []*voc.Certificate) []voc.ResourceID {
	certificateIDs := []voc.ResourceID{}
	for _, c := range certs {
		certificateIDs = append(certificateIDs, c.GetID())
	}
	return certificateIDs
}

// isActive determines whether the key vault is being actively used. Measuring is done by examining the API traffic of
// the key vault (API hits via Azure Monitoring). The number of required API hits and the time period measured are
// defined by NumberOfAPIHits and PeriodOfAPIHits, respectively.
func (d *azureKeyVaultDiscovery) isActive(kv *armkeyvault.Vault) (isActive bool, err error) {
	// We need the client for doing metric queries to Azure Monitor
	err = d.initMetricsClient()
	if err != nil {
		return false, fmt.Errorf("could not create Azure Metrics Client (Azure Monitor): %v", err)
	}
	metrics, err := d.metricsClient.QueryResource(context.TODO(), util.Deref(kv.ID),
		&azquery.MetricsClientQueryResourceOptions{
			Aggregation:     nil,
			Filter:          nil,
			Interval:        util.Ref("P1D"), // TODO(lebogg): For testing. In the end we probably want to use timespan to increase this number (max allowed is 1 day)
			MetricNames:     util.Ref("ServiceApiHit"),
			MetricNamespace: nil,
			OrderBy:         nil,
			ResultType:      nil,
			Timespan:        nil,
			Top:             nil,
		})
	if err != nil {
		// TODO(lebogg): To Test: Maybe there are resources (in this case, key vaults) where no API Hit is defined -> Then it is not an error but, e.g., false?
		return false, fmt.Errorf("could not query resource for metric (Monitoring): %v", err)
	}
	// TODO(lebogg): To test: Can this even happen?
	if metrics.Value == nil {
		return false, fmt.Errorf("something went wrong. There are no value(s) for this metric")
	}
	// TODO(lebogg): We only asked for one metric, so we should only get one value ?!
	if l := len(metrics.Value); l != 1 {
		return false, fmt.Errorf("we got %d metrics. But should be one", l)
	}
	metric := metrics.Value[0]
	// TODO(lebogg): If timeseries or data is nil nothing is tracked -> No API Hit or error?
	if metric.TimeSeries[0] == nil || metric.TimeSeries[0].Data[0] == nil {
		return false, nil
	}
	if util.Deref(metric.TimeSeries[0].Data[0].Count) >= 5 {
		return true, nil
	} else {
		return false, nil
	}

}

// Todo(lebogg): What happens with different versions of a key
func (d *azureKeyVaultDiscovery) getKeys(kv *armkeyvault.Vault) ([]*voc.Key, error) {
	var (
		keys []*voc.Key
		c    *armkeyvault.KeysClient
	)
	c = d.clients.keysClient
	if c == nil {
		return nil, errors.New("keys client is empty")
	}
	pager := c.NewListPager(util.Deref(d.azureDiscovery.rg), util.Deref(kv.Name), &armkeyvault.KeysClientListOptions{})
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("could not page next page (paging error): %v", err)
		}
		for _, k := range page.Value {
			// We have to request each single key because this lazy NewListPager doesn't fill out all key information
			res, err := c.Get(context.Background(), util.Deref(d.rg), util.Deref(kv.Name),
				util.Deref(k.Name), &armkeyvault.KeysClientGetOptions{})
			if err != nil {
				return nil, fmt.Errorf("could not get key: %v", err)
			}
			k = util.Ref(res.Key) // maybe not the most beautiful thing to re-use var `k`
			key := &voc.Key{
				Resource: discovery.NewResource(d,
					voc.ResourceID(util.Deref(k.ID)),
					util.Deref(k.Name),
					util.Ref(time.Unix(util.Deref(k.Properties.Attributes.Created), 0)),
					voc.GeoLocation{Region: util.Deref(k.Location)},
					labels(k.Tags),
					voc.ResourceID(util.Deref(kv.ID)),
					voc.KeyType,
					k, kv),
				Enabled:        util.Deref(k.Properties.Attributes.Enabled),
				ActivationDate: util.Deref(k.Properties.Attributes.NotBefore),
				ExpirationDate: util.Deref(k.Properties.Attributes.Expires),
				KeyType:        getKeyType(k.Properties.Kty),
				KeySize:        int(util.Deref(k.Properties.KeySize)),
				NumberOfUsages: len(keyUsage[util.Deref(k.Properties.KeyURI)]), // TODO(lebogg): Test this!
			}
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func (d *azureKeyVaultDiscovery) getSecrets(kv *armkeyvault.Vault) ([]*voc.Secret, error) {
	var (
		secrets []*voc.Secret
		c       *armkeyvault.SecretsClient
	)
	c = d.clients.secretsClient
	if c == nil {
		return nil, errors.New("secrets client is empty")
	}
	pager := c.NewListPager(util.Deref(d.azureDiscovery.rg), util.Deref(kv.Name), &armkeyvault.SecretsClientListOptions{})
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("could not page next page (paging error): %v", err)
		}
		for _, s := range page.Value {
			// We have to request each single key because this lazy NewListPager doesn't fill out all key information
			res, err := d.clients.secretsClient.Get(context.Background(), util.Deref(d.rg), util.Deref(kv.Name),
				util.Deref(s.Name), &armkeyvault.SecretsClientGetOptions{})
			if err != nil {
				return nil, fmt.Errorf("could not get key: %v", err)
			}
			s = util.Ref(res.Secret) // maybe not the most beautiful thing to re-use var `k`
			secret := &voc.Secret{
				Resource: discovery.NewResource(d,
					voc.ResourceID(util.Deref(s.ID)),
					util.Deref(s.Name),
					s.Properties.Attributes.Created,
					voc.GeoLocation{Region: util.Deref(s.Location)},
					labels(s.Tags),
					voc.ResourceID(util.Deref(kv.ID)),
					voc.SecretType,
					s, kv),
				Enabled:        util.Deref(s.Properties.Attributes.Enabled),
				ActivationDate: convertTime(s.Properties.Attributes.NotBefore),
				ExpirationDate: convertTime(s.Properties.Attributes.Expires),
				NumberOfUsages: len(secretUsage[util.Deref(s.Properties.SecretURI)]), // TODO(lebogg): Test this!
			}
			secrets = append(secrets, secret)
		}
	}
	return secrets, nil
}

func (d *azureKeyVaultDiscovery) getCertificates(kv *armkeyvault.Vault) (certs []*voc.Certificate, err error) {
	certs = []*voc.Certificate{}
	// initialize certificates client
	if err = d.initCertificatesClient(kv); err != nil {
		err = fmt.Errorf("could not initialize certificates client: %v", err)
		return
	}

	pager := d.clients.certificationsClient.NewListCertificatePropertiesPager(&azcertificates.ListCertificatePropertiesOptions{})
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			log.Errorf("Could not page next page (paging error), probably you do not have the right permissions: %v", err)
			return certs, nil
		}
		for _, c := range page.Value {
			cert := &voc.Certificate{
				Resource: discovery.NewResource(d,
					voc.ResourceID(util.Deref(c.ID)),
					getCertificateName(c.ID),
					c.Attributes.Created,
					voc.GeoLocation{Region: util.Deref(kv.Location)},
					labels(c.Tags),
					voc.ResourceID(util.Deref(kv.ID)),
					voc.CertificateType,
					c, kv),
				Enabled:        util.Deref(c.Attributes.Enabled),
				ActivationDate: convertTime(c.Attributes.NotBefore),
				ExpirationDate: convertTime(c.Attributes.Expires),
			}
			certs = append(certs, cert)
		}
	}
	return

}

func getCertificateName(id *azcertificates.ID) (certName string) {
	certID := string(util.Deref(id))
	splits := strings.Split(certID, "/certificates/")
	if len(splits) < 2 {
		log.Errorf("Could not split certificate with ID '%s'. Probably it is wrongly formatted. Using ID instead.", certID)
		certName = certID
		return
	}
	certName = splits[len(splits)-1]
	return
}

// convertTime converts a time pointer (given via Azure SDK) to a Unix time - handling nil pointers appropriately
func convertTime(t *time.Time) int64 {
	if t == nil {
		return -1
	} else {
		return util.Deref(t).Unix()
	}
}

// TODO(lebogg): How to define the range/scope of key types in the ontology?
// TODO(lebogg): Extract these returned options in a type or const to internal/api s.t. all discoverers use the same values
func getKeyType(kt *armkeyvault.JSONWebKeyType) string {
	switch util.Deref(kt) {
	case armkeyvault.JSONWebKeyTypeEC:
		return "EC"
	case armkeyvault.JSONWebKeyTypeECHSM:
		return "EC"
	case armkeyvault.JSONWebKeyTypeRSA:
		return "RSA"
	case armkeyvault.JSONWebKeyTypeRSAHSM:
		return "RSA"
	}
	// In the future, there could be new types not handled so far. Return it anyway but warn in console.
	keyType := string(util.Deref(kt))
	log.Warnf("This key is not supported yet: '%s'. Probably, metrics won't work properly.", keyType)
	return keyType
}
