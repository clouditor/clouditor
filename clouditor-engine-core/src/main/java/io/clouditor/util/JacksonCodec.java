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

package io.clouditor.util;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.mongodb.MongoException;
import io.clouditor.rest.ObjectMapperResolver.DatabaseOnly;
import java.io.IOException;
import org.bson.BsonReader;
import org.bson.BsonWriter;
import org.bson.RawBsonDocument;
import org.bson.codecs.Codec;
import org.bson.codecs.DecoderContext;
import org.bson.codecs.EncoderContext;
import org.bson.codecs.configuration.CodecRegistry;

public class JacksonCodec<T> implements Codec<T> {

  private final Class<T> clazz;
  private ObjectMapper mapper;
  private Codec<RawBsonDocument> codec;

  JacksonCodec(ObjectMapper mapper, Class<T> clazz, CodecRegistry registry) {
    this.mapper = mapper;
    this.clazz = clazz;
    this.codec = registry.get(RawBsonDocument.class);
  }

  @Override
  public T decode(BsonReader reader, DecoderContext decoderContext) {
    RawBsonDocument doc = codec.decode(reader, decoderContext);
    try {
      return mapper
          .readerWithView(DatabaseOnly.class)
          .forType(this.clazz)
          .readValue(doc.getByteBuffer().array());
    } catch (IOException e) {
      throw new MongoException(e.getMessage());
    }
  }

  @Override
  public void encode(BsonWriter writer, T value, EncoderContext encoderContext) {
    try {
      byte[] data = mapper.writerWithView(DatabaseOnly.class).writeValueAsBytes(value);
      codec.encode(writer, new RawBsonDocument(data), encoderContext);
    } catch (JsonProcessingException e) {
      throw new MongoException(e.getMessage());
    }
  }

  @Override
  public Class<T> getEncoderClass() {
    return this.clazz;
  }
}
