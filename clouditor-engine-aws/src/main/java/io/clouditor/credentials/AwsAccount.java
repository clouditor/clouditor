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

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonTypeName;
import java.io.IOException;
import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.Id;
import javax.persistence.Table;
import software.amazon.awssdk.auth.credentials.AwsCredentials;
import software.amazon.awssdk.auth.credentials.AwsCredentialsProvider;
import software.amazon.awssdk.auth.credentials.DefaultCredentialsProvider;
import software.amazon.awssdk.core.exception.SdkClientException;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.sts.StsClient;
import software.amazon.awssdk.services.sts.model.StsException;

@Entity(name = "aws_account")
@Table(name = "aws_account")
@JsonTypeName(value = "AWS")
public class AwsAccount extends CloudAccount<AwsCredentials>
    implements AwsCredentials, AwsCredentialsProvider {

  private static final DefaultCredentialsProvider DEFAULT_PROVIDER =
      DefaultCredentialsProvider.create();

  private static final long serialVersionUID = 1928775323719265066L;

  @Column(name = "access_key_id")
  @JsonProperty
  private String accessKeyId;

  @Column(name = "secret_access_key")
  @JsonProperty
  private String secretAccessKey;

  @Column(name = "region")
  @JsonProperty
  private String region;

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

  /**
   * Discovers an AWS account.
   *
   * @return null, if no account was discovered. Otherwise the discovered {@link AwsAccount}.
   */
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

  @Id //  enable the access to the property accessKeyId through the getter method by default
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

  public void setAccessKeyId(String accessKeyId) {
    this.accessKeyId = accessKeyId;
  }

  public void setSecretAccessKey(String secretAccessKey) {
    this.secretAccessKey = secretAccessKey;
  }
}
