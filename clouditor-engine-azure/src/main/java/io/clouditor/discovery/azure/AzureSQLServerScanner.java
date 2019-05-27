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

import com.microsoft.azure.management.sql.SqlActiveDirectoryAdministrator;
import com.microsoft.azure.management.sql.SqlEncryptionProtector;
import com.microsoft.azure.management.sql.SqlFirewallRule;
import com.microsoft.azure.management.sql.SqlServer;
import com.microsoft.azure.management.sql.SqlServerSecurityAlertPolicy;
import com.microsoft.azure.management.sql.implementation.ServerBlobAuditingPolicyInner;
import io.clouditor.discovery.Asset;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.List;

@ScannerInfo(assetType = "SQLServer", group = "Azure", service = "SQL", assetIcon = "fas fa-server")
public class AzureSQLServerScanner extends AzureScanner<SqlServer> {

  public AzureSQLServerScanner() {
    super(SqlServer::id, SqlServer::name);
  }

  @Override
  protected List<SqlServer> list() {
    return this.resourceGroup != null
        ? this.api.azure().sqlServers().listByResourceGroup(this.resourceGroup)
        : this.api.azure().sqlServers().list();
  }

  @Override
  protected Asset transform(SqlServer server) throws ScanException {
    var asset = super.transform(server);

    enrich(
        asset,
        "securityAlertPolicy",
        server,
        x -> x.serverSecurityAlertPolicies().get(),
        SqlServerSecurityAlertPolicy::id,
        SqlServerSecurityAlertPolicy::name);

    enrich(
        asset,
        "encryptionProtectors",
        server,
        x -> x.encryptionProtectors().get(),
        SqlEncryptionProtector::id,
        SqlEncryptionProtector::serverKeyName);

    enrich(
        asset,
        "activeDirectoryAdmin",
        server,
        x -> server.getActiveDirectoryAdministrator(),
        SqlActiveDirectoryAdministrator::id,
        SqlActiveDirectoryAdministrator::signInName);

    enrichList(
        asset,
        "firewallRules",
        server,
        x -> x.firewallRules().list(),
        SqlFirewallRule::id,
        SqlFirewallRule::name);

    enrich(
        asset,
        "auditingPolicy",
        server,
        x ->
            this.api
                .azure()
                .sqlServers()
                .manager()
                .inner()
                .serverBlobAuditingPolicies()
                .get(server.resourceGroupName(), server.name()),
        ServerBlobAuditingPolicyInner::id,
        ServerBlobAuditingPolicyInner::name);

    return asset;
  }
}
