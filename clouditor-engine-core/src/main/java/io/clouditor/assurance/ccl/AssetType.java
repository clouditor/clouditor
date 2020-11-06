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

package io.clouditor.assurance.ccl;

import com.fasterxml.jackson.annotation.JsonTypeInfo;
import io.clouditor.data_access_layer.PersistentObject;
import io.clouditor.discovery.AssetProperties;
import java.io.Serializable;
import javax.persistence.*;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;

@JsonTypeInfo(use = JsonTypeInfo.Id.CLASS, property = "@class")
@Entity(name = "asset_type")
@Table(name = "asset_type")
@Inheritance(strategy = InheritanceType.JOINED)
public class AssetType implements Serializable, PersistentObject<String> {

  private static final long serialVersionUID = 5799337361597734304L;

  @Id
  @Column(name = "type_value", nullable = false)
  private String value;

  public String getValue() {
    return this.value;
  }

  public void setValue(String value) {
    this.value = value;
  }

  public boolean evaluate(AssetProperties properties) {
    return true;
  }

  @Override
  public String getId() {
    return value;
  }

  @Override
  public String toString() {
    return new ToStringBuilder(this).append("value", value).toString();
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    AssetType assetType = (AssetType) o;

    return new EqualsBuilder().append(value, assetType.value).isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37).append(value).toHashCode();
  }
}
