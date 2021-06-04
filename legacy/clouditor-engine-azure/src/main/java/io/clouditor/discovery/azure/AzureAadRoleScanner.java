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

package io.clouditor.discovery.azure;

import com.microsoft.azure.management.graphrbac.RoleDefinition;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.List;

@ScannerInfo(assetType = "ActiveDirectoryRole", group = "Azure", service = "Active Directory")
public class AzureAadRoleScanner extends AzureScanner<RoleDefinition> {

  public AzureAadRoleScanner() {
    super(RoleDefinition::id, RoleDefinition::roleName);
  }

  @Override
  protected List<RoleDefinition> list() {
    return this.api.azure().accessManagement().roleDefinitions().listByScope("");
  }

  @Override
  protected Asset transform(RoleDefinition role) throws ScanException {
    var asset = super.transform(role);

    var hasGlobalScope = false;
    var isAdminRole = false;

    for (var scope : role.assignableScopes()) {
      if (scope.equals("/") || scope.contains("subscription")) {
        hasGlobalScope = true;
      }
    }

    if (hasGlobalScope) {
      for (var permission : role.permissions()) {
        for (var action : permission.actions()) {
          if (action.equals("*")) {
            isAdminRole = true;
          }
        }
      }
    }

    asset.setProperty("customAdminRole", isAdminRole);

    return asset;
  }
}
