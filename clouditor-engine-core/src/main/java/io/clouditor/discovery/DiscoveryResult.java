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

package io.clouditor.discovery;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.assurance.ccl.AssetType;
import java.io.Serializable;
import java.time.Instant;
import java.util.HashMap;
import java.util.Map;
import javax.persistence.*;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.commons.lang3.builder.ToStringStyle;

@Entity(name = "discovery_result")
@Table(name = "discovery_result")
public class DiscoveryResult implements Serializable {

  private static final long serialVersionUID = -7032902561471865653L;

  @Id
  @Column(name = "time_stamp")
  private Instant timestamp;

  @ManyToMany @Embedded private Map<String, Asset> discoveredAssets = new HashMap<>();

  @Column(name = "failed")
  private boolean failed = false;

  @Column(name = "error")
  private String error;

  @JsonProperty
  @OneToOne(cascade = CascadeType.ALL)
  private final Scan scanId;

  public void setTimestamp(Instant timestamp) {
    this.timestamp = timestamp;
  }

  public Instant getTimestamp() {
    return timestamp;
  }

  @JsonCreator
  public DiscoveryResult(@JsonProperty(value = "scanId") AssetType scanId) {
    this.timestamp = Instant.now();

    this.scanId = new Scan();
    this.scanId.setAssetType(scanId);
  }

  public void setDiscoveredAssets(Map<String, Asset> discoveredAssets) {
    this.discoveredAssets = discoveredAssets;
  }

  public Map<String, Asset> getDiscoveredAssets() {
    return discoveredAssets;
  }

  public Asset get(String assetId) {
    return this.discoveredAssets.get(assetId);
  }

  public boolean isFailed() {
    return failed;
  }

  public void setFailed(boolean failed) {
    this.failed = failed;
  }

  public String getError() {
    return error;
  }

  public void setError(String error) {
    this.error = error;
  }

  @Override
  public String toString() {
    return new ToStringBuilder(this, ToStringStyle.JSON_STYLE)
        .append("timestamp", timestamp)
        .append("discoveredAssets", discoveredAssets)
        .append("failed", failed)
        .append("error", error)
        .toString();
  }

  public Scan getScanId() {
    return this.scanId;
  }

  public void setScanId(AssetType scanId) {
    this.scanId.setAssetType(scanId);
  }
}
