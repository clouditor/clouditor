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

package io.clouditor.credentials;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;
import com.microsoft.azure.AzureEnvironment;
import com.microsoft.azure.credentials.ApplicationTokenCredentials;
import com.microsoft.azure.credentials.AzureCliCredentials;
import com.microsoft.azure.credentials.AzureTokenCredentials;
import com.microsoft.azure.management.Azure;
import java.io.File;
import java.io.IOException;

@JsonTypeName("Azure")
public class AzureAccount extends CloudAccount<AzureTokenCredentials> {

  @JsonProperty private String clientId;
  @JsonProperty private String tenantId;
  @JsonProperty private String domain;
  // TODO: might be needed again if an account has multiple subscriptions to find the correct one
  // @JsonProperty private String subscriptionId;
  @JsonProperty private String clientSecret;

  public static AzureAccount discover() throws IOException {
    var account = new AzureAccount();

    // fetch credentials from default credential chain
    var credentials = defaultCredentialProviderChain();

    var azure = Azure.authenticate(credentials).withDefaultSubscription();

    account.setAccountId(azure.getCurrentSubscription().displayName());
    account.setAutoDiscovered(true);
    account.setDomain(credentials.domain());

    return account;
  }

  private static AzureTokenCredentials defaultCredentialProviderChain() throws IOException {
    // check if the default credentials-file exists
    var credentialsFile = new File(defaultAuthFile());

    if (credentialsFile.exists()) {
      LOGGER.info("Using default credentials file {}", credentialsFile);
      return ApplicationTokenCredentials.fromFile(credentialsFile);
    } else {
      // otherwise, use default locations
      LOGGER.info("Did not find default credentials. Trying to use AzureCLI credentials instead.");
      return AzureCliCredentials.create();
    }
  }

  @Override
  public void validate() throws IOException {
    var credentials = this.resolveCredentials();

    try {
      var azure = Azure.authenticate(credentials).withDefaultSubscription();

      this.setAccountId(azure.getCurrentSubscription().displayName());
    } catch (RuntimeException ex) {
      throw new IOException(ex.getCause());
    }
  }

  public AzureTokenCredentials resolveCredentials() throws IOException {
    if (this.isAutoDiscovered()) {
      return AzureAccount.defaultCredentialProviderChain();
    } else {
      return new ApplicationTokenCredentials(
          clientId, tenantId, clientSecret, AzureEnvironment.AZURE);
    }
  }

  private static String defaultAuthFile() {
    return System.getenv()
        .getOrDefault(
            "AZURE_AUTH_LOCATION", System.getProperty("user.home") + "/.azure/clouditor.azureauth");
  }

  public void setClientId(String clientId) {
    this.clientId = clientId;
  }

  public void setTenantId(String tenantId) {
    this.tenantId = tenantId;
  }

  public void setClientSecret(String clientSecret) {
    this.clientSecret = clientSecret;
  }

  @JsonProperty
  public String getAuthFile() {
    return defaultAuthFile();
  }

  public void setDomain(String domain) {
    this.domain = domain;
  }
}
