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

package io.clouditor.discovery.azure;

import com.microsoft.azure.management.resources.fluentcore.arm.models.HasId;
import com.microsoft.azure.management.resources.fluentcore.arm.models.HasName;
import com.microsoft.azure.management.resources.fluentcore.model.HasInner;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.AssetProperties;
import io.clouditor.discovery.MixInIgnore;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.Scanner;
import io.clouditor.discovery.ScannerPostProcessor;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Collection;
import java.util.Map.Entry;
import java.util.function.Function;
import java.util.stream.Collectors;

public abstract class AzureScanner<T extends HasInner> extends Scanner<AzureClients, T> {

  public static final String PROPERTIES_KEY_PREFIX = "properties.";
  protected String resourceGroup;

  AzureScanner(Function<T, String> idGenerator, Function<T, String> nameGenerator) {
    super(null, idGenerator, nameGenerator);

    this.postProcessor = HasInner::inner;
  }

  @Override
  public void init() throws IOException {
    super.init();

    MAPPER.addMixIn(com.microsoft.azure.SubResource.class, MixInIgnore.class);

    this.api = new AzureClients();

    this.api.init();
  }

  <O, I> void enrich(
      Asset asset,
      String key,
      O outer,
      Function<O, I> supplier,
      Function<I, String> idGenerator,
      Function<I, String> nameGenerator) {

    var object = supplier.apply(outer);

    if (object == null) {
      return;
    }

    AssetProperties tmp;
    if (object instanceof HasInner) {
      tmp = MAPPER.convertValue(((HasInner) object).inner(), AssetProperties.class);
    } else {
      tmp = MAPPER.convertValue(object, AssetProperties.class);
    }

    // TODO: find a nicer way to do that
    var properties = new AssetProperties();
    properties.putAll(
        tmp.entrySet().stream()
            .collect(
                Collectors.toMap(
                    e -> e.getKey().replace(PROPERTIES_KEY_PREFIX, ""), Entry::getValue)));

    properties.put("id", idGenerator.apply(object));
    properties.put("name", nameGenerator.apply(object));

    asset.setProperty(key, properties);
  }

  <O, I extends HasInner & HasName> void enrich(
      Asset asset, String key, O outer, Function<O, I> supplier, Function<I, String> idGenerator) {
    enrich(asset, key, outer, supplier, idGenerator, HasName::name);
  }

  <O, I extends HasInner & HasId & HasName> void enrich(
      Asset asset, String key, O outer, Function<O, I> supplier) {
    enrich(asset, key, outer, supplier, HasId::id, HasName::name);
  }

  <O, T extends HasInner> void enrichList(
      Asset asset,
      String key,
      O outer,
      Function<O, Collection<T>> listSupplier,
      Function<T, String> idGenerator,
      Function<T, String> nameGenerator) {
    this.enrichList(asset, key, outer, listSupplier, HasInner::inner, idGenerator, nameGenerator);
  }

  <O, S> void enrichList(
      Asset asset,
      String key,
      O outer,
      Function<O, Collection<S>> listSupplier,
      ScannerPostProcessor<?, S> postProcessor,
      Function<S, String> idGenerator,
      Function<S, String> nameGenerator) {
    var list = new ArrayList<>();

    for (S object : listSupplier.apply(outer)) {
      AssetProperties tmp;

      if (postProcessor != null) {
        tmp = MAPPER.convertValue(postProcessor.handle(object), AssetProperties.class);
      } else {
        tmp = MAPPER.convertValue(object, AssetProperties.class);
      }

      // TODO: find a nicer way to do that
      var properties = new AssetProperties();
      properties.putAll(
          tmp.entrySet().stream()
              .collect(
                  Collectors.toMap(e -> e.getKey().replace("properties.", ""), Entry::getValue)));

      properties.put("id", idGenerator.apply(object));
      properties.put("name", nameGenerator.apply(object));

      list.add(properties);
    }

    asset.setProperty(key, list);
  }

  @Override
  protected Asset transform(T object) throws ScanException {
    var asset = super.transform(object);

    // TODO: find a nicer way to do that
    var properties = new AssetProperties();
    properties.putAll(
        asset.getProperties().entrySet().stream()
            .collect(
                Collectors.toMap(e -> e.getKey().replace("properties.", ""), Entry::getValue)));

    asset.setProperties(properties);

    return asset;
  }
}
