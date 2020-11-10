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

package io.clouditor.discovery;

import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

import io.clouditor.AbstractEngineUnitTest;
import io.clouditor.assurance.RuleService;
import java.util.concurrent.TimeUnit;
import org.awaitility.Awaitility;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

class DiscoveryServiceTest extends AbstractEngineUnitTest {

  private static final String OBJECT_ID_FAKE = "fake";

  @BeforeEach
  protected void setUp() {
    super.setUp();
    this.engine.setDBName("DiscoveryServiceTestDB");
    this.engine.initDB();
  }

  @Test
  void testScanning() {
    var scanService = this.engine.getService(DiscoveryService.class);

    assertNotNull(scanService);

    scanService.init();

    var scans = scanService.getScans();

    // subscribe with RuleService
    var ruleServices = this.engine.getService(RuleService.class);

    assertNotNull(ruleServices);

    scanService.subscribe(ruleServices);

    Assertions.assertFalse(scans.isEmpty());

    var scan = scanService.getScan("fake");

    scanService.enableScan(scan);

    assertTrue(scan.isEnabled());
    Awaitility.await()
        .atMost(5, TimeUnit.SECONDS)
        .until(() -> this.engine.getService(AssetService.class).get(OBJECT_ID_FAKE) != null);
  }
}
