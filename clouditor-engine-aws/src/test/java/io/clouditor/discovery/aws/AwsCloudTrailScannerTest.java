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
import software.amazon.awssdk.services.cloudtrail.CloudTrailClient;
import software.amazon.awssdk.services.cloudtrail.model.DescribeTrailsResponse;
import software.amazon.awssdk.services.cloudtrail.model.Trail;

class AwsCloudTrailScannerTest extends AwsScannerTest {

  private static final String ENCRYPTED_TRAIL =
      "arn:aws:cloudtrail:eu-central-1:575163611729:trail/encrypted-trail";
  private static final String UNENCRYPTED_TRAIL =
      "arn:aws:cloudtrail:eu-central-1:575163611729:trail/unencrypted-trail";

  @BeforeAll
  static void setUpOnce() throws IOException {
    discoverAssets(
        CloudTrailClient.class,
        AwsCloudTrailScanner::new,
        api ->
            when(api.describeTrails())
                .thenReturn(
                    DescribeTrailsResponse.builder()
                        .trailList(
                            Trail.builder()
                                .trailARN(ENCRYPTED_TRAIL)
                                .kmsKeyId("some-key-arn")
                                .build(),
                            Trail.builder().trailARN(UNENCRYPTED_TRAIL).build())
                        .build()));
  }

  @Test
  void testEncryption() throws IOException {
    var encryptedTrail = assets.get(ENCRYPTED_TRAIL);

    assertNotNull(encryptedTrail);

    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/cloudtrail/trail-encrypted.yaml"));

    assertNotNull(rule);

    assertTrue(rule.evaluate(encryptedTrail).isOk());

    var unencryptedTrail = assets.get(UNENCRYPTED_TRAIL);

    assertNotNull(unencryptedTrail);

    assertFalse(rule.evaluate(unencryptedTrail).isOk());
  }
}
