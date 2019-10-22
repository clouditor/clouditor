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

package io.clouditor.util;

import static com.mongodb.client.model.Filters.eq;
import static io.clouditor.rest.AbstractAPI.sanitize;
import static org.bson.codecs.configuration.CodecRegistries.fromProviders;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.module.SimpleModule;
import com.mongodb.MongoClient;
import com.mongodb.MongoException;
import com.mongodb.MongoTimeoutException;
import com.mongodb.ServerAddress;
import com.mongodb.client.FindIterable;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoDatabase;
import com.mongodb.client.model.UpdateOptions;
import de.undercouch.bson4jackson.BsonFactory;
import io.clouditor.rest.ObjectMapperResolver;
import java.time.Instant;
import org.bson.Document;
import org.bson.codecs.configuration.CodecRegistries;
import org.bson.codecs.configuration.CodecRegistry;
import org.bson.conversions.Bson;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class PersistenceManager {

  public static final String FIELD_ID = "_id";

  protected static final Logger LOGGER = LoggerFactory.getLogger(PersistenceManager.class);

  private static PersistenceManager instance;

  private String host;

  private int port;

  private MongoClient mongo;
  private MongoDatabase mongoDatabase;
  private CodecRegistry codecRegistry;

  private boolean initialized = false;

  private PersistenceManager() {
    var factory = new BsonFactory();

    var module = new SimpleModule();
    // the default Jackson Java 8 time (de)serializer are not compatible with MongoDB
    module.addSerializer(Instant.class, new BsonInstantSerializer());
    module.addDeserializer(Instant.class, new BsonInstantDeserializer());

    var mapper = new ObjectMapper(factory);
    ObjectMapperResolver.configureObjectMapper(mapper);

    mapper.registerModule(module);

    this.codecRegistry =
        CodecRegistries.fromRegistries(
            MongoClient.getDefaultCodecRegistry(), fromProviders(new JacksonCodecProvider(mapper)));
  }

  public static synchronized PersistenceManager getInstance() {
    if (instance == null) {
      instance = new PersistenceManager();
    }
    return instance;
  }

  public <T> FindIterable<T> find(Class<T> clazz) {
    var coll = this.getCollection(clazz);

    return coll.find();
  }

  public <T> FindIterable<T> find(Class<T> clazz, Bson filter) {
    var coll = this.getCollection(clazz);

    return coll.find(filter);
  }

  public <T> long count(Class<T> clazz) {
    var coll = this.getCollection(clazz);

    return coll.count();
  }

  public <T> long count(Class<T> clazz, Bson filter) {
    var coll = this.getCollection(clazz);

    return coll.count(filter);
  }

  public <T> T getById(Class<T> clazz, String id) {
    return this.find(clazz, eq(FIELD_ID, id)).limit(1).first();
  }

  public String getHost() {
    return this.host;
  }

  public int getPort() {
    return this.port;
  }

  public void init(String dbName, String dbHost, int dbPort) {
    this.host = dbHost;
    this.port = dbPort;

    this.init(dbName, new MongoClient(new ServerAddress(dbHost, dbPort)));
  }

  public void destroy() {
    if (this.mongo != null) {
      this.mongo.close();
    }
  }

  public void persist(PersistentObject object) {
    try {
      MongoCollection<PersistentObject> coll = this.getCollection(object);

      coll.replaceOne(
          new Document(FIELD_ID, object.getId()), object, new UpdateOptions().upsert(true));
    } catch (MongoException e) {
      LOGGER.error("Error while saving into database: {}", e.getMessage());
    }
  }

  private <T> MongoCollection<T> getCollection(Class<T> clazz) {
    return this.mongoDatabase
        .getCollection(getCollectionName(clazz), clazz)
        .withCodecRegistry(this.codecRegistry);
  }

  private MongoCollection<PersistentObject> getCollection(PersistentObject object) {
    return this.mongoDatabase
        .getCollection(getCollectionName(object.getClass()), PersistentObject.class)
        .withCodecRegistry(this.codecRegistry);
  }

  private String getCollectionName(Class clazz) {
    String collectionName;
    var collection = (Collection) clazz.getAnnotation(Collection.class);
    if (collection != null) {
      collectionName = collection.value();
    } else {
      collectionName = clazz.getSimpleName().toLowerCase() + "s";
    }
    return collectionName;
  }

  public <T> void delete(Class<T> clazz, String id) {
    var coll = getCollection(clazz);

    var result = coll.deleteOne(new Document(FIELD_ID, id));

    if (result.wasAcknowledged()) {
      LOGGER.info("Deleted id {} (Type: {}) from database", sanitize(id), clazz.getSimpleName());
    } else {
      LOGGER.error(
          "Could not delete id {} (Type: {}) from database", sanitize(id), clazz.getSimpleName());
    }
  }

  public boolean isInitialized() {
    return initialized;
  }

  public void close() {
    this.initialized = false;

    this.mongo.close();
  }

  public void init(String dbName, MongoClient client) {
    this.mongo = client;

    this.mongoDatabase = this.mongo.getDatabase(dbName);

    // wait for the DB
    try {
      var address = mongo.getAddress();
      LOGGER.info("Connected to MongoDB @ {}/{}.", address, dbName);

      this.initialized = true;
    } catch (MongoTimeoutException ex) {
      LOGGER.error("Fatal error. Could not connect to the MongoDB: {}", ex);
    }
  }
}
