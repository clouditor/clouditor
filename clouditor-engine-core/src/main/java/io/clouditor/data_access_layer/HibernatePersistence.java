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

import java.io.Serializable;
import java.util.*;
import java.util.function.Consumer;
import java.util.function.Function;
import javax.persistence.criteria.CriteriaQuery;
import org.hibernate.NonUniqueObjectException;
import org.hibernate.Session;
import org.hibernate.Transaction;

/**
 * An implementation of the of the PersistenceManager.
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
public class HibernatePersistence implements PersistenceManager {

  @Override
  public <T> void saveOrUpdate(final T toSave) {
    Objects.requireNonNull(toSave);
    execConsumer(
        session -> {
          try {
            session.saveOrUpdate(toSave);
          } catch (NonUniqueObjectException exception) {
            // There are different objects (by reference)
            // with the same PrimaryKey  (@javax.persistence.Id)
            // then the session must be cleared to avoid collisions.
            session.clear();
            session.saveOrUpdate(toSave);
          }
        });
  }

  @Override
  public <T> Optional<T> get(final Class<T> resultType, final Serializable primaryKey) {
    Objects.requireNonNull(resultType);
    Objects.requireNonNull(primaryKey);
    return exec(session -> session.get(resultType, primaryKey));
  }

  @Override
  public <T> List<T> listAll(final Class<T> type) {
    Objects.requireNonNull(type);
    return exec(session -> {
          // Create a JPA CriteriaQuery wit the result type T.
          final CriteriaQuery<T> criteriaQuery = session.getCriteriaBuilder().createQuery(type);
          // Set the datasource to the Entity relation to T.
          criteriaQuery.from(type);
          return session.createQuery(criteriaQuery).getResultList();
        })
        .orElseThrow();
  }

  @Override
  public <T> void delete(final T toDelete) {
    Objects.requireNonNull(toDelete);
    execConsumer(session -> session.delete(toDelete));
  }

  @Override
  public <T> void delete(final Class<T> deleteType, final Serializable id) {
    Objects.requireNonNull(deleteType);
    Objects.requireNonNull(id);
    execConsumer(
        session ->
            session.delete(
                // Get a proxied instance of the object to delete.
                session.load(deleteType, id)));
  }

  @Override
  public <T> int count(final Class<T> countType) {
    Objects.requireNonNull(countType);
    return listAll(countType).size();
  }

  /**
   * Wrapper for the exec method, to be able to execute methods without a return value lieke save,
   * update or delete.
   *
   * @param consumer the operation to execute. The input value is an open Hibernate <code>Session
   *     </code>. See: https://docs.jboss.org/hibernate/orm/5.0/javadocs.
   * @throws NullPointerException if the <code>consumer</code> is null.
   */
  private void execConsumer(final Consumer<Session> consumer) {
    exec(
        session -> {
          consumer.accept(session);
          return new Object(); // The result is ignored.
        });
  }

  /**
   * Executes some operation on the database in a transaction with an open session.
   *
   * @param function the operation to execute. The input value is an open Hibernate <code>Session
   *     </code>. See: https://docs.jboss.org/hibernate/orm/5.0/javadocs.
   * @param <T> The result type of the operation.
   * @return the result of the operation.
   * @throws NullPointerException if the <code>function</code> is null.
   */
  private <T> Optional<T> exec(final Function<Session, T> function) {
    synchronized (HibernateUtils.LOCK) {
      final Session session = HibernateUtils.getSession();
      Transaction transaction = session.getTransaction();
      if (!transaction.isActive()) transaction = session.beginTransaction();
      // Execute the function.
      final Optional<T> result = Optional.ofNullable(function.apply(session));
      session.flush();
      transaction.commit();
      return result;
    }
  }
}
