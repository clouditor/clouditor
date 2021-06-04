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

import io.clouditor.discovery.Asset;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.List;
import java.util.stream.Collectors;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.SqsClientBuilder;
import software.amazon.awssdk.services.sqs.model.GetQueueAttributesRequest;
import software.amazon.awssdk.services.sqs.model.QueueAttributeName;

@ScannerInfo(assetType = "Queue", group = "AWS", service = "SQS")
public class AwsSqsScanner extends AwsScanner<SqsClient, SqsClientBuilder, Queue> {

  private static final String ARN_PREFIX_SQS = "arn:aws:sqs";

  private static final String ARN_RESOURCE_TYPE_QUEUE = "queue";

  public AwsSqsScanner() {
    super(
        SqsClient::builder,
        queue ->
            ARN_PREFIX_SQS
                + ARN_SEPARATOR
                + ARN_RESOURCE_TYPE_QUEUE
                + RESOURCE_TYPE_SEPARATOR
                + queue.url(),
        Queue::url);
  }

  @Override
  protected List<Queue> list() {
    return this.api.listQueues().queueUrls().stream()
        .map(url -> Queue.builder().url(url).build())
        .collect(Collectors.toList());
  }

  @Override
  protected Asset transform(Queue queue) throws ScanException {
    var asset = super.transform(queue);

    asset.setProperty(
        "queueAttributes",
        this.api
            .getQueueAttributes(
                GetQueueAttributesRequest.builder()
                    .queueUrl(queue.url())
                    .attributeNames(QueueAttributeName.ALL)
                    .build())
            .attributesAsStrings());

    return asset;
  }
}
