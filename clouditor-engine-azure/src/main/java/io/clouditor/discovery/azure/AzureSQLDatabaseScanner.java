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
