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

import io.clouditor.auth.UserContext;
import io.clouditor.credentials.AccountService;
import io.clouditor.credentials.CloudAccount;
import java.io.IOException;
import java.util.Map;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.BadRequestException;
import javax.ws.rs.GET;
import javax.ws.rs.NotFoundException;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.Response.Status;

@Path("accounts")
@RolesAllowed(UserContext.ROLE_USERS)
public class AccountsResource {

  private AccountService service;

  @Inject
  public AccountsResource(AccountService service) {
    this.service = service;
  }

  @GET
  public Map<String, CloudAccount> getAccounts() {
    return this.service.getAccounts();
  }

  @GET
  @Path("{provider}")
  public CloudAccount getAccount(@PathParam("provider") String provider) {
    var account = this.service.getAccount(provider);

    if (account == null) {
      throw new NotFoundException();
    }

    return account;
  }

  @POST
  @Path("discover/{provider}")
  public CloudAccount discover(@PathParam("provider") String provider) {
    var account = this.service.discover(provider);

    if (account == null) {
      throw new NotFoundException();
    }

    return account;
  }

  @PUT
  @Path("{provider}")
  public void putAccount(@PathParam("provider") String provider, CloudAccount account) {
    try {
      this.service.addAccount(provider, account);
    } catch (IOException ex) {
      throw new BadRequestException(
          Response.status(Status.BAD_REQUEST).entity(ex.getMessage()).build());
    }
  }
}
