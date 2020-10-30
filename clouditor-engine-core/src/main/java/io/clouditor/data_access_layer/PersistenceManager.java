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
import java.util.List;
import java.util.Optional;

/**
 * The persistence manager provides functions to interact with the underlying database.
 *
 * @author Andreas Hager, andreas.hager@aisec.fraunhofer.de
 */
public interface PersistenceManager {

  /**
   * Saves or updates the data from the in the database.
   *
   * @param toSave the object containing the data to store or to update in the database.
   * @param <T> the generic type of the object to store.
   * @throws NullPointerException if the <code>toSave</code> is null.
   */
  <T> void saveOrUpdate(final T toSave);

  /**
   * Loads a stored object from the database.
   *
   * @param resultType the result type of the object to load.
   * @param primaryKey the primaryKey to identify the dataset.
   * @param <T> the generic type of the object to load.
   * @return a new instance of the object of type T,
   *          filled with the data from the database.
   * @throws NullPointerException if <code>resultType</code> or <code>primaryKey</code> is null.
   */
  <T> Optional<T> get(final Class<T> resultType, final Serializable primaryKey);

  /**
   * Loads all stored objects of the type T.
   *
   * @param resultType the type of the stored objects.
   * @param <T> the generic type of the object to load.
   * @return a list containing new instances of the type t,
   *          filled with the data from the database.
   * @throws NullPointerException if the <code>resultType</code> is null.
   */
  <T> List<T> listAll(final Class<T> resultType);

  /**
   * Deletes a stored object form the database.
   *
   * @param toDelete the object to delete.
   * @param <T> the generic type of the object to delete.
   * @throws NullPointerException if <code>toDelete</code> is null.
   */
  <T> void delete(final T toDelete);

  /**
   * Deletes a stored object form the database.
   *
   * @param deleteType the type of the object to delete.
   * @param id the id of the dataset.
   * @param <T> the generic type of the object to delete.
   * @throws NullPointerException if <code>deleteType</code> or <code>id</code> is null.
   */
  <T> void delete(
      final Class<T> deleteType, final Serializable id);

  /**
   * Counts the datasets of the type T that are stored in the database.
   *
   * @param countType the type of the dataset to count.
   * @param <T> the generic type of the object to count.
   * @return the count of the datasets of the type T
   * @throws NullPointerException if <code>countType</code> is null.
   */
  <T> int count(final Class<T> countType);
}
