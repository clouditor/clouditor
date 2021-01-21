package io.clouditor.rest;

import static org.junit.jupiter.api.Assertions.*;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;
import io.clouditor.Engine;
import io.clouditor.auth.AuthenticationService;
import io.clouditor.auth.User;
import io.clouditor.discovery.DiscoveryService;
import java.util.List;
import javax.ws.rs.client.Entity;
import javax.ws.rs.core.Application;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import org.glassfish.jersey.test.JerseyTest;
import org.glassfish.jersey.test.TestProperties;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

public class UsersResourceTest extends JerseyTest {
  private static final Engine engine = new Engine();
  private String token;
  private static final String prefix = "/users/";
  private final String clouditorUserName = "clouditor";
  private static String mockUserName = "MockUser";
  private static String mockUser2Name = "MockUser2";

  /** Tests */
  @Test
  void givenGetUsers() {
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
  void givenGetUser_whenUserNotExist_thenStatusNotFound() {
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
  void givenGetUser_whenUserExist_thenStatusNotFound() {
    // Request
    Response response =
        target(prefix + mockUserName)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    User actualUser = response.readEntity(User.class);
    assertEquals(mockUserName, actualUser.getName());
  }

  @Test
  void givenUpdateUser() {
    // Preparation
    ObjectMapper objectMapper = new ObjectMapper();
    ObjectNode userAsJson = objectMapper.createObjectNode();
    String mockUserEmail = "mock@email.de";
    userAsJson.put("email", mockUserEmail);
    AuthenticationService authenticationService = engine.getService(AuthenticationService.class);
    //        User clouditor = authenticationService.getUser(clouditorUserName);
    String mockUserMailBeforeUpdate = authenticationService.getUser(mockUserName).getEmail();

    // Request
    Response response =
        target(prefix + mockUserName)
            .request(MediaType.APPLICATION_JSON_TYPE)
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .put(Entity.json(userAsJson));

    // Assertions
    assertEquals(Response.Status.NO_CONTENT.getStatusCode(), response.getStatus());
    User mockUser = authenticationService.getUser(mockUserName);
    assertEquals(mockUserEmail, mockUser.getEmail());
    assertNotEquals(mockUserMailBeforeUpdate, mockUser.getEmail());
  }

  @Test
  void givenDeleteUser() {
    // Pre assertion
    AuthenticationService authenticationService = engine.getService(AuthenticationService.class);
    assertNotNull(authenticationService.getUser(mockUser2Name));

    // Request
    Response response =
        target(prefix + mockUser2Name)
            .request(MediaType.APPLICATION_JSON_TYPE)
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .delete();

    // Post assertion
    assertEquals(Response.Status.NO_CONTENT.getStatusCode(), response.getStatus());
    assertNull(authenticationService.getUser(mockUser2Name));
  }

  @Test
  void givenCreateUser_whenUserAlreadyExists_StatusBadRequest() {
    // Preparation
    ObjectMapper objectMapper = new ObjectMapper();
    ObjectNode userAsJson = objectMapper.createObjectNode();
    userAsJson.put("username", mockUserName);
    userAsJson.put("password", mockUserName);

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
  void givenCreateUser_whenUserNotAlreadyExists_StatusOk() {
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

    // Creating a second user
    AuthenticationService authenticationService = engine.getService(AuthenticationService.class);
    User mockUser = new User(mockUserName, mockUserName);
    authenticationService.createUser(mockUser);
    User mockUser2 = new User(mockUser2Name, mockUser2Name);
    authenticationService.createUser(mockUser2);
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
      this.token = engine.authenticateAPI(target(), clouditorUserName, "clouditor");
    }
  }

  @Override
  protected Application configure() {
    // Find first available port.
    forceSet(TestProperties.CONTAINER_PORT, "0");
    return new EngineAPI(engine);
  }
}
