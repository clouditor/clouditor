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
import static org.mockito.Mockito.when;

import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import software.amazon.awssdk.services.ec2.Ec2Client;
import software.amazon.awssdk.services.ec2.model.DescribeStaleSecurityGroupsRequest;
import software.amazon.awssdk.services.ec2.model.DescribeStaleSecurityGroupsResponse;
import software.amazon.awssdk.services.ec2.model.DescribeVpcsResponse;
import software.amazon.awssdk.services.ec2.model.StaleSecurityGroup;
import software.amazon.awssdk.services.ec2.model.Vpc;

class AwsEc2VpcScannerTest extends AwsScannerTest {

  private static final String VPC_ARN = "arn:aws:ec2:vpc/vpc-1";

  @BeforeAll
  static void setUpOnce() throws IOException {
    discoverAssets(
        Ec2Client.class,
        AwsEc2VpcScanner::new,
        api -> {
          when(api.describeVpcs())
              .thenReturn(
                  DescribeVpcsResponse.builder()
                      .vpcs(Vpc.builder().vpcId("vpc-1").build())
                      .build());

          when(api.describeStaleSecurityGroups(
                  DescribeStaleSecurityGroupsRequest.builder().vpcId("vpc-1").build()))
              .thenReturn(
                  DescribeStaleSecurityGroupsResponse.builder()
                      .staleSecurityGroupSet(
                          StaleSecurityGroup.builder().groupId("some-group").build())
                      .build());
        });
  }

  @Test
  void testStaleSecurityGroups() throws IOException {
    var vpc = assets.get(VPC_ARN);

    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/ec2/vpc-stale-security-groups.md"));

    assertFalse(rule.evaluate(vpc).isOk());
  }
}
