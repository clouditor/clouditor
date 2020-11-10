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

import com.fasterxml.jackson.annotation.JsonProperty;
import io.clouditor.Engine;
import io.clouditor.assurance.Certification;
import io.clouditor.assurance.CertificationImporter;
import io.clouditor.assurance.CertificationService;
import io.clouditor.assurance.Control;
import java.util.Map;
import javax.annotation.security.RolesAllowed;
import javax.inject.Inject;
import javax.ws.rs.GET;
import javax.ws.rs.NotFoundException;
import javax.ws.rs.POST;
import javax.ws.rs.Path;
import javax.ws.rs.PathParam;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * A resource end-point for {@link Certification}.
 *
 * @author Christian Banse
 */
@Path("certification")
@RolesAllowed(ROLE_USER)
public class CertificationResource {

  private static final Logger LOGGER = LoggerFactory.getLogger(CertificationResource.class);

  private final Engine engine;
  private final CertificationService service;

  /**
   * Constructs a new resource.
   *
   * @param engine the Clouditor Engine
   */
  @Inject
  public CertificationResource(Engine engine, CertificationService service) {
    this.engine = engine;
    this.service = service;
  }

  @POST
  @Path("{certificationId}/{controlId}/status")
  public void modifyControlStatus(
      @PathParam("certificationId") String certificationId,
      @PathParam("controlId") String controlId,
      ControlStatusRequest request) {
    certificationId = sanitize(certificationId);
    controlId = sanitize(controlId);

    var certification = this.getCertification(certificationId);

    if (certification == null) {
      throw new NotFoundException();
    }

    String finalControlId = controlId;
    var first =
        certification.getControls().stream()
            .filter(control -> control.getControlId().equals(finalControlId))
            .findFirst();

    if (first.isEmpty()) {
      throw new NotFoundException();
    }

    var control = first.get();

    if (!request.status && control.isActive()) {
      this.service.stopMonitoring(control);
    } else if (request.status && !control.isActive()) {
      this.service.startMonitoring(control);
    }
  }

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  public Map<String, Certification> getCertifications() {
    return this.service.getCertifications();
  }

  @GET
  @Produces(MediaType.APPLICATION_JSON)
  @Path("{id}/")
  public Certification getCertification(@PathParam(value = "id") String certificationId) {
    certificationId = sanitize(certificationId);

    var certifications = this.service.getCertifications();

    var certification = certifications.get(certificationId);

    if (certification == null) {
      throw new NotFoundException();
    }

    return certification;
  }

  @GET
  @Path("{certificationId}/{controlId}")
  public Control getControl(
      @PathParam("certificationId") String certificationId,
      @PathParam("controlId") String controlId) {
    certificationId = sanitize(certificationId);
    controlId = sanitize(controlId);

    var certification = getCertification(certificationId);

    String finalControlId = controlId;
    var any =
        certification.getControls().stream()
            .filter(control -> control.getControlId().equals(finalControlId))
            .findAny();

    if (any.isEmpty()) {
      throw new NotFoundException();
    }

    return any.get();
  }

  @POST
  @Path("import/{certificationId}")
  public void importCertification(@PathParam("certificationId") String certificationId) {
    certificationId = sanitize(certificationId);

    var certification = this.service.load(certificationId);

    if (certification == null) {
      throw new NotFoundException();
    }

    this.service.modifyCertification(certification);
  }

  @GET
  @Path("importers")
  public Map<String, CertificationImporter> getImporters() {
    return this.service.getImporters();
  }

  public static class ControlStatusRequest {

    @JsonProperty private boolean status;

    public ControlStatusRequest() {}
  }
}
