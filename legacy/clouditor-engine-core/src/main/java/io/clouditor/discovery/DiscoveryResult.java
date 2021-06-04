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
import java.io.Serializable;
import java.time.Instant;
import java.util.HashMap;
import java.util.Map;
import javax.persistence.*;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.commons.lang3.builder.ToStringStyle;
import org.hibernate.annotations.LazyCollection;
import org.hibernate.annotations.LazyCollectionOption;

@Entity(name = "discovery_result")
@Table(name = "discovery_result")
public class DiscoveryResult implements Serializable {

  private static final long serialVersionUID = -7032902561471865653L;

  @Id
  @Column(nullable = false)
  private Instant timestamp;

  @ManyToMany
  @LazyCollection(LazyCollectionOption.FALSE)
  @Embedded
  private Map<String, Asset> discoveredAssets = new HashMap<>();

  @Column() private boolean failed = false;

  @Column() private String error;

  @JsonProperty
  @Column(name = "scan_id")
  private String scanId;

  public DiscoveryResult() {}

  @JsonCreator
  public DiscoveryResult(@JsonProperty(value = "scanId") String scanId) {
    // round the Instant to epoch milli seconds (To be consistent with the PostgreSQL DB)
    this.timestamp = Instant.ofEpochMilli(Instant.now().toEpochMilli());
    this.scanId = scanId;
  }

  public Instant getTimestamp() {
    return timestamp;
  }

  public void setTimestamp(final Instant timestamp) {
    // round the Instant to epoch milli seconds (To be consistent with the PostgreSQL DB)
    this.timestamp = Instant.ofEpochMilli(timestamp.toEpochMilli());
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

  public String getScanId() {
    return this.scanId;
  }

  public void setScanId(final String scanId) {
    this.scanId = scanId;
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

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    DiscoveryResult that = (DiscoveryResult) o;

    return new EqualsBuilder()
        .append(failed, that.failed)
        .append(timestamp, that.timestamp)
        .append(discoveredAssets, that.discoveredAssets)
        .append(error, that.error)
        .append(scanId, that.scanId)
        .isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37)
        .append(timestamp)
        .append(discoveredAssets)
        .append(failed)
        .append(error)
        .append(scanId)
        .toHashCode();
  }
}
