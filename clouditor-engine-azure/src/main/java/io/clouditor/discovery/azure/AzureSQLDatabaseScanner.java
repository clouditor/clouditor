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

import com.microsoft.azure.management.sql.ReplicationLink;
import com.microsoft.azure.management.sql.SqlDatabase;
import com.microsoft.azure.management.sql.implementation.DatabaseBlobAuditingPolicyInner;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.ArrayList;
import java.util.List;

@ScannerInfo(
    assetType = "SQLDatabase",
    group = "Azure",
    service = "SQL",
    assetIcon = "fas fa-database")
// TODO: would be nice, if we could get the results of the AzureSQLServerScanner instead of
// re-scanning servers again
public class AzureSQLDatabaseScanner extends AzureScanner<SqlDatabase> {

  public AzureSQLDatabaseScanner() {
    super(SqlDatabase::id, SqlDatabase::name);
  }

  @Override
  protected List<SqlDatabase> list() {
    var servers =
        new ArrayList<>(
            this.resourceGroup != null
                ? this.api.azure().sqlServers().listByResourceGroup(this.resourceGroup)
                : this.api.azure().sqlServers().list());

    List<SqlDatabase> databases = new ArrayList<>();

    for (var server : servers) {
      server
          .databases()
          .list()
          .forEach(
              database -> {
                // filter out master databases, since they are internal Azure structures
                if (!database.name().equals("master")) {
                  databases.add(database);
                }
              });
    }

    return databases;
  }

  @Override
  protected Asset transform(SqlDatabase database) throws ScanException {
    var asset = super.transform(database);

    enrich(asset, "transparentDataEncryption", database, SqlDatabase::getTransparentDataEncryption);

    enrich(
        asset,
        "auditingPolicy",
        database,
        x ->
            this.api
                .azure()
                .sqlServers()
                .manager()
                .inner()
                .databaseBlobAuditingPolicies()
                .get(database.resourceGroupName(), database.sqlServerName(), database.name()),
        DatabaseBlobAuditingPolicyInner::id,
        DatabaseBlobAuditingPolicyInner::name);

    enrichList(
        asset,
        "replicationLinks",
        database,
        x -> x.listReplicationLinks().values(),
        ReplicationLink::id,
        ReplicationLink::name);

    return asset;
  }
}
