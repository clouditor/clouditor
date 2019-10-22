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

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.mockito.ArgumentMatchers;
import software.amazon.awssdk.services.cloudtrail.CloudTrailClient;
import software.amazon.awssdk.services.cloudtrail.model.DescribeTrailsResponse;
import software.amazon.awssdk.services.cloudtrail.model.EventSelector;
import software.amazon.awssdk.services.cloudtrail.model.GetEventSelectorsRequest;
import software.amazon.awssdk.services.cloudtrail.model.GetEventSelectorsResponse;
import software.amazon.awssdk.services.cloudtrail.model.GetTrailStatusRequest;
import software.amazon.awssdk.services.cloudtrail.model.GetTrailStatusResponse;
import software.amazon.awssdk.services.cloudtrail.model.ReadWriteType;
import software.amazon.awssdk.services.cloudtrail.model.Trail;

class AwsCloudTrailAccountScannerTest extends AwsScannerTest {

  @BeforeAll
  static void setUpOnce() throws IOException {
    discoverAssets(
        CloudTrailClient.class,
        AwsCloudTrailAccountScanner::new,
        api -> {
          when(api.describeTrails())
              .thenReturn(
                  DescribeTrailsResponse.builder()
                      .trailList(
                          Trail.builder()
                              .trailARN(
                                  "arn:aws:cloudtrail:eu-central-1:123456789:trail/management-trail-1")
                              .name("management-trail-1")
                              .isMultiRegionTrail(true)
                              .build(),
                          Trail.builder()
                              .trailARN(
                                  "arn:aws:cloudtrail:eu-central-1:123456789:trail/management-trail-2")
                              .name("management-trail-2")
                              .isMultiRegionTrail(true)
                              .build())
                      .build());

          when(api.getTrailStatus(ArgumentMatchers.any(GetTrailStatusRequest.class)))
              .thenReturn(GetTrailStatusResponse.builder().isLogging(true).build());

          when(api.getEventSelectors(
                  GetEventSelectorsRequest.builder().trailName("management-trail-1").build()))
              .thenReturn(
                  GetEventSelectorsResponse.builder()
                      .eventSelectors(
                          EventSelector.builder()
                              .includeManagementEvents(true)
                              .readWriteType(ReadWriteType.WRITE_ONLY)
                              .build())
                      .build());

          when(api.getEventSelectors(
                  GetEventSelectorsRequest.builder().trailName("management-trail-2").build()))
              .thenReturn(
                  GetEventSelectorsResponse.builder()
                      .eventSelectors(
                          EventSelector.builder()
                              .includeManagementEvents(true)
                              .readWriteType(ReadWriteType.READ_ONLY)
                              .build())
                      .build());
        });
  }

  @Test
  void testTrailManagementInAllRegionsManagementPassSplit() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/cloudtrail/account-cloud-trail-active.md"));

    assertNotNull(rule);

    var account = assets.get("current");

    assertNotNull(account);
    assertTrue(rule.evaluate(account).isOk());
  }
}
