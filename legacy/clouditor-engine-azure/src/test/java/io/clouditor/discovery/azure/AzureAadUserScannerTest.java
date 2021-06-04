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
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

import com.microsoft.azure.management.graphrbac.ActiveDirectoryUser;
import com.microsoft.azure.management.graphrbac.implementation.UserInner;
import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

class AzureAadUserScannerTest extends AzureScannerTest {

  @BeforeAll
  static void setUpOnce() {
    discoverAssets(
        AzureAadUserScanner::new,
        api -> {
          var user1 =
              createWithId(
                  ActiveDirectoryUser.class,
                  "1",
                  new UserInner().withUserPrincipalName("#EXT#name"));

          var user2 =
              createWithId(
                  ActiveDirectoryUser.class, "2", new UserInner().withUserPrincipalName("name"));

          when(api.azure.accessManagement().activeDirectoryUsers().list())
              .thenReturn(MockedPagedList.of(user1, user2));
        });
  }

  @Test
  void testGuestUsers() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/azure/aad/aad-no-guest-user.md"));

    var user = assets.get("1");

    assertNotNull(user);
    assertFalse(rule.evaluate(user).isOk());

    user = assets.get("2");

    assertNotNull(user);
    assertTrue(rule.evaluate(user).isOk());
  }
}
