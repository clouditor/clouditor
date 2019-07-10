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
import org.mockito.ArgumentMatchers;
import software.amazon.awssdk.services.kms.KmsClient;
import software.amazon.awssdk.services.kms.model.DescribeKeyRequest;
import software.amazon.awssdk.services.kms.model.DescribeKeyResponse;
import software.amazon.awssdk.services.kms.model.GetKeyPolicyRequest;
import software.amazon.awssdk.services.kms.model.GetKeyPolicyResponse;
import software.amazon.awssdk.services.kms.model.GetKeyRotationStatusRequest;
import software.amazon.awssdk.services.kms.model.GetKeyRotationStatusResponse;
import software.amazon.awssdk.services.kms.model.KeyListEntry;
import software.amazon.awssdk.services.kms.model.KeyManagerType;
import software.amazon.awssdk.services.kms.model.KeyMetadata;
import software.amazon.awssdk.services.kms.model.ListKeysResponse;
import software.amazon.awssdk.services.kms.model.OriginType;

class AwsKmsScannerTest extends AwsScannerTest {

  @BeforeAll
  static void setUpOnce() throws IOException {
    discoverAssets(
        KmsClient.class,
        AwsKmsScanner::new,
        api -> {
          when(api.listKeys())
              .thenReturn(
                  ListKeysResponse.builder()
                      .keys(
                          KeyListEntry.builder().keyArn("key1").keyId("key1").build(),
                          KeyListEntry.builder().keyArn("key2").keyId("key2").build(),
                          KeyListEntry.builder().keyArn("key3").keyId("key3").build())
                      .build());

          when(api.describeKey(DescribeKeyRequest.builder().keyId("key1").build()))
              .thenReturn(
                  DescribeKeyResponse.builder()
                      .keyMetadata(
                          KeyMetadata.builder()
                              .keyId("key1")
                              .arn("key1")
                              .origin(OriginType.EXTERNAL)
                              .build())
                      .build());

          when(api.getKeyRotationStatus(
                  GetKeyRotationStatusRequest.builder().keyId("key1").build()))
              .thenReturn(GetKeyRotationStatusResponse.builder().keyRotationEnabled(true).build());

          when(api.describeKey(DescribeKeyRequest.builder().keyId("key2").build()))
              .thenReturn(
                  DescribeKeyResponse.builder()
                      .keyMetadata(
                          KeyMetadata.builder()
                              .keyId("key2")
                              .arn("key2")
                              .origin(OriginType.AWS_KMS)
                              .build())
                      .build());

          when(api.getKeyRotationStatus(
                  GetKeyRotationStatusRequest.builder().keyId("key2").build()))
              .thenReturn(GetKeyRotationStatusResponse.builder().keyRotationEnabled(false).build());

          when(api.describeKey(DescribeKeyRequest.builder().keyId("key3").build()))
              .thenReturn(
                  DescribeKeyResponse.builder()
                      .keyMetadata(
                          KeyMetadata.builder()
                              .keyId("key3")
                              .arn("key3")
                              .origin(OriginType.AWS_KMS)
                              .keyManager(KeyManagerType.AWS)
                              .build())
                      .build());

          when(api.getKeyRotationStatus(
                  GetKeyRotationStatusRequest.builder().keyId("key3").build()))
              .thenReturn(GetKeyRotationStatusResponse.builder().keyRotationEnabled(false).build());

          when(api.getKeyPolicy(ArgumentMatchers.any(GetKeyPolicyRequest.class)))
              .thenReturn(GetKeyPolicyResponse.builder().policy("my-policy").build());
        });
  }

  @Test
  void testExternalOrigin() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/kms/key-origin-external.md"));

    assertNotNull(rule);

    var key1 = assets.get("key1");

    assertNotNull(key1);
    assertTrue(rule.evaluate(key1).isOk());

    var key2 = assets.get("key2");

    assertNotNull(key2);
    assertFalse(rule.evaluate(key2).isOk());
  }

  @Test
  void testKMSOrigin() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/kms/key-origin-kms.md"));

    assertNotNull(rule);

    var key1 = assets.get("key1");

    assertNotNull(key1);
    assertFalse(rule.evaluate(key1).isOk());

    var key2 = assets.get("key2");

    assertNotNull(key2);
    assertTrue(rule.evaluate(key2).isOk());
  }

  @Test
  void testKeyRotationEnabled() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/kms/key-origin-kms.md"));

    assertNotNull(rule);

    var key1 = assets.get("key1");

    assertNotNull(key1);
    assertFalse(rule.evaluate(key1).isOk());

    var key2 = assets.get("key2");

    assertNotNull(key2);
    assertTrue(rule.evaluate(key2).isOk());
  }
}
