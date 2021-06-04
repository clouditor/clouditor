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
import software.amazon.awssdk.services.ec2.model.DescribeStaleSecurityGroupsRequest;
import software.amazon.awssdk.services.ec2.model.DescribeStaleSecurityGroupsResponse;
import software.amazon.awssdk.services.ec2.model.Vpc;

@ScannerInfo(assetType = "VPC", group = "AWS", service = "EC2", assetIcon = "fas fa-network-wired")
public class AwsEc2VpcScanner extends AwsEc2Scanner<Vpc> {

  private static final String ARN_RESOURCE_TYPE_VPC = "vpc";

  public AwsEc2VpcScanner() {
    // TODO: getScanners name from tags
    super(
        vpc ->
            ARN_PREFIX_EC2
                + AwsScanner.ARN_SEPARATOR
                + ARN_RESOURCE_TYPE_VPC
                + AwsScanner.RESOURCE_TYPE_SEPARATOR
                + vpc.vpcId(),
        Vpc::vpcId);
  }

  @Override
  protected List<Vpc> list() {
    return this.api.describeVpcs().vpcs();
  }

  @Override
  protected Asset transform(Vpc vpc) throws ScanException {
    var asset = super.transform(vpc);

    enrichList(
        asset,
        "staleSecurityGroups",
        this.api::describeStaleSecurityGroups,
        DescribeStaleSecurityGroupsResponse::staleSecurityGroupSet,
        DescribeStaleSecurityGroupsRequest.builder().vpcId(vpc.vpcId()).build());

    return asset;
  }
}
