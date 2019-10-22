/*
 * Copyright 2016-2019 Fraunhofer AISEC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *            $$\                           $$\ $$\   $$\
 *            $$ |                          $$ |\__|  $$ |
 *   $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 *  $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 *  $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ |  \__|
 *  $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 *  \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *   \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 */

package io.clouditor.rest;

import io.clouditor.Component;
import io.clouditor.auth.AuthenticationService;
import io.clouditor.auth.User;
import io.clouditor.auth.UserContext;
import java.util.Objects;
import javax.annotation.Priority;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.ForbiddenException;
import javax.ws.rs.NotAuthorizedException;
import javax.ws.rs.Priorities;
import javax.ws.rs.container.ContainerRequestContext;
import javax.ws.rs.container.ContainerRequestFilter;
import javax.ws.rs.container.ResourceInfo;
import javax.ws.rs.core.Context;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Priority(Priorities.AUTHENTICATION)
public class AuthenticationFilter implements ContainerRequestFilter {

  public static final String HEADER_AUTHORIZATION = "Authorization";
  protected static final Logger LOGGER = LoggerFactory.getLogger(AuthenticationFilter.class);
  @Inject private Component component;
  @Inject private AuthenticationService authenticationService;

  @Context private ResourceInfo resourceInfo;

  public static String createAuthorization(String token) {
    return "Bearer " + token;
  }

  @Override
  public void filter(ContainerRequestContext requestContext) {
    // ignore filter for classes that do not have @RolesAllowed
    var rolesAllowed = resourceInfo.getResourceClass().getAnnotation(RolesAllowed.class);

    if (rolesAllowed == null) {
      return;
    }

    // ignore filter for OPTIONS requests (pre-flight requests)
    if (Objects.equals(requestContext.getMethod(), "OPTIONS")) {
      return;
    }

    String authorization = requestContext.getHeaderString(HEADER_AUTHORIZATION);

    if (authorization == null || authorization.isEmpty()) {
      // try cookies
      var cookie = requestContext.getCookies().get("authentication");
      if (cookie != null) {
        authorization = cookie.getValue();
      }
    }

    if (authorization == null || !authorization.startsWith("Bearer")) {
      throw new NotAuthorizedException("No token was specified");
    }

    String[] rr = authorization.split(" ");

    if (rr.length != 2) {
      throw new NotAuthorizedException("Invalid authentication format");
    }

    String token = rr[1];

    try {
      User user = authenticationService.verifyToken(token);

      LOGGER.debug(
          "Authenticated API access to {} as {}",
          requestContext.getUriInfo().getPath(),
          user.getName());

      var ctx = new UserContext(user, requestContext.getSecurityContext().isSecure());

      requestContext.setSecurityContext(ctx);

      var authorized = false;

      for (var role : rolesAllowed.value()) {
        if (ctx.isUserInRole(role)) {
          authorized = true;
          break;
        }
      }

      if (!authorized) {
        throw new ForbiddenException(
            "User " + user.getName() + " does not have appropriate role to view resource.");
      }

    } catch (NotAuthorizedException | ForbiddenException ex) {
      // log the error
      LOGGER.error(
          "API access to {} was denied: {}",
          requestContext.getUriInfo().getPath(),
          ex.getMessage());

      // re-throw it
      throw ex;
    }
  }
}
