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
