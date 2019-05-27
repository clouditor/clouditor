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

package io.clouditor.discovery;

import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.assurance.EvaluationResult;
import java.util.ArrayList;
import java.util.List;
import org.apache.commons.lang3.builder.ToStringBuilder;
import org.apache.commons.lang3.builder.ToStringStyle;

public class Asset {

  private List<EvaluationResult> evaluationResults = new ArrayList<>();

  private AssetProperties properties = new AssetProperties();

  private String id;
  private String name;
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
}
