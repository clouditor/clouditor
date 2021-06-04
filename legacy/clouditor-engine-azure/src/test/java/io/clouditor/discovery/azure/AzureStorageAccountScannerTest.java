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
 */

package io.clouditor.discovery.azure;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

import com.microsoft.azure.management.storage.StorageService;
import com.microsoft.azure.management.storage.implementation.StorageAccountInner;
import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import java.util.Map;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

class AzureStorageAccountScannerTest extends AzureScannerTest {

  @BeforeAll
  static void setUpOnce() {
    discoverAssets(
        AzureStorageAccountScanner::new,
        api -> {
          var account =
              createStorageAccount(
                  "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/mock/providers/Microsoft.Storage/storageAccounts/account1",
                  new StorageAccountInner(),
                  Map.of(StorageService.BLOB, true));

          when(api.azure.storageAccounts().list()).thenReturn(MockedPagedList.of(account));
        });
  }

  @Test
  void testEnforceSecureTransfer() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource(
                        "rules/azure/storage/storage-account-enforce-secure-transfer.md"));

    var account =
        assets.get(
            "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/mock/providers/Microsoft.Storage/storageAccounts/account1");

    assertNotNull(account);
    assertFalse(rule.evaluate(account).isOk());
  }

  @Test
  void testEncryption() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/azure/storage/storage-account-blob-encrypted.md"));

    var account =
        assets.get(
            "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/mock/providers/Microsoft.Storage/storageAccounts/account1");

    assertNotNull(account);
    assertTrue(rule.evaluate(account).isOk());
  }
}
