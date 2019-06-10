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

package io.clouditor.rest;

import static io.clouditor.auth.AuthenticationService.ROLE_USER;

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
  @Path("assets/{assetType}")
  public Set<Rule> getRules(@PathParam("assetType") String assetType) {
    var rules = this.ruleService.getRules().get(assetType);

    if (rules == null) {
      rules = new HashSet<>();
    }

    return rules;
  }

  @Produces(MediaType.APPLICATION_JSON)
  @GET
  @Path("{ruleId}")
  public RuleEvaluation get(@PathParam("ruleId") String ruleId) {
    var rule = this.ruleService.getWithId(ruleId);

    if (rule == null) {
      throw new NotFoundException();
    }

    return this.ruleService.getStatus(rule);
  }
}
