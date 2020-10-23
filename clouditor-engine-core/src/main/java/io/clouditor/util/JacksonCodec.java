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
  private final ObjectMapper mapper;
  private final Codec<RawBsonDocument> codec;

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
