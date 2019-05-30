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
