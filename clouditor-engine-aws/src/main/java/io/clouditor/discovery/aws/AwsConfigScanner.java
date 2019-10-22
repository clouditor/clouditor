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
