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

import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.Engine;
import io.clouditor.assurance.Certification;
import io.clouditor.assurance.CertificationService;
import io.clouditor.assurance.Control;
import io.clouditor.discovery.DiscoveryService;
import io.clouditor.discovery.Scan;
import java.util.List;
import java.util.stream.Collectors;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.GET;
import javax.ws.rs.Path;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

@Path("statistics")
@RolesAllowed(ROLE_USER)
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
