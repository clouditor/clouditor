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

package io.clouditor.rest;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

import io.clouditor.Engine;
import io.clouditor.assurance.Certification;
import io.clouditor.assurance.CertificationService;
import io.clouditor.assurance.Rule;
import io.clouditor.assurance.RuleService;
import io.clouditor.auth.LoginRequest;
import io.clouditor.auth.User;
import io.clouditor.discovery.*;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.Set;
import javax.ws.rs.client.Entity;
import javax.ws.rs.core.Application;
import javax.ws.rs.core.GenericType;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.glassfish.jersey.test.JerseyTest;
import org.glassfish.jersey.test.TestProperties;
import org.junit.jupiter.api.*;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

class EngineAPIResourcesTest extends JerseyTest {

  private static final Engine engine = new Engine();

  private static final Logger LOGGER = LogManager.getLogger();
  private static final String ASSET_TYPE = "fake";

  private String token;

  @AfterAll
  static void cleanUpOnce() {
    engine.shutdown();
  }

  @BeforeAll
  static void startUpOnce() {
    engine.setDbInMemory(true);

    engine.setDBName("EngineAPIResorcesTestDB");

    // init db
    engine.initDB();

    // initialize every else
    engine.init();

    // start the DiscoveryService
    engine.getService(DiscoveryService.class).start();

    /*Job<?, ?> job = new Job<>(new TrueCheck());
    job.setId(EXPECTED_JOB_ID);

    job.setIterations(1);
    job.setInterval(0);

    LegacyAssetManager.getInstance()
        .update(null, LegacyAsset.of(IdGenerator.nextId(), "name", Object.class));*/
  }

  @BeforeEach
  public void setUp() throws Exception {
    super.setUp();

    client().register(ObjectMapperResolver.class);

    if (this.token == null) {
      this.token = engine.authenticateAPI(target(), "clouditor", "clouditor");
    }
  }

  @Override
  protected Application configure() {
    // Find first available port.
    forceSet(TestProperties.CONTAINER_PORT, "0");

    return new EngineAPI(engine);
  }

  @Test
  void testAuthenticate() {
    var fail =
        target("authenticate").request().post(Entity.json(new LoginRequest("wrong", "password")));

    assertEquals(401, fail.getStatus());

    var success =
        target("authenticate")
            .request()
            .post(
                Entity.json(
                    new LoginRequest(
                        engine.getDefaultApiUsername(), engine.getDefaultApiPassword())));

    assertEquals(200, success.getStatus());
  }

  @ParameterizedTest
  @ValueSource(
      strings = {"certification", "assets/" + ASSET_TYPE, "accounts", "discovery", "rules"})
  void testGetNotAuthenticated(String endpoint) {
    var response = target(endpoint).request().get();

    LOGGER.debug("Endpoint {} returned response code {}", endpoint, response.getStatus());

    assertEquals(401, response.getStatus());
  }

  @Test
  void testCertification() {
    // add some fake certification
    var mockCert = new Certification();
    mockCert.setId("mock-cert");

    var service = engine.getService(CertificationService.class);

    service.modifyCertification(mockCert);

    var certifications =
        target("certification")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .get(new GenericType<Map<String, Certification>>() {});

    assertEquals(service.getCertifications().get("mock-cert"), certifications.get("mock-cert"));

    assertTrue(certifications.containsKey(mockCert.getId()));

    var cert = certifications.get(mockCert.getId());

    assertEquals(mockCert, cert);
  }

  @Test
  void testGetAssets() {
    var service = engine.getService(AssetService.class);

    var asset = new Asset(ASSET_TYPE, "some-id", "some-name", new AssetProperties());

    service.update(asset);

    var assets =
        target("assets")
            .path(ASSET_TYPE)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .get(new GenericType<Set<Asset>>() {});

    assertNotNull(assets);
    assertFalse(assets.isEmpty());
  }

  @Test
  void testGetScans() {
    var service = engine.getService(DiscoveryService.class);

    var fakeScan = service.getScan(ASSET_TYPE);

    service.enableScan(fakeScan);

    var scans =
        target("discovery")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .get(new GenericType<List<Scan>>() {});

    assertNotNull(scans);
    assertFalse(scans.isEmpty());

    var scan = scans.get(0);

    assertTrue(scan.isEnabled());
  }

  @Test
  void testRules() throws IOException {
    var service = engine.getService(RuleService.class);

    service.load(FileSystemManager.getInstance().getPathForResource("rules/test"));

    var rules =
        target("rules")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .get(new GenericType<Map<String, Set<Rule>>>() {});

    assertNotNull(rules);

    var rule = rules.get("Asset").toArray()[0];

    assertNotNull(rule);
  }

  @Test
  void testUsers() {
    var users =
        target("users")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .get(new GenericType<List<User>>() {});

    assertNotNull(users);

    assertFalse(users.isEmpty());

    // store number of users (we need it later)
    var numberOfUsers = users.size();

    // add a user
    var guest = new User();
    guest.setPassword("test");
    guest.setUsername("test");

    var resp =
        target("users")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .post(Entity.json(guest));

    assertEquals(200, resp.getStatus());

    // retrieve user
    var guest2 =
        target("users")
            .path("test")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .get(User.class);

    assertNotNull(guest2);

    assertEquals(guest, guest2);

    // delete user again
    target("users")
        .path("test")
        .request()
        .header(
            AuthenticationFilter.HEADER_AUTHORIZATION,
            AuthenticationFilter.createAuthorization(this.token))
        .delete();

    // number of users should be back to beginning
    var users2 =
        target("users")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .get(new GenericType<List<User>>() {});

    assertNotNull(users2);
    assertEquals(numberOfUsers, users2.size());
  }
}
