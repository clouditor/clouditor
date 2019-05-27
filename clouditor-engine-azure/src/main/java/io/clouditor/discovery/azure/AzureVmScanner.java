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

import com.microsoft.azure.management.compute.VirtualMachine;
import com.microsoft.azure.management.resources.fluentcore.model.HasInner;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.HashMap;
import java.util.List;
import java.util.stream.Collectors;

@ScannerInfo(assetType = "VirtualMachine", group = "Azure", service = "Compute")
public class AzureVmScanner extends AzureScanner<VirtualMachine> {

  public AzureVmScanner() {
    super(VirtualMachine::id, VirtualMachine::name);
  }

  @Override
  protected List<VirtualMachine> list() {
    return this.resourceGroup != null
        ? this.api.azure().virtualMachines().listByResourceGroup(this.resourceGroup)
        : this.api.azure().virtualMachines().list();
  }

  @Override
  protected Asset transform(VirtualMachine vm) throws ScanException {
    var asset = super.transform(vm);

    asset.setProperty(
        "extensions",
        vm.listExtensions().values().stream()
            .map(HasInner::inner)
            .map(inner -> MAPPER.convertValue(inner, HashMap.class))
            .collect(Collectors.toList()));

    asset.setProperty(
        "osDiskEncryption", vm.diskEncryption().getMonitor().osDiskStatus().toString());

    if (!vm.dataDisks().isEmpty()) {
      asset.setProperty(
          "dataDiskEncryption", vm.diskEncryption().getMonitor().dataDiskStatus().toString());
    }

    return asset;
  }
}
