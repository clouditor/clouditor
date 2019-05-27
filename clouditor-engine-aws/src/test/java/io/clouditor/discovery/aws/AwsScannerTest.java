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

package io.clouditor.discovery.aws;

import static org.mockito.Mockito.mock;

import io.clouditor.Engine;
import io.clouditor.discovery.DiscoveryResult;
import java.util.function.Consumer;
import java.util.function.Supplier;
import software.amazon.awssdk.awscore.client.builder.AwsClientBuilder;
import software.amazon.awssdk.core.SdkClient;
import software.amazon.awssdk.utils.builder.ToCopyableBuilder;

public class AwsScannerTest {

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

    assets = scanner.scan();
  }
}
