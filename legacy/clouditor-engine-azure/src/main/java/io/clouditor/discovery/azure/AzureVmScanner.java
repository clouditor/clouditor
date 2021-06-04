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
