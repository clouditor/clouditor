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
import io.clouditor.oauth.OAuthResource;
import javax.validation.constraints.NotNull;
import javax.ws.rs.core.UriBuilder;
import org.glassfish.grizzly.http.server.HttpServer;
import org.glassfish.grizzly.http.server.NetworkListener;
import org.glassfish.grizzly.http.server.StaticHttpHandler;
import org.glassfish.grizzly.servlet.WebappContext;
import org.glassfish.jersey.grizzly2.httpserver.GrizzlyHttpServerFactory;
import org.glassfish.jersey.server.ResourceConfig;
import org.glassfish.jersey.servlet.ServletContainer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** Abstract base class for a REST API server. */
public abstract class AbstractAPI<C extends Component> extends ResourceConfig {

  /** The logger. */
  protected static final Logger LOGGER = LoggerFactory.getLogger(AbstractAPI.class);

  /** The grizzly HTTP server. */
  protected HttpServer httpServer;

  /** The API port. */


  private final int port;

  /** The context path */
  @NotNull private final String contextPath;

  /** The associated Clouditor component */
  @NotNull private final C component;

  AbstractAPI(C component, int port, String contextPath) {
    this.component = component;
    this.port = port;
    this.contextPath = contextPath;

    // set the component service locator
    InjectionBridge.setComponentServiceLocator(this.component.getServiceLocator());

    // bridges the service locator of the engine with the one of the REST web service
    this.register(InjectionBridge.class);

    this.register(CORSResponseFilter.class);
    this.register(ObjectMapperResolver.class);
    this.register(AuthenticationFilter.class);

    this.packages(this.getClass().getPackage().toString());
  }

  /**
   * Sanitizes input for several factors: a) to not break the log file pattern.
   *
   * @param input the untrusted input
   * @return the sanitized output
   */
  public static String sanitize(String input) {
    return input == null ? null : input.replaceAll("[\n|\r\t]", "_");
  }

  /** Starts the API. */
  public void start() {
    LOGGER.info("Starting {}...", this.getClass().getSimpleName());

    this.httpServer =
        GrizzlyHttpServerFactory.createHttpServer(
            UriBuilder.fromUri(
                    "http://" + NetworkListener.DEFAULT_NETWORK_HOST + "/" + this.contextPath)
                .port(this.port)
                .build(),
            this);

    LOGGER.info("{} successfully started.", this.getClass().getSimpleName());

    // update the associated with the real port used, if port 0 was specified.
    if (this.port == 0) {
      component.setAPIPort(this.httpServer.getListener("grizzly").getPort());
    }

    var config = new ResourceConfig();
    config.register(OAuthResource.class);
    config.register(InjectionBridge.class);

    var context = new WebappContext("WebappContext", "/oauth2");
    var registration = context.addServlet("OAuth2 Client", new ServletContainer(config));
    registration.addMapping("/*");
    context.deploy(httpServer);

    this.httpServer.getServerConfiguration().addHttpHandler(new StaticHttpHandler("html"), "/");
  }

  /** Stops the API. */
  public void stop() {
    if (this.httpServer != null) {
      this.httpServer.shutdown();
    }
  }
}
