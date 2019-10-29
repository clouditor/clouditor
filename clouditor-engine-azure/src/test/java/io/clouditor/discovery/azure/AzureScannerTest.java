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

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.mockito.Mockito.RETURNS_DEEP_STUBS;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

import com.microsoft.azure.Page;
import com.microsoft.azure.PagedList;
import com.microsoft.azure.keyvault.models.KeyBundle;
import com.microsoft.azure.management.Azure;
import com.microsoft.azure.management.compute.EncryptionStatus;
import com.microsoft.azure.management.compute.VirtualMachine;
import com.microsoft.azure.management.compute.implementation.VirtualMachineInner;
import com.microsoft.azure.management.keyvault.Key;
import com.microsoft.azure.management.monitor.implementation.MonitorManager;
import com.microsoft.azure.management.resources.ResourceGroup;
import com.microsoft.azure.management.resources.Subscription;
import com.microsoft.azure.management.resources.fluentcore.arm.models.HasId;
import com.microsoft.azure.management.resources.fluentcore.arm.models.HasName;
import com.microsoft.azure.management.resources.fluentcore.model.HasInner;
import com.microsoft.azure.management.resources.implementation.ResourceGroupInner;
import com.microsoft.azure.management.resources.implementation.SubscriptionInner;
import com.microsoft.azure.management.sql.SqlDatabase;
import com.microsoft.azure.management.sql.SqlServer;
import com.microsoft.azure.management.sql.implementation.DatabaseInner;
import com.microsoft.azure.management.sql.implementation.ServerInner;
import com.microsoft.azure.management.sql.implementation.SqlServerManager;
import com.microsoft.azure.management.storage.StorageAccount;
import com.microsoft.azure.management.storage.StorageAccountEncryptionStatus;
import com.microsoft.azure.management.storage.StorageService;
import com.microsoft.azure.management.storage.implementation.StorageAccountInner;
import com.microsoft.rest.RestException;
import io.clouditor.Engine;
import io.clouditor.discovery.DiscoveryResult;
import java.util.Arrays;
import java.util.Map;
import java.util.Map.Entry;
import java.util.function.Consumer;
import java.util.function.Supplier;
import java.util.stream.Collectors;
import org.joda.time.DateTime;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public abstract class AzureScannerTest {

  private static final String MOCK_SUBSCRIPTION_ID = "00000000-1111-2222-3333-444444444444";
  private static final String MOCK_RESOURCE_GROUP_ID = "mock";
  private static final String MOCK_RESOURCE_GROUP_NAME =
      "/subscriptions/" + MOCK_SUBSCRIPTION_ID + "/resourceGroups/" + MOCK_RESOURCE_GROUP_ID;

  private static final Logger LOGGER = LoggerFactory.getLogger(AzureScannerTest.class);

  Engine engine = new Engine();

  static DiscoveryResult assets;

  static <T extends HasInner> void discoverAssets(
      Supplier<AzureScanner<T>> supplier, Consumer<AzureClients> configurator) {

    var scanner = supplier.get();

    // scanner.init(); don't call init
    scanner.setInitialized(true);

    var api = new AzureClients();
    api.azure = mock(Azure.class, RETURNS_DEEP_STUBS);
    api.monitor = mock(MonitorManager.class, RETURNS_DEEP_STUBS);

    // set up a subscription
    var subscription = createSubscription(MOCK_SUBSCRIPTION_ID, "mock", new SubscriptionInner());

    when(api.azure.getCurrentSubscription()).thenReturn(subscription);
    when(api.azure.subscriptions().list()).thenReturn(MockedPagedList.of(subscription));

    // set up a resource group
    var resourceGroup = createResourceGroup(MOCK_RESOURCE_GROUP_ID, MOCK_RESOURCE_GROUP_NAME);

    when(api.azure.resourceGroups().list()).thenReturn(MockedPagedList.of(resourceGroup));

    // some helpers for sql
    when(api.azure.sqlServers().manager())
        .thenReturn(mock(SqlServerManager.class, RETURNS_DEEP_STUBS));

    configurator.accept(api);

    scanner.setApi(api);

    assets = scanner.scan(null);

    assertNotNull(assets);

    LOGGER.info("Assets: {}", assets);

    assertFalse(assets.getDiscoveredAssets().isEmpty());
  }

  private static ResourceGroup createResourceGroup(String id, String name) {
    return createWithIdAndName(ResourceGroup.class, id, name, new ResourceGroupInner());
  }

  private static Subscription createSubscription(String id, String name, SubscriptionInner inner) {
    var subscription = mock(Subscription.class, RETURNS_DEEP_STUBS);

    when(subscription.subscriptionId()).thenReturn(id);
    when(subscription.displayName()).thenReturn(name);
    when(subscription.inner()).thenReturn(inner);

    return subscription;
  }

  static class MockedPagedList<E> extends PagedList<E> {

    @Override
    public Page<E> nextPage(String nextPageLink) throws RestException {
      return null;
    }

    static <E> MockedPagedList<E> of(E element) {
      var list = new MockedPagedList<E>();

      list.add(element);

      return list;
    }

    @SafeVarargs
    static <E> MockedPagedList<E> of(E... elements) {
      var list = new MockedPagedList<E>();

      list.addAll(Arrays.asList(elements));

      return list;
    }
  }

  static VirtualMachine createVirtualMachine(
      String id,
      VirtualMachineInner inner,
      EncryptionStatus osDiskStatus,
      EncryptionStatus dataDiskStatus) {
    var vm = createWithId(VirtualMachine.class, id, inner);

    when(vm.diskEncryption().getMonitor().osDiskStatus()).thenReturn(osDiskStatus);
    when(vm.diskEncryption().getMonitor().dataDiskStatus()).thenReturn(dataDiskStatus);

    return vm;
  }

  static StorageAccount createStorageAccount(
      String id, StorageAccountInner inner, Map<StorageService, Boolean> encryptionStatuses) {
    var storage = createWithId(StorageAccount.class, id, inner);

    Map<StorageService, StorageAccountEncryptionStatus> y =
        encryptionStatuses.entrySet().stream()
            .collect(
                Collectors.toMap(
                    Entry::getKey,
                    e ->
                        new StorageAccountEncryptionStatus() {
                          @Override
                          public StorageService storageService() {
                            return e.getKey();
                          }

                          @Override
                          public boolean isEnabled() {
                            return e.getValue();
                          }

                          @Override
                          public DateTime lastEnabledTime() {
                            return DateTime.now();
                          }
                        }));

    when(storage.encryptionStatuses()).thenReturn(y);

    return storage;
  }

  static SqlDatabase createSqlDatabase(String id, String name, DatabaseInner inner) {
    var db = mock(SqlDatabase.class);

    when(db.id()).thenReturn(id);
    when(db.name()).thenReturn(name);
    when(db.inner()).thenReturn(inner);

    return db;
  }

  static SqlServer createSqlServer(String id, ServerInner inner) {
    return createWithId(SqlServer.class, id, inner);
  }

  static <T extends HasInner<InnerT> & HasId, InnerT> T createWithId(
      Class<T> clazz, String id, InnerT inner) {
    var resource = mock(clazz, RETURNS_DEEP_STUBS);

    when(resource.id()).thenReturn(id);
    when(resource.inner()).thenReturn(inner);

    return resource;
  }

  static Key createKey(String id, String name, KeyBundle inner) {
    return createWithIdAndName(Key.class, id, name, inner);
  }

  static <T extends HasInner<InnerT> & HasId & HasName, InnerT> T createWithIdAndName(
      Class<T> clazz, String id, String name, InnerT inner) {
    var resource = mock(clazz, RETURNS_DEEP_STUBS);

    when(resource.id()).thenReturn(id);
    when(resource.name()).thenReturn(name);
    when(resource.inner()).thenReturn(inner);

    return resource;
  }
}
