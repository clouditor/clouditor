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
import io.clouditor.discovery.AssetProperties;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import software.amazon.awssdk.services.iam.model.GetAccessKeyLastUsedRequest;
import software.amazon.awssdk.services.iam.model.GetAccessKeyLastUsedResponse;
import software.amazon.awssdk.services.iam.model.ListAccessKeysRequest;
import software.amazon.awssdk.services.iam.model.ListAccessKeysResponse;
import software.amazon.awssdk.services.iam.model.ListGroupsForUserRequest;
import software.amazon.awssdk.services.iam.model.ListGroupsForUserResponse;
import software.amazon.awssdk.services.iam.model.ListMfaDevicesRequest;
import software.amazon.awssdk.services.iam.model.ListMfaDevicesResponse;
import software.amazon.awssdk.services.iam.model.User;

@ScannerInfo(assetType = "User", group = "AWS", service = "IAM")
public class AwsIamUserScanner extends AwsIamScanner<User> {

  public AwsIamUserScanner() {
    super(User::arn, User::userName);
  }

  @Override
  public List<User> list() {
    return this.api.listUsers().users();
  }

  @Override
  public Asset transform(User user) throws ScanException {
    var asset = super.transform(user);

    enrichList(
        asset,
        "mfaDevices",
        this.api::listMFADevices,
        ListMfaDevicesResponse::mfaDevices,
        ListMfaDevicesRequest.builder().userName(user.userName()).build());

    enrichList(
        asset,
        "groups",
        this.api::listGroupsForUser,
        ListGroupsForUserResponse::groups,
        ListGroupsForUserRequest.builder().userName(user.userName()).build());

    enrichList(
        asset,
        "accessKeys",
        this.api::listAccessKeys,
        ListAccessKeysResponse::accessKeyMetadata,
        ListAccessKeysRequest.builder().userName(user.userName()).build());

    // TODO: this should probably be in a separate scanner
    var keys =
        (ArrayList<AssetProperties>)
            asset.getProperties().getOrDefault("accessKeys", Collections.emptyList());

    for (var key : keys) {
      enrich(
          key,
          "accessKeyLastUsed",
          this.api::getAccessKeyLastUsed,
          GetAccessKeyLastUsedResponse::accessKeyLastUsed,
          GetAccessKeyLastUsedRequest.builder()
              .accessKeyId(String.valueOf(key.get("accessKeyId")))
              .build());
    }

    return asset;
  }
}
