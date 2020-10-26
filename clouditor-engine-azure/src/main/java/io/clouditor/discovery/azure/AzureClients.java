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

import com.microsoft.azure.management.Azure;
import com.microsoft.azure.management.monitor.implementation.MonitorManager;
import com.microsoft.rest.RestClient.Builder;
import io.clouditor.credentials.AzureAccount;
import io.clouditor.util.PersistenceManager;
import java.io.IOException;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/** Helper class for Azure and Azure Monitor API */
public class AzureClients {

  private static final Logger LOGGER = LoggerFactory.getLogger(AzureClients.class);

  Azure azure;

  MonitorManager monitor;

  protected Builder builder;

  AzureClients() {}

  public void init() throws IOException {
    authenticate();
  }

  /**
   * Authenticates, either using an authentication file or the azure cli access tokens. See
   * https://github.com/Azure/azure-sdk-for-java/blob/master/AUTH.md for details.
   */
  private void authenticate() throws IOException {
    // fetch information about the account
    var account = PersistenceManager.getInstance().getById(AzureAccount.class, "Azure");

    if (account == null) {
      throw new IOException("Azure not configured");
    }

    var credentials = account.resolveCredentials();

    this.azure = Azure.authenticate(credentials).withDefaultSubscription();

    this.monitor =
        MonitorManager.authenticate(
            credentials, this.azure.getCurrentSubscription().subscriptionId());

    LOGGER.info("Successfully created Azure client");
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }

    return new EqualsBuilder().appendSuper(super.equals(o)).isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37).appendSuper(super.hashCode()).toHashCode();
  }

  public Azure azure() {
    return azure;
  }

  public MonitorManager monitor() {
    return monitor;
  }
}
