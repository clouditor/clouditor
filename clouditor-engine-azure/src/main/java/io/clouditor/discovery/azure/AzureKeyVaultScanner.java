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

import com.microsoft.azure.management.keyvault.Key;
import com.microsoft.azure.management.keyvault.Secret;
import com.microsoft.azure.management.keyvault.Vault;
import com.microsoft.azure.management.monitor.DiagnosticSetting;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.AssetProperties;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.stream.Collectors;

@ScannerInfo(assetType = "KeyVault", group = "Azure", service = "Key Vaults")
public class AzureKeyVaultScanner extends AzureScanner<Vault> {

  public AzureKeyVaultScanner() {
    super(Vault::id, Vault::name);
  }

  @Override
  protected List<Vault> list() {
    return this.listVaultsBySubscription();
  }

  private List<Vault> listVaultsBySubscription() {
    // for some reason Vaults does not directly expose the listBySubscription of VaultsImpl, so we
    // have to loop over all resource groups
    if (this.resourceGroup == null) {
      List<Vault> vaults = new ArrayList<>();
      for (var group : this.api.azure().resourceGroups().list()) {
        vaults.addAll(this.api.azure().vaults().listByResourceGroup(group.name()));
      }

      return vaults;
    } else {
      return this.api.azure().vaults().listByResourceGroup(this.resourceGroup);
    }
  }

  @Override
  protected Asset transform(Vault vault) throws ScanException {
    var asset = super.transform(vault);

    enrichList(asset, "keys", vault, x -> vault.keys().list(), Key::id, Key::name);

    enrichList(asset, "secrets", vault, x -> vault.secrets().list(), Secret::id, Secret::name);

    asset.setProperty(
        "logs",
        this.api.monitor().diagnosticSettings().listByResource(vault.id()).stream()
            .map(DiagnosticSetting::logs)
            .flatMap(Collection::stream)
            .map(log -> MAPPER.convertValue(log, AssetProperties.class))
            .collect(Collectors.toList()));

    return asset;
  }
}
