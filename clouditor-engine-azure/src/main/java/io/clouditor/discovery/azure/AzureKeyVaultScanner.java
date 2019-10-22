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
