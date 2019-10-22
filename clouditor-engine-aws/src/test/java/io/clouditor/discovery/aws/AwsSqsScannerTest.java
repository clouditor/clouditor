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
import java.util.Map;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.GetQueueAttributesRequest;
import software.amazon.awssdk.services.sqs.model.GetQueueAttributesResponse;
import software.amazon.awssdk.services.sqs.model.ListQueuesResponse;
import software.amazon.awssdk.services.sqs.model.QueueAttributeName;

class AwsSqsScannerTest extends AwsScannerTest {

  @BeforeAll
  static void setUpOnce() throws IOException {
    discoverAssets(
        SqsClient.class,
        AwsSqsScanner::new,
        api -> {
          when(api.listQueues())
              .thenReturn(
                  ListQueuesResponse.builder()
                      .queueUrls("123456789012/MyQueue1", "123456789012/MyQueue2")
                      .build());

          when(api.getQueueAttributes(
                  GetQueueAttributesRequest.builder()
                      .queueUrl("123456789012/MyQueue1")
                      .attributeNames(QueueAttributeName.ALL)
                      .build()))
              .thenReturn(
                  GetQueueAttributesResponse.builder()
                      .attributesWithStrings(Map.of("VisibilityTimeout", "30"))
                      .build());

          when(api.getQueueAttributes(
                  GetQueueAttributesRequest.builder()
                      .queueUrl("123456789012/MyQueue2")
                      .attributeNames(QueueAttributeName.ALL)
                      .build()))
              .thenReturn(
                  GetQueueAttributesResponse.builder()
                      .attributesWithStrings(Map.of("VisibilityTimeout", "60"))
                      .build());
        });
  }

  @Test
  void testVisibilityTimeout() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/sqs/queue-visibility-timeout.md"));

    assertNotNull(rule);

    var queue1 = assets.get("arn:aws:sqs:queue/123456789012/MyQueue1");

    assertNotNull(queue1);
    assertTrue(rule.evaluate(queue1).isOk());

    var queue2 = assets.get("arn:aws:sqs:queue/123456789012/MyQueue2");

    assertNotNull(queue2);
    assertFalse(rule.evaluate(queue2).isOk());
  }
}
