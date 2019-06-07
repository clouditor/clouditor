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

package io.clouditor;

import io.clouditor.auth.LoginResponse;
import io.clouditor.auth.User;
import javax.ws.rs.client.Entity;
import javax.ws.rs.client.WebTarget;
import org.glassfish.hk2.api.ServiceLocator;
import org.glassfish.hk2.utilities.ServiceLocatorUtilities;
import org.jvnet.hk2.annotations.Contract;
import org.kohsuke.args4j.CmdLineException;
import org.kohsuke.args4j.CmdLineParser;
import org.kohsuke.args4j.Option;
import org.kohsuke.args4j.spi.ExplicitBooleanOptionHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.bridge.SLF4JBridgeHandler;

/**
 * The base class for different Clouditor components, such as the Engine or Explorer.
 *
 * @author Banse, Christian
 */
@Contract
public abstract class Component {

  public static final String CONFIG_API_PORT = "api.port";
  public static final String CONFIG_API_USERNAME = "api.username";
  public static final String CONFIG_API_PW = "api.password";
  public static final String CONFIG_API_ALLOWED_ORIGIN = "api.allowedOrigin";

  protected static final Logger LOGGER = LoggerFactory.getLogger(Component.class);

  private static final boolean DEFAULT_IS_DAEMON = true;
  private static final String DEFAULT_API_USERNAME = "clouditor";
  private static final String DEFAULT_API_PW = "clouditor";
  private static final String DEFAULT_API_SECRET = "changeme";
  private static final String DEFAULT_API_ALLOWED_ORIGIN = "*";

  /** Specifies whether the component is running as a daemon. */
  @Option(
      name = "-D",
      usage = "specifies whether to run as daemon",
      handler = ExplicitBooleanOptionHandler.class)
  boolean isDaemon = DEFAULT_IS_DAEMON;

  /** Specifies the API username. */
  @Option(name = "--api-default-user", usage = "specifies the API username")
  String defaultApiUsername = DEFAULT_API_USERNAME;

  /** Specifies the API password. */
  @Option(name = "--api-default-password", usage = "specifies the API password")
  String defaultApiPw = DEFAULT_API_PW;

  /** Specifies the secret used by API tokens. */
  @Option(name = "--api-secret", usage = "specifies the secret used by API tokens")
  String apiSecret = DEFAULT_API_SECRET;

  /** Specifies the allowed origin for API requests. */
  @Option(name = "--api-allowed-origin", usage = "specifies the allowed origin for API requests")
  String apiAllowedOrigin = DEFAULT_API_ALLOWED_ORIGIN;

  // TODO: somehow use the one that is already there in the api
  /** The service locator from HK2 */
  private ServiceLocator locator;

  public Component() {
    this.locator = ServiceLocatorUtilities.createAndPopulateServiceLocator();
    ServiceLocatorUtilities.addOneConstant(this.locator, this);

    // Optionally remove existing handlers attached to j.u.l root logger
    SLF4JBridgeHandler.removeHandlersForRootLogger();

    // add SLF4JBridgeHandler to java.util.logging's root logger
    SLF4JBridgeHandler.install();
  }

  public String getDefaultApiUsername() {
    return this.defaultApiUsername;
  }

  public String getDefaultApiPassword() {
    return this.defaultApiPw;
  }

  public String getAPIAllowedOrigin() {
    return this.apiAllowedOrigin;
  }

  public abstract int getAPIPort();

  public abstract void setAPIPort(int port);

  public abstract void init();

  public boolean isDaemon() {
    return this.isDaemon;
  }

  public void setDaemon(boolean daemon) {
    this.isDaemon = daemon;
  }

  public boolean parseArgs(String[] args) {
    var parser = new CmdLineParser(this);
    try {
      parser.parseArgument(args);
      return true;
    } catch (CmdLineException e) {
      LOGGER.error("Could not parse command line arguments: {}", e.getLocalizedMessage());
      return false;
    }
  }

  public abstract void shutdown();

  public abstract void start(String[] args) throws InterruptedException;

  public abstract void startAPI();

  public abstract void stopAPI();

  public String authenticateAPI(WebTarget target, String username, String password) {
    var response =
        target
            .path("authenticate")
            .request()
            .post(Entity.json(new User(username, password)), LoginResponse.class);

    if (response != null) {
      return response.getToken();
    }

    return null;
  }

  public <T> T getService(Class<T> clazz) {
    return this.locator.getService(clazz);
  }

  public ServiceLocator getServiceLocator() {
    return this.locator;
  }

  public String getApiSecret() {
    return this.apiSecret;
  }
}
