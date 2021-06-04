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

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.mockito.Mockito.when;

import com.microsoft.azure.management.sql.implementation.DatabaseInner;
import com.microsoft.azure.management.sql.implementation.ServerInner;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

class AzureSQLDatabaseScannerTest extends AzureScannerTest {

  @BeforeAll
  static void setUpOnce() {
    discoverAssets(
        AzureSQLDatabaseScanner::new,
        api -> {
          var server = createSqlServer("id", new ServerInner());
          when(server.name()).thenReturn("myServer");

          when(api.azure.sqlServers().list()).thenReturn(MockedPagedList.of(server));

          var db = createSqlDatabase("id", "name", new DatabaseInner());
          when(db.name()).thenReturn("myDatabase");

          when(server.databases().list()).thenReturn(MockedPagedList.of(db));
        });
  }

  @Test
  void testScan() {
    assertEquals(1, assets.getDiscoveredAssets().size());

    var server = assets.get("id");

    assertNotNull(server);
    assertEquals("myDatabase", server.getName());
  }
}
