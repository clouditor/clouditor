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
import java.util.HashMap;
import java.util.List;
import java.util.Map;
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

@ScannerInfo(assetType = "Bucket", group = "AWS", service = "S3")
public class AwsS3BucketScanner extends AwsScanner<S3Client, S3ClientBuilder, Bucket> {

  private Map<String, S3Client> regionClients = new HashMap<>();

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

        client =
            regionClients.getOrDefault(
                o.get(), S3Client.builder().region(Region.of(region)).build());
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
