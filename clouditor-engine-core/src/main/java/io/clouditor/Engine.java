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

package io.clouditor;

import io.clouditor.assurance.CertificationService;
import io.clouditor.assurance.RuleService;
import io.clouditor.auth.AuthenticationService;
import io.clouditor.data_access_layer.HibernateUtils;
import io.clouditor.discovery.DiscoveryService;
import io.clouditor.rest.EngineAPI;
import io.clouditor.util.FileSystemManager;
import org.jvnet.hk2.annotations.Service;
import org.kohsuke.args4j.Option;

/**
 * The main Clouditor Engine class.
 *
 * @author Philipp Stephanow
 */
@Service
public class Engine extends Component {

  private static final String DEFAULT_DB_USER_NAME = "postgres";

  private static final String DEFAULT_DB_PASSWORD = "postgres";

  private static final String DEFAULT_DB_HOST = "localhost";

  private static final String DEFAULT_DB_NAME = "clouditor";

  private static final int DEFAULT_DB_PORT = 5432;

  private static final short DEFAULT_API_PORT = 9999;

  private static final boolean DEFAULT_DB_IN_MEMORY = true;

  @Option(name = "--db-user-name", usage = "provides user name of database")
  private String dbUserName = DEFAULT_DB_USER_NAME;

  @Option(name = "--db-password", usage = "provides password of database")
  private String dbPassword = DEFAULT_DB_PASSWORD;

  @Option(name = "--db-host", usage = "provides address of database")
  private String dbHost = DEFAULT_DB_HOST;

  @Option(name = "--db-name", usage = "provides name of database")
  private String dbName = DEFAULT_DB_NAME;

  @Option(name = "--db-port", usage = "provides port for database")
  private int dbPort = DEFAULT_DB_PORT;

  @Option(
      name = "--db-in-memory",
      usage = "uses an in-memory database which is not persisted at all")
  private boolean dbInMemory = DEFAULT_DB_IN_MEMORY;

  @Option(
      name = "-p",
      aliases = {"--port"},
      usage = "specifies port for REST API")
  private int apiPort = DEFAULT_API_PORT;

  @Option(name = "--server-base-url")
  private String baseUrl = "http://localhost:" + DEFAULT_API_PORT;

  @Option(name = "--oauth-client-id")
  private String oAuthClientId;

  @Option(name = "--oauth-client-secret")
  private String oAuthClientSecret;

  @Option(name = "--oauth-token-url")
  private String oAuthTokenUrl;

  @Option(name = "--oauth-auth-url")
  private String oAuthAuthUrl;

  @Option(name = "--oauth-jwt-secret")
  private String oAuthJwtSecret;

  @Option(name = "--oauth-jwt-issuer")
  private String oAuthJwtIssuer;

  /** The web api. */
  private EngineAPI api;

  public Engine() {
    // Nothing to do
  }

  public int getAPIPort() {
    return this.apiPort;
  }

  public void setAPIPort(int port) {
    this.apiPort = port;
  }

  public void setDbInMemory(boolean dbInMemory) {
    this.dbInMemory = dbInMemory;
  }

  public String getOAuthClientId() {
    return oAuthClientId;
  }

  public String getOAuthClientSecret() {
    return oAuthClientSecret;
  }

  public String getOAuthTokenUrl() {
    return oAuthTokenUrl;
  }

  public String getOAuthJwtSecret() {
    return oAuthJwtSecret;
  }

  public String getOAuthAuthUrl() {
    return this.oAuthAuthUrl;
  }

  /**
   * Returns the Clouditor Engine REST API.
   *
   * @return the Clouditor Engine REST API.
   */
  public EngineAPI getAPI() {
    return this.api;
  }

  /** Initializes the Clouditor Engine, i.e. loads all configuration files. */
  @Override
  public void init() {
    // init user service
    this.getService(AuthenticationService.class).init();

    // load the certificate importers
    this.getService(CertificationService.class).loadImporters();

    loadRules();
    initDiscoveryService();
    loadCertifications();
    loadSubscribers();

    this.getService(CertificationService.class).updateCertification();
  }

  private void initDiscoveryService() {
    var service = this.getService(DiscoveryService.class);
    service.init();
  }

  public void loadRules() {
    var ruleService = this.getService(RuleService.class);
    ruleService.loadAll();
  }

  public void loadSubscribers() {
    // subscribe rule service to scan service as an asset subscriber
    this.getService(DiscoveryService.class).subscribe(this.getService(RuleService.class));

    this.getService(CertificationService.class).loadSubscribers();
  }

  private void loadCertifications() {
    this.getService(CertificationService.class).loadCertifications();
  }

  public void initDB() {
    if (this.dbInMemory) HibernateUtils.init(this.dbName, this.dbUserName, this.dbPassword);
    else {
      HibernateUtils.init(this.dbHost, this.dbPort, this.dbName, this.dbUserName, this.dbPassword);
      Runtime.getRuntime().addShutdownHook(new Thread(HibernateUtils::close));
    }
  }

  /** Shuts down the Clouditor Engine */
  @Override
  public void shutdown() {
    LOGGER.info("Shutting down...");

    // clean up file systems
    FileSystemManager.getInstance().cleanup();

    // stop the API
    this.stopAPI();
  }

  /**
   * Starts the Clouditor Engine. Assumes, that all configuration values are set correctly. This
   * call will block until all tasks are done.
   */
  @Override
  public void start(String[] args) {
    // parse command line args
    this.parseArgs(args);

    // init db
    this.initDB();

    // initialize every else
    this.init();

    // start the DiscoveryService
    this.getService(DiscoveryService.class).start();

    // start the REST API
    this.startAPI();
  }

  /** Starts the Clouditor Engine REST API. */
  @Override
  public void startAPI() {
    // start the web api
    this.api = new EngineAPI(this);
    this.api.start();
  }

  /** Stops the Clouditor Engine REST API. */
  @Override
  public void stopAPI() {
    if (this.api != null) {
      this.api.stop();
    }
  }

  public String getBaseUrl() {
    return baseUrl;
  }

  public String getoAuthJwtIssuer() {
    return oAuthJwtIssuer;
  }
}
