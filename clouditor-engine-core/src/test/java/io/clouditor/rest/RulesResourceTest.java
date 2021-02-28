package io.clouditor.rest;

import static org.junit.jupiter.api.Assertions.*;

import io.clouditor.Engine;
import io.clouditor.assurance.Rule;
import io.clouditor.assurance.RuleEvaluation;
import io.clouditor.assurance.RuleService;
import io.clouditor.discovery.DiscoveryService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import java.util.Iterator;
import java.util.Map;
import java.util.Set;
import javax.ws.rs.core.Application;
import javax.ws.rs.core.Response;
import org.glassfish.jersey.test.JerseyTest;
import org.glassfish.jersey.test.TestProperties;
import org.junit.jupiter.api.*;

@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
public class RulesResourceTest extends JerseyTest {
  private static final Engine engine = new Engine();
  private String token;
  private static final String targetPrefix = "/rules/";

  /** Tests */

  // ToDo: UI check if removing rules is possible (otherwise test can fail due to other tests)
  @Test
  @Order(1)
  public void testGetRules_whenNoRulesAvailable_thenStatusOkAndResponseEmpty() {
    // Request
    Response response =
        target(targetPrefix)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    assertTrue(response.readEntity(Map.class).isEmpty());
  }

  @Test
  public void testGetRules_thenAmountOfRulesIsEqual() {
    // Preparation
    RuleService ruleService = engine.getService(RuleService.class);
    try {
      ruleService.load(FileSystemManager.getInstance().getPathForResource("rules/test"));
    } catch (IOException e) {
      e.printStackTrace();
    }

    // Request
    Response response =
        target(targetPrefix)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    System.out.println(ruleService.getRules().size());
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    assertEquals(ruleService.getRules().size(), response.readEntity(Map.class).size());
  }

  @Test
  public void testGetRules_whenNoRulesWithAssetTypeAvailable_thenStatusOkAndResponseEmpty() {
    // Request
    Response response =
        target(targetPrefix + "assets/NoAssetWithThisName")
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
  public void testGetRules_whenRulesWithAssetTypeAvailable_thenStatusOkAndResponseEqual() {
    // Preparation
    RuleService ruleService = engine.getService(RuleService.class);
    try {
      ruleService.load(FileSystemManager.getInstance().getPathForResource("rules/test"));
    } catch (IOException e) {
      e.printStackTrace();
    }
    Iterator<?> iter;
    iter = ruleService.get("Asset").iterator();
    Rule rule = (Rule) iter.next();
    String expectedId = rule.getId();

    // Request
    Response response =
        target(targetPrefix + "assets/Asset")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    Set<?> responseSet = response.readEntity(Set.class);
    iter = responseSet.iterator();
    Map<?, ?> actualRule = (Map<?, ?>) iter.next();
    assertEquals(expectedId, actualRule.get("_id"));
  }

  @Test
  public void testGet_whenNoRuleWithIdAvailable_thenStatusNotFound() {
    // Request
    Response response =
        target(targetPrefix + "No Id With This Name")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());
    assertNull(response.readEntity(RuleEvaluation.class));
  }

  @Test
  public void testGet_whenRuleWithIdAvailable_thenStatusOkAndRespondIt() {
    // Preparation
    RuleService ruleService = engine.getService(RuleService.class);
    try {
      ruleService.load(FileSystemManager.getInstance().getPathForResource("rules/test"));
    } catch (IOException e) {
      e.printStackTrace();
    }
    Iterator<?> iter;
    iter = ruleService.get("Asset").iterator();
    Rule rule = (Rule) iter.next();
    String expectedId = rule.getId();

    // Request
    Response response =
        target(targetPrefix + expectedId)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    assertNotNull(response.readEntity(RuleEvaluation.class));
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
