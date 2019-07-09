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
