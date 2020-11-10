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

import io.clouditor.credentials.AwsAccount;
import io.clouditor.data_access_layer.HibernatePersistence;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import software.amazon.awssdk.core.exception.SdkClientException;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.S3ClientBuilder;
import software.amazon.awssdk.services.s3.model.Bucket;
import software.amazon.awssdk.services.s3.model.GetBucketEncryptionRequest;
import software.amazon.awssdk.services.s3.model.GetBucketEncryptionResponse;
import software.amazon.awssdk.services.s3.model.GetBucketLifecycleConfigurationRequest;
import software.amazon.awssdk.services.s3.model.GetBucketLifecycleConfigurationResponse;
import software.amazon.awssdk.services.s3.model.GetBucketReplicationRequest;
import software.amazon.awssdk.services.s3.model.GetBucketReplicationResponse;
import software.amazon.awssdk.services.s3.model.GetPublicAccessBlockRequest;
import software.amazon.awssdk.services.s3.model.GetPublicAccessBlockResponse;
import software.amazon.awssdk.services.s3.model.HeadBucketRequest;
import software.amazon.awssdk.services.s3.model.ListBucketsRequest;
import software.amazon.awssdk.services.s3.model.S3Exception;

@ScannerInfo(assetType = "Bucket", group = "AWS", service = "S3", assetIcon = "fas fa-archive")
public class AwsS3BucketScanner extends AwsScanner<S3Client, S3ClientBuilder, Bucket> {

  private final Map<String, S3Client> regionClients = new HashMap<>();

  static final String ARN_PREFIX_S3 = "arn:aws:s3:::";

  public AwsS3BucketScanner() {
    super(S3Client::builder, bucket -> ARN_PREFIX_S3 + bucket.name(), Bucket::name);
  }

  @Override
  public List<Bucket> list() {
    return this.api.listBuckets().buckets();
  }

  @Override
  public Asset transform(Bucket bucket) throws ScanException {
    var map = super.transform(bucket);

    this.api.listBuckets(ListBucketsRequest.builder().build());

    var client = this.api;
    try {
      this.api.headBucket(HeadBucketRequest.builder().bucket(bucket.name()).build());
    } catch (S3Exception ex) {
      var o = ex.awsErrorDetails().sdkHttpResponse().firstMatchingHeader("x-amz-bucket-region");

      if (o.isPresent()) {
        var region = o.get();

        // needed, as long as https://github.com/aws/aws-sdk-java-v2/issues/52 is not fixed
        LOGGER.info("Switching to region-specific S3 client ({})", region);

        var account =
            new HibernatePersistence()
                .get(AwsAccount.class, "AWS")
                .orElseThrow(() -> SdkClientException.create("AWS account not configured"));

        client =
            regionClients.getOrDefault(
                region,
                S3Client.builder().credentialsProvider(account).region(Region.of(region)).build());
      }
    }

    enrich(
        map,
        "bucketEncryption",
        client::getBucketEncryption,
        GetBucketEncryptionResponse::serverSideEncryptionConfiguration,
        GetBucketEncryptionRequest.builder().bucket(bucket.name()).build());

    enrich(
        map,
        "publicAccessBlockConfiguration",
        client::getPublicAccessBlock,
        GetPublicAccessBlockResponse::publicAccessBlockConfiguration,
        GetPublicAccessBlockRequest.builder().bucket(bucket.name()).build());

    enrich(
        map,
        "bucketReplication",
        client::getBucketReplication,
        GetBucketReplicationResponse::replicationConfiguration,
        GetBucketReplicationRequest.builder().bucket(bucket.name()).build());

    enrichList(
        map,
        "lifecycleConfiguration",
        client::getBucketLifecycleConfiguration,
        GetBucketLifecycleConfigurationResponse::rules,
        GetBucketLifecycleConfigurationRequest.builder().bucket(bucket.name()).build());

    return map;
  }
}
