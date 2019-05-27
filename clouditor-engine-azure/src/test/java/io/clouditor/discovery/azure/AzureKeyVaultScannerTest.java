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
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.when;

import com.microsoft.azure.keyvault.models.KeyAttributes;
import com.microsoft.azure.keyvault.models.KeyBundle;
import com.microsoft.azure.management.keyvault.Vault;
import com.microsoft.azure.management.keyvault.implementation.VaultInner;
import com.microsoft.azure.management.monitor.DiagnosticSetting;
import com.microsoft.azure.management.monitor.LogSettings;
import com.microsoft.azure.management.monitor.RetentionPolicy;
import com.microsoft.azure.management.monitor.implementation.DiagnosticSettingsResourceInner;
import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import java.util.List;
import org.joda.time.DateTime;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

class AzureKeyVaultScannerTest extends AzureScannerTest {

  @BeforeAll
  static void setUpOnce() {
    discoverAssets(
        AzureKeyVaultScanner::new,
        api -> {
          var vault1 = createWithId(Vault.class, "vault-with-expiry", new VaultInner());

          var key =
              createKey(
                  "key",
                  "key-name",
                  new KeyBundle()
                      .withAttributes(
                          (KeyAttributes)
                              new KeyAttributes().withExpires(new DateTime().plusWeeks(30))));

          when(vault1.keys().list()).thenReturn(MockedPagedList.of(key));

          var vault2 = createWithId(Vault.class, "vault-without-expiry", new VaultInner());

          key = createKey("key", "key-name", new KeyBundle());

          when(vault2.keys().list()).thenReturn(MockedPagedList.of(key));

          when(api.azure.vaults().listByResourceGroup(anyString()))
              .thenReturn(MockedPagedList.of(vault1, vault2));

          var settings =
              createDiagnosticsSetting(
                  "some-id",
                  "some-name",
                  new DiagnosticSettingsResourceInner()
                      .withLogs(
                          List.of(
                              new LogSettings()
                                  .withEnabled(true)
                                  .withRetentionPolicy(
                                      new RetentionPolicy().withEnabled(true).withDays(270)))));

          when(api.monitor().diagnosticSettings().listByResource(anyString()))
              .thenReturn(MockedPagedList.of(settings));
        });
  }

  private static DiagnosticSetting createDiagnosticsSetting(
      String id, String name, DiagnosticSettingsResourceInner inner) {
    var settings = createWithIdAndName(DiagnosticSetting.class, id, name, inner);

    when(settings.logs()).thenReturn(inner.logs());

    return settings;
  }

  @Test
  void testKeyExpiry() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .parseFromMarkDown(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/azure/keyvault/key-expiry.md"));

    var vault1 = assets.get("vault-with-expiry");

    assertTrue(rule.evaluate(vault1).isOk());

    var vault2 = assets.get("vault-without-expiry");

    assertFalse(rule.evaluate(vault2).isOk());
  }

  @Test
  void testLogging() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .parseFromMarkDown(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/azure/keyvault/vault-logging.md"));

    var vault1 = assets.get("vault-with-expiry");

    assertTrue(rule.evaluate(vault1).isOk());
  }
}
