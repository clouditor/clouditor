package io.clouditor.rest;

import static org.junit.jupiter.api.Assertions.*;

import io.clouditor.Engine;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.AssetService;
import io.clouditor.discovery.DiscoveryService;
import java.util.Iterator;
import java.util.Map;
import java.util.Set;
import javax.ws.rs.core.Application;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import org.glassfish.jersey.test.JerseyTest;
import org.glassfish.jersey.test.TestProperties;
import org.junit.jupiter.api.*;

public class AssetsResourceTest extends JerseyTest {
  private static final Engine engine = new Engine();
  private String token;
  private static final String prefix = "/assets/";

  /** Tests */
  @Test
  void givenGetAssetsWithType_whenNoAssetAvailable_then() {
    // Request
    Response response =
        target(prefix + "This Asset Type Does Not Exist")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    assertTrue(response.readEntity(Set.class).isEmpty());
  }

  @Test
  void givenGetAssetsWithType() {
    // Preparation
    AssetService assetService = engine.getService(AssetService.class);
    Asset mockAsset = new Asset();
    mockAsset.setId("Mock Asset");
    mockAsset.setType("Mock Asset Type");
    assetService.update(mockAsset);

    // Request
    Response response =
        target(prefix + "Mock Asset Type")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    Iterator<?> iterator = response.readEntity(Set.class).iterator();
    Map<?, ?> actualAsset = (Map<?, ?>) iterator.next();
    assertEquals(mockAsset.getId(), actualAsset.get("_id"));
  }

  // ToDo: Request Response is 500 (but full coverage)
  @Test
  void givenGetServerSentEvents() {
    // Preparation
    AssetService assetService = engine.getService(AssetService.class);
    Asset mockAsset = new Asset();
    mockAsset.setId("Mock Asset");
    mockAsset.setType("Mock Asset Type");
    assetService.update(mockAsset);

    // Request
    Response response =
        target(prefix + "Mock Asset Type" + "/subscribe")
            .request(MediaType.SERVER_SENT_EVENTS_TYPE)
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
  }

  /** Test Settings */
  @BeforeAll
  static void startUpOnce() {
    engine.setDbInMemory(true);

    engine.setDBName("CertificationResourceTestDB");

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
}
