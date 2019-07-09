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
