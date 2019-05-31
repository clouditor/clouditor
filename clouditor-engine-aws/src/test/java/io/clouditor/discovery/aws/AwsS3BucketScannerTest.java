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
import org.mockito.ArgumentMatchers;
import software.amazon.awssdk.awscore.exception.AwsServiceException;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.Bucket;
import software.amazon.awssdk.services.s3.model.GetBucketEncryptionRequest;
import software.amazon.awssdk.services.s3.model.GetBucketEncryptionResponse;
import software.amazon.awssdk.services.s3.model.GetBucketLifecycleConfigurationRequest;
import software.amazon.awssdk.services.s3.model.GetBucketReplicationRequest;
import software.amazon.awssdk.services.s3.model.GetBucketReplicationResponse;
import software.amazon.awssdk.services.s3.model.GetPublicAccessBlockRequest;
import software.amazon.awssdk.services.s3.model.GetPublicAccessBlockResponse;
import software.amazon.awssdk.services.s3.model.ListBucketsResponse;
import software.amazon.awssdk.services.s3.model.PublicAccessBlockConfiguration;
import software.amazon.awssdk.services.s3.model.ReplicationConfiguration;
import software.amazon.awssdk.services.s3.model.ServerSideEncryption;
import software.amazon.awssdk.services.s3.model.ServerSideEncryptionByDefault;
import software.amazon.awssdk.services.s3.model.ServerSideEncryptionConfiguration;
import software.amazon.awssdk.services.s3.model.ServerSideEncryptionRule;

class AwsS3BucketScannerTest extends AwsScannerTest {

  @BeforeAll
  static void setUpOnce() {
    discoverAssets(
        S3Client.class,
        AwsS3BucketScanner::new,
        api -> {
          when(api.listBuckets())
              .thenReturn(
                  ListBucketsResponse.builder()
                      .buckets(
                          Bucket.builder().name("Bucket-A").build(),
                          Bucket.builder().name("Bucket-B").build(),
                          Bucket.builder().name("Bucket-C").build())
                      .build());

          when(api.getBucketEncryption(
                  GetBucketEncryptionRequest.builder().bucket("Bucket-A").build()))
              .thenReturn(
                  GetBucketEncryptionResponse.builder()
                      .serverSideEncryptionConfiguration(
                          ServerSideEncryptionConfiguration.builder()
                              .rules(
                                  ServerSideEncryptionRule.builder()
                                      .applyServerSideEncryptionByDefault(
                                          ServerSideEncryptionByDefault.builder()
                                              .kmsMasterKeyID("key")
                                              .sseAlgorithm(ServerSideEncryption.AES256)
                                              .build())
                                      .build())
                              .build())
                      .build());

          when(api.getBucketEncryption(
                  GetBucketEncryptionRequest.builder().bucket("Bucket-B").build()))
              .thenReturn(
                  GetBucketEncryptionResponse.builder()
                      .serverSideEncryptionConfiguration(
                          ServerSideEncryptionConfiguration.builder()
                              .rules(
                                  ServerSideEncryptionRule.builder()
                                      .applyServerSideEncryptionByDefault(
                                          ServerSideEncryptionByDefault.builder()
                                              .kmsMasterKeyID("key")
                                              .sseAlgorithm(ServerSideEncryption.AWS_KMS)
                                              .build())
                                      .build())
                              .build())
                      .build());

          when(api.getBucketEncryption(
                  GetBucketEncryptionRequest.builder().bucket("Bucket-C").build()))
              .thenThrow(AwsServiceException.builder().statusCode(404).build());

          when(api.getPublicAccessBlock(ArgumentMatchers.any(GetPublicAccessBlockRequest.class)))
              .thenReturn(
                  GetPublicAccessBlockResponse.builder()
                      .publicAccessBlockConfiguration(
                          PublicAccessBlockConfiguration.builder().build())
                      .build());

          when(api.getBucketReplication(ArgumentMatchers.any(GetBucketReplicationRequest.class)))
              .thenReturn(
                  GetBucketReplicationResponse.builder()
                      .replicationConfiguration(ReplicationConfiguration.builder().build())
                      .build());

          when(api.getBucketLifecycleConfiguration(
                  (GetBucketLifecycleConfigurationRequest) ArgumentMatchers.any()))
              .thenThrow(AwsServiceException.builder().statusCode(404).build());
        });
  }

  @Test
  void testBucketDefaultEncryptionCheck() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/s3/bucket-default-encryption.yaml"));

    assertNotNull(rule);

    var bucketA = assets.get("arn:aws:s3:::Bucket-A");

    assertNotNull(bucketA);
    assertTrue(rule.evaluate(bucketA).isOk());

    var bucketB = assets.get("arn:aws:s3:::Bucket-B");

    assertNotNull(bucketB);
    assertTrue(rule.evaluate(bucketB).isOk());

    var bucketC = assets.get("arn:aws:s3:::Bucket-C");

    assertNotNull(bucketC);
    assertFalse(rule.evaluate(bucketC).isOk());
  }
}
