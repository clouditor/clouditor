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

  private List<Control> process(BufferedReader reader) throws IOException {
    List<Control> controls = new ArrayList<>();

    /*
        // first line is the name (TODO: Skip the number)
        var line = reader.readLine();

        // TODO: only works for the first one currently
        Domain domain = new Domain();
        domain.setName(line);

        // next lines are the description, until the next sub-chapter
        var description = readUntilBeforeSection(reader);

        domain.setDescription(description);
    */
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
