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
                    .getPathForResource("/rules/aws/ec2/volume-encryption.yaml"));

    assertFalse(rule.evaluate(volume).isOk());

    volume = assets.get(ENCRYPTED_VOLUME_ARN);

    assertNotNull(volume);

    assertTrue(rule.evaluate(volume).isOk());
  }
}
