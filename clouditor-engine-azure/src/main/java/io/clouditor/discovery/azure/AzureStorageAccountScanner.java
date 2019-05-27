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

import com.microsoft.azure.management.resources.fluentcore.arm.models.HasId;
import com.microsoft.azure.management.resources.fluentcore.arm.models.HasName;
import com.microsoft.azure.management.storage.StorageAccount;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.AssetProperties;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.List;

@ScannerInfo(assetType = "StorageAccount", group = "Azure", service = "Storage")
public class AzureStorageAccountScanner extends AzureScanner<StorageAccount> {

  public AzureStorageAccountScanner() {
    super(HasId::id, HasName::name);
  }

  @Override
  protected List<StorageAccount> list() {
    return this.api.azure().storageAccounts().list();
  }

  @Override
  protected Asset transform(StorageAccount account) throws ScanException {
    var asset = super.transform(account);

    // TODO: support maps instead of just lists
    asset.setProperty(
        "encryptionStatuses",
        MAPPER.convertValue(account.encryptionStatuses(), AssetProperties.class));

    var isKeyRegenerated = false;

    /*List<EventData> accountLogs =
        this.api
            .monitor()
            .activityLogs()
            .defineQuery()
            .startingFrom(new DateTime().minus(90 * 24 * 3600 * 1000L))
            .endsBefore(new DateTime())
            .withAllPropertiesInResponse()
            .filterByResource(account.id())
            .execute();

    for (var data : accountLogs) {
      if (data.operationName()
          .value()
          .equals("Microsoft.Storage/storageAccounts/regenerateKey/action")) {
        isKeyRegenerated = true;
        break;
      }
    }

    map.put("keyRegenerated", isKeyRegenerated);*/

    return asset;
  }
}
