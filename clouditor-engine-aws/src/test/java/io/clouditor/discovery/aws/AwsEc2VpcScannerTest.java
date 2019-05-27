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
                    .getPathForResource("rules/aws/ec2/vpc-stale-security-groups.yaml"));

    assertFalse(rule.evaluate(vpc).isOk());
  }
}
