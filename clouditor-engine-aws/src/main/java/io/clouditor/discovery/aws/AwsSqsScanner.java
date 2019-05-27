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
