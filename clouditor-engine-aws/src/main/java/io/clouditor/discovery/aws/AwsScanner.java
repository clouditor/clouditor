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

package io.clouditor.discovery.aws;

import io.clouditor.credentials.AwsAccount;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.AssetProperties;
import io.clouditor.discovery.Scanner;
import io.clouditor.util.PersistenceManager;
import java.io.IOException;
import java.util.List;
import java.util.function.Function;
import java.util.function.Supplier;
import java.util.stream.Collectors;
import software.amazon.awssdk.awscore.AwsRequest;
import software.amazon.awssdk.awscore.AwsResponse;
import software.amazon.awssdk.awscore.client.builder.AwsClientBuilder;
import software.amazon.awssdk.awscore.exception.AwsServiceException;
import software.amazon.awssdk.core.SdkClient;
import software.amazon.awssdk.core.exception.SdkClientException;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.iam.IamClientBuilder;
import software.amazon.awssdk.utils.builder.ToCopyableBuilder;

public abstract class AwsScanner<
        C extends SdkClient, B extends AwsClientBuilder<B, C>, T extends ToCopyableBuilder>
    extends Scanner<C, T> {

  static final String ARN_SEPARATOR = ":";
  static final String RESOURCE_TYPE_SEPARATOR = "/";

  private Supplier<AwsClientBuilder<B, C>> builderSupplier;

  public AwsScanner(
      Supplier<AwsClientBuilder<B, C>> builderSupplier,
      Function<T, String> idGenerator,
      Function<T, String> nameGenerator) {
    super(null, idGenerator, nameGenerator);

    this.builderSupplier = builderSupplier;
    this.postProcessor = ToCopyableBuilder::toBuilder;
  }

  <O extends AwsResponse, R extends AwsRequest, S> void enrichSimple(
      Asset asset, String key, Function<R, O> function, Function<O, S> valueFunction, R builder) {
    try {
      var response = function.apply(builder);

      asset.setProperty(key, valueFunction.apply(response));
    } catch (AwsServiceException ex) {
      // ignore if error is 404, since this just means that the value is empty
      if (ex.awsErrorDetails().sdkHttpResponse().statusCode() != 404) {
        LOGGER.info(
            "Got exception from {} while retrieving info for {}: {}",
            this.getClass(),
            key,
            ex.getMessage());
      }
    }
  }

  <O extends AwsResponse, R extends AwsRequest, S extends ToCopyableBuilder> void enrich(
      Asset asset,
      String key,
      Function<R, O> listFunction,
      Function<O, S> valueFunction,
      R builder) {
    this.enrich(asset.getProperties(), key, listFunction, valueFunction, builder);
  }

  <O extends AwsResponse, R extends AwsRequest, S extends ToCopyableBuilder> void enrich(
      AssetProperties properties,
      String key,
      Function<R, O> listFunction,
      Function<O, S> valueFunction,
      R builder) {

    try {
      var response = listFunction.apply(builder);

      properties.put(
          key,
          MAPPER.convertValue(valueFunction.apply(response).toBuilder(), AssetProperties.class));
    } catch (AwsServiceException ex) {
      // ignore if error is 404, since this just means that the value is empty
      if (ex.statusCode() != 404) {
        LOGGER.info(
            "Got exception from {} while retrieving info for {}: {}",
            this.getClass(),
            key,
            ex.getMessage());
      }
    }
  }

  <O extends AwsResponse, R extends AwsRequest, S extends ToCopyableBuilder> void enrichList(
      Asset asset,
      String key,
      Function<R, O> listFunction,
      Function<O, List<S>> valueFunction,
      R builder) {
    try {
      var response = listFunction.apply(builder);

      var list = valueFunction.apply(response);

      asset.setProperty(
          key,
          list.stream()
              .map(x -> MAPPER.convertValue(x.toBuilder(), AssetProperties.class))
              .collect(Collectors.toList()));
    } catch (AwsServiceException ex) {
      // ignore if error is 404, since this just means that the value is empty
      if (ex.statusCode() != 404) {
        LOGGER.info(
            "Got exception from {} while retrieving info for {}: {}",
            this.getClass(),
            key,
            ex.getMessage());
      }
    }
  }

  @Override
  public void init() throws IOException {
    super.init();

    var builder = this.builderSupplier.get();

    var account = PersistenceManager.getInstance().getById(AwsAccount.class, "AWS");

    if (account == null) {
      throw SdkClientException.create("AWS account not configured");
    }

    // TODO: find a generic way to find which client is global
    if (!account.isAutoDiscovered() && !(builder instanceof IamClientBuilder)) {
      builder.region(Region.of(account.getRegion()));
    }

    builder.credentialsProvider(account);

    this.api = builder.build();
  }

  public C getApi() {
    return api;
  }

  @Override
  public void setApi(C api) {
    this.api = api;
  }
}
