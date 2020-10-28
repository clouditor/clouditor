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

import io.clouditor.auth.AuthenticationService;
import io.clouditor.auth.LoginRequest;
import io.clouditor.auth.LoginResponse;
import io.clouditor.auth.User;
import io.clouditor.data_access_layer.HibernatePersistence;
import javax.inject.Inject;
import javax.ws.rs.Consumes;
import javax.ws.rs.NotAuthorizedException;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.NewCookie;
import javax.ws.rs.core.Response;

/**
 * A resource end-point to authenticate.
 *
 * @author Christian Banse
 */
@Path("authenticate")
public class AuthenticateResource {

  private final AuthenticationService service;

  @Inject
  public AuthenticateResource(AuthenticationService service) {
    this.service = service;
  }

  @POST
  @Consumes(MediaType.APPLICATION_JSON)
  @Produces(MediaType.APPLICATION_JSON)
  public Response login(LoginRequest request) {
    var payload = new LoginResponse();

    if (!service.verifyLogin(request)) {
      throw new NotAuthorizedException("Invalid user and/or password");
    }

    var user = new HibernatePersistence().get(User.class, request.getUsername()).orElseThrow();

    payload.setToken(service.createToken(user));

    // TODO: max age, etc.
    return Response.ok(payload).cookie(new NewCookie("authorization", payload.getToken())).build();
  }
}
