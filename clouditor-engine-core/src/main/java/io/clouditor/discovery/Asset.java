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

import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.assurance.EvaluationResult;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;
import javax.persistence.*;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.commons.lang3.builder.ToStringStyle;

@Entity(name = "asset")
@Table(name = "asset")
@Inheritance(strategy = InheritanceType.JOINED)
public class Asset implements Serializable {

  private static final long serialVersionUID = -2382328140875227410L;

  @ManyToMany(targetEntity = EvaluationResult.class)
  private List<EvaluationResult> evaluationResults = new ArrayList<>();

  @CollectionTable(name = "asset_properties", joinColumns = @JoinColumn(name = "key_id"))
  @MapKeyColumn(name = "mapKey")
  @Column(name = "asset_properties")
  private AssetProperties properties = new AssetProperties();

  @Id
  @Column(name = "asset_id")
  private String id;

  @Column(name = "asset_name")
  private String name;

  @Column(name = "type_value")
  private String type;

  public Asset() {
    // nothing to do
  }

  public Asset(String type, String id, String name, AssetProperties properties) {
    this.type = type;
    this.id = id;
    this.name = name;
    this.properties = properties;
  }

  public void setEvaluationResults(List<EvaluationResult> evaluationResults) {
    this.evaluationResults = evaluationResults;
  }

  public AssetProperties getProperties() {
    return properties;
  }

  public void setProperties(AssetProperties properties) {
    this.properties = properties;
  }

  public String getId() {
    return id;
  }

  public void setId(String id) {
    this.id = id;
  }

  public String getName() {
    return name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public String getType() {
    return type;
  }

  public void setType(String type) {
    this.type = type;
  }

  @JsonProperty
  public boolean isCompliant() {
    return this.evaluationResults.stream().allMatch(EvaluationResult::isOk);
  }

  public void addEvaluationResult(EvaluationResult result) {
    this.evaluationResults.add(result);
  }

  public List<EvaluationResult> getEvaluationResults() {
    return evaluationResults;
  }

  public <T> void setProperty(String key, T value) {
    this.properties.put(key, value);
  }

  @Override
  public String toString() {
    return new ToStringBuilder(this, ToStringStyle.JSON_STYLE)
        .append("evaluationResults", evaluationResults)
        .append("properties", properties)
        .append("id", id)
        .append("name", name)
        .append("type", type)
        .toString();
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;
    if (o == null || getClass() != o.getClass()) return false;
    Asset asset = (Asset) o;
    return Objects.equals(getEvaluationResults(), asset.getEvaluationResults())
        && Objects.equals(getProperties(), asset.getProperties())
        && Objects.equals(getId(), asset.getId())
        && Objects.equals(getName(), asset.getName())
        && Objects.equals(getType(), asset.getType());
  }

  @Override
  public int hashCode() {
    return Objects.hash(getEvaluationResults(), getProperties(), getId(), getName(), getType());
  }
}
