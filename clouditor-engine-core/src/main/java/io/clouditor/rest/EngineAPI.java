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

import io.clouditor.Engine;
import io.clouditor.auth.UserContext;
import io.swagger.annotations.ApiKeyAuthDefinition;
import io.swagger.annotations.ApiKeyAuthDefinition.ApiKeyLocation;
import io.swagger.annotations.SecurityDefinition;
import io.swagger.annotations.SwaggerDefinition;
import io.swagger.jaxrs.config.BeanConfig;
import io.swagger.jaxrs.listing.ApiListingResource;
import io.swagger.jaxrs.listing.SwaggerSerializers;
import javax.annotation.security.RolesAllowed;
import javax.ws.rs.ApplicationPath;

/**
 * The Engine REST API.
 *
 * @author Banse, Christian
 */
@ApplicationPath(EngineAPI.CONTEXT_PATH)
@RolesAllowed(UserContext.ROLE_USERS)
@SwaggerDefinition(securityDefinition = @SecurityDefinition(apiKeyAuthDefintions = @ApiKeyAuthDefinition(key = "token", name = "Authorization", in = ApiKeyLocation.HEADER)))
public class EngineAPI extends AbstractAPI<Engine> {

  static final String CONTEXT_PATH = "engine";

  /**
   * Constructs a new {@link EngineAPI} from an {@link Engine}.
   *
   * @param engine The Clouditor Engine
   */
  public EngineAPI(Engine engine) {
    super(engine, engine.getAPIPort(), CONTEXT_PATH);

    var beanConfig = new BeanConfig();
    beanConfig.setVersion("1.0.0");
    beanConfig.setSchemes(new String[] { "http" });
    beanConfig.setHost("clouditor.io");
    beanConfig.setBasePath(EngineAPI.CONTEXT_PATH);
    beanConfig.setResourcePackage(EngineAPI.class.getPackage().getName());
    beanConfig.setScan(true);

    this.register(ApiListingResource.class);
    this.register(SwaggerSerializers.class);
  }
}
