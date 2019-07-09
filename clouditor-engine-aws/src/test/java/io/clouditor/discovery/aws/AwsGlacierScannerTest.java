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

package io.clouditor.discovery.aws;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import software.amazon.awssdk.awscore.exception.AwsServiceException;
import software.amazon.awssdk.services.glacier.GlacierClient;
import software.amazon.awssdk.services.glacier.model.DescribeVaultOutput;
import software.amazon.awssdk.services.glacier.model.GetVaultNotificationsRequest;
import software.amazon.awssdk.services.glacier.model.GetVaultNotificationsResponse;
import software.amazon.awssdk.services.glacier.model.ListVaultsResponse;
import software.amazon.awssdk.services.glacier.model.VaultNotificationConfig;

class AwsGlacierScannerTest extends AwsScannerTest {

  @BeforeAll
  static void setUpOnce() throws IOException {
    discoverAssets(
        GlacierClient.class,
        AwsGlacierScanner::new,
        api -> {
          when(api.listVaults())
              .thenReturn(
                  ListVaultsResponse.builder()
                      .vaultList(
                          DescribeVaultOutput.builder()
                              .vaultARN("arn:aws:glacier:us-west-2:012345678901:vaults/vault1")
                              .vaultName("vault1")
                              .build(),
                          DescribeVaultOutput.builder()
                              .vaultARN("arn:aws:glacier:us-west-2:012345678901:vaults/vault2")
                              .vaultName("vault2")
                              .build())
                      .build());

          when(api.getVaultNotifications(
                  GetVaultNotificationsRequest.builder().vaultName("vault1").build()))
              .thenReturn(
                  GetVaultNotificationsResponse.builder()
                      .vaultNotificationConfig(
                          VaultNotificationConfig.builder().snsTopic("some-topic").build())
                      .build());

          // better would be a catch-all
          when(api.getVaultNotifications(
                  GetVaultNotificationsRequest.builder().vaultName("vault2").build()))
              .thenThrow(AwsServiceException.builder().statusCode(404).build());
        });
  }

  @Test
  void testNotifications() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/glacier/vault-notifications.md"));

    assertNotNull(rule);

    var vault = assets.get("arn:aws:glacier:us-west-2:012345678901:vaults/vault1");

    assertTrue(rule.evaluate(vault).isOk());

    vault = assets.get("arn:aws:glacier:us-west-2:012345678901:vaults/vault2");

    assertFalse(rule.evaluate(vault).isOk());
  }
}
