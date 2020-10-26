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

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import io.clouditor.assurance.ccl.BinaryComparison;
import io.clouditor.assurance.ccl.CCLDeserializer;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.AssetProperties;
import io.clouditor.rest.ObjectMapperResolver;
import java.io.IOException;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.Test;

class CCLDeserializerTest {

  @Test
  void testFromYAML() throws IOException {
    var mapper = new ObjectMapper(new YAMLFactory());
    ObjectMapperResolver.configureObjectMapper(mapper);

    var rule = mapper.readValue("conditions: \n - \"User has field == true\"", Rule.class);

    assertNotNull(rule);

    var conditions = rule.getConditions();

    assertNotNull(conditions);

    assertTrue(conditions.get(0).getExpression() instanceof BinaryComparison);
  }

  @Test
  void testEqualComparison() {
    var ccl = new CCLDeserializer();
    var condition = ccl.parse("Volume has encrypted == true");

    assertNotNull(condition);

    var asset = new AssetProperties();
    asset.put("encrypted", true);

    assertTrue(condition.evaluate(asset));
  }

  @Test
  void testEmptyExpression() {
    var ccl = new CCLDeserializer();
    var condition = ccl.parse("User has not empty mfaDevices");

    var asset = new AssetProperties();
    asset.put("mfaDevices", new ArrayList<>());

    // empty can work on a list
    assertFalse(condition.evaluate(asset));

    asset.clear();
    asset.put("mfaDevices", List.of(Map.of("deviceId", "5"), Map.of("deviceId", "6")));

    assertTrue(condition.evaluate(asset));

    // empty can also check, if a string value is null or 'empty'
    condition = ccl.parse("User has empty name");

    asset.put("name", null);

    assertTrue(condition.evaluate(asset));

    asset.put("name", "");

    assertTrue(condition.evaluate(asset));

    asset.put("name", "Name");

    assertFalse(condition.evaluate(asset));
  }

  @Test
  void testInExpression() {
    var ccl = new CCLDeserializer();
    var condition = ccl.parse("User has deviceId == \"5\" in any mfaDevices");

    var asset = new AssetProperties();
    asset.put("mfaDevices", List.of(Map.of("deviceId", "5"), Map.of("deviceId", "6")));

    assertTrue(condition.evaluate(asset));

    condition = ccl.parse("User has deviceId == \"5\" in all mfaDevices");

    assertFalse(condition.evaluate(asset));
  }

  @Test
  void testIsBeforeComparison() {
    var ccl = new CCLDeserializer();
    var condition = ccl.parse("AccessKey has expiry before 10 days");

    var asset = new AssetProperties();
    asset.put("expiry", Map.of("epochSecond", Instant.now().getEpochSecond(), "nano", 1));

    assertTrue(condition.evaluate(asset));

    asset.clear();

    asset.put(
        "expiry",
        Map.of("epochSecond", Instant.now().plus(20, ChronoUnit.DAYS).getEpochSecond(), "nano", 1));

    assertFalse(condition.evaluate(asset));
  }

  @Test
  void testAfterComparison() {
    var ccl = new CCLDeserializer();
    var condition = ccl.parse("AccessKey has expiry after 10 days");

    var asset = new AssetProperties();
    asset.put(
        "expiry",
        Map.of("epochSecond", Instant.now().plus(20, ChronoUnit.DAYS).getEpochSecond(), "nano", 1));

    assertTrue(condition.evaluate(asset));

    asset.clear();

    asset.put("expiry", Map.of("epochSecond", Instant.now().getEpochSecond(), "nano", 1));

    assertFalse(condition.evaluate(asset));
  }

  @Test
  void testYoungerComparison() {
    var ccl = new CCLDeserializer();
    var condition = ccl.parse("User has createDate younger 10 days");

    var asset = new AssetProperties();
    asset.put("createDate", Map.of("epochSecond", Instant.now().getEpochSecond(), "nano", 1));

    assertTrue(condition.evaluate(asset));

    asset.clear();
    // something in the future
    asset.put(
        "createDate",
        Map.of(
            "epochSecond", Instant.now().minus(20, ChronoUnit.DAYS).getEpochSecond(), "nano", 1));

    assertFalse(condition.evaluate(asset));
  }

  @Test
  void testOlderComparison() {
    var ccl = new CCLDeserializer();
    var condition = ccl.parse("User has createDate older 10 days");

    var asset = new AssetProperties();

    asset.put(
        "createDate",
        Map.of(
            "epochSecond", Instant.now().minus(20, ChronoUnit.DAYS).getEpochSecond(), "nano", 1));

    assertTrue(condition.evaluate(asset));

    asset.clear();

    asset.put("createDate", Map.of("epochSecond", Instant.now().getEpochSecond(), "nano", 1));

    assertFalse(condition.evaluate(asset));
  }

  @Test
  void testWithinExpression() {
    var ccl = new CCLDeserializer();
    var condition =
        ccl.parse("BaseConfigurationItem  has awsRegion within \"eu-central1\", \"eu-west1\"");

    var asset = new AssetProperties();
    asset.put("awsRegion", "eu-west1");

    assertTrue(condition.evaluate(asset));

    asset.put("awsRegion", "us-east1");

    assertFalse(condition.evaluate(asset));
  }

  @Test
  void testSubFields() {
    var ccl = new CCLDeserializer();
    var condition = ccl.parse("LegacyAsset has description.name == \"test\"");

    var asset = new AssetProperties();
    asset.put("description", Map.of("name", "test"));

    assertTrue(condition.evaluate(asset));
  }

  @Test
  void testContains() {
    var ccl = new CCLDeserializer();
    var condition = ccl.parse("User has name contains \"name\"");

    var asset = new AssetProperties();
    asset.put("name", "Some name");

    assertTrue(condition.evaluate(asset));
  }

  @Test
  void testAmbiguous() {
    var ccl = new CCLDeserializer();
    var rule = new Rule();
    rule.setConditions(
        List.of(ccl.parse("Storage has (not empty policy.algorithm) in any encryption.rules")));

    // asset does
    // - not have field 'encryption.rules' at all
    var asset = new Asset(null, null, null, AssetProperties.of());
    assertFalse(rule.evaluate(asset).isOk());

    // asset does
    // - have field 'encryption'
    // - not have field 'encryption.rules'
    asset = new Asset(null, null, null, AssetProperties.of("encryption", AssetProperties.of()));
    assertFalse(rule.evaluate(asset).isOk());

    // asset does
    // - have field 'encryption'
    // - have field 'encryption.rules'
    // - not have inner field 'policy'
    asset =
        new Asset(
            null,
            null,
            null,
            AssetProperties.of("encryption", AssetProperties.of("rules", List.of())));
    assertFalse(rule.evaluate(asset).isOk());

    // asset does
    // - have field 'encryption'
    // - have field 'encryption.rules'
    // - have inner field 'policy'
    // - not have inner field 'policy.algorithm'
    asset =
        new Asset(
            null,
            null,
            null,
            AssetProperties.of(
                "encryption",
                AssetProperties.of(
                    "rules", List.of(AssetProperties.of("policy", AssetProperties.of())))));
    assertFalse(rule.evaluate(asset).isOk());

    // asset does
    // - have field 'encryption'
    // - have field 'encryption.rules'
    // - have inner field 'policy'
    // - have inner field 'policy.algorithm' (with a value)
    asset =
        new Asset(
            null,
            null,
            null,
            AssetProperties.of(
                "encryption",
                AssetProperties.of(
                    "rules",
                    List.of(
                        AssetProperties.of(
                            "policy", AssetProperties.of("algorithm", "AES-265"))))));
    assertTrue(rule.evaluate(asset).isOk());
  }
}
