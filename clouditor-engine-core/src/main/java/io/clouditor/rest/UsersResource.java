package io.clouditor.rest;

import static io.clouditor.auth.AuthenticationService.ROLE_ADMIN;

import io.clouditor.auth.AuthenticationService;
import io.clouditor.auth.User;
import java.util.List;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

@Path("users")
@RolesAllowed(ROLE_ADMIN)
public class UsersResource {

  private AuthenticationService service;

  @Inject
  public UsersResource(AuthenticationService service) {
    this.service = service;
  }

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  public List<User> getUsers() {
    return this.service.getUsers();
  }
}
