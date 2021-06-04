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

package io.clouditor.rest;

import static io.clouditor.auth.AuthenticationService.ROLE_USER;
import static io.clouditor.rest.AbstractAPI.sanitize;

import io.clouditor.assurance.Rule;
import io.clouditor.assurance.RuleEvaluation;
import io.clouditor.assurance.RuleService;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.GET;
import javax.ws.rs.NotFoundException;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

@Path("rules")
@RolesAllowed(ROLE_USER)
public class RulesResource {

  private final RuleService ruleService;

  @Inject
  public RulesResource(RuleService ruleService) {
    this.ruleService = ruleService;
  }

  @Produces(MediaType.APPLICATION_JSON)
  @GET
  public Map<String, Set<Rule>> getRules() {
    return this.ruleService.getRules();
  }

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  @Path("assets/{assetType}")
  public Set<Rule> getRules(@PathParam("assetType") String assetType) {
    assetType = sanitize(assetType);

    var rules = this.ruleService.getRules().get(assetType);

    if (rules == null) {
      rules = new HashSet<>();
    }

    return rules;
  }

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  @Path("{ruleId}")
  public RuleEvaluation get(@PathParam("ruleId") String ruleId) {
    ruleId = sanitize(ruleId);

    var rule = this.ruleService.getWithId(ruleId);

    if (rule == null) {
      throw new NotFoundException();
    }

    return this.ruleService.getStatus(rule);
  }
}
