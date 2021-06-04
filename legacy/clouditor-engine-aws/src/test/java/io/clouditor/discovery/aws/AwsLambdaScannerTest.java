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
