/*
 * Copyright (c) 2016-2019, Fraunhofer AISEC. All rights reserved.
 *
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
 *
 * Clouditor Community Edition is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Clouditor Community Edition is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * long with Clouditor Community Edition.  If not, see <https://www.gnu.org/licenses/>
 */

package io.clouditor.rest;

import io.clouditor.Component;
import io.clouditor.auth.User;
import io.clouditor.auth.UserContext;
import io.clouditor.auth.UserService;
import java.util.Objects;
import javax.annotation.Priority;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
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
  @Inject private UserService userService;

  @Context private ResourceInfo resourceInfo;

  public static String createAuthorization(String token) {
    return "Bearer " + token;
  }

  @Override
  public void filter(ContainerRequestContext requestContext) {
    // ignore filter for classes that do not have @RolesAllowed
    RolesAllowed rolesAllowed = resourceInfo.getResourceClass().getAnnotation(RolesAllowed.class);

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
      User user = userService.verifyToken(token);

      LOGGER.debug(
          "Authenticated API access to {} as {}",
          requestContext.getUriInfo().getPath(),
          user.getName());

      requestContext.setSecurityContext(
          new UserContext(user, requestContext.getSecurityContext().isSecure()));
    } catch (NotAuthorizedException ex) {
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
