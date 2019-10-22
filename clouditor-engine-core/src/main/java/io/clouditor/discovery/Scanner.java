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

package io.clouditor.discovery;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.databind.ObjectMapper;
import io.clouditor.assurance.Rule;
import io.clouditor.rest.ObjectMapperResolver;
import java.io.IOException;
import java.util.HashMap;
import java.util.List;
import java.util.function.Function;
import java.util.function.Supplier;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * This class represents a Scanner, which discovers a certain asset type. The properties of the
 * asset will then later queried by a {@link Rule}.
 *
 * @param <C> The API type.
 * @param <T> The discovered object type.
 */
public abstract class Scanner<C, T> {

  protected static final Logger LOGGER = LoggerFactory.getLogger(Scanner.class);

  protected static final ObjectMapper MAPPER = new ObjectMapper();

  static {
    ObjectMapperResolver.configureObjectMapper(MAPPER);
  }

  protected ScannerPostProcessor<?, T> postProcessor;

  private final Supplier<C> supplier;
  private final Function<T, String> idGenerator;
  private final Function<T, String> nameGenerator;

  protected C api;
  private boolean initialized = false;

  public Scanner(
      Supplier<C> supplier, Function<T, String> idGenerator, Function<T, String> nameGenerator) {
    this.supplier = supplier;
    this.idGenerator = idGenerator;
    this.nameGenerator = nameGenerator;
  }

  public void init() throws IOException {
    if (this.supplier != null) {
      this.api = supplier.get();
    }
  }

  protected abstract List<T> list() throws ScanException;

  @JsonIgnore
  public DiscoveryResult scan() {
    var result = new DiscoveryResult();

    try {

      // initialize the scanner, if not done already
      if (!this.initialized) {
        this.init();

        // we need to set that here and not within init() because a lot of scanners overwrite
        // init()
        this.initialized = true;
      }
      var assets = new HashMap<String, Asset>();

      LOGGER.info("Scanner {} is now scanning", this.getId());

      for (var object : list()) {
        var asset = transform(object);

        assets.put(asset.getId(), asset);
      }

      result.setDiscoveredAssets(assets);
    } catch (Exception e) {
      // it is important to catch all exceptions here, since that means a scan interval is failed
      // however, we might be able to recover in the next interval, so we do NOT want our future to
      // get cancelled
      LOGGER.info("Exception during scan", e);

      // mark the scan result as a failure
      result.setFailed(true);
      result.setError(e.getMessage() != null ? e.getMessage() : e.getClass().getSimpleName());
    }

    return result;
  }

  protected String getIdForObject(T object) {
    return this.idGenerator.apply(object);
  }

  protected String getNameForObject(T object) {
    return this.nameGenerator.apply(object);
  }

  public String getId() {
    return this.getClass().getName();
  }

  /**
   * Transforms an object to an {@link Asset}.
   *
   * @param object the object to transform
   * @return the asset.
   */
  protected Asset transform(T object) throws ScanException {
    Asset asset = new Asset();

    if (postProcessor != null) {
      asset.setProperties(MAPPER.convertValue(postProcessor.handle(object), AssetProperties.class));
    } else {
      asset.setProperties(MAPPER.convertValue(object, AssetProperties.class));
    }

    // TODO: not really the most efficient way to do this. better would be if the Scanner would have
    // access to the Scan object
    var info = this.getClass().getAnnotation(ScannerInfo.class);
    if (info != null) {
      asset.setType(info.assetType());
    } else {
      asset.setType(object.getClass().getSimpleName());
    }
    asset.setId(this.idGenerator.apply(object));
    asset.setName(this.nameGenerator.apply(object));

    return asset;
  }

  // TODO: just for mocking debug

  public void setApi(C api) {
    this.api = api;
  }

  public void setInitialized(boolean initialized) {
    this.initialized = initialized;
  }

  public boolean getInitialized() {
    return initialized;
  }

  public ScannerInfo getInfo() {
    return this.getClass().getAnnotation(ScannerInfo.class);
  }
}
