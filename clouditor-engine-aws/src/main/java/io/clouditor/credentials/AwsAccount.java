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
import java.io.IOException;
import software.amazon.awssdk.auth.credentials.AwsCredentials;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProvider;
import software.amazon.awssdk.auth.credentials.DefaultCredentialsProvider;
import software.amazon.awssdk.core.exception.SdkClientException;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.sts.StsClient;
import software.amazon.awssdk.services.sts.model.StsException;

@JsonTypeName(value = "AWS")
public class AwsAccount extends CloudAccount<AwsCredentials>
    implements AwsCredentials, AwsCredentialsProvider {

  private static final DefaultCredentialsProvider DEFAULT_PROVIDER =
      DefaultCredentialsProvider.create();

  @JsonProperty private String accessKeyId;
  @JsonProperty private String secretAccessKey;
  @JsonProperty private String region;

  @Override
  public void validate() throws IOException {
    try {
      // use STS to find account id and user

      var builder = StsClient.builder();

      if (!this.isAutoDiscovered()) {
        builder.region(Region.of(this.region));
        builder.credentialsProvider(() -> this);
      }

      var stsClient = builder.build();

      var identity = stsClient.getCallerIdentity();

      this.accountId = identity.account();
      this.user = identity.arn();

      LOGGER.info("Account {} validated with user {}.", this.accountId, this.user);
    } catch (SdkClientException | StsException ex) {
      // TODO: log error, etc.
      throw new IOException(ex.getMessage());
    }
  }

  public static AwsAccount discover() {
    try {
      var account = new AwsAccount();

      // use STS to find account id using the default provider
      var stsClient = StsClient.builder().credentialsProvider(DEFAULT_PROVIDER).build();

      var identity = stsClient.getCallerIdentity();

      account.setAutoDiscovered(true);
      account.setAccountId(identity.account());
      account.setUser(identity.arn());

      return account;
    } catch (SdkClientException ex) {
      // TODO: log error, etc.
      return null;
    }
  }

  @Override
  public AwsCredentials resolveCredentials() {
    // check, if account is auto-discovered
    if (this.isAutoDiscovered()) {
      // then, hand it down to the default AWS provider chain
      return DEFAULT_PROVIDER.resolveCredentials();
    }

    // otherwise, we need to specify the stored credentials
    return this;
  }

  @Override
  public String accessKeyId() {
    return this.accessKeyId;
  }

  @Override
  public String secretAccessKey() {
    return this.secretAccessKey;
  }

  public String getRegion() {
    return this.region;
  }
}
