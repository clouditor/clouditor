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
import software.amazon.awssdk.services.ec2.Ec2Client;
import software.amazon.awssdk.services.ec2.model.DescribeVolumesResponse;
import software.amazon.awssdk.services.ec2.model.Volume;
import software.amazon.awssdk.services.ec2.model.VolumeState;

class AwsEc2VolumeScannerTest extends AwsScannerTest {

  private static final String NOT_ENCRYPTED_VOLUME_ID = "vol-06d882eadbcee7968";
  private static final String NOT_ENCRYPTED_VOLUME_ARN =
      "arn:aws:ec2:volume/" + NOT_ENCRYPTED_VOLUME_ID;
  private static final String ENCRYPTED_VOLUME_ID = "vol-06d882eadbcee7967";
  private static final String ENCRYPTED_VOLUME_ARN = "arn:aws:ec2:volume/" + ENCRYPTED_VOLUME_ID;

  @BeforeAll
  static void setUpOnce() throws IOException {
    discoverAssets(
        Ec2Client.class,
        AwsEc2VolumeScanner::new,
        api ->
            when(api.describeVolumes())
                .thenReturn(
                    DescribeVolumesResponse.builder()
                        .volumes(
                            Volume.builder()
                                .volumeId(NOT_ENCRYPTED_VOLUME_ID)
                                .encrypted(false)
                                .state(VolumeState.AVAILABLE)
                                .build(),
                            Volume.builder()
                                .volumeId(ENCRYPTED_VOLUME_ID)
                                .encrypted(true)
                                .state(VolumeState.AVAILABLE)
                                .build())
                        .build()));
  }

  @Test
  void testEncryption() throws IOException {
    var volume = assets.get(NOT_ENCRYPTED_VOLUME_ARN);

    assertNotNull(volume);

    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("/rules/aws/ec2/volume-encryption.md"));

    assertFalse(rule.evaluate(volume).isOk());

    volume = assets.get(ENCRYPTED_VOLUME_ARN);

    assertNotNull(volume);

    assertTrue(rule.evaluate(volume).isOk());
  }
}
