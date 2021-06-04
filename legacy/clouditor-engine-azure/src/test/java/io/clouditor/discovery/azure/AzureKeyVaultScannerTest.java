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
            .loadRule(
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
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/azure/keyvault/vault-logging.md"));

    var vault1 = assets.get("vault-with-expiry");

    assertTrue(rule.evaluate(vault1).isOk());
  }
}
