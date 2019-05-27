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

  protected final Supplier<C> supplier;
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
}
