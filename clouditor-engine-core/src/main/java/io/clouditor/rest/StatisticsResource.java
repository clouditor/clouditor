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

import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.Engine;
import io.clouditor.assurance.Certification;
import io.clouditor.assurance.CertificationService;
import io.clouditor.assurance.Control;
import io.clouditor.auth.UserContext;
import io.clouditor.discovery.DiscoveryService;
import io.clouditor.discovery.Scan;
import io.swagger.annotations.Api;
import io.swagger.annotations.Authorization;
import java.util.List;
import java.util.stream.Collectors;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

@Path("statistics")
@Api(authorizations = @Authorization(value = "token"))
@RolesAllowed(UserContext.ROLE_USERS)
public class StatisticsResource {

  @Inject private Engine engine;

  @Inject private DiscoveryService discoveryService;

  @Produces(MediaType.APPLICATION_JSON)
  @GET
  public Statistics getStatistics() {
    var stats = new Statistics();
    stats.numAssets = 0;
    stats.assetCoverage = 0;

    var certifications = this.engine.getService(CertificationService.class).getCertifications();

    stats.numActiveScanners =
        (int) discoveryService.getScans().values().stream().map(Scan::isEnabled).count();
    stats.numCertifications = certifications.size();
    stats.numPassedControls =
        certifications.values().stream()
            .map(Certification::getControls)
            .map(controls -> controls.stream().filter(Control::isGood).collect(Collectors.toList()))
            .mapToInt(List::size)
            .sum();
    stats.numFailedControls =
        certifications.values().stream()
            .map(Certification::getControls)
            .map(
                controls ->
                    controls.stream().filter(Control::hasWarning).collect(Collectors.toList()))
            .mapToInt(List::size)
            .sum();

    return stats;
  }

  static class Statistics {

    @JsonProperty int numAssets;
    @JsonProperty double assetCoverage;
    @JsonProperty int numActiveScanners;
    @JsonProperty int numCertifications;
    @JsonProperty int numFailedControls;
    @JsonProperty int numPassedControls;
  }
}
