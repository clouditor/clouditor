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
import io.clouditor.discovery.Asset;
import java.util.HashMap;
import java.util.Map;
import javax.validation.constraints.NotNull;

public class RuleEvaluation {

  /** The rule. */
  @JsonProperty @NotNull private Rule rule;

  @JsonProperty private final Map<String, Boolean> compliance = new HashMap<>();

  @JsonCreator
  public RuleEvaluation(@JsonProperty("rule") Rule rule) {
    this.rule = rule;
  }

  void addCompliant(@NotNull Asset asset) {
    this.compliance.put(asset.getId(), true);
  }

  void addNonCompliant(@NotNull Asset asset) {
    this.compliance.put(asset.getId(), false);
  }

  public boolean isOk() {
    return this.compliance.values().stream().allMatch(status -> status);
  }

  @JsonProperty
  public long getNumberOfCompliant() {
    return this.compliance.values().stream().filter(status -> status).count();
  }

  @JsonProperty
  public long getNumberOfNonCompliant() {
    return this.compliance.values().stream().filter(status -> !status).count();
  }
}
