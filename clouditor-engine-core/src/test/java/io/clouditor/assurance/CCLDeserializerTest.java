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

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.yaml.YAMLFactory;
import io.clouditor.assurance.ccl.BinaryComparison;
import io.clouditor.assurance.ccl.CCLDeserializer;
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

    var rule = mapper.readValue("condition: \"User has field == true\"", Rule.class);

    assertNotNull(rule);

    var condition = rule.getCondition();

    assertNotNull(condition);

    assertTrue(condition.getExpression() instanceof BinaryComparison);
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
    var condition = ccl.parse("User has createDate before now");

    var asset = new AssetProperties();
    asset.put("createDate", Map.of("epochSecond", 1550131042, "nano", 1));

    assertTrue(condition.evaluate(asset));

    asset.clear();
    // something in the future
    asset.put("createDate", Map.of("epochSecond", Integer.MAX_VALUE, "nano", 1));

    assertFalse(condition.evaluate(asset));

    condition = ccl.parse("User has createDate after now");

    assertTrue(condition.evaluate(asset));

    // something older than 90 days
    asset.put(
        "createDate",
        Map.of(
            "epochSecond", Instant.now().minus(100, ChronoUnit.DAYS).getEpochSecond(), "nano", 0));

    // user create date should not be older than 90 days
    condition = ccl.parse("User has not createDate before 90 days");

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
}
