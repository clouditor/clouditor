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

package io.clouditor.assurance;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;
import io.clouditor.assurance.ccl.*;
import io.clouditor.data_access_layer.PersistentObject;
import io.clouditor.discovery.Asset;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;
import javax.persistence.*;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.hibernate.annotations.LazyCollection;
import org.hibernate.annotations.LazyCollectionOption;

@Entity(name = "rule")
@Table(name = "rule")
public class Rule implements PersistentObject<String> {

  private static final long serialVersionUID = -5934273783785749037L;

  @JsonDeserialize(contentUsing = CCLDeserializer.class)
  @JsonSerialize(contentUsing = CCLSerializer.class)
  @Embedded
  @ElementCollection(targetClass = Condition.class)
  @LazyCollection(LazyCollectionOption.FALSE)
  private List<Condition> conditions = new ArrayList<>();

  @Column(name = "active")
  private boolean active;

  @JsonProperty
  @Column(name = "rule_name")
  private String name;

  @JsonProperty
  @Column(name = "rule_description")
  private String description;

  @JsonProperty
  @Column(name = "icon")
  private final String icon = "far fa-file-alt";

  @JsonProperty
  @Embedded
  @ElementCollection(targetClass = String.class)
  @LazyCollection(LazyCollectionOption.FALSE)
  private List<String> controls = new ArrayList<>();

  @JsonProperty
  @Id
  @Column(name = "rule_id", nullable = false)
  private String id;

  public boolean isAssetFiltered(Asset asset) {
    if (this.conditions != null && !this.conditions.isEmpty()) {
      // Theoretically, a list of conditions could have different asset types each, but we assume
      // they are all equal
      if (this.conditions.get(0).getAssetType() instanceof FilteredAssetType) {
        return !this.conditions.get(0).getAssetType().evaluate(asset.getProperties());
      }
    }
    return false;
  }

  public EvaluationResult evaluate(Asset asset) {
    var eval = new EvaluationResult(this, asset.getProperties());

    if (!this.conditions.isEmpty()) {
      // get those conditions which evaluate to false
      eval.setFailedConditions(
          this.conditions.stream()
              .filter(c -> !c.evaluate(asset.getProperties()))
              .collect(Collectors.toList()));
    }

    return eval;
  }

  public String getAssetType() {
    if (this.conditions != null
        && !this.conditions.isEmpty()
        && this.conditions.get(0).getAssetType() != null) {
      return this.conditions.get(0).getAssetType().getValue();
    }

    // no asset type found, we cannot really use this rule then
    return null;
  }

  public boolean isActive() {
    return active;
  }

  public void setActive(boolean active) {
    this.active = active;
  }

  public String getName() {
    return this.name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public void setId(String id) {
    this.id = id;
  }

  public List<String> getControls() {
    return this.controls;
  }

  public void setControls(final List<String> controls) {
    this.controls = controls;
  }

  public void addControls(final String... controls) {
    final List<String> controlList = List.of(controls);
    this.controls.addAll(controlList);
  }

  public boolean containsControl(final String controlId) {
    return getControls().stream().anyMatch(c -> c.equals(controlId));
  }

  @Override
  public String getId() {
    return this.id;
  }

  public String getDescription() {
    return description;
  }

  public void setDescription(String description) {
    this.description = description;
  }

  public void addCondition(Condition condition) {
    this.conditions.add(condition);
  }

  public List<Condition> getConditions() {
    return this.conditions;
  }

  public void setConditions(List<Condition> conditions) {
    this.conditions = conditions;
  }

  @Override
  public String toString() {
    return new ToStringBuilder(this)
        .append("conditions", conditions)
        .append("active", active)
        .append("name", name)
        .append("description", description)
        .append("icon", icon)
        .append("controls", controls)
        .append("id", id)
        .toString();
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    Rule rule = (Rule) o;

    return new EqualsBuilder()
        .append(active, rule.active)
        .append(new ArrayList<>(conditions), new ArrayList<>(rule.conditions))
        .append(name, rule.name)
        .append(description, rule.description)
        .append(icon, rule.icon)
        .append(new ArrayList<>(controls), new ArrayList<>(rule.controls))
        .append(id, rule.id)
        .isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37)
        .append(conditions)
        .append(active)
        .append(name)
        .append(description)
        .append(icon)
        .append(controls)
        .append(id)
        .toHashCode();
  }
}
