/*
 * Copyright (c) 2016-2019, Fraunhofer AISEC. All rights reserved.
 *
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
 *
 * Clouditor Community Edition is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Clouditor Community Edition is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * long with Clouditor Community Edition.  If not, see <https://www.gnu.org/licenses/>
 */

package io.clouditor.discovery.azure;

import com.microsoft.azure.management.monitor.implementation.LogProfileResourceInner;
import com.microsoft.azure.management.network.ProvisioningState;
import com.microsoft.azure.management.resources.Subscription;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.AssetProperties;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.List;

@ScannerInfo(assetType = "Subscription", group = "Azure", service = "Account")
public class AzureSubscriptionScanner extends AzureScanner<Subscription> {

  AzureSubscriptionScanner() {
    super(Subscription::subscriptionId, Subscription::displayName);
  }

  @Override
  protected List<Subscription> list() {
    return this.api.azure().subscriptions().list();
  }

  @Override
  protected Asset transform(Subscription subscription) throws ScanException {
    var asset = super.transform(subscription);

    enrichList(
        asset,
        "logProfiles",
        subscription,
        x -> this.api.monitor().inner().logProfiles().list(),
        null,
        LogProfileResourceInner::id,
        LogProfileResourceInner::name);

    var regions = new AssetProperties();

    // Get the available locations
    for (var location : this.api.azure().getCurrentSubscription().listLocations()) {
      var properties = new AssetProperties();
      properties.put("enabled", false);

      regions.put(location.region().name(), properties);
    }

    // These regions are not selectable in the dashboard, but they are listed as available
    // regions for the subscription, so they need to be removed for the network watcher check
    regions.remove("francesouth");
    regions.remove("southafricawest");
    regions.remove("australiacentral");
    regions.remove("australiacentral2");

    // Get all active Network Watchers
    for (var networkWatcher : this.api.azure().networkWatchers().list()) {
      if (networkWatcher.inner().provisioningState().equals(ProvisioningState.SUCCEEDED)) {
        var properties = new AssetProperties();
        properties.put("enabled", true);

        regions.put(networkWatcher.regionName(), properties);
      }
    }

    asset.setProperty("watchers", regions);

    return asset;
  }
}
