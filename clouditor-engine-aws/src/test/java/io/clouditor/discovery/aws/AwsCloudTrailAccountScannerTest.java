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
                    .getPathForResource("rules/aws/cloudtrail/account-cloud-trail-active.yaml"));

    assertNotNull(rule);

    var account = assets.get("current");

    assertNotNull(account);
    assertTrue(rule.evaluate(account).isOk());
  }
}
