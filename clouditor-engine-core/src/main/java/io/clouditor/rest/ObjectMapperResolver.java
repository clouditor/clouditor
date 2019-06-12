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
