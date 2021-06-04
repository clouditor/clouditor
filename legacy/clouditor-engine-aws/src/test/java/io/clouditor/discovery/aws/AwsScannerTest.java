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

import static org.mockito.Mockito.mock;

import io.clouditor.Engine;
import io.clouditor.discovery.DiscoveryResult;
import java.util.function.Consumer;
import java.util.function.Supplier;
import software.amazon.awssdk.awscore.client.builder.AwsClientBuilder;
import software.amazon.awssdk.core.SdkClient;
import software.amazon.awssdk.utils.builder.ToCopyableBuilder;

public abstract class AwsScannerTest {

  Engine engine = new Engine();

  static DiscoveryResult assets;

  static <B extends SdkClient, C extends AwsClientBuilder<C, B>, T extends ToCopyableBuilder>
      void discoverAssets(
          Class<B> apiClass, Supplier<AwsScanner<B, C, T>> supplier, Consumer<B> configurator) {
    var scanner = supplier.get();

    System.setProperty("aws.region", "mock");
    System.setProperty("aws.accessKeyId", "mock");
    System.setProperty("aws.secretAccessKey", "mock");

    // mock the client
    var api = mock(apiClass);
    // scanner.init(); don't init
    scanner.setInitialized(true);

    configurator.accept(api);

    // force the api
    scanner.setApi(api);

    assets = scanner.scan(null);
  }
}
