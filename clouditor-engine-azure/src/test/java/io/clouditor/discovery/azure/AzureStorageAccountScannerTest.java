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
