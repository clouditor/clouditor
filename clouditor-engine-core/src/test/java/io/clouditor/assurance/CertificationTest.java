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

    var certificationService = this.engine.getService(CertificationService.class);

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
    certificationService.startMonitoring(control);

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
