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

import java.io.BufferedReader;
import java.io.IOException;
import java.io.StringReader;
import java.net.URI;
import java.net.URISyntaxException;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.text.PDFTextStripper;
import org.apache.pdfbox.text.PDFTextStripperByArea;

public class AzureCISBenchmarkImporter extends CertificationImporter {
  private static final String BENCHMARK_PREFIX = "Ensure that ";
  private static final String SECTION_REGEX = "(\\d+[\\.\\d+]+) " + BENCHMARK_PREFIX + "(.*)";

  // TODO: also retrieve this from PDF:
  private static Map<String, String> domains = new HashMap<>();

  static {
    domains.put("Azure 1", "Identity and Access Management");
    domains.put("Azure 2", "Security Center");
    domains.put("Azure 3", "Storage Accounts");
    domains.put("Azure 4.1", "SQL Servers");
    domains.put("Azure 4.2", "SQL Databases");
    domains.put("Azure 5", "Logging and Monitoring");
    domains.put("Azure 6", "Networking");
    domains.put("Azure 7", "Virtual Machines");
    domains.put("Azure 8", "Other Security Considerations");
  }

  public String getName() {
    return "CIS Microsoft Azure Foundations Benchmark";
  }

  public Certification load() {
    var certification = new Certification();

    certification.setId(this.getName());
    certification.setPublisher("Center for Internet Security");
    certification.setDescription(
        "The CIS Microsoft Azure Foundations Benchmark contains a set of technical controls to assure the security of a Azure-based Cloud workload.");
    certification.setWebsite("https://www.cisecurity.org/benchmark/azure/");

    certification.setControls(importCIS());

    return certification;
  }

  private List<Control> importCIS() {
    try {
      System.setProperty("sun.java2d.cmm", "sun.java2d.cmm.kcms.KcmsServiceProvider");

      var url =
          "https://azure.microsoft.com/mediahandler/files/resourcefiles/cis-microsoft-azure-foundations-security-benchmark/CIS_Microsoft_Azure_Foundations_Benchmark_v1.0.0.pdf";

      LOGGER.info("Fetching Azure CIS Benchmark from {}...", url);

      var document = PDDocument.load(URI.create(url).toURL().openStream());

      var stripper = new PDFTextStripperByArea();
      stripper.setSortByPosition(true);

      var tStripper = new PDFTextStripper();

      var pdfFileInText = tStripper.getText(document);

      return processText(pdfFileInText);
    } catch (IOException | URISyntaxException e) {
      LOGGER.error("An error occurred while importing the CIS benchmark: {}", e.getMessage());
    }

    return Collections.emptyList();
  }

  private List<Control> processText(String pdfFileInText) throws IOException, URISyntaxException {
    var reader = new BufferedReader(new StringReader(pdfFileInText));

    // skip forward to "Recommendations"
    skipUntilAfter(reader, "Recommendations");

    // lets controls
    return process(reader);
  }

  private void skipUntilAfter(BufferedReader reader, String needle) throws IOException {
    while (true) {
      var line = reader.readLine();
      if (line.trim().matches(needle)) {
        return;
      }
    }
  }

  private String readUntilBeforeSection(BufferedReader reader) throws IOException {
    return this.readUntilBefore(reader, "([0-9]+\\.)+[0-9]+ .*");
  }

  private String readUntilBefore(BufferedReader reader, String needle) throws IOException {
    var buffer = new StringBuilder();
    while (true) {
      // mark the beginning of the line, so we can return to it
      reader.mark(2048);

      // read the line
      var line = reader.readLine();

      if (line == null) {
        return null;
      }

      // see if it matches a section
      if (line.trim().matches(needle)) {
        // reset to mark
        reader.reset();

        // return buffer
        return buffer.toString().trim();
      }

      buffer.append(line);
    }
  }

  private List<Control> process(BufferedReader reader) throws IOException, URISyntaxException {
    List<Control> controls = new ArrayList<>();

    // first line is the name (TODO: Skip the number)
    var line = reader.readLine();

    // TODO: only works for the first one currently
    // Domain domain = new Domain();
    // domain.setName(line);

    // next lines are the description, until the next sub-chapter
    var description = readUntilBeforeSection(reader);

    // domain.setDescription(description);

    // now lets read the individual controls
    // TODO: sometimes we have a sub-group
    while (true) {
      var control = processControl(reader);

      if (control == null) {
        break;
      }

      // check, if one of the domains matches the beginning of the control id, not really very
      // efficient but works
      domains.forEach(
          (key, value) -> {
            if (control.getControlId().startsWith(key)) {
              Domain domain = new Domain();
              domain.setName(value);
              control.setDomain(domain);
            }
          });

      // read, but ignore -> skip
      readUntilBeforeSection(reader);

      controls.add(control);
    }

    return controls;
  }

  private Control processControl(BufferedReader reader) throws IOException {
    var control = new Control();

    var nameAndId = this.readUntilBefore(reader, "Profile Applicability:");

    if (nameAndId == null) {
      return null;
    }

    var m = Pattern.compile(SECTION_REGEX).matcher(nameAndId);

    if (!m.find()) {
      return null;
    }

    control.setControlId("Azure " + nameAndId.substring(m.start(1), m.end(1)));

    var name = BENCHMARK_PREFIX + nameAndId.substring(m.start(2), m.end(2));
    // remove (Scored) and (Not Scored)
    name = name.replace("(Scored)", "").trim();
    name = name.replace("(Not Scored)", "").trim();

    control.setName(name);

    this.skipUntilAfter(reader, "Description:");
    var description = this.readUntilBefore(reader, "Rationale:");

    control.setDescription(description);

    this.skipUntilAfter(reader, "CIS Controls:");

    return control;
  }
}
