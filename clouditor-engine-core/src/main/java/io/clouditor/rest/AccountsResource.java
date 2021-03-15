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

import static io.clouditor.auth.AuthenticationService.ROLE_ADMIN;
import static io.clouditor.rest.AbstractAPI.sanitize;

import io.clouditor.credentials.AccountService;
import io.clouditor.credentials.CloudAccount;
import java.io.IOException;
import java.util.Map;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.*;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.Response.Status;

@Path("accounts")
@RolesAllowed(ROLE_ADMIN)
public class AccountsResource {

  private AccountService service;

  @Inject
  public AccountsResource(AccountService service) {
    this.service = service;
  }

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  public Map<String, CloudAccount> getAccounts() {
    return this.service.getAccounts();
  }

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  @Path("{provider}")
  public CloudAccount getAccount(@PathParam("provider") String provider) {
    provider = sanitize(provider);

    var account = this.service.getAccount(provider);

    if (account == null) {
      throw new NotFoundException();
    }

    return account;
  }

  @POST
  @Produces(MediaType.APPLICATION_JSON)
  @Path("discover/{provider}")
  public CloudAccount discover(@PathParam("provider") String provider) {
    provider = sanitize(provider);

    var account = this.service.discover(provider);

    if (account == null) {
      throw new NotFoundException();
    }

    return account;
  }

  @PUT
  @Consumes(MediaType.APPLICATION_JSON)
  @Path("{provider}")
  public void putAccount(@PathParam("provider") String provider, CloudAccount account) {
    provider = sanitize(provider);

    try {
      this.service.addAccount(provider, account);
    } catch (IOException ex) {
      throw new BadRequestException(
          Response.status(Status.BAD_REQUEST).entity(ex.getMessage()).build());
    }
  }
}
