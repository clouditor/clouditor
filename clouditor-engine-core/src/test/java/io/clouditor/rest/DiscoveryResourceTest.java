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

public class DiscoveryResourceTest extends JerseyTest {
  private static final Engine engine = new Engine();
  private String token;
  private static final String targetPrefix = "/discovery/";

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
  public void testGetScans_whenOneScannerAvailable_thenStatusOkAndResponseNotEmpty() {
    // Request
    Response response =
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
  public void testGetScan_whenRequestedScannerAvailable_thenStatusOkAndRespondIt() {
    // Preparation
    String id = "fake";
    DiscoveryService discoveryService = engine.getService(DiscoveryService.class);

    // Request
    Response response =
        target(targetPrefix + id)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    Assertions.assertNotNull(discoveryService.getScan(id));
    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    Scan scan = response.readEntity(Scan.class);
    Assertions.assertEquals("fake", scan.getId());
  }

  @Test
  public void testGetScan_whenRequestedScannerNotAvailable_thenStatusOkAndRespondIt() {
    // Preparation
    String id = "I Am Not There";
    DiscoveryService discoveryService = engine.getService(DiscoveryService.class);

    // Request
    Response response =
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
  public void testEnable_whenScannerIsAvailable_thenScanEnabledStatusNoContent() {
    // Preparation
    String id = "fake";
    DiscoveryService discoveryService = engine.getService(DiscoveryService.class);
    Scan scan = discoveryService.getScan(id);
    discoveryService.disableScan(scan);

    // Request
    Response response =
        target(targetPrefix + id + "/enable")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(Entity.json("{}"));

    // Assertions
    scan = discoveryService.getScan(id);
    Assertions.assertTrue(scan.isEnabled());
    Assertions.assertEquals(Response.Status.NO_CONTENT.getStatusCode(), response.getStatus());
  }

  @Test
  public void testEnable_whenScannerIsNotAvailable_thenStatusNotFound() {
    // Request
    String id = "I am Not There";
    Response response =
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
  public void testDisable_whenScannerIsAvailable_thenStatusNoContent() {
    // Preparation
    String id = "fake";
    DiscoveryService discoveryService = engine.getService(DiscoveryService.class);
    Scan scan = discoveryService.getScan(id);
    discoveryService.enableScan(scan);

    // Request
    Response response =
        target(targetPrefix + id + "/disable")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(Entity.json("{}"));

    // Assertions
    scan = discoveryService.getScan(id);
    Assertions.assertFalse(scan.isEnabled());
    Assertions.assertEquals(Response.Status.NO_CONTENT.getStatusCode(), response.getStatus());
  }

  @Test
  public void testDisable_whenScannerIsNotAvailable_thenStatusNotFound() {
    // Request
    String id = "I am Not There";
    Response response =
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
