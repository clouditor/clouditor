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

package io.clouditor.credentials;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

import org.junit.jupiter.api.Test;

class AwsAccountTest {

  @Test
  void testResolveCredentials() {
    var account = new AwsAccount();
    account.setAccessKeyId("my-key-id");
    account.setSecretAccessKey("my-secret");
    account.setAutoDiscovered(false);

    var credentials = account.resolveCredentials();

    assertNotNull(credentials);

    assertEquals("my-key-id", credentials.accessKeyId());
    assertEquals("my-secret", credentials.secretAccessKey());

    account = new AwsAccount();
    account.setAutoDiscovered(true);

    // system properties are the first in the discovery chain so we override them
    System.setProperty("aws.accessKeyId", "my-discovered-key-id");
    System.setProperty("aws.secretAccessKey", "my-discovered-secret");

    assertTrue(account.isAutoDiscovered());

    credentials = account.resolveCredentials();

    assertNotNull(credentials);

    assertEquals("my-discovered-key-id", credentials.accessKeyId());
    assertEquals("my-discovered-secret", credentials.secretAccessKey());
  }
}
