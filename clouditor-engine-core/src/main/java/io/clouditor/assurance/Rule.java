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

package io.clouditor.assurance;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;
import io.clouditor.assurance.ccl.*;
import io.clouditor.discovery.Asset;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

public class Rule {

  @JsonDeserialize(using = CCLDeserializer.class)
  @JsonSerialize(using = CCLSerializer.class)
  private Condition condition;

  @JsonDeserialize(contentUsing = CCLDeserializer.class)
  @JsonSerialize(contentUsing = CCLSerializer.class)
  private List<Condition> conditions = new ArrayList<>();

  private boolean active;

  @JsonProperty private String name;
  @JsonProperty private String description;
  @JsonProperty private String icon = "far fa-file-alt";
  @JsonProperty private List<String> controls = new ArrayList<>();
  @JsonProperty private String id;

  public boolean evaluateApplicability(Asset asset) {

    if (this.condition != null) {
      if (this.condition.getAssetType() instanceof GroupedAsset) {
        return this.condition.getAssetType().evaluate(asset.getProperties());
      }
    } else if (this.conditions != null) {
      if (this.conditions.get(0).getAssetType() instanceof GroupedAsset) {
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
      return this.condition.getAssetType().getField();
    }

    // multiple conditions
    if (this.conditions != null
        && !this.conditions.isEmpty()
        && this.conditions.get(0).getAssetType() != null) {
      // take the first one
      return this.conditions.get(0).getAssetType().getField();
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

  public List<String> getControls() {
    return this.controls;
  }

  public void setControls(List<String> controls) {
    this.controls = controls;
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
}
