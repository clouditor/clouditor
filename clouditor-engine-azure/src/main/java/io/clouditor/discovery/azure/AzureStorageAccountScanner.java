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
