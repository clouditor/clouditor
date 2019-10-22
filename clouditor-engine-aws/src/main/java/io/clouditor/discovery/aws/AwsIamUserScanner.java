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

@ScannerInfo(assetType = "User", group = "AWS", service = "IAM", assetIcon = "fas fa-user")
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
