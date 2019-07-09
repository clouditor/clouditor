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
