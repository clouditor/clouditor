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
import software.amazon.awssdk.services.cloudtrail.CloudTrailClient;
import software.amazon.awssdk.services.cloudtrail.CloudTrailClientBuilder;
import software.amazon.awssdk.services.cloudtrail.model.GetEventSelectorsRequest;
import software.amazon.awssdk.services.cloudtrail.model.GetTrailStatusRequest;

@ScannerInfo(
    assetType = "Account",
    group = "AWS",
    service = "CloudTrail",
    assetIcon = "fas fa-sitemap")
public class AwsCloudTrailAccountScanner
    extends AwsScanner<CloudTrailClient, CloudTrailClientBuilder, Account> {

  public AwsCloudTrailAccountScanner() {
    // TODO: do we need a global client?
    super(CloudTrailClient::builder, Account::id, Account::id);
  }

  @Override
  protected List<Account> list() {
    // TODO: getScanners the account id somehow
    return List.of(Account.builder().id("current").build());
  }

  @Override
  protected Asset transform(Account account) throws ScanException {
    var asset = super.transform(account);

    asset.setProperty("cloudTrailActive", hasCloudTrail());

    return asset;
  }

  private boolean hasCloudTrail() {
    var hasActiveAllRegionManagementWriteTrail = false;
    var hasActiveAllRegionManagementReadTrail = false;

    var trails = this.api.describeTrails().trailList();

    var found = false;

    for (var trail : trails) {
      boolean isActive =
          this.api
              .getTrailStatus(GetTrailStatusRequest.builder().name(trail.name()).build())
              .isLogging();

      /*
       We are NOT interested in trails that are
       - inactive or
       - are not "multiRegion", i.e. is not active in all regions
      */
      if (!isActive || Boolean.FALSE.equals(trail.isMultiRegionTrail())) {
        continue;
      }

      for (var selector :
          this.api
              .getEventSelectors(GetEventSelectorsRequest.builder().trailName(trail.name()).build())
              .eventSelectors()) {
        if (Boolean.TRUE.equals(selector.includeManagementEvents())) {
          switch (selector.readWriteType()) {
            case READ_ONLY:
              hasActiveAllRegionManagementReadTrail = true;
              break;

            case WRITE_ONLY:
              hasActiveAllRegionManagementWriteTrail = true;
              break;

            case ALL:
              hasActiveAllRegionManagementWriteTrail = true;
              hasActiveAllRegionManagementReadTrail = true;
              break;
          }
        }
      }

      if (hasActiveAllRegionManagementWriteTrail && hasActiveAllRegionManagementReadTrail) {
        found = true;
        break;
      }
    }

    return found;
  }
}
