package io.clouditor.rest;

import io.clouditor.Engine;
import io.clouditor.assurance.Certification;
import io.clouditor.assurance.CertificationService;
import io.clouditor.discovery.DiscoveryService;
import java.util.Map;
import javax.ws.rs.core.Application;
import javax.ws.rs.core.Response;
import org.glassfish.jersey.test.JerseyTest;
import org.glassfish.jersey.test.TestProperties;
import org.junit.jupiter.api.*;

public class StatisticsResourceTest extends JerseyTest {
  private static final Engine engine = new Engine();
  private String token;
  private static final String targetPrefix = "/statistics/";

  /*
    Tests
  */
  @Test
  public void givenGetStatistic_when_then() {
    // Preparation: Add a certification
    CertificationService certificationService = engine.getService(CertificationService.class);
    Certification mockCertification = new Certification();
    mockCertification.setId("1");
    mockCertification.setDescription("I am a Mock Certification");
    certificationService.modifyCertification(mockCertification);
    Map<String, Certification> certifications = certificationService.getCertifications();

    // Request
    Response response =
        target(targetPrefix)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    System.out.println(response.getStatus());
    System.out.println(response);

    // Assertions
    Assertions.assertFalse(certifications.isEmpty());
    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
  }

  /*
  Test Settings
   */
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
