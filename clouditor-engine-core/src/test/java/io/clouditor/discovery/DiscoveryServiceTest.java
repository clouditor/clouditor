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
