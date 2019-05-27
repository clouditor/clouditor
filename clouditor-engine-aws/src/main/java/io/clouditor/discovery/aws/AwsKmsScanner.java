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

import io.clouditor.discovery.Asset;
import io.clouditor.discovery.ScanException;
import io.clouditor.discovery.ScannerInfo;
import java.util.List;
import java.util.stream.Collectors;
import software.amazon.awssdk.services.kms.KmsClient;
import software.amazon.awssdk.services.kms.KmsClientBuilder;
import software.amazon.awssdk.services.kms.model.DescribeKeyRequest;
import software.amazon.awssdk.services.kms.model.GetKeyPolicyRequest;
import software.amazon.awssdk.services.kms.model.GetKeyRotationStatusRequest;
import software.amazon.awssdk.services.kms.model.KeyManagerType;
import software.amazon.awssdk.services.kms.model.KeyMetadata;

@ScannerInfo(assetType = "Key", group = "AWS", service = "KMS")
public class AwsKmsScanner extends AwsScanner<KmsClient, KmsClientBuilder, KeyMetadata> {

  public AwsKmsScanner() {
    // TODO: name from tags?
    super(KmsClient::builder, KeyMetadata::arn, KeyMetadata::keyId);
  }

  @Override
  protected List<KeyMetadata> list() {
    /*
     * Filter out "master keys", since they are managed by AWS and no properties can be set for them.
     * An AWS master key can be identified as such, if the keyManager type of a key is "AWS".
     */
    return this.api.listKeys().keys().stream()
        .map(
            keyListEntry ->
                this.api
                    .describeKey(DescribeKeyRequest.builder().keyId(keyListEntry.keyId()).build())
                    .keyMetadata())
        .filter(keyMetadata -> keyMetadata.keyManager() != KeyManagerType.AWS)
        .collect(Collectors.toList());
  }

  @Override
  protected Asset transform(KeyMetadata keyMetadata) throws ScanException {
    var asset = super.transform(keyMetadata);

    asset.setProperty(
        "keyRotationStatus",
        this.api
            .getKeyRotationStatus(
                GetKeyRotationStatusRequest.builder().keyId(keyMetadata.keyId()).build())
            .keyRotationEnabled());

    asset.setProperty(
        "keyPolicy",
        this.api
            .getKeyPolicy(GetKeyPolicyRequest.builder().keyId(keyMetadata.keyId()).build())
            .policy());

    return asset;
  }
}
