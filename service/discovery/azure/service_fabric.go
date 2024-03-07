package azure

import (
	"context"
	"fmt"

	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/internal/util"
	"clouditor.io/clouditor/voc"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicefabric/armservicefabric"
)

type azureServiceFabricDiscovery struct {
	*azureDiscovery
	defenderProperties map[string]*defenderProperties
}

func NewFabricServiceDiscovery(opts ...DiscoveryOption) discovery.Discoverer {
	d := &azureServiceFabricDiscovery{
		&azureDiscovery{
			discovererComponent: ServiceFabricComponent,
			csID:                discovery.DefaultCloudServiceID,
			backupMap:           make(map[string]*backup),
		},
		make(map[string]*defenderProperties),
	}

	// Apply options
	for _, opt := range opts {
		opt(d.azureDiscovery)
	}

	return d
}

func (*azureServiceFabricDiscovery) Name() string {
	return "Service Fabric"
}

func (*azureServiceFabricDiscovery) Description() string {
	return "Discovery for Service Fabric Clusters."
}

// List service fabric resources
func (d *azureServiceFabricDiscovery) List() (list []voc.IsCloudResource, err error) {
	if err = d.authorize(); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrCouldNotAuthenticate, err)
	}

	log.Info("Discover Azure Service Fabric resources")
	// Discover Azure Service Fabric Clusters
	storage, err := d.discoverClusters()
	if err != nil {
		return nil, fmt.Errorf("could not discover service fabric clusters: %w", err)
	}
	list = append(list, storage...)

	// Add backup block storages
	if d.backupMap[DataSourceTypeDisc] != nil && d.backupMap[DataSourceTypeDisc].backupStorages != nil {
		list = append(list, d.backupMap[DataSourceTypeDisc].backupStorages...)
	}

	return
}

func (d *azureServiceFabricDiscovery) initClusterClient() (err error) {
	d.clients.fabricsServiceClusterClient, err = initClient(d.clients.fabricsServiceClusterClient, d.azureDiscovery,
		armservicefabric.NewClustersClient)
	return
}

func (d *azureServiceFabricDiscovery) discoverClusters() ([]voc.IsCloudResource, error) {
	var (
		list []voc.IsCloudResource
	)

	// initialize backup policies client
	if err := d.initClusterClient(); err != nil {
		return nil, err
	}

	var clusters []*armservicefabric.Cluster
	if util.Deref(d.rg) != "" {
		response, err := d.clients.fabricsServiceClusterClient.List(context.Background(), &armservicefabric.ClustersClientListOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not get fabric service clusters: %v", err)
		}
		clusters = append(clusters, response.Value...)
	} else {
		response, err := d.clients.fabricsServiceClusterClient.ListByResourceGroup(context.Background(),
			util.Deref(d.rg), &armservicefabric.ClustersClientListByResourceGroupOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not get fabric service clusters: %v", err)
		}
		clusters = append(clusters, response.Value...)
	}
	for _, c := range clusters {
		var r = &voc.Redundancy{}

		if c.Properties.VmssZonalUpgradeMode != nil {
			r.Zone = true
		}
		list = append(list, &voc.Cluster{
			Resource: discovery.NewResource(d,
				voc.ResourceID(util.Deref(c.ID)),
				util.Deref(c.Name),
				// No creation time available
				nil,
				voc.GeoLocation{
					Region: util.Deref(c.Location),
				},
				labels(c.Tags),
				resourceGroupID(c.ID),
				voc.ClusterType,
				c),
			Redundancy: r,
		})
	}

	return list, nil
}
