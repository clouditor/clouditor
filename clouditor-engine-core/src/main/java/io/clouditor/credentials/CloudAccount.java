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

import com.fasterxml.jackson.annotation.JsonTypeInfo;
import com.fasterxml.jackson.annotation.JsonTypeInfo.Id;
import com.fasterxml.jackson.annotation.JsonTypeName;
import io.clouditor.util.Collection;
import io.clouditor.util.PersistentObject;
import java.io.IOException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@JsonTypeInfo(use = Id.NAME, property = "provider")
@Collection("accounts")
public abstract class CloudAccount<T> implements PersistentObject {

  protected static final Logger LOGGER = LoggerFactory.getLogger(CloudAccount.class);

  protected String accountId;
  protected String user;

  /**
   * Specifies that this account was auto-discovered and that credentials are provided by the
   * default chain provided by the Cloud provider client API library. Thus credential fields are
   * ignored.
   */
  private boolean autoDiscovered;

  public String getAccountId() {
    return accountId;
  }

  public void setAccountId(String accountId) {
    this.accountId = accountId;
  }

  public String getUser() {
    return user;
  }

  public void setUser(String user) {
    this.user = user;
  }

  public void setAutoDiscovered(boolean autoDiscovered) {
    this.autoDiscovered = autoDiscovered;
  }

  public boolean isAutoDiscovered() {
    return autoDiscovered;
  }

  public String getProvider() {
    var typeName = this.getClass().getAnnotation(JsonTypeName.class);

    return typeName != null ? typeName.value() : null;
  }

  public String getId() {
    return this.getProvider();
  }

  /**
   * Validates this account by issuing a simple call which should return additional information,
   * such as the connected user.
   *
   * <p>Additionally, it should update the {@link CloudAccount#accountId} and {@link
   * CloudAccount#user} fields with recent information.
   *
   * @throws IOException
   */
  public abstract void validate() throws IOException;

  public abstract T resolveCredentials() throws IOException;
}
