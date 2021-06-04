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
                    .getPathForResource("rules/aws/cloudtrail/trail-encrypted.md"));

    assertNotNull(rule);

    assertTrue(rule.evaluate(encryptedTrail).isOk());

    var unencryptedTrail = assets.get(UNENCRYPTED_TRAIL);

    assertNotNull(unencryptedTrail);

    assertFalse(rule.evaluate(unencryptedTrail).isOk());
  }
}
