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

import java.util.Objects;
import java.util.Set;
import javax.persistence.Entity;
import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.hibernate.cfg.Configuration;
import org.reflections.Reflections;

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
          .setProperty("hibernate.enable_lazy_load_no_trans", "true")
          .setProperty("hibernate.dialect", "org.hibernate.dialect.PostgreSQL94Dialect")
          // Enable an automatic generation of the data model.
          .setProperty("hibernate.hbm2ddl.auto", "create");

  private static SessionFactory sessionFactory;
  private static Session session;

  static final Object LOCK = new Object();

  /**
   * Getter for the current session. If there is no current session, it opens a new session.
   *
   * @return the current session
   * @throws IllegalStateException if the database was not initialized.
   */
  public static Session getSession() {
    if (session == null) setSession(getSessionFactory().openSession());
    return session;
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
    finishSessionFactory(userName, password);
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
  public static void initInMemoryH2(
      final String dbName, final String userName, final String password) {
    Objects.requireNonNull(dbName);
    Objects.requireNonNull(userName);
    Objects.requireNonNull(password);
    CONFIGURATION
        .setProperty("hibernate.connection.driver_class", "org.h2.Driver")
        .setProperty("hibernate.connection.url", "jdbc:h2:~/" + dbName);
    finishSessionFactory(userName, password);
  }

  /** Closes the session factory. */
  public static synchronized void close() {
    getSession().close();
    setSession(null);
  }

  private static void finishSessionFactory(final String userName, final String password) {
    setUserNameAndPassword(userName, password);
    setAnnotatedClasses();
    buildSessionFactory();
  }

  private static void setUserNameAndPassword(final String userName, final String password) {
    CONFIGURATION
        .setProperty("hibernate.connection.username", userName)
        .setProperty("hibernate.connection.password", password);
  }

  private static void setAnnotatedClasses() {
    final Reflections reflections = new Reflections("io.clouditor");
    final Set<Class<?>> entities = reflections.getTypesAnnotatedWith(Entity.class);
    for (final Class<?> clazz : entities) {
      CONFIGURATION.addAnnotatedClass(clazz);
    }
  }

  private static void buildSessionFactory() {
    setSessionFactory(CONFIGURATION.buildSessionFactory());
  }

  private static void setSession(final Session sessionToSet) {
    session = sessionToSet;
  }

  private static SessionFactory getSessionFactory() {
    if (sessionFactory == null)
      throw new IllegalStateException("The Database Connection is not initialized.");
    return sessionFactory;
  }

  private static void setSessionFactory(final SessionFactory sessionFactoryToSet) {
    sessionFactory = sessionFactoryToSet;
  }
}
