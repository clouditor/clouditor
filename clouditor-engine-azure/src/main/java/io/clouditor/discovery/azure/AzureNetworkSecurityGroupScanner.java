/*
 * Copyright 2016-2019 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *            $$\                           $$\ $$\   $$\
 *            $$ |                          $$ |\__|  $$ |
 *   $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 *  $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 *  $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ |  \__|
 *  $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 *  \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *   \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package io.clouditor.discovery.azure;

import com.microsoft.azure.management.network.NetworkSecurityGroup;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.List;

@ScannerInfo(assetType = "NetworkSecurityGroup", group = "Azure", service = "Networking")
public class AzureNetworkSecurityGroupScanner extends AzureScanner<NetworkSecurityGroup> {

  public AzureNetworkSecurityGroupScanner() {
    super(NetworkSecurityGroup::id, NetworkSecurityGroup::name);
  }

  @Override
  protected List<NetworkSecurityGroup> list() {
    return this.resourceGroup != null
        ? this.api.azure().networkSecurityGroups().listByResourceGroup(this.resourceGroup)
        : this.api.azure().networkSecurityGroups().list();
  }

  @Override
  protected Asset transform(NetworkSecurityGroup nsg) throws ScanException {
    var asset = super.transform(nsg);

    var watcher =
        this.api
            .azure()
            .networkWatchers()
            .getById(
                "/subscriptions/"
                    + this.api.azure().subscriptionId()
                    + "/resourceGroups/NetworkWatcherRG/providers/Microsoft.Network/networkWatchers/NetworkWatcher_"
                    + nsg.regionName());

    if (watcher != null) {
      // this needs the Network Contributor role!
      enrich(
          asset,
          "flowLogSettings",
          nsg,
          x -> watcher.getFlowLogSettings(nsg.id()),
          x -> null,
          x -> null);
    }

    return asset;
  }
}
