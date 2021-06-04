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
