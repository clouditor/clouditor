package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/graphrbac/graphrbac"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/voc"
)

type azureIdentityDiscovery struct {
	*azureDiscovery
}

func (*azureIdentityDiscovery) Name() string {
	return "Azure Identity"
}

func (*azureIdentityDiscovery) Description() string {
	return "Discovery Azure identities."
}

func NewAzureIdentityDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureIdentityDiscovery{
		&azureDiscovery{
			discovererComponent: IdentityComponent,
			csID:                discovery.DefaultCloudServiceID,
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(d.azureDiscovery)
	}

	return d
}

func (d *azureIdentityDiscovery) List() (list []voc.IsCloudResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	log.Info("Discover Azure identities")
	identities, err := d.discoverIdentities()

	if err != nil {
		return nil, fmt.Errorf("could not discover identities: %w", err)
	}

	list = append(list, identities...)

	return
}

func (d *azureIdentityDiscovery) discoverIdentities() ([]voc.IsCloudResource, error) {

	var list []voc.IsCloudResource

	// initialize the identity client
	if err := d.initIdentityClient(); err != nil {
		return nil, err
	}

	result, err := d.clients.identityClient.List(context.TODO(), "", "")

	for _, user := range result.Values() {
		log.Infof("Adding user '%s'", user.DisplayName)
		list = append(list, d.handleIdentities(user))
	}

	if err != nil {
		return nil, err
	}

	return list, nil
}

func (d *azureIdentityDiscovery) handleIdentities(identity graphrbac.User) voc.IsCloudResource {
	return &voc.Identity{
		Identifiable: &voc.Identifiable{
			Resource: &voc.Resource{
				ID:           "",
				ServiceID:    "",
				Name:         *identity.DisplayName,
				CreationTime: 0,
				Type:         nil,
				GeoLocation:  voc.GeoLocation{},
				Labels:       nil,
			},
			Authenticity:  nil,
			Authorization: nil,
			Activated:     false,
		},
		Authenticities:        nil,
		Privileged:            false,
		LastActivity:          LastActivity(identity),
		DisablePasswordPolicy: false,
	}
}

func LastActivity(i graphrbac.User) time.Time {
	//this need Azure AD Premium, P1 or P2
	//log.Infof(i.AdditionalProperties["signInActivity"].(string))

	return time.Time{}
}

func (d *azureIdentityDiscovery) initIdentityClient() (err error) {
	d.clients.identityClient = graphrbac.NewUsersClient("")

	return
}
