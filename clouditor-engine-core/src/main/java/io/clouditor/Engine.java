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

import de.bwaldvogel.mongo.MongoServer;
import de.bwaldvogel.mongo.backend.memory.MemoryBackend;
import io.clouditor.assurance.CertificationService;
import io.clouditor.assurance.RuleService;
import io.clouditor.discovery.DiscoveryService;
import io.clouditor.rest.EngineAPI;
import io.clouditor.util.FileSystemManager;
import io.clouditor.util.PersistenceManager;
import org.jvnet.hk2.annotations.Service;
import org.kohsuke.args4j.Option;

/**
 * The main Clouditor Engine class.
 *
 * @author Philipp Stephanow
 */
@Service
public class Engine extends Component {

  private static final String DEFAULT_DB_HOST = "localhost";

  private static final String DEFAULT_DB_NAME = "clouditor";

  private static final int DEFAULT_DB_PORT = 27017;

  private static final short DEFAULT_API_PORT = 9999;

  private static final boolean DEFAULT_DB_IN_MEMORY = false;

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

  public String getDbHost() {
    return dbHost;
  }

  public void setDbHost(String dbHost) {
    this.dbHost = dbHost;
  }

  public String getDbName() {
    return dbName;
  }

  public void setDbName(String dbName) {
    this.dbName = dbName;
  }

  public int getDbPort() {
    return dbPort;
  }

  public void setDbPort(int dbPort) {
    this.dbPort = dbPort;
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
    // load the certificate importers
    this.getService(CertificationService.class).loadImporters();

    loadRules();
    loadScans();
    loadCertifications();
    loadSubscribers();

    this.getService(CertificationService.class).updateCertification();
  }

  private void loadScans() {
    var service = this.getService(DiscoveryService.class);
    service.load();
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
    if (this.dbInMemory) {
      var server = new MongoServer(new MemoryBackend());

      var address = server.bind();

      LOGGER.info("Starting database purely in-memory...");

      this.dbHost = address.getHostName();
      this.dbPort = address.getPort();

      Runtime.getRuntime().addShutdownHook(new Thread(server::shutdown));
    }

    PersistenceManager.getInstance().init(this.dbName, this.dbHost, this.dbPort);
  }

  /** Shuts down the Clouditor Engine */
  @Override
  public void shutdown() {
    LOGGER.info("Shutting down...");

    // clean up file systems
    FileSystemManager.getInstance().cleanup();

    // stop the API
    this.stopAPI();

    // destroy the persistence manager
    PersistenceManager.getInstance().destroy();
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
}
