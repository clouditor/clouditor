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

import static io.clouditor.auth.AuthenticationService.ROLE_USER;

import io.clouditor.discovery.DiscoveryService;
import io.clouditor.discovery.Scan;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.GET;
import javax.ws.rs.NotFoundException;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Path("discovery")
@RolesAllowed(ROLE_USER)
public class DiscoveryResource {

  private static final Logger LOGGER = LoggerFactory.getLogger(DiscoveryResource.class);

  private final DiscoveryService service;

  @Inject
  public DiscoveryResource(DiscoveryService service) {
    this.service = service;
  }

  @Produces(MediaType.APPLICATION_JSON)
  @GET
  public List<Scan> getScans() {
    return new ArrayList<>(this.service.getScans().values());
  }

  @GET
  @Path("{id}")
  public Scan getScan(@PathParam("id") String id) {
    return this.service.getScan(id);
  }

  @POST
  @Path("{id}/enable")
  public void enable(@PathParam("id") String id) throws IOException {
    var scan = service.getScan(id);

    if (scan == null) {
      LOGGER.error("Could not find scan with id {}", id);
      throw new NotFoundException("Could not find scan with id " + id);
    }

    service.enableScan(scan);
  }

  @POST
  @Path("{id}/disable")
  public void disable(@PathParam("id") String id) {
    var scan = service.getScan(id);

    if (scan == null) {
      LOGGER.error("Could not find scan with id {}", id);
      throw new NotFoundException("Could not find scan with id " + id);
    }

    service.disableScan(scan);
  }
}
