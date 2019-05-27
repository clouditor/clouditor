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
    for (var scan : this.discoveryService.getScans().values()) {
      var scanner = scan.getScanner();

      if (scanner.getInitialized() && Objects.equals(scan.getGroup(), provider)) {
        LOGGER.info("Forcing scanner {} / {} to re-authenticate.", scan.getService(), scan.getId());
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
