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
import io.clouditor.data_access_layer.PersistentObject;
import io.clouditor.discovery.AssetService;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Objects;
import java.util.stream.Collectors;
import javax.persistence.Column;
import javax.persistence.Embedded;
import javax.persistence.Entity;
import javax.persistence.EnumType;
import javax.persistence.Enumerated;
import javax.persistence.Id;
import javax.persistence.ManyToMany;
import javax.persistence.Table;
import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.glassfish.hk2.api.ServiceLocator;
import org.hibernate.annotations.LazyCollection;
import org.hibernate.annotations.LazyCollectionOption;

@Entity(name = "control")
@Table(name = "control")
public class Control implements PersistentObject<String> {

  private static final long serialVersionUID = -6926507274525122348L;
  /**
   * The rules associated with this control. This is actually redundant a little bit since its
   * already stored in the {@link Rule}.
   */
  @ManyToMany
  @LazyCollection(LazyCollectionOption.FALSE)
  private List<Rule> rules = new ArrayList<>();

  /** The last evaluation results */
  @ManyToMany
  @LazyCollection(LazyCollectionOption.FALSE)
  private final List<EvaluationResult> results = new ArrayList<>();

  /** The id of the control this objective is referring to, i.e. a CCM control id. */
  @JsonProperty
  @Id
  @Column(name = "control_id", nullable = false)
  private String controlId;

  /** A short description. */
  @JsonProperty
  @Column(name = "control_description", length = 65535)
  private String description;

  /** Is this control ok or not? By default we start in the NOT_EVALUATED state. */
  @JsonProperty
  @Enumerated(EnumType.ORDINAL)
  @Column(name = "fulfillment_value")
  private Fulfillment fulfilled = Fulfillment.NOT_EVALUATED;

  @Embedded private Domain domain;

  @JsonProperty
  @Column(name = "control_name")
  private String name;

  /** Describes, whether the control can be automated or not. */
  @JsonProperty
  @Column(name = "automated")
  private boolean automated;

  /** Is the control actively monitored? */
  @JsonProperty
  @Column(name = "active")
  private boolean active = false;

  @JsonProperty
  @Column(name = "violations")
  private int violations = 0;

  @Override
  public String getId() {
    return this.controlId;
  }

  public void evaluate(ServiceLocator locator) {
    // clear old results
    this.results.clear();

    if (this.rules.isEmpty()) {
      this.fulfilled = Fulfillment.NOT_EVALUATED;

      return;
    }

    this.fulfilled = Fulfillment.GOOD;

    // retrieve assets that belong to a rule within the control
    for (var rule : this.rules) {
      // TODO: use the new function in RulesService#get
      var assets = locator.getService(AssetService.class).getAssetsWithType(rule.getAssetType());

      for (var asset : assets) {
        this.results.addAll(
            asset.getEvaluationResults().stream()
                .filter(result -> Objects.equals(result.getRule().getId(), rule.getId()))
                .collect(Collectors.toList()));
      }

      if (this.results.stream().anyMatch(EvaluationResult::hasFailedConditions)) {
        this.fulfilled = Fulfillment.WARNING;
      }
    }

    // we should handle this better
    if (this.results.isEmpty()) {
      this.fulfilled = Fulfillment.NOT_EVALUATED;
    }

    this.violations = 0;
  }

  public String getControlId() {
    return controlId;
  }

  public void setControlId(String controlId) {
    this.controlId = controlId;
  }

  public Fulfillment getFulfilled() {
    return fulfilled;
  }

  public void setFulfilled(Fulfillment fulfilled) {
    this.fulfilled = fulfilled;
  }

  public void setDomain(Domain domain) {
    this.domain = domain;
  }

  public void setName(String name) {
    this.name = name;
  }

  public void setDescription(String description) {
    this.description = description;
  }

  public boolean isActive() {
    return this.active;
  }

  public void setActive(boolean active) {
    this.active = active;
  }

  public boolean isAutomated() {
    return this.automated;
  }

  public void setAutomated(boolean automated) {
    this.automated = automated;
  }

  public boolean isGood() {
    return this.active && this.fulfilled == Fulfillment.GOOD;
  }

  public boolean hasWarning() {
    return this.active && this.fulfilled == Fulfillment.WARNING;
  }

  public List<Rule> getRules() {
    return rules;
  }

  public void setRules(List<Rule> rules) {
    this.rules = rules;
  }

  public List<EvaluationResult> getResults() {
    return Collections.unmodifiableList(this.results);
  }

  public void setResults(final EvaluationResult... evaluationResults) {
    final List<EvaluationResult> evaluationResultList = List.of(evaluationResults);
    this.results.addAll(evaluationResultList);
  }

  public void removeResults(final EvaluationResult... evaluationResults) {
    final List<EvaluationResult> evaluationResultList = List.of(evaluationResults);
    this.results.removeAll(evaluationResultList);
  }

  @Override
  public String toString() {
    return new ToStringBuilder(this)
        .append("rules", rules)
        .append("results", results)
        .append("controlId", controlId)
        .append("description", description)
        .append("fulfilled", fulfilled)
        .append("domain", domain)
        .append("name", name)
        .append("automated", automated)
        .append("active", active)
        .append("violations", violations)
        .toString();
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;

    if (o == null || getClass() != o.getClass()) return false;

    Control control = (Control) o;

    return new EqualsBuilder()
        .append(automated, control.automated)
        .append(active, control.active)
        .append(violations, control.violations)
        .append(new ArrayList<>(rules), new ArrayList<>(control.rules))
        .append(results, control.results)
        .append(controlId, control.controlId)
        .append(description, control.description)
        .append(fulfilled, control.fulfilled)
        .append(domain, control.domain)
        .append(name, control.name)
        .isEquals();
  }

  @Override
  public int hashCode() {
    return new HashCodeBuilder(17, 37)
        .append(rules)
        .append(results)
        .append(controlId)
        .append(description)
        .append(fulfilled)
        .append(domain)
        .append(name)
        .append(automated)
        .append(active)
        .append(violations)
        .toHashCode();
  }

  public Domain getDomain() {
    return this.domain;
  }

  public enum Fulfillment {
    NOT_EVALUATED,
    WARNING,
    GOOD
  }
}
