package io.clouditor.rest;

import static org.junit.jupiter.api.Assertions.assertThrows;

import com.fasterxml.jackson.annotation.JsonTypeName;
import io.clouditor.Engine;
import io.clouditor.credentials.AccountService;
import io.clouditor.credentials.CloudAccount;
import io.clouditor.discovery.DiscoveryService;
import java.io.IOException;
import java.util.Map;
import javax.persistence.Entity;
import javax.persistence.Table;
import javax.ws.rs.NotFoundException;
import javax.ws.rs.core.Application;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import org.glassfish.jersey.test.JerseyTest;
import org.glassfish.jersey.test.TestProperties;
import org.junit.jupiter.api.*;

public class AccountsResourceTest extends JerseyTest {
  private static final Engine engine = new Engine();
  private String token;
  private static final String accountsPrefix = "/accounts/";

  @Test
  public void givenGetAccounts_whenNoAccountsAvailable_thenStatusOkAndResponseEmpty() {
    // Request
    Response response =
        target(accountsPrefix)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    Map<?, ?> accountsResource = response.readEntity(Map.class);
    Assertions.assertTrue(accountsResource.isEmpty());
  }

  /*
  ToDo: responded provider is "Mock Cloud" instead of "Mock Provider"
   */
  @Test
  public void givenGetAccounts_whenOneAccountAvailable_thenRespondWithAccount() {
    // Create account
    AccountService accService = engine.getService(AccountService.class);
    CloudAccount mockCloudAccount = new MockCloudAccount();
    String provider = "Mock Provider";
    try {
      accService.addAccount(provider, mockCloudAccount);
    } catch (IOException e) {
      e.printStackTrace();
    }

    // Request
    Response response =
        target(accountsPrefix)
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    Map<?, Map> accounts = response.readEntity(Map.class);
    Assertions.assertEquals(mockCloudAccount.getId(), accounts.get("Mock Cloud").get("_id"));
    // Assertions.assertEquals(provider, accounts.get("Mock Cloud").get("provider"));
    Assertions.assertEquals(
        mockCloudAccount.isAutoDiscovered(), accounts.get("Mock Cloud").get("autoDiscovered"));
  }

  @Test
  public void givenGetAccount_whenNoAccountAvailable_then404AndNull() {
    // Request
    Response response =
        target(accountsPrefix + "provider")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    // Assertions
    Assertions.assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());
    Assertions.assertNull(response.readEntity(CloudAccount.class));

    assertThrows(
        NotFoundException.class,
        () ->
            target(accountsPrefix + "provider")
                .request()
                .header(
                    AuthenticationFilter.HEADER_AUTHORIZATION,
                    AuthenticationFilter.createAuthorization(token))
                .get(CloudAccount.class));
  }

  @Test
  public void givenGetAccount_whenOneAccountAvailable_then200AndResponseWithAccount() {
    // Create account
    AccountService accService = engine.getService(AccountService.class);
    CloudAccount mockCloudAccount = new MockCloudAccount();
    String provider = "Mock Provider";
    try {
      accService.addAccount(provider, mockCloudAccount);
    } catch (IOException e) {
      e.printStackTrace();
    }

    // Request
    Response response =
        target(accountsPrefix + "Mock Cloud")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .get();

    CloudAccount responseCloudAccount = response.readEntity(MockCloudAccount.class);
    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
    Assertions.assertEquals(mockCloudAccount.getId(), responseCloudAccount.getId());
  }

  /*
  ToDo: AssertThrow
   */
  @Test
  public void givenDiscover_whenNoAccountAvailable_Then404AndNull() {
    // Request
    Response response =
        target(accountsPrefix + "discover/Mock Cloud")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(javax.ws.rs.client.Entity.json("{}"));

    // Assertions
    Assertions.assertNull(response.readEntity(CloudAccount.class));
    Assertions.assertEquals(Response.Status.NOT_FOUND.getStatusCode(), response.getStatus());

    //    assertThrows(
    //        NotFoundException.class,
    //        () ->
    //            target(accountsPrefix + "discover/Mock Cloud")
    //                .request()
    //                .header(
    //                    AuthenticationFilter.HEADER_AUTHORIZATION,
    //                    AuthenticationFilter.createAuthorization(token))
    //                .post(javax.ws.rs.client.Entity.json("{}")));
  }

  /*
  ToDo: Discover AWS. Problem: ClassNotFound, because it is in other module?
   */
  @Test
  public void givenDiscover_whenProviderAvailable_thenResponse() {

    // Request
    Response response =
        target(accountsPrefix + "discover/AWS")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .post(javax.ws.rs.client.Entity.json("{}"));

    System.out.println(response);
    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
  }

  /*
  ToDo: Passing on the request parameter CloudAccount object does not work. (With other object types it works)
   */
  @Test
  public void givenPutAccount() {
    // Create Account
    CloudAccount mockCloudAccount = new MockCloudAccount();

    // Request with account and provider as PathParam
    Response response =
        target(accountsPrefix + "AWS")
            .request()
            .header(
                AuthenticationFilter.HEADER_AUTHORIZATION,
                AuthenticationFilter.createAuthorization(token))
            .put(
                javax.ws.rs.client.Entity.entity(
                    mockCloudAccount, MediaType.APPLICATION_JSON_TYPE));

    Assertions.assertEquals(Response.Status.OK.getStatusCode(), response.getStatus());
  }

  // Helper methods and classes

  @Table(name = "mock_account")
  @Entity(name = "mock_account")
  @JsonTypeName(value = "Mock Cloud")
  private static class MockCloudAccount extends CloudAccount {

    @Override
    public void validate() {
      System.out.println("Mock Cloud Account validated.");
    }

    @Override
    public Object resolveCredentials() {
      return null;
    }
  }

  // Test Settings
  @Override
  protected Application configure() {
    // CONTAINER_PORT = 0 means first available port is used
    forceSet(TestProperties.CONTAINER_PORT, "0");
    return new EngineAPI(engine);
  }

  @BeforeEach
  public void startUp() {
    // Init DB
    engine.setDbInMemory(true);
    engine.setDBName("AccountsResourceTestDB");
    engine.initDB();

    // Init everything else
    engine.init();

    // Start DiscoveryService
    engine.getService(DiscoveryService.class).start();
  }

  @AfterEach
  public void cleanUp() {
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
}
