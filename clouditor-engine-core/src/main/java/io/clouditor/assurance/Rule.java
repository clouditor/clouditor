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
import io.clouditor.discovery.Asset;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;
import javax.persistence.*;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;

@Entity(name = "rule")
@Table(name = "rule")
public class Rule implements Serializable {

  @JsonDeserialize(using = CCLDeserializer.class)
  @JsonSerialize(using = CCLSerializer.class)
  @ManyToOne(cascade = CascadeType.ALL)
  @JoinTable(
      name = "condition_to_rule",
      joinColumns = @JoinColumn(name = "rule_id", referencedColumnName = "rule_id"),
      inverseJoinColumns = {
        @JoinColumn(name = "source", referencedColumnName = "source"),
        @JoinColumn(name = "type_value", referencedColumnName = "type_value"),
      })
  private Condition condition;

  @JsonDeserialize(contentUsing = CCLDeserializer.class)
  @JsonSerialize(contentUsing = CCLSerializer.class)
  @ManyToMany
  @JoinTable(
      name = "condition_to_many_rules",
      joinColumns = @JoinColumn(name = "rule_id", referencedColumnName = "rule_id"),
      inverseJoinColumns = {
        @JoinColumn(name = "source", referencedColumnName = "source"),
        @JoinColumn(name = "type_value", referencedColumnName = "type_value"),
      })
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
  @ManyToMany
  @JoinTable(
      name = "rule_to_control",
      joinColumns = @JoinColumn(name = "rule_id", referencedColumnName = "rule_id"),
      inverseJoinColumns = @JoinColumn(name = "control_id", referencedColumnName = "control_id"))
  private final List<Control> controls = new ArrayList<>();

  @JsonProperty
  @Id
  @Column(name = "rule_id")
  private String id;

  /**
   * TOTO: what does it do?
   *
   * @param asset
   * @return
   */
  public boolean evaluateApplicability(Asset asset) {

    if (this.condition != null) {
      if (this.condition.getAssetType() instanceof FilteredAssetType) {
        return this.condition.getAssetType().evaluate(asset.getProperties());
      }
    } else if (this.conditions != null) {
      if (this.conditions.get(0).getAssetType() instanceof FilteredAssetType) {
        return this.conditions.get(0).getAssetType().evaluate(asset.getProperties());
      }
    }
    return true;
  }

  public EvaluationResult evaluate(Asset asset) {
    var eval = new EvaluationResult(this, asset.getProperties());

    if (!this.conditions.isEmpty()) {
      // get those conditions with evaluate as false
      eval.setFailedConditions(
          this.conditions.stream()
              .filter(c -> !c.evaluate(asset.getProperties()))
              .collect(Collectors.toList()));
    } else {
      if (!this.condition.evaluate(asset.getProperties())) {
        eval.setFailedConditions(Collections.singletonList(this.condition));
      }
    }

    return eval;
  }

  public void setCondition(Condition condition) {
    this.condition = condition;
  }

  public String getAssetType() {
    // single condition
    if (this.condition != null && this.condition.getAssetType() != null) {
      return this.condition.getAssetType().getValue();
    }

    // multiple conditions
    if (this.conditions != null
        && !this.conditions.isEmpty()
        && this.conditions.get(0).getAssetType() != null) {
      // take the first one
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

  public Condition getCondition() {
    return this.condition;
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

  public List<Control> getControls() {
    return this.controls;
  }

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

  public boolean containsControl(final String controlId) {
    return getControls().stream().anyMatch(control -> control.getControlId().equals(controlId));
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    Rule rule = (Rule) o;

    return new EqualsBuilder()
        .append(isActive(), rule.isActive())
        .append(getCondition(), rule.getCondition())
        .append(getConditions(), rule.getConditions())
        .append(getName(), rule.getName())
        .append(getDescription(), rule.getDescription())
        .append(icon, rule.icon)
        .append(getControls(), rule.getControls())
        .append(getId(), rule.getId())
        .isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37)
        .append(getCondition())
        .append(getConditions())
        .append(isActive())
        .append(getName())
        .append(getDescription())
        .append(icon)
        .append(getControls())
        .append(getId())
        .toHashCode();
  }
}
