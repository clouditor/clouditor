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
import io.swagger.jaxrs.config.BeanConfig;
import io.swagger.jaxrs.listing.ApiListingResource;
import io.swagger.jaxrs.listing.SwaggerSerializers;
import javax.validation.constraints.NotNull;
import javax.ws.rs.core.UriBuilder;
import org.glassfish.grizzly.http.server.HttpServer;
import org.glassfish.grizzly.http.server.NetworkListener;
import org.glassfish.grizzly.http.server.StaticHttpHandler;
import org.glassfish.jersey.grizzly2.httpserver.GrizzlyHttpServerFactory;
import org.glassfish.jersey.server.ResourceConfig;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** Abstract base class for a REST API server. */
public abstract class AbstractAPI<C extends Component> extends ResourceConfig {

  /** The logger. */
  protected static final Logger LOGGER = LoggerFactory.getLogger(AbstractAPI.class);

  /** The grizzly HTTP server. */
  protected HttpServer httpServer;

  /** The API port. */
  private int port;

  /** The context path */
  @NotNull
  private String contextPath;

  /** The associated Clouditor component */
  @NotNull
  private C component;

  AbstractAPI(C component, int port, String contextPath) {
    this.component = component;
    this.port = port;
    this.contextPath = contextPath;

    BeanConfig beanConfig = new BeanConfig();
    beanConfig.setVersion("1.0.0");
    beanConfig.setSchemes(new String[] { "http" });
    beanConfig.setHost("clouditor.io");
    beanConfig.setBasePath(contextPath);
    beanConfig.setResourcePackage(this.getClass().getPackage().getName());
    beanConfig.setScan(true);

    // set the component service locator
    InjectionBridge.setComponentServiceLocator(this.component.getServiceLocator());

    this.register(ApiListingResource.class);
    this.register(SwaggerSerializers.class);
    this.register(CORSResponseFilter.class);
    this.register(ObjectMapperResolver.class);
    this.register(AuthenticationFilter.class);

    // registers the component itself as a service in the service locator
    // TODO: this might be obsolete now
    this.register(new ComponentFeature(component));
    this.register(InjectionBridge.class);

    this.packages(this.getClass().getPackage().toString());
  }

  /** Starts the API. */
  public void start() {
    LOGGER.info("Starting {}...", this.getClass().getSimpleName());

    this.httpServer = GrizzlyHttpServerFactory.createHttpServer(UriBuilder
        .fromUri("http://" + NetworkListener.DEFAULT_NETWORK_HOST + "/" + this.contextPath).port(this.port).build(),
        this);

    LOGGER.info("{} successfully started.", this.getClass().getSimpleName());

    // update the associated with the real port used, if port 0 was specified.
    if (this.port == 0) {
      component.setAPIPort(this.httpServer.getListener("grizzly").getPort());
    }

    this.httpServer.getServerConfiguration().addHttpHandler(new StaticHttpHandler("html"), "/");
  }

  /** Stops the API. */
  public void stop() {
    if (this.httpServer != null) {
      this.httpServer.shutdown();
    }
  }
}
