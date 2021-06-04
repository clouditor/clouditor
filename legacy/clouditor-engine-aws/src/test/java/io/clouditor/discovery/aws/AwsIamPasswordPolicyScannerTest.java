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

import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import software.amazon.awssdk.services.iam.IamClient;
import software.amazon.awssdk.services.iam.model.GetAccountPasswordPolicyResponse;
import software.amazon.awssdk.services.iam.model.PasswordPolicy;

class AwsIamPasswordPolicyScannerTest extends AwsScannerTest {

  @BeforeAll
  static void setUpOnce() {
    discoverAssets(
        IamClient.class,
        AwsIamPasswordPolicyScanner::new,
        api ->
            when(api.getAccountPasswordPolicy())
                .thenReturn(
                    GetAccountPasswordPolicyResponse.builder()
                        .passwordPolicy(
                            PasswordPolicy.builder()
                                .requireUppercaseCharacters(true)
                                .requireSymbols(true)
                                .expirePasswords(false)
                                .passwordReusePrevention(24)
                                .requireLowercaseCharacters(true)
                                .maxPasswordAge(90)
                                .hardExpiry(false)
                                .requireNumbers(true)
                                .minimumPasswordLength(14)
                                .build())
                        .build()));
  }

  @Test
  void testPasswordPolicy() throws IOException {
    var policy = assets.get(AwsIamPasswordPolicyScanner.ARN_AWS_IAM_PW_POLICY);

    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/iam/password-policy.md"));

    assertTrue(rule.evaluate(policy).isOk());
  }
}
