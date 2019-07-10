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
import software.amazon.awssdk.services.lambda.LambdaClient;
import software.amazon.awssdk.services.lambda.model.FunctionConfiguration;
import software.amazon.awssdk.services.lambda.model.GetPolicyRequest;
import software.amazon.awssdk.services.lambda.model.GetPolicyResponse;
import software.amazon.awssdk.services.lambda.model.ListFunctionsResponse;

class AwsLambdaScannerTest extends AwsScannerTest {

  @BeforeAll
  static void setUpOnce() {
    discoverAssets(
        LambdaClient.class,
        AwsLambdaScanner::new,
        api -> {
          when(api.listFunctions())
              .thenReturn(
                  ListFunctionsResponse.builder()
                      .functions(
                          FunctionConfiguration.builder()
                              .functionArn(
                                  "arn:aws:lambda:eu-central-1:123456789:function:function-1")
                              .functionName("function-1")
                              .kmsKeyArn("some-key")
                              .build(),
                          FunctionConfiguration.builder()
                              .functionArn(
                                  "arn:aws:lambda:eu-central-1:123456789:function:function-2")
                              .functionName("function-2")
                              .build())
                      .build());

          when(api.getPolicy(GetPolicyRequest.builder().functionName("function-1").build()))
              .thenReturn(GetPolicyResponse.builder().policy("*").build());

          when(api.getPolicy(GetPolicyRequest.builder().functionName("function-2").build()))
              .thenReturn(GetPolicyResponse.builder().policy("no-wildcard").build());
        });
  }

  @Test
  void testEnvVariablesEncryption() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/lambda/function-env-encryption.md"));

    assertNotNull(rule);

    var function1 = assets.get("arn:aws:lambda:eu-central-1:123456789:function:function-1");

    assertNotNull(function1);
    assertTrue(rule.evaluate(function1).isOk());

    var function2 = assets.get("arn:aws:lambda:eu-central-1:123456789:function:function-2");

    assertNotNull(function2);
    assertFalse(rule.evaluate(function2).isOk());
  }

  @Test
  void testLambdaPolicyCheck() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/lambda/function-policy-wildcard.md"));

    assertNotNull(rule);

    var function1 = assets.get("arn:aws:lambda:eu-central-1:123456789:function:function-1");

    assertNotNull(function1);
    assertFalse(rule.evaluate(function1).isOk());

    var function2 = assets.get("arn:aws:lambda:eu-central-1:123456789:function:function-2");

    assertNotNull(function2);
    assertTrue(rule.evaluate(function2).isOk());
  }
}
