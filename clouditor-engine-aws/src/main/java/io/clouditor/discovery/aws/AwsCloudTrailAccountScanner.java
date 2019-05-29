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
      if (!isActive || !trail.isMultiRegionTrail()) {
        continue;
      }

      for (var selector :
          this.api
              .getEventSelectors(GetEventSelectorsRequest.builder().trailName(trail.name()).build())
              .eventSelectors()) {
        if (selector.includeManagementEvents()) {
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
