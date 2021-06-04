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
 */

package io.clouditor.discovery.azure;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

import com.microsoft.azure.management.compute.EncryptionStatus;
import com.microsoft.azure.management.compute.OSProfile;
import com.microsoft.azure.management.compute.implementation.VirtualMachineInner;
import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

class AzureVmScannerTest extends AzureScannerTest {

  @BeforeAll
  static void setUpOnce() {
    discoverAssets(
        AzureVmScanner::new,
        api -> {
          var vm1 =
              createVirtualMachine(
                  "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/test/providers/Microsoft.Compute/virtualMachines/centvm",
                  new VirtualMachineInner()
                      .withOsProfile(new OSProfile().withAdminPassword("password")),
                  EncryptionStatus.ENCRYPTED,
                  EncryptionStatus.ENCRYPTED);

          var vm2 =
              createVirtualMachine(
                  "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/test/providers/Microsoft.Compute/virtualMachines/unencryptedvm",
                  new VirtualMachineInner()
                      .withOsProfile(new OSProfile().withAdminPassword("password")),
                  EncryptionStatus.NOT_ENCRYPTED,
                  EncryptionStatus.NOT_ENCRYPTED);

          when(api.azure.virtualMachines().list()).thenReturn(MockedPagedList.of(vm1, vm2));
        });
  }

  @Test
  void testDataEncryption() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/azure/compute/vm-data-encryption.md"));

    var vm1 =
        assets.get(
            "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/test/providers/Microsoft.Compute/virtualMachines/centvm");

    assertTrue(rule.evaluate(vm1).isOk());

    var vm2 =
        assets.get(
            "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/test/providers/Microsoft.Compute/virtualMachines/unencryptedvm");

    assertFalse(rule.evaluate(vm2).isOk());
  }

  @Test
  void testOsEncryption() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/azure/compute/vm-os-encryption.md"));

    var vm1 =
        assets.get(
            "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/test/providers/Microsoft.Compute/virtualMachines/centvm");

    assertTrue(rule.evaluate(vm1).isOk());

    var vm2 =
        assets.get(
            "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/test/providers/Microsoft.Compute/virtualMachines/unencryptedvm");

    assertFalse(rule.evaluate(vm2).isOk());
  }
}
