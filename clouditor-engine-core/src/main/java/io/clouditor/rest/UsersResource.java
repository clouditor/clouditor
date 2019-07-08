package io.clouditor.rest;

import static io.clouditor.auth.AuthenticationService.ROLE_ADMIN;
import static io.clouditor.rest.AbstractAPI.sanitize;

import io.clouditor.auth.AuthenticationService;
import io.clouditor.auth.User;
import java.util.List;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.Consumes;
import javax.ws.rs.DELETE;
import javax.ws.rs.GET;
import javax.ws.rs.NotFoundException;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.Response.Status;

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

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  @Path("{id}")
  public User getUser(@PathParam("id") String id) {
    id = sanitize(id);

    var user = this.service.getUser(id);

    if (user == null) {
      throw new NotFoundException("User does not exist");
    }

    return user;
  }

  @PUT
  @Consumes(MediaType.APPLICATION_JSON)
  @Path("{id}")
  public void updateUser(@PathParam("id") String id, User user) {
    id = sanitize(id);

    this.service.updateUser(id, user);
  }

  @DELETE
  @Consumes(MediaType.APPLICATION_JSON)
  @Path("{id}")
  public void deleteUser(@PathParam("id") String id) {
    id = sanitize(id);

    this.service.deleteUser(id);
  }

  @POST
  @Consumes(MediaType.APPLICATION_JSON)
  public Response createUser(User user) {
    if (!this.service.createUser(user)) {
      return Response.status(Status.BAD_REQUEST).entity("User already exists").build();
    }

    return Response.ok().build();
  }
}
