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

import com.vladmihalcea.hibernate.type.json.JsonBinaryType;  // TODO: check for licence: not gpl

import io.clouditor.assurance.ccl.Condition;
import io.clouditor.discovery.AssetProperties;
import org.hibernate.annotations.Type;
import org.hibernate.annotations.TypeDef;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import javax.persistence.*;
import javax.validation.constraints.NotNull;

@Entity(name = "evaluation_result")
@Table(name = "evaluation_result")
@TypeDef(name = "jsonb", typeClass = JsonBinaryType.class)
public class EvaluationResult implements Serializable {

  private static final long serialVersionUID = 7255742076812915308L;

  @Id
  @Column(name = "time_stamp")
  private final String timeStamp = new Date().toString();

  /** The rule according to which this was evaluated. */
  @NotNull
  @JsonProperty
  @Id
  @ManyToOne
  @JoinColumn(name = "rule_id")
  private final Rule rule;

  @JsonProperty
  @Type(type = "jsonb")
  @Column(name = "asset_properties")
  private final AssetProperties evaluatedProperties;

  @ManyToMany
  @JoinTable(name="condition_to_evaluation_result",
          joinColumns = {
                  @JoinColumn(name = "time_stamp", referencedColumnName = "time_stamp"),
                  @JoinColumn(name = "rule_id", referencedColumnName = "rule_id")
          },
          inverseJoinColumns = {
                  @JoinColumn(name = "source", referencedColumnName = "source"),
                  @JoinColumn(name = "type_value", referencedColumnName = "type_value"),
          }
  )
  private List<Condition> failedConditions = new ArrayList<>();

  @JsonCreator
  public EvaluationResult(
      @JsonProperty("rule") Rule rule,
      @JsonProperty("evaluatedProperties") AssetProperties evaluatedProperties) {
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
}
