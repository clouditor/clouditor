package azure

import (
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
)

type azureKeyVaultDiscovery struct {
	*azureDiscovery
	// TODO(lebogg): Don't know if we need these defenderProperties here as well
	//defenderProperties map[string]*defenderProperties
}

func NewKeyVaultDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureKeyVaultDiscovery{
		azureDiscovery: &azureDiscovery{
			clientOptions: arm.ClientOptions{},
			// Todo(all): What do we need this for?
			discovererComponent: KeyVaultComponent,
			// Todo(lebogg): In service/discovery/discovery.go:257 csID is set anyway. Maybe for testing?
			csID:      discovery.DefaultCloudServiceID,
			backupMap: make(map[string]*backup),
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
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	// Discover Key Vaults
	log.Info("Discover Azure Key Vaults")
	keyVaults, err := d.discoverKeyVaults()
	if err != nil {
		return nil, fmt.Errorf("could not discover Key Vaults: %w", err)
	}
	list = append(list, keyVaults...)

	return
}

// TODO(lebogg): Finished here last time. Not tested, yet. Next: Implement handleKeyVault
func (d *azureKeyVaultDiscovery) discoverKeyVaults() (list []voc.IsCloudResource, err error) {
	// initialize key vault client
	if err = d.initKeyVaultClient(); err != nil {
		return nil, err
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
			keyVault, err := d.handleKeyVault(kv)
			if err != nil {
				return fmt.Errorf("could not handle key vault: %w", err)
			}

			log.Infof("Adding block storage '%s'", keyVault.GetName())

			list = append(list, keyVault)
			return nil
		})
	if err != nil {
		list = nil
		return
	}
	return
}

func (d *azureKeyVaultDiscovery) initKeyVaultClient() (err error) {
	d.clients.keyVaultClient, err = initClient(d.clients.keyVaultClient, d.azureDiscovery, armkeyvault.NewVaultsClient)
	return
}

func (d *azureKeyVaultDiscovery) handleKeyVault(kv *armkeyvault.Vault) (*voc.KeyVault, error) {
	//TODO
	return nil, nil
}
