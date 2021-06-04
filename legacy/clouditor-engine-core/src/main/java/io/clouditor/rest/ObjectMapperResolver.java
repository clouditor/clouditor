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

import com.fasterxml.jackson.annotation.JsonInclude.Include;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.InjectableValues;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import io.clouditor.credentials.CloudAccount;
import javax.ws.rs.ext.ContextResolver;
import org.reflections.Reflections;
import org.reflections.scanners.SubTypesScanner;
import org.reflections.util.ClasspathHelper;
import org.reflections.util.ConfigurationBuilder;

/** A {@link ContextResolver} used to resolve {@link ObjectMapper} for Jersey server and clients. */
public class ObjectMapperResolver implements ContextResolver<ObjectMapper> {

  public static class DatabaseOnly {}

  public static class ApiOnly {}

  private static final Reflections REFLECTIONS_SUBTYPE_SCANNER =
      new Reflections(
          new ConfigurationBuilder()
              .addUrls(ClasspathHelper.forPackage(CloudAccount.class.getPackage().getName()))
              .setScanners(new SubTypesScanner()));

  private static final ObjectMapper MAPPER = new ObjectMapper();

  public ObjectMapperResolver() {
    configureObjectMapper(MAPPER);
  }

  public static void configureObjectMapper(ObjectMapper mapper) {
    mapper.registerModule(new JavaTimeModule());
    mapper.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
    mapper.disable(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES);
    mapper.disable(SerializationFeature.FAIL_ON_EMPTY_BEANS);
    mapper.enable(SerializationFeature.INDENT_OUTPUT);
    mapper.setSerializationInclusion(Include.NON_NULL);
    mapper.setConfig(mapper.getSerializationConfig().withView(ApiOnly.class));

    // add all sub types of CloudAccount
    for (var type : REFLECTIONS_SUBTYPE_SCANNER.getSubTypesOf(CloudAccount.class)) {
      mapper.registerSubtypes(type);
    }

    // set injectable value to null
    var values = new InjectableValues.Std();
    values.addValue("hash", null);

    mapper.setInjectableValues(values);
  }

  @Override
  public ObjectMapper getContext(Class<?> type) {
    return MAPPER;
  }
}
