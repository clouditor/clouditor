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

import io.clouditor.discovery.ScannerInfo;
import java.util.Collections;
import java.util.List;
import software.amazon.awssdk.services.iam.model.PasswordPolicy;

@ScannerInfo(assetType = "PasswordPolicy", group = "AWS", service = "IAM")
public class AwsIamPasswordPolicyScanner extends AwsIamScanner<PasswordPolicy> {

  /**
   * This is not a really ARN but rather a simple hack to give the password policy an asset
   * identifier.
   */
  static final String ARN_AWS_IAM_PW_POLICY = "arn:aws:iam:::password-policy";

  public AwsIamPasswordPolicyScanner() {
    super(policy -> ARN_AWS_IAM_PW_POLICY, policy -> "Password Policy");
  }

  @Override
  protected List<PasswordPolicy> list() {
    return Collections.singletonList(this.api.getAccountPasswordPolicy().passwordPolicy());
  }
}
