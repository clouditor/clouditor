package io.clouditor.rest;

import static org.junit.jupiter.api.Assertions.*;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;
import io.clouditor.Engine;
import io.clouditor.auth.AuthenticationService;
import io.clouditor.auth.User;
import io.clouditor.data_access_layer.HibernatePersistence;
import io.clouditor.discovery.DiscoveryService;
import java.util.List;
import javax.ws.rs.client.Entity;
import javax.ws.rs.core.Application;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import org.glassfish.jersey.test.JerseyTest;
import org.glassfish.jersey.test.TestProperties;
import org.junit.jupiter.api.*;

class UsersResourceTest extends JerseyTest {
  private static final Engine engine = new Engine();
  private String token;
  private static final String prefix = "/users/";
  private static AuthenticationService authenticationService;
  private final String clouditorUserName = "clouditor";
  private static final String MOCK_USER_NAME = "MockUser";
  private static final String MOCK_USER_2_NAME = "MockUser2";

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

    // Creating a second user
    authenticationService = engine.getService(AuthenticationService.class);
    User mockUser = new User(MOCK_USER_NAME, MOCK_USER_NAME);
    mockUser.setEmail("");
    authenticationService.createUser(mockUser);
    User mockUser2 = new User(MOCK_USER_2_NAME, MOCK_USER_2_NAME);
    authenticationService.createUser(mockUser2);
    mockUser2.setEmail("");
  }

  // Delete users from DB so that they do not potentially affect other test suites
  @AfterAll
  static void cleanUpOnce() {
    if (authenticationService.getUser(MOCK_USER_NAME) != null) {
      new HibernatePersistence().delete(User.class, MOCK_USER_NAME);
    }
    if (authenticationService.getUser(MOCK_USER_2_NAME) != null) {
      new HibernatePersistence().delete(User.class, MOCK_USER_2_NAME);
    }
  }

  @BeforeEach
  public void setUp() throws Exception {
    super.setUp();

    client().register(ObjectMapperResolver.class);

    if (this.token == null) {
      this.token = engine.authenticateAPI(target(), clouditorUserName, "clouditor");
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
  void testGetUsers() {
    AuthenticationService authenticationService = engine.getService(AuthenticationService.class);

    // Request
    Response response =
        target(prefix)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    List<?> users = response.readEntity(List.class);
    assertEquals(authenticationService.getUsers().size(), users.size());
  }

  @Test
  void testGetUser_whenUserNotExist_thenStatusNotFound() {
    // Request
    Response response =
        target(prefix + "Non-Existent User")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());
  }

  @Test
  void testGetUser_whenUserExist_thenStatusNotFound() {
    // Request
    Response response =
        target(prefix + MOCK_USER_NAME)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    User actualUser = response.readEntity(User.class);
    assertEquals(MOCK_USER_NAME, actualUser.getName());
  }

  @Test
  void testUpdateUser() {
    // Preparation
    String mockUserEmail = "mock@mail.io";
    User user = new User();
    user.setEmail(mockUserEmail);
    AuthenticationService authenticationService = engine.getService(AuthenticationService.class);
    String mockUserMailBeforeUpdate = authenticationService.getUser(MOCK_USER_NAME).getEmail();

    // Request
    Response response =
        target(prefix + MOCK_USER_NAME)
            .request(MediaType.APPLICATION_JSON_TYPE)
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .put(Entity.json(user));

    // Assertions
    assertEquals(Response.Status.NO_CONTENT.getStatusCode(), response.getStatus());
    User mockUser = authenticationService.getUser(MOCK_USER_NAME);
    assertEquals(mockUserEmail, mockUser.getEmail());
    assertNotEquals(mockUserMailBeforeUpdate, mockUser.getEmail());
  }

  @Test
  void testDeleteUser() {
    // Pre assertion
    AuthenticationService authenticationService = engine.getService(AuthenticationService.class);
    assertNotNull(authenticationService.getUser(MOCK_USER_2_NAME));

    // Request
    Response response =
        target(prefix + MOCK_USER_2_NAME)
            .request(MediaType.APPLICATION_JSON_TYPE)
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .delete();

    // Post assertion
    assertEquals(Response.Status.NO_CONTENT.getStatusCode(), response.getStatus());
    assertNull(authenticationService.getUser(MOCK_USER_2_NAME));
  }

  @Test
  void testCreateUser_whenUserAlreadyExists_StatusBadRequest() {
    // Preparation
    ObjectMapper objectMapper = new ObjectMapper();
    ObjectNode userAsJson = objectMapper.createObjectNode();
    userAsJson.put("username", MOCK_USER_NAME);
    userAsJson.put("password", MOCK_USER_NAME);

    // Request
    Response response =
        target(prefix)
            .request(MediaType.APPLICATION_JSON_TYPE)
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(Entity.json(userAsJson));

    // Assertions
    assertEquals(Response.Status.BAD_REQUEST.getStatusCode(), response.getStatus());
  }

  @Test
  void testCreateUser_whenUserNotAlreadyExists_StatusOk() {
    // Preparation
    ObjectMapper objectMapper = new ObjectMapper();
    ObjectNode userAsJson = objectMapper.createObjectNode();
    String newUserUsername = "Complete New Mock User Name";
    userAsJson.put("username", newUserUsername);
    userAsJson.put("password", "Complete New Mock User Password");
    AuthenticationService authenticationService = engine.getService(AuthenticationService.class);
    // Pre Assertion
    assertNull(authenticationService.getUser(newUserUsername));

    // Request
    Response response =
        target(prefix)
            .request(MediaType.APPLICATION_JSON_TYPE)
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(Entity.json(userAsJson));

    // Assertions
    assertNotNull(authenticationService.getUser(newUserUsername));
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
  }
}
