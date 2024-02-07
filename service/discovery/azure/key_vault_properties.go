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
	"clouditor.io/clouditor/internal/util"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
)

// TODO(lebogg): Add property functions used by the discover and handle parts here, e.g. `getKeys`

// TODO(lebogg): Test it
func hasPublicAccess(kv *armkeyvault.Vault) (hasPublicAccess bool) {
	if kv.Properties != nil {
		hasPublicAccess = util.Deref(kv.Properties.PublicNetworkAccess) == "Enabled"
	}
	return
}

// TODO(lebogg): Test it and resolve TODOs
func (d *azureDiscovery) isKeyVaultActive(kv *armkeyvault.Vault) (bool, error) {
	// Query Metric "API Hits" to derive the Key Vault's amount of traffic
	metrics, err := d.clients.metricsClient.QueryResource(context.TODO(), util.Deref(kv.ID),
		&azquery.MetricsClientQueryResourceOptions{
			Interval:    util.Ref("P1D"),
			MetricNames: util.Ref("ServiceApiHit"),
		})
	if err != nil {
		// TODO(lebogg): To Test: Maybe there are resources (in this case, key vaults) where no API Hit is defined -> Then it is not an error but, e.g., false?
		return false, fmt.Errorf("could not query resource for metric (Monitoring): %v", err)
	}
	// Check if value is non-nil (shouldn't be)
	if metrics.Value == nil {
		return false, fmt.Errorf("something went wrong. There are no value(s) for this metric")
	}
	// We only asked for one metric, so we should only get one value
	if l := len(metrics.Value); l != 1 {
		return false, fmt.Errorf("we got %d metrics. But should be one", l)
	}

	// Determine if there is enough traffic s.t. Key Vault is active
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
