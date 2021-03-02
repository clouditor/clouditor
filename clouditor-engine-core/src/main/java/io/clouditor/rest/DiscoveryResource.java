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

import static io.clouditor.auth.AuthenticationService.ROLE_USER;
import static io.clouditor.rest.AbstractAPI.sanitize;

import io.clouditor.discovery.DiscoveryService;
import io.clouditor.discovery.Scan;
import java.util.ArrayList;
import java.util.List;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.*;
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

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  public List<Scan> getScans() {
    return new ArrayList<>(this.service.getScans().values());
  }

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  @Path("{id}")
  public Scan getScan(@PathParam("id") String id) {
    id = sanitize(id);

    return this.service.getScan(id);
  }

  @POST
  @Path("{id}/enable")
  public void enable(@PathParam("id") String id) {
    id = sanitize(id);

    var scan = service.getScan(id);

    if (scan == null) {
      var sanitizedId = sanitize(id);
      LOGGER.error("Could not find scan with id {}", sanitizedId);
      throw new NotFoundException("Could not find scan with id " + id);
    }

    service.enableScan(scan);
  }

  @POST
  @Path("{id}/disable")
  public void disable(@PathParam("id") String id) {
    id = sanitize(id);

    var scan = service.getScan(id);

    if (scan == null) {
      LOGGER.error("Could not find scan with id {}", id);
      throw new NotFoundException("Could not find scan with id " + id);
    }

    service.disableScan(scan);
  }
}
