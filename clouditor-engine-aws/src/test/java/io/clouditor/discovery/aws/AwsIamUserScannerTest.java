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
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.when;

import io.clouditor.assurance.RuleService;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import java.time.Instant;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import software.amazon.awssdk.services.iam.IamClient;
import software.amazon.awssdk.services.iam.model.AccessKeyLastUsed;
import software.amazon.awssdk.services.iam.model.AccessKeyMetadata;
import software.amazon.awssdk.services.iam.model.GetAccessKeyLastUsedRequest;
import software.amazon.awssdk.services.iam.model.GetAccessKeyLastUsedResponse;
import software.amazon.awssdk.services.iam.model.Group;
import software.amazon.awssdk.services.iam.model.ListAccessKeysRequest;
import software.amazon.awssdk.services.iam.model.ListAccessKeysResponse;
import software.amazon.awssdk.services.iam.model.ListGroupsForUserRequest;
import software.amazon.awssdk.services.iam.model.ListGroupsForUserResponse;
import software.amazon.awssdk.services.iam.model.ListMfaDevicesRequest;
import software.amazon.awssdk.services.iam.model.ListMfaDevicesResponse;
import software.amazon.awssdk.services.iam.model.ListUsersResponse;
import software.amazon.awssdk.services.iam.model.MFADevice;
import software.amazon.awssdk.services.iam.model.StatusType;
import software.amazon.awssdk.services.iam.model.User;

class AwsIamUserScannerTest extends AwsScannerTest {

  private static final String USER1_ARN =
      "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/engineering/USER1";
  private static final String USER2_ARN =
      "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/engineering/USER2";
  private static final String USER3_ARN =
      "arn:aws:iam::123456789012:user/division_abc/subdivision_xyz/engineering/USER3";

  @BeforeAll
  static void setUpOnce() throws IOException {
    discoverAssets(
        IamClient.class,
        AwsIamUserScanner::new,
        api -> {
          when(api.listUsers())
              .thenReturn(
                  ListUsersResponse.builder()
                      .users(
                          User.builder().arn(USER1_ARN).userName("USER1").build(),
                          User.builder().arn(USER2_ARN).userName("USER2").build(),
                          User.builder().arn(USER3_ARN).userName("USER3").build())
                      .build());

          when(api.listMFADevices(ListMfaDevicesRequest.builder().userName("USER1").build()))
              .thenReturn(ListMfaDevicesResponse.builder().build());

          when(api.listGroupsForUser(ListGroupsForUserRequest.builder().userName("USER1").build()))
              .thenReturn(
                  ListGroupsForUserResponse.builder()
                      .groups(Group.builder().groupName("some-group").build())
                      .build());

          when(api.listAccessKeys(ListAccessKeysRequest.builder().userName("USER1").build()))
              .thenReturn(
                  ListAccessKeysResponse.builder()
                      .accessKeyMetadata(
                          AccessKeyMetadata.builder().accessKeyId("some-key").build())
                      .build());

          when(api.listMFADevices(ListMfaDevicesRequest.builder().userName("USER2").build()))
              .thenReturn(
                  ListMfaDevicesResponse.builder()
                      .mfaDevices(MFADevice.builder().serialNumber("1234556").build())
                      .build());

          when(api.listGroupsForUser(ListGroupsForUserRequest.builder().userName("USER2").build()))
              .thenReturn(
                  ListGroupsForUserResponse.builder()
                      .groups(Group.builder().groupName("some-group").build())
                      .build());

          when(api.listAccessKeys(ListAccessKeysRequest.builder().userName("USER2").build()))
              .thenReturn(
                  ListAccessKeysResponse.builder()
                      .accessKeyMetadata(AccessKeyMetadata.builder().accessKeyId("old-key").build())
                      .build());

          when(api.listMFADevices(ListMfaDevicesRequest.builder().userName("USER3").build()))
              .thenReturn(ListMfaDevicesResponse.builder().build());

          when(api.listGroupsForUser(ListGroupsForUserRequest.builder().userName("USER3").build()))
              .thenReturn(ListGroupsForUserResponse.builder().build());

          when(api.listAccessKeys(ListAccessKeysRequest.builder().userName("USER3").build()))
              .thenReturn(
                  ListAccessKeysResponse.builder()
                      .accessKeyMetadata(
                          AccessKeyMetadata.builder()
                              .accessKeyId("some-key")
                              .status(StatusType.INACTIVE)
                              .build())
                      .build());

          when(api.getAccessKeyLastUsed(
                  GetAccessKeyLastUsedRequest.builder().accessKeyId("some-key").build()))
              .thenReturn(
                  GetAccessKeyLastUsedResponse.builder()
                      .accessKeyLastUsed(
                          AccessKeyLastUsed.builder().lastUsedDate(Instant.now()).build())
                      .build());

          when(api.getAccessKeyLastUsed(
                  GetAccessKeyLastUsedRequest.builder().accessKeyId("old-key").build()))
              .thenReturn(
                  GetAccessKeyLastUsedResponse.builder()
                      .accessKeyLastUsed(
                          AccessKeyLastUsed.builder().lastUsedDate(Instant.MIN).build())
                      .build());
        });
  }

  @Test
  void testMFA() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(FileSystemManager.getInstance().getPathForResource("rules/aws/iam/mfa.md"));

    // user1 has no MFA
    assertFalse(rule.evaluate(assets.get(USER1_ARN)).isOk());

    // user2 should have an active MFA
    assertTrue(rule.evaluate(assets.get(USER2_ARN)).isOk());
  }

  @Test
  void testUsersInGroup() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/iam/users-in-groups.md"));

    // user3 has no groups
    assertFalse(rule.evaluate(assets.get(USER3_ARN)).isOk());
  }

  @Test
  void testAccessKeyRotation() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/iam/access-key-rotation.md"));

    // user2 has a very old access key
    assertFalse(rule.evaluate(assets.get(USER2_ARN)).isOk());
  }

  @Test
  void testInactiveAccessKeys() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(
                FileSystemManager.getInstance()
                    .getPathForResource("rules/aws/iam/inactive-access-keys.md"));

    // user3 has inactive access keys
    assertFalse(rule.evaluate(assets.get(USER3_ARN)).isOk());
  }
}
