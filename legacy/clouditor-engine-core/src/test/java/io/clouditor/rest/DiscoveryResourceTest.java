package io.clouditor.rest;

import io.clouditor.Engine;
import io.clouditor.discovery.DiscoveryService;
import io.clouditor.discovery.Scan;
import java.util.List;
import javax.ws.rs.client.Entity;
import javax.ws.rs.core.Application;
import javax.ws.rs.core.Response;
import org.glassfish.jersey.test.JerseyTest;
import org.glassfish.jersey.test.TestProperties;
import org.junit.jupiter.api.*;

class DiscoveryResourceTest extends JerseyTest {
  private static final Engine engine = new Engine();
  private String token;
  private static final String targetPrefix = "/discovery/";
  private static final String FAKE_ID = "fake";
  private static DiscoveryService discoveryService;

  /* Test Settings */
  @BeforeAll
  static void startUpOnce() {
    engine.setDbInMemory(true);

    engine.setDBName("DiscoveryResourceTestDB");

    // init db
    engine.initDB();

    // initialize every else
    engine.init();

    // start the DiscoveryService
    engine.getService(DiscoveryService.class).start();

    // Initialize discoveryService
    discoveryService = engine.getService(DiscoveryService.class);
  }

  @AfterAll
  static void cleanUpOnce() {
    engine.shutdown();
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

  /* Tests */
  @Test
  void testGetScansWhenOneScannerAvailableThenStatusOkAndResponseNotEmpty() {
    // Request
    var response =
        target(targetPrefix)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    List<Scan> scans = response.readEntity(List.class);
    Assertions.assertFalse(scans.isEmpty());
    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
  }

  @Test
  void testGetScanWhenRequestedScannerAvailableThenStatusOkAndRespondIt() {
    // Request
    var response =
        target(targetPrefix + FAKE_ID)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    Assertions.assertNotNull(discoveryService.getScan(FAKE_ID));
    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    Scan scan = response.readEntity(Scan.class);
    Assertions.assertEquals(FAKE_ID, scan.getId());
  }

  @Test
  void testGetScanWhenRequestedScannerNotAvailableThenStatusOkAndRespondIt() {
    // Preparation
    String id = "I Am Not There";

    // Request
    var response =
        target(targetPrefix + id)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    Assertions.assertNull(discoveryService.getScan(id));
    Assertions.assertEquals(Response.Status.NO_CONTENT.getStatusCode(), response.getStatus());
    Scan scan = response.readEntity(Scan.class);
    Assertions.assertNull(scan);
  }

  @Test
  void testEnableWhenScannerIsAvailableAndNotEnabledThenScanEnabledStatusNoContent() {
    // Preparation
    Scan scan = discoveryService.getScan(FAKE_ID);
    scan.setEnabled(false);

    // Assert that scan is disabled before request
    Assertions.assertFalse(scan.isEnabled());

    // Request
    var response =
        target(targetPrefix + FAKE_ID + "/enable")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(Entity.json("{}"));

    // Assertions
    Assertions.assertTrue(scan.isEnabled());
    Assertions.assertEquals(Response.Status.NO_CONTENT.getStatusCode(), response.getStatus());
  }

  @Test
  void testEnableWhenScannerIsNotAvailableThenStatusNotFound() {
    // Request
    String id = "I am Not There";
    var response =
        target(targetPrefix + id + "/enable")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(Entity.json("{}"));

    // Assertions
    Assertions.assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());
  }

  @Test
  void testDisableWhenScannerIsAvailableAndEnabledThenStatusNoContent() {
    // Preparation
    Scan scan = discoveryService.getScan(FAKE_ID);
    scan.setEnabled(true);

    // Request
    var response =
        target(targetPrefix + FAKE_ID + "/disable")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(Entity.json("{}"));

    // Assertions
    Assertions.assertFalse(scan.isEnabled());
    Assertions.assertEquals(Response.Status.NO_CONTENT.getStatusCode(), response.getStatus());
  }

  @Test
  void testDisableWhenScannerIsNotAvailableThenStatusNotFound() {
    // Request
    String id = "I am Not There";
    var response =
        target(targetPrefix + id + "/disable")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(Entity.json("{}"));

    // Assertions
    Assertions.assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());
  }
}
