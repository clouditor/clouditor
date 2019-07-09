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

package io.clouditor.discovery.aws;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;
import software.amazon.awssdk.services.dynamodb.model.DescribeTableRequest;
import software.amazon.awssdk.services.dynamodb.model.DescribeTableResponse;
import software.amazon.awssdk.services.dynamodb.model.ListTablesResponse;
import software.amazon.awssdk.services.dynamodb.model.SSEDescription;
import software.amazon.awssdk.services.dynamodb.model.SSEStatus;
import software.amazon.awssdk.services.dynamodb.model.TableDescription;

class AwsDynamoDbScannerTest extends AwsScannerTest {

  @BeforeAll
  static void setUpOnce() {
    discoverAssets(
        DynamoDbClient.class,
        AwsDynamoDbScanner::new,
        api -> {
          when(api.listTables())
              .thenReturn(
                  ListTablesResponse.builder()
                      .tableNames(
                          "enabled_encryption",
                          "enabling_encryption",
                          "disabled_encryption",
                          "disabling_encryption")
                      .build());

          when(api.describeTable(
                  DescribeTableRequest.builder().tableName("enabled_encryption").build()))
              .thenReturn(
                  DescribeTableResponse.builder()
                      .table(
                          TableDescription.builder()
                              .sseDescription(
                                  SSEDescription.builder().status(SSEStatus.ENABLED).build())
                              .tableArn(
                                  "arn:aws:dynamodb:eu-central-1:123456789:table/encryption-enabled-table")
                              .build())
                      .build());

          when(api.describeTable(
                  DescribeTableRequest.builder().tableName("enabling_encryption").build()))
              .thenReturn(
                  DescribeTableResponse.builder()
                      .table(
                          TableDescription.builder()
                              .sseDescription(
                                  SSEDescription.builder().status(SSEStatus.ENABLING).build())
                              .tableArn(
                                  "arn:aws:dynamodb:eu-central-1:123456789:table/encryption-enabling-table")
                              .build())
                      .build());

          when(api.describeTable(
                  DescribeTableRequest.builder().tableName("disabling_encryption").build()))
              .thenReturn(
                  DescribeTableResponse.builder()
                      .table(
                          TableDescription.builder()
                              .sseDescription(
                                  SSEDescription.builder().status(SSEStatus.DISABLING).build())
                              .tableArn(
                                  "arn:aws:dynamodb:eu-central-1:123456789:table/encryption-disabling-table")
                              .build())
                      .build());

          when(api.describeTable(
                  DescribeTableRequest.builder().tableName("disabled_encryption").build()))
              .thenReturn(
                  DescribeTableResponse.builder()
                      .table(
                          TableDescription.builder()
                              .sseDescription(
                                  SSEDescription.builder().status(SSEStatus.DISABLED).build())
                              .tableArn(
                                  "arn:aws:dynamodb:eu-central-1:123456789:table/encryption-disabled-table")
                              .build())
                      .build());
        });
  }

  @Test
  void testEncryption() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/dynamodb/table-encrypted.md"));

    assertNotNull(rule);

    var table =
        assets.get("arn:aws:dynamodb:eu-central-1:123456789:table/encryption-enabling-table");

    assertNotNull(table);
    assertFalse(rule.evaluate(table).isOk());

    table = assets.get("arn:aws:dynamodb:eu-central-1:123456789:table/encryption-disabling-table");

    assertNotNull(table);
    assertFalse(rule.evaluate(table).isOk());

    table = assets.get("arn:aws:dynamodb:eu-central-1:123456789:table/encryption-enabled-table");

    assertNotNull(table);
    assertTrue(rule.evaluate(table).isOk());

    table = assets.get("arn:aws:dynamodb:eu-central-1:123456789:table/encryption-disabled-table");

    assertNotNull(table);
    assertFalse(rule.evaluate(table).isOk());
  }
}
