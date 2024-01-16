package azure

import (
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

// handleResourceGroup returns a [voc.Account] out of an existing [armsubscription.Subscription].
func (d *azureDiscovery) handleResourceGroup(rg *armresources.ResourceGroup) voc.IsCloudResource {
	return &voc.ResourceGroup{
		Resource: discovery.NewResource(
			d,
			voc.ResourceID(util.Deref(rg.ID)),
			util.Deref(rg.Name),
			nil,
			voc.GeoLocation{
				Region: util.Deref(rg.Location),
			},
			labels(rg.Tags),
			voc.ResourceID(*d.sub.ID),
			voc.ResourceGroupType,
		),
	}
}

// handleSubscription returns a [voc.Account] out of an existing [armsubscription.Subscription].
func (d *azureDiscovery) handleSubscription(s *armsubscription.Subscription) *voc.Account {
	return &voc.Account{
		Resource: discovery.NewResource(
			d,
			voc.ResourceID(util.Deref(s.ID)),
			util.Deref(s.DisplayName),
			// subscriptions do not have a creation date
			nil,
			// subscriptions are global
			voc.GeoLocation{},
			// subscriptions do not have labels,
			nil,
			// subscriptions are the top-most item and have no parent,
			"",
			voc.AccountType,
		),
	}
}
