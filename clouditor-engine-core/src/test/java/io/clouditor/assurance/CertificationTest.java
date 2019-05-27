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

package io.clouditor.assurance;

import io.clouditor.AbstractEngineUnitTest;
import io.clouditor.assurance.ccl.CCLDeserializer;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.AssetProperties;
import io.clouditor.discovery.DiscoveryResult;
import io.clouditor.discovery.DiscoveryService;
import io.clouditor.discovery.Scan;
import java.util.List;
import java.util.Map;
import java.util.Set;
import org.junit.jupiter.api.Test;

class CertificationTest extends AbstractEngineUnitTest {

  @Test
  void testEvaluate() {
    // build a simple rule
    var rule = new Rule();
    rule.setCondition(new CCLDeserializer().parse("MockAsset has property == true"));

    var ruleService = this.engine.getService(RuleService.class);
    ruleService.getRules().put("MockAsset", Set.of(rule));

    // put some assets into our asset service
    // TODO: inject service into test
    var discoveryService = this.engine.getService(DiscoveryService.class);

    var properties = new AssetProperties();
    properties.put("property", true);

    var asset = new Asset("MockAsset", "id", "name", properties);
    var discovery = new DiscoveryResult();
    discovery.setDiscoveredAssets(Map.of(asset.getId(), asset));

    // pipe it through the discovery pipeline
    discoveryService.submit(new Scan(), discovery);

    Certification cert = new Certification();
    cert.setId("some-cert");

    var control = new Control();
    control.setAutomated(true);
    control.setControlId("good-control-id");
    control.setRules(List.of(rule));
    this.engine.startMonitoring(control);

    control.evaluate(this.engine.getServiceLocator());

    var results = control.getResults();

    /*cert.setControls(Arrays.asList(goodControl, controlWithWarnings));

    this.engine.modifyCertification(cert);

    assertEquals(Fulfillment.GOOD, goodControl.getFulfilled());
    assertEquals(Fulfillment.WARNING, controlWithWarnings.getFulfilled());

    var assets = this.engine.getNonCompliantAssets("some-cert", "warning-control-id");

    assertEquals(1, assets.size());

    var detail = (ResultDetail) assets.values().toArray()[0];

    assertNotNull(detail);*/
  }

  /*@Test
  void testEqual() {
    var cert = new Certification();
    var control = new Control();
    control.setAutomated(true);
    control.setDomain(new Domain("Some Domain"));
    control.setObjectives(
        Collections.singletonList(new Objective(URI.create("test"), "true", "1")));
    cert.setControls(Arrays.asList(control, new Control()));

    // compare with self
    assertEquals(cert, cert);

    // compare with null
    assertNotEquals(cert, null);

    var other = new Certification();
    var otherControl = new Control();
    otherControl.setDomain(new Domain("Some Other Domain"));
    otherControl.setObjectives(
        Collections.singletonList(new Objective(URI.create("test"), "true", "1")));
    other.setControls(Collections.singletonList(otherControl));

    // compare with other
    assertNotEquals(cert, other);

    // compare with wrong class
    assertNotEquals(cert, new Object());
  }*/
}
