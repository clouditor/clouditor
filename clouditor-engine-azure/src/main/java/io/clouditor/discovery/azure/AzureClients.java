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

package io.clouditor.discovery.azure;

import com.microsoft.azure.management.Azure;
import com.microsoft.azure.management.monitor.implementation.MonitorManager;
import com.microsoft.rest.RestClient;
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

  // exposing restClient and baseUrl is necessary, since BYOK is not exposed in the Java REST API
  private RestClient restClient;

  MonitorManager monitor;

  protected Builder builder;

  private String baseUrl;

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
