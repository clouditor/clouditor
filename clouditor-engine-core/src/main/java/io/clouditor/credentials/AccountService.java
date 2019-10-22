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

import io.clouditor.discovery.DiscoveryService;
import io.clouditor.util.PersistenceManager;
import java.io.IOException;
import java.lang.reflect.InvocationTargetException;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import java.util.function.Consumer;
import javax.inject.Inject;
import org.jvnet.hk2.annotations.Service;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@Service
public class AccountService {

  private DiscoveryService discoveryService;

  @Inject
  public AccountService(DiscoveryService discoveryService) {
    this.discoveryService = discoveryService;
  }

  protected static final Logger LOGGER = LoggerFactory.getLogger(AccountService.class);

  public CloudAccount discover(String provider) {
    LOGGER.info("Trying to discover accounts for provider {}", provider);

    try {

      Class<?> c;
      switch (provider) {
        case "AWS":
          c = Class.forName("io.clouditor.credentials.AwsAccount");
          break;

        case "Azure":
          c = Class.forName("io.clouditor.credentials.AzureAccount");
          break;
        default:
          throw new IOException("Provider not supported");
      }

      var m = c.getMethod("discover");

      return (CloudAccount) m.invoke(null);
    } catch (ClassNotFoundException
        | NoSuchMethodException
        | IllegalAccessException
        | InvocationTargetException
        | IOException e) {
      LOGGER.error("Could not discover {} account: {}", provider, e.getCause());
      return null;
    }
  }

  public CloudAccount getAccount(String provider) {
    return PersistenceManager.getInstance().getById(CloudAccount.class, provider);
  }

  public void addAccount(String provider, CloudAccount account) throws IOException {
    LOGGER.info("Trying to validate account for provider {}...", provider);

    account.validate();

    LOGGER.info("Adding account for provider {} with id {}", provider, account.getId());

    // TODO: check, if something actually has changed

    PersistenceManager.getInstance().persist(account);

    // since we changed the account (potentially), we need to make sure the scanners associated with
    // this provider re-authenticate properly
    for (var scanner : this.discoveryService.getScanners()) {
      var info = scanner.getInfo();

      if (scanner.getInitialized() && Objects.equals(info.group(), provider)) {
        LOGGER.info("Forcing scanner {} to re-authenticate.", scanner.getId());
        scanner.setInitialized(false);
      }
    }
  }

  public Map<String, CloudAccount> getAccounts() {
    var accounts = new HashMap<String, CloudAccount>();

    PersistenceManager.getInstance()
        .find(CloudAccount.class)
        .forEach(
            (Consumer<? super CloudAccount>)
                account -> accounts.put(account.getProvider(), account));

    return accounts;
  }
}
