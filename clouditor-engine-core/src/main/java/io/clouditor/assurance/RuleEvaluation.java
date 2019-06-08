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

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.discovery.Asset;
import java.util.HashMap;
import java.util.Map;
import javax.validation.constraints.NotNull;

public class RuleEvaluation {

  /** The rule. */
  @JsonProperty @NotNull private Rule rule;

  @JsonProperty private Map<String, Boolean> compliance = new HashMap<>();

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
