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
import software.amazon.awssdk.services.kinesis.KinesisClient;
import software.amazon.awssdk.services.kinesis.model.DescribeStreamRequest;
import software.amazon.awssdk.services.kinesis.model.DescribeStreamResponse;
import software.amazon.awssdk.services.kinesis.model.EncryptionType;
import software.amazon.awssdk.services.kinesis.model.ListStreamsResponse;
import software.amazon.awssdk.services.kinesis.model.StreamDescription;

class AwsKinesisScannerTest extends AwsScannerTest {

  @BeforeAll
  static void setUpOnce() throws IOException {
    discoverAssets(
        KinesisClient.class,
        AwsKinesisScanner::new,
        api -> {
          when(api.listStreams())
              .thenReturn(
                  ListStreamsResponse.builder()
                      .streamNames("stream-encrypted", "stream-not-encrypted")
                      .build());

          when(api.describeStream(
                  DescribeStreamRequest.builder().streamName("stream-encrypted").build()))
              .thenReturn(
                  DescribeStreamResponse.builder()
                      .streamDescription(
                          StreamDescription.builder()
                              .streamARN("arn:aws:kinesis:us-east-1:111122223333:encrypted")
                              .encryptionType(EncryptionType.KMS)
                              .build())
                      .build());

          when(api.describeStream(
                  DescribeStreamRequest.builder().streamName("stream-not-encrypted").build()))
              .thenReturn(
                  DescribeStreamResponse.builder()
                      .streamDescription(
                          StreamDescription.builder()
                              .streamARN("arn:aws:kinesis:us-east-1:111122223333:unencrypted")
                              .encryptionType(EncryptionType.NONE)
                              .build())
                      .build());
        });
  }

  @Test
  void testStreamEncryption() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/kinesis/stream-encryption.md"));

    assertNotNull(rule);

    var encrypted = assets.get("arn:aws:kinesis:us-east-1:111122223333:encrypted");

    assertTrue(rule.evaluate(encrypted).isOk());

    var unencrypted = assets.get("arn:aws:kinesis:us-east-1:111122223333:unencrypted");

    assertFalse(rule.evaluate(unencrypted).isOk());
  }
}
