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

import io.clouditor.discovery.ScannerInfo;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;
import software.amazon.awssdk.services.config.ConfigClient;
import software.amazon.awssdk.services.config.ConfigClientBuilder;
import software.amazon.awssdk.services.config.model.BaseConfigurationItem;
import software.amazon.awssdk.services.config.model.BatchGetResourceConfigRequest;
import software.amazon.awssdk.services.config.model.ListDiscoveredResourcesRequest;
import software.amazon.awssdk.services.config.model.ResourceKey;
import software.amazon.awssdk.services.config.model.ResourceType;

@ScannerInfo(assetType = "Configuration Item", group = "AWS", service = "ConfigService")
public class AwsConfigScanner
    extends AwsScanner<ConfigClient, ConfigClientBuilder, BaseConfigurationItem> {

  public AwsConfigScanner() {
    super(ConfigClient::builder, BaseConfigurationItem::arn, BaseConfigurationItem::resourceName);
  }

  @Override
  protected List<BaseConfigurationItem> list() {
    var resourceItems = new ArrayList<BaseConfigurationItem>();

    // for each resource type add the available resources to the asset list
    for (var type : ResourceType.values()) {
      var request = ListDiscoveredResourcesRequest.builder();
      request.resourceType(type);

      var r = api.listDiscoveredResources(request.build());

      var resourceKeyList =
          r.resourceIdentifiers().stream()
              .map(
                  id ->
                      ResourceKey.builder()
                          .resourceId(id.resourceId())
                          .resourceType(type.toString())
                          .build())
              .collect(Collectors.toList());

      if (!resourceKeyList.isEmpty()) {
        // TODO: is this the best way? wouldn't it be better to first gather all resource keys and
        // then do ONE batch request?
        var bres =
            api.batchGetResourceConfig(
                BatchGetResourceConfigRequest.builder().resourceKeys(resourceKeyList).build());

        resourceItems.addAll(bres.baseConfigurationItems());
      }
    }

    return resourceItems;
  }
}
