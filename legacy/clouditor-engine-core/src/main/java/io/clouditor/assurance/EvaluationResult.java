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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.assurance.ccl.Condition;
import io.clouditor.data_access_layer.PersistentObject;
import io.clouditor.discovery.AssetProperties;
import java.time.Instant;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;
import javax.persistence.CollectionTable;
import javax.persistence.Column;
import javax.persistence.ElementCollection;
import javax.persistence.Embedded;
import javax.persistence.Entity;
import javax.persistence.Id;
import javax.persistence.JoinColumn;
import javax.persistence.ManyToOne;
import javax.persistence.MapKeyColumn;
import javax.persistence.Table;
import javax.validation.constraints.NotNull;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.hibernate.annotations.LazyCollection;
import org.hibernate.annotations.LazyCollectionOption;
import org.jetbrains.annotations.Nullable;

@Entity(name = "evaluation_result")
@Table(name = "evaluation_result")
public class EvaluationResult implements PersistentObject<String> {

  private static final long serialVersionUID = 7255742076812915308L;

  @NotNull @Id private String id;

  @Column(nullable = false)
  private Instant timestamp;

  /** The rule according to which this was evaluated. */
  @NotNull
  @JsonProperty
  @ManyToOne
  @JoinColumn(name = "rule_id", nullable = false)
  private final Rule rule;

  @JsonProperty
  @CollectionTable(name = "asset_properties", joinColumns = @JoinColumn(name = "key_id"))
  @MapKeyColumn(name = "mapKey")
  @Column(name = "asset_properties")
  private final AssetProperties evaluatedProperties;

  @Embedded
  @ElementCollection(targetClass = Condition.class)
  @LazyCollection(LazyCollectionOption.FALSE)
  private List<Condition> failedConditions = new ArrayList<>();

  public EvaluationResult(Rule rule, AssetProperties evaluatedProperties) {
    this(rule, evaluatedProperties, null);
  }

  @JsonCreator
  public EvaluationResult(
      @JsonProperty("rule") Rule rule,
      @JsonProperty("evaluatedProperties") AssetProperties evaluatedProperties,
      @Nullable @JsonProperty("id") String id) {
    if (id == null) {
      this.id = UUID.randomUUID().toString();
    }

    this.timestamp = Instant.now();
    this.rule = rule;
    this.evaluatedProperties = evaluatedProperties;
  }

  public void setFailedConditions(List<Condition> failedConditions) {
    this.failedConditions = failedConditions;
  }

  public List<Condition> getFailedConditions() {
    return this.failedConditions;
  }

  public boolean isOk() {
    return this.failedConditions.isEmpty();
  }

  public Rule getRule() {
    return rule;
  }

  public boolean hasFailedConditions() {
    return !this.failedConditions.isEmpty();
  }

  public Instant getTimestamp() {
    return timestamp;
  }

  public void setTimestamp(final Instant timestamp) {
    this.timestamp = Instant.ofEpochMilli(timestamp.toEpochMilli());
  }

  @Override
  public String getId() {
    return this.id;
  }

  @Override
  public String toString() {
    return new ToStringBuilder(this)
        .append("timeStamp", timestamp)
        .append("rule", rule)
        .append("evaluatedProperties", evaluatedProperties)
        .append("failedConditions", failedConditions)
        .toString();
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    EvaluationResult that = (EvaluationResult) o;

    return new EqualsBuilder()
        .append(timestamp, that.timestamp)
        .append(rule, that.rule)
        .append(evaluatedProperties, that.evaluatedProperties)
        .append(new ArrayList<>(failedConditions), new ArrayList<>(that.failedConditions))
        .isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37)
        .append(timestamp)
        .append(rule)
        .append(evaluatedProperties)
        .append(failedConditions)
        .toHashCode();
  }
}
