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

package io.clouditor.data_access_layer;

import io.clouditor.assurance.*;
import io.clouditor.assurance.ccl.AssetType;
import io.clouditor.assurance.ccl.Condition;
import io.clouditor.assurance.ccl.FilteredAssetType;
import io.clouditor.auth.User;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.DiscoveryResult;
import io.clouditor.discovery.Scan;
import java.util.Objects;
import org.hibernate.SessionFactory;
import org.hibernate.cfg.Configuration;

/**
 * An Utility class to initialize and configure Hibernate.
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
public class HibernateUtils {

  private HibernateUtils() {}

  private static final Configuration CONFIGURATION =
      new Configuration()
          .configure()
          // Add all Entities.
          .addAnnotatedClass(Domain.class)
          .addAnnotatedClass(Certification.class)
          .addAnnotatedClass(Control.class)
          .addAnnotatedClass(Rule.class)
          .addAnnotatedClass(EvaluationResult.class)
          .addAnnotatedClass(AssetType.class)
          .addAnnotatedClass(FilteredAssetType.class)
          .addAnnotatedClass(Asset.class)
          .addAnnotatedClass(Condition.class)
          .addAnnotatedClass(Scan.class)
          .addAnnotatedClass(DiscoveryResult.class)
          .addAnnotatedClass(User.class)
          .setProperty("hibernate.enable_lazy_load_no_trans", "true")
          .setProperty("hibernate.dialect", "org.hibernate.dialect.PostgreSQL94Dialect")
          // Enable an automatic generation of the data model.
          .setProperty("hibernate.hbm2ddl.auto", "create-drop");

  private static SessionFactory sessionFactory;

  /**
   * Getter for the session factory
   *
   * @return the session factory. Not null.
   * @throws IllegalStateException if the session factory was not initialized by one of the init
   *     functions.
   */
  public static SessionFactory getSessionFactory() {
    if (sessionFactory == null)
      throw new IllegalStateException("The Database Connection is not initialized.");
    return sessionFactory;
  }

  /**
   * Connects Hibernate to the persistent PostgreSQL database at the address:
   * "jdbc:postgresql://{host}:{port}/{dbName}".
   *
   * @param host the host of the database.
   * @param port the port of the database should be in range [1024; 49151].
   * @param dbName the database name.
   * @param userName the user name.
   * @param password the password of the user
   * @throws IllegalArgumentException if the <code>port</code> is not in range [1024; 49151].
   * @throws NullPointerException if <code>host</code>, <code>dbName</code>, <code>userName</code>
   *     or <code>password</code> is null.
   */
  public static void init(
      final String host,
      final int port,
      final String dbName,
      final String userName,
      final String password) {
    Objects.requireNonNull(host);
    Objects.requireNonNull(dbName);
    Objects.requireNonNull(userName);
    Objects.requireNonNull(password);
    if (port < 1024 || port > 49151)
      throw new IllegalArgumentException(
          "The given port: " + port + ", was not in range [1024; 49151].");
    CONFIGURATION
        .setProperty("hibernate.connection.driver_class", "org.postgresql.Driver")
        .setProperty(
            "hibernate.connection.url", "jdbc:postgresql://" + host + ":" + port + "/" + dbName);
    setUserNameAndPassword(userName, password);
    buildSessionFactory();
  }

  /**
   * Connects to a automatically created H2 in memory database at the address: "jdbc:h2:~/{dbName}".
   *
   * @param dbName the database name.
   * @param userName the user name.
   * @param password the password of the user
   * @throws IllegalArgumentException if the <code>port</code> is not in range [1024; 49151].
   * @throws NullPointerException if <code>dbName</code>, <code>userName</code> or <code>password
   *     </code> is null.
   */
  public static void init(final String dbName, final String userName, final String password) {
    Objects.requireNonNull(dbName);
    Objects.requireNonNull(userName);
    Objects.requireNonNull(password);
    CONFIGURATION
        .setProperty("hibernate.connection.driver_class", "org.h2.Driver")
        .setProperty("hibernate.connection.url", "jdbc:h2:~/" + dbName);
    setUserNameAndPassword(userName, password);
    buildSessionFactory();
  }

  /** Closes the session factory. */
  public static void close() {
    if (sessionFactory != null) sessionFactory.close();
  }

  private static void setUserNameAndPassword(final String userName, final String password) {
    CONFIGURATION
        .setProperty("hibernate.connection.username", userName)
        .setProperty("hibernate.connection.password", password);
  }

  private static void buildSessionFactory() {
    sessionFactory = CONFIGURATION.buildSessionFactory();
  }
}
