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

package io.clouditor.discovery.aws;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.doReturn;
import static org.mockito.Mockito.when;

import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentMatchers;
import software.amazon.awssdk.services.config.ConfigClient;
import software.amazon.awssdk.services.config.model.BaseConfigurationItem;
import software.amazon.awssdk.services.config.model.BatchGetResourceConfigRequest;
import software.amazon.awssdk.services.config.model.BatchGetResourceConfigResponse;
import software.amazon.awssdk.services.config.model.ListDiscoveredResourcesRequest;
import software.amazon.awssdk.services.config.model.ListDiscoveredResourcesResponse;
import software.amazon.awssdk.services.config.model.ResourceIdentifier;
import software.amazon.awssdk.services.config.model.ResourceKey;

class AwsConfigScannerTest extends AwsScannerTest {

  @BeforeAll
  static void setUpOnce() throws IOException {
    discoverAssets(
        ConfigClient.class,
        AwsConfigScanner::new,
        api -> {
          when(api.listDiscoveredResources(
                  ArgumentMatchers.any(ListDiscoveredResourcesRequest.class)))
              .thenReturn(ListDiscoveredResourcesResponse.builder().build());

          doReturn(
                  ListDiscoveredResourcesResponse.builder()
                      .resourceIdentifiers(ResourceIdentifier.builder().resourceId("rds-1").build())
                      .build())
              .when(api)
              .listDiscoveredResources(
                  ListDiscoveredResourcesRequest.builder()
                      .resourceType("AWS::RDS::DBInstance")
                      .build());

          doReturn(
                  ListDiscoveredResourcesResponse.builder()
                      .resourceIdentifiers(ResourceIdentifier.builder().resourceId("vm-1").build())
                      .build())
              .when(api)
              .listDiscoveredResources(
                  ListDiscoveredResourcesRequest.builder().resourceType("AWS::EC2::Host").build());

          doReturn(
                  ListDiscoveredResourcesResponse.builder()
                      .resourceIdentifiers(
                          ResourceIdentifier.builder().resourceId("user-1").build())
                      .build())
              .when(api)
              .listDiscoveredResources(
                  ListDiscoveredResourcesRequest.builder().resourceType("AWS::IAM::User").build());

          doReturn(
                  ListDiscoveredResourcesResponse.builder()
                      .resourceIdentifiers(
                          ResourceIdentifier.builder().resourceId("bucket-1").build())
                      .build())
              .when(api)
              .listDiscoveredResources(
                  ListDiscoveredResourcesRequest.builder().resourceType("AWS::S3::Bucket").build());

          when(api.batchGetResourceConfig(
                  ArgumentMatchers.any(BatchGetResourceConfigRequest.class)))
              .thenReturn(BatchGetResourceConfigResponse.builder().build());

          doReturn(
                  BatchGetResourceConfigResponse.builder()
                      .baseConfigurationItems(
                          BaseConfigurationItem.builder()
                              .awsRegion("eu-central-1")
                              .arn("eu-arn")
                              .build())
                      .build())
              .when(api)
              .batchGetResourceConfig(
                  BatchGetResourceConfigRequest.builder()
                      .resourceKeys(
                          ResourceKey.builder()
                              .resourceId("rds-1")
                              .resourceType("AWS::RDS::DBInstance")
                              .build())
                      .build());

          doReturn(
                  BatchGetResourceConfigResponse.builder()
                      .baseConfigurationItems(
                          BaseConfigurationItem.builder()
                              .awsRegion("eu-west-1")
                              .arn("eu-west-arn")
                              .build())
                      .build())
              .when(api)
              .batchGetResourceConfig(
                  BatchGetResourceConfigRequest.builder()
                      .resourceKeys(
                          ResourceKey.builder()
                              .resourceId("vm-1")
                              .resourceType("AWS::EC2::Host")
                              .build())
                      .build());

          doReturn(
                  BatchGetResourceConfigResponse.builder()
                      .baseConfigurationItems(
                          BaseConfigurationItem.builder()
                              .awsRegion("global")
                              .arn("some-global-arn")
                              .build())
                      .build())
              .when(api)
              .batchGetResourceConfig(
                  BatchGetResourceConfigRequest.builder()
                      .resourceKeys(
                          ResourceKey.builder()
                              .resourceId("user-1")
                              .resourceType("AWS::IAM::User")
                              .build())
                      .build());

          doReturn(
                  BatchGetResourceConfigResponse.builder()
                      .baseConfigurationItems(
                          BaseConfigurationItem.builder()
                              .awsRegion("us-east-1")
                              .arn("us-arn")
                              .build())
                      .build())
              .when(api)
              .batchGetResourceConfig(
                  BatchGetResourceConfigRequest.builder()
                      .resourceKeys(
                          ResourceKey.builder()
                              .resourceId("bucket-1")
                              .resourceType("AWS::S3::Bucket")
                              .build())
                      .build());
        });
  }

  @Test
  void testAllowedRegions() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/config/resources-region-eu.md"));

    assertNotNull(rule);

    var resource = assets.get("eu-arn");

    assertNotNull(resource);
    assertTrue(rule.evaluate(resource).isOk());

    resource = assets.get("eu-west-arn");

    assertNotNull(resource);
    assertTrue(rule.evaluate(resource).isOk());

    resource = assets.get("some-global-arn");

    assertNotNull(resource);
    assertTrue(rule.evaluate(resource).isOk());

    resource = assets.get("us-arn");

    assertNotNull(resource);
    assertFalse(rule.evaluate(resource).isOk());
  }
}
