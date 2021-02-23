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

import io.clouditor.data_access_layer.HibernatePersistence;
import io.clouditor.data_access_layer.PersistenceManager;
import io.clouditor.events.CertificationSubscriber;
import io.clouditor.events.SubscriptionManager;
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Modifier;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.SubmissionPublisher;
import java.util.stream.Collectors;
import javax.inject.Inject;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.glassfish.hk2.api.ServiceLocator;
import org.jvnet.hk2.annotations.Service;
import org.reflections.Reflections;
import org.reflections.scanners.SubTypesScanner;
import org.reflections.util.ClasspathHelper;
import org.reflections.util.ConfigurationBuilder;

@Service
public class CertificationService {

  private static final Logger LOGGER = LogManager.getLogger();

  private static final Reflections REFLECTIONS_SUBTYPE_SCANNER =
      new Reflections(
          new ConfigurationBuilder()
              .addUrls(
                  ClasspathHelper.forPackage(CertificationService.class.getPackage().getName()))
              .setScanners(new SubTypesScanner()));

  private final Map<String, CertificationImporter> importers = new HashMap<>();

  private final SubmissionPublisher<Certification> certificationPublisher =
      new SubmissionPublisher<>();

  /** The certifications objective */
  private final Map<String, Certification> certifications = new HashMap<>();

  private final ServiceLocator locator;

  @Inject
  public CertificationService(ServiceLocator locator) {
    this.locator = locator;
  }

  public void loadImporters() {
    try {
      // load importers via reflection
      List<Class<? extends CertificationImporter>> importerClasses =
          REFLECTIONS_SUBTYPE_SCANNER.getSubTypesOf(CertificationImporter.class).stream()
              .filter(
                  clazz -> !Modifier.isAbstract(clazz.getModifiers()) && !clazz.isAnonymousClass())
              .collect(Collectors.toList());

      for (Class<? extends CertificationImporter> clazz : importerClasses) {
        CertificationImporter importer = clazz.getDeclaredConstructor().newInstance();
        this.importers.put(importer.getName(), importer);
      }
    } catch (NoSuchMethodException
        | InstantiationException
        | IllegalAccessException
        | InvocationTargetException ex) {
      LOGGER.error("Could not load certification importers: {}", ex.getMessage());
    }
  }

  public Certification load(String certificationId) {
    CertificationImporter importer = this.importers.get(certificationId);

    if (importer == null) {
      return null;
    }

    var ruleService = this.locator.getService(RuleService.class);

    var certification = importer.load();

    // loop through all controls and
    // a) look for associated rules
    // b) start monitoring controls which have associated rules

    for (var control : certification.getControls()) {
      // find associated rules
      control.setRules(
          ruleService.getRulesForControl(certificationId + "/" + control.getControlId()));

      if (!control.getRules().isEmpty()) {
        control.setAutomated(true);

        startMonitoring(control);
      }
    }

    return certification;
  }

  public void startMonitoring(Control control) {
    if (!control.isAutomated()) {
      LOGGER.error("Non-automated control {} cannot be enabled. Ignoring.", control.getControlId());
      return;
    }

    control.setActive(true);

    // previous results could already be there, try to update the control
    this.updateCertification(Collections.singletonList(control.getControlId()));
  }

  public void stopMonitoring(Control control) {
    // TODO: actually stop all associated jobs
    control.setActive(false);
  }

  public Map<String, CertificationImporter> getImporters() {
    return this.importers;
  }

  public void loadCertifications() {
    LOGGER.info("Loading certifications and controls...");

    for (var certification : new HibernatePersistence().listAll(Certification.class)) {
      loadCertification(certification);
    }

    LOGGER.info("Clouditor Engine loaded {} certifications", this.certifications.size());
  }

  private void loadCertification(Certification certification) {
    this.certifications.put(certification.getId(), certification);
  }

  public void modifyCertification(Certification certification) {
    // load it
    this.loadCertification(certification);

    // make sure, controls are active
    /*for (Control control : certification.getControls()) {
      this.startMonitoring(control);
    }*/

    // update
    this.updateCertification();
  }

  public Map<String, Certification> getCertifications() {
    return certifications;
  }

  public void updateCertification() {
    this.updateCertification(Collections.emptyList());
  }

  public void updateCertification(List<String> controlIds) {
    final PersistenceManager persistenceManager = new HibernatePersistence();
    for (var certification : this.certifications.values()) {
      for (var control : certification.getControls()) {
        // only update certain controls, or all if the list is empty
        if (!controlIds.isEmpty() && !controlIds.contains(control.getControlId())) {
          continue;
        }

        // skip non-active controls
        if (!control.isActive()) {
          continue;
        }

        control.evaluate(this.locator);

        control.getResults().forEach(persistenceManager::saveOrUpdate);
        persistenceManager.saveOrUpdate(control);

        LOGGER.info(
            "Evaluated fulfillment of control {} as {}",
            control.getControlId(),
            control.getFulfilled());

        LOGGER.debug("Control {} is now {}", control.getControlId(), control);
      }

      persistenceManager.saveOrUpdate(certification);
      this.certificationPublisher.submit(certification);
    }
  }

  public void loadSubscribers() {
    // load Certification subscribers
    for (var subscriber :
        SubscriptionManager.getInstance().loadSubscribers(CertificationSubscriber.class)) {
      this.certificationPublisher.subscribe(subscriber);
    }
  }

  // For testing purposes (remove certifications when previous test suites added some)
  public void removeAllCertifications() {
    for (Certification certification : getCertifications().values()) {
      new HibernatePersistence().delete(certification);
      certifications.remove(certification.getId());
    }
  }
}
