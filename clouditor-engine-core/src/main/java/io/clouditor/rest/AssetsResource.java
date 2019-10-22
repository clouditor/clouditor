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

import io.clouditor.discovery.Asset;
import io.clouditor.discovery.AssetService;
import io.clouditor.discovery.DiscoveryService;
import java.util.Set;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import org.glassfish.jersey.media.sse.EventOutput;
import org.glassfish.jersey.media.sse.SseFeature;

@Path("assets")
@RolesAllowed(ROLE_USER)
public class AssetsResource {

  @Inject private AssetService service;
  @Inject private DiscoveryService discoveryService;

  @GET
  @Path("{assetType}")
  @Produces(MediaType.APPLICATION_JSON)
  public Set<Asset> getAssetsWithType(@PathParam("assetType") String assetType) {
    assetType = sanitize(assetType);

    return this.service.getAssetsWithType(assetType);
  }

  @Path("{assetType}/subscribe")
  @GET
  @Produces(SseFeature.SERVER_SENT_EVENTS)
  public EventOutput getServerSentEvents(@PathParam("assetType") String assetType) {
    assetType = sanitize(assetType);

    return this.discoveryService.subscribeToEvents(assetType);
  }
}
