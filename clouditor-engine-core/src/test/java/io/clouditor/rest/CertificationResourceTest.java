package io.clouditor.rest;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.Engine;
import io.clouditor.assurance.*;
import io.clouditor.discovery.DiscoveryService;
import java.util.*;
import javax.ws.rs.NotFoundException;
import javax.ws.rs.client.Entity;
import javax.ws.rs.core.Application;
import javax.ws.rs.core.HttpHeaders;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import org.glassfish.jersey.test.JerseyTest;
import org.glassfish.jersey.test.TestProperties;
import org.junit.jupiter.api.*;

@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
public class CertificationResourceTest extends JerseyTest {
  private static final Engine engine = new Engine();
  private String token;

  /* Test Settings */
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

    // Clean all certificates from previous test suites
    engine.getService(CertificationService.class).removeAllCertifications();
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
  public void testCertificationResource_constructor() {
    target("certification")
        .request()
        .header(
            AuthenticationFilter.HEADER_AUTHORIZATION,
            AuthenticationFilter.createAuthorization(token))
        .get();
  }

  @Test
  @Order(1)
  public void
      testGetCertifications_whenNoCertificationsAvailable_thenStatusOKButEmptyResponseContent() {
    // Execute request
    Response response =
        target("/certification")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();
    // Check conditions
    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    Assertions.assertEquals(
        MediaType.APPLICATION_JSON, response.getHeaderString(HttpHeaders.CONTENT_TYPE));
    var responseCertifications = response.readEntity(Map.class);
    Assertions.assertTrue(responseCertifications.isEmpty());
  }

  @Test
  public void testModifyControlStatus_whenCertificationIsNullAndControlIsNull_thenNotFound() {
    // Tests
    Response response =
        target("certification/1/1/status")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(Entity.json("{}"));

    // Assertions
    Assertions.assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());
  }

  @Test
  public void
      testModifyControlStatus_whenCertificationNotNullAndControlIsNull_thenThrowException() {

    // Request
    CertificationService certService = engine.getService(CertificationService.class);
    Certification mockCertification = new Certification();
    mockCertification.setId("1");
    certService.modifyCertification(mockCertification);

    Response response =
        target("certification/1/1/status")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(Entity.json("{}"));

    // Assertions
    Assertions.assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());
  }

  @Test
  public void testModifyControlStatus_whenCertificationNotNullAndControlNotNull_thenStatus204() {
    // Create mock of certification  with id=1 (the one that will be searched for in
    // testModifyControlStatus())
    CertificationService certService = engine.getService(CertificationService.class);
    Certification mockCertification = new Certification();
    mockCertification.setId("1");
    // Print for debugging purposes

    // Create mock of control with id=2 (will be attached to the mocked certification)
    Control mockControl = new Control();
    // mockControl.setAutomated(true);
    mockControl.setControlId("2");
    mockControl.setDomain(new Domain("TestDomain"));

    // Add mocked control (as list of one control) to mocked certification
    List<Control> oneControlList = new ArrayList<>();
    oneControlList.add(mockControl);
    mockCertification.setControls(oneControlList);
    // Check if control is inside

    // Update the certificate in the certification service
    certService.modifyCertification(mockCertification);

    Response response =
        target("certification/1/2/status")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .post(Entity.entity("{}", MediaType.APPLICATION_JSON_TYPE));
    Assertions.assertEquals(
        Response.Status.NO_CONTENT.getStatusCode(),
        response.getStatus(),
        "HTTP Response should be 204, no content: ");
  }

  @Test
  public void testModifyControlStatus_whenStatusFalseAndControlActiveTrue_thenStopMonitoring() {
    // Create mock of certification  with id=1 (the one that will be searched for in
    // testModifyControlStatus())
    CertificationService certService = engine.getService(CertificationService.class);
    Certification mockCertification = new Certification();
    mockCertification.setId("1");

    // Create mock of control with id=2 (will be attached to the mocked certification)
    Control mockControl = new Control();
    // mockControl.setAutomated(true);
    mockControl.setControlId("2");
    mockControl.setDomain(new Domain("TestDomain"));
    mockControl.setActive(true);

    // Add mocked control (as list of one control) to mocked certification
    List<Control> oneControlList = new ArrayList<>();
    oneControlList.add(mockControl);
    mockCertification.setControls(oneControlList);

    // Update the certificate in the certification service
    certService.modifyCertification(mockCertification);

    // Create ControlStatusRequest
    CertificationResource.ControlStatusRequest controlStatusRequest =
        new CertificationResource.ControlStatusRequest();

    Assertions.assertTrue(mockControl.isActive());
    Response response =
        target("certification/1/2/status")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .post(Entity.entity(controlStatusRequest, MediaType.APPLICATION_JSON_TYPE));
    Assertions.assertEquals(
        Response.Status.NO_CONTENT.getStatusCode(),
        response.getStatus(),
        "HTTP Response should be 204: ");
    Assertions.assertFalse(mockControl.isActive());
  }

  @Test
  public void testModifyControlStatus_whenStatusTrueAndControlActiveFalse_thenStartMonitoring() {
    // Create mock of certification  with id=1 (the one that will be searched for in
    // testModifyControlStatus())
    CertificationService certService = engine.getService(CertificationService.class);
    Certification mockCertification = new Certification();
    mockCertification.setId("1");

    // Create mock of control with id=2 (will be attached to the mocked certification)
    Control mockControl = new Control();
    // mockControl.setAutomated(true);
    mockControl.setControlId("2");
    mockControl.setDomain(new Domain("TestDomain"));
    mockControl.setActive(false);
    mockControl.setAutomated(true);

    // Add mocked control (as list of one control) to mocked certification
    List<Control> oneControlList = new ArrayList<>();
    oneControlList.add(mockControl);
    mockCertification.setControls(oneControlList);

    // Update the certificate in the certification service
    certService.modifyCertification(mockCertification);

    // Create ControlStatusRequest
    // Using the innerclass of CertificationResource not possible since variable status is
    // unreachable
    CertificationResourceTest.ControlStatusRequest controlStatusRequest =
        new CertificationResourceTest.ControlStatusRequest(true);
    Assertions.assertTrue(controlStatusRequest.status);

    // Get the control from the certification service and assert it is not active
    Map<String, Certification> certifications = certService.getCertifications();
    Certification certification = certifications.get("1");
    List<Control> controls = certification.getControls();
    Control control = controls.get(0);
    Assertions.assertFalse(control.isActive());

    target("certification/1/2/status")
        .request()
        .header(
            AuthenticationFilter.HEADER_AUTHORIZATION,
            AuthenticationFilter.createAuthorization(this.token))
        .post(Entity.entity(controlStatusRequest, MediaType.APPLICATION_JSON_TYPE));
    // Assert the control now is active
    Assertions.assertTrue(control.isActive());
  }

  @Test
  public void testGetCertifications_whenCorrectRequest_thenResponseIsOkAndContainsCertificate() {
    CertificationService certService = engine.getService(CertificationService.class);
    Certification mockCertification = new Certification();
    mockCertification.setId("Mock1");

    certService.modifyCertification(mockCertification);

    Response response =
        target("/certification/Mock1")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();
    Assertions.assertEquals(
        MediaType.APPLICATION_JSON,
        response.getHeaderString(HttpHeaders.CONTENT_TYPE),
        "Http Content-Type should be: APPLICATION_JSON");
  }

  @Test
  public void testGetCertifications_whenTwoCertificationsAvailable_thenStatusOkAndReturnBoth() {
    // Create two mocks of certification  with id=1 and id=2
    var id1 = "1";
    var id2 = "2";
    CertificationService certService = engine.getService(CertificationService.class);
    Certification mockCertification1 = new Certification();
    mockCertification1.setId(id1);
    // mockCertification1.setDescription("Description1");
    certService.modifyCertification(mockCertification1);
    Certification mockCertification2 = new Certification();
    mockCertification2.setId(id2);
    certService.modifyCertification(mockCertification2);

    Response response =
        target("/certification")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    var certifications = response.readEntity(Map.class);
    var certificationsOfCertService = certService.getCertifications();
    Assertions.assertFalse(certifications.isEmpty());
    Assertions.assertTrue(certifications.containsKey(id1));
    Assertions.assertTrue(certifications.containsKey(id2));
    Assertions.assertTrue(
        compareExpectedAndActualCertification(
            certificationsOfCertService.get(id1), (HashMap<?, ?>) certifications.get(id1)));
    Assertions.assertTrue(
        compareExpectedAndActualCertification(
            certificationsOfCertService.get(id2), (HashMap<?, ?>) certifications.get(id2)));
  }

  @Test
  public void testGetCertification_whenCertificationAvailable_thenStatusOkAndReturnIt() {
    // Create one mock of certification  with id=1
    var id = "1";
    var description = "Test Description";
    CertificationService certService = engine.getService(CertificationService.class);
    Certification mockCertification1 = new Certification();
    mockCertification1.setId(id);
    mockCertification1.setDescription(description);
    certService.modifyCertification(mockCertification1);

    // Execute get request
    Response response =
        target("/certification/1")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    Certification certification1 = response.readEntity(Certification.class);
    assertEquals(mockCertification1, certification1);
  }

  @Test
  public void
      testGetCertification_whenNoCertificationAvailableWithGivenId_thenStatus404AndThrowException() {

    // Execute first get request asserting the status code
    Response response =
        target("/certification/66")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();
    Assertions.assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());

    assertThrows(
        NotFoundException.class,
        () -> {
          {
            // Execute second get request asserting the error
            target("/certification/66")
                .request()
                .header(
                    AuthenticationFilter.HEADER_AUTHORIZATION,
                    AuthenticationFilter.createAuthorization(token))
                .get(Certification.class);
          }
        });
  }

  @Test
  public void testGetControl_whenNoControlAvailableWithGivenId_thenStatus404AndThrowException() {
    Certification mockCertification = new Certification();
    mockCertification.setId("1");
    CertificationService certService = engine.getService(CertificationService.class);
    certService.modifyCertification(mockCertification);

    // Execute first get request (get status code)
    Response response =
        target("/certification/1/2")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();
    Assertions.assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());

    // Execute second get request (get content)
    assertThrows(
        NotFoundException.class,
        () -> {
          // Execute second get request
          target("/certification/1/2")
              .request()
              .header(
                  AuthenticationFilter.HEADER_AUTHORIZATION,
                  AuthenticationFilter.createAuthorization(token))
              .get(Control.class);
        });
  }

  @Test
  public void testGetControl_whenControlAvailableWithGivenId_thenStatusOkAndReturnIt() {
    // Create mock of certification  with id=1 (the one that will be searched for in
    // testModifyControlStatus())
    var certificationId = "1";
    var controlId = "5";

    CertificationService certService = engine.getService(CertificationService.class);
    Certification mockCertification = new Certification();
    mockCertification.setId(certificationId);

    // Create mock of control with id=2 (will be attached to the mocked certification)
    Control mockControl = new Control();
    // mockControl.setAutomated(true);
    mockControl.setControlId(controlId);
    mockControl.setDomain(new Domain("TestDomain"));

    // Add mocked control (as list of one control) to mocked certification
    List<Control> oneControlList = new ArrayList<>();
    oneControlList.add(mockControl);
    mockCertification.setControls(oneControlList);

    // Update the certificate in the certification service
    certService.modifyCertification(mockCertification);

    Response response =
        target(String.format("certification/%s/%s", certificationId, controlId))
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();
    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    assertEquals(response.readEntity(Control.class).getId(), controlId);
  }

  @Test
  public void
      testImportCertification_whenNoCertificationAvailable_thenStatus404AndThrowException() {
    // Execute first Post Request (for status)
    Response response =
        target("certification/import/1")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .post(Entity.json("{}"));
    Assertions.assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());
  }

  @Test
  public void
      testImportCertification_whenCertificationAvailableWithGivenId_thenStatusOkAndCertificateThere() {
    // Get CertificationService
    CertificationService certService = engine.getService(CertificationService.class);

    // Verify that there are no certifications currently available (by checking the hash map of
    // certifications)
    Assertions.assertNull(certService.getCertifications().get("BSI C5"));
    Response response =
        target("certification/import/BSI C5")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .post(Entity.json("{}"));
    // Check if status is 204, No Content
    Assertions.assertEquals(Response.Status.NO_CONTENT.getStatusCode(), response.getStatus());
    // Check if certification hashMap is not empty any more
    Assertions.assertNotNull(certService.getCertifications().get("BSI C5"));
    // Check if the certification from firstCertificationImporter was properly imported
    //    Assertions.assertNotNull(
    //        certService.getCertifications().get(firstCertificationImporter.getKey()));
  }

  @Test
  public void testGetImporters_when_then() {
    // Get CertificationService
    CertificationService certService = engine.getService(CertificationService.class);
    // Get Importers
    var importers = certService.getImporters();
    Response response =
        target("certification/importers")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(this.token))
            .get();
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());

    Map<?, ?> certificationImporter = response.readEntity(Map.class);
    Assertions.assertFalse(certificationImporter.isEmpty());
    Assertions.assertTrue(certificationImporter.containsKey("BSI C5"));
  }

  /** Helper classes and methods */

  // The corresponding inner class of CertificationResource does not allow to change
  // the status
  public static class ControlStatusRequest {

    @JsonProperty private final boolean status;

    public ControlStatusRequest(boolean status) {
      this.status = status;
    }
  }

  private boolean compareExpectedAndActualCertification(
      Certification expectedCertification, HashMap<?, ?> actualCertification) {
    String expectedDescription = expectedCertification.getDescription();
    String expectedPublisher = expectedCertification.getPublisher();
    String expectedWebsite = expectedCertification.getWebsite();
    List<Control> expectedControls = expectedCertification.getControls();

    if (!(expectedCertification.getId().equals(actualCertification.get("_id")))) {
      return false;
    } else if (expectedDescription != null
        && !(expectedDescription.equals(actualCertification.get("description")))) {
      return false;
    } else if (expectedPublisher != null
        && !(expectedPublisher.equals(actualCertification.get("publisher")))) {
      return false;
    } else if (expectedWebsite != null
        && !(expectedWebsite.equals(actualCertification.get("publisher")))) {
      return false;
    } else
      return expectedControls.isEmpty() || expectedControls == actualCertification.get("controls");
  }
}
