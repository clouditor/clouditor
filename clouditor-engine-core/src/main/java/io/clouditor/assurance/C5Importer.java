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

import java.io.IOException;
import java.net.URL;
import java.util.ArrayList;
import org.apache.poi.xssf.usermodel.XSSFWorkbook;

public class C5Importer extends CertificationImporter {

  @Override
  public Certification load() {
    var url =
        "https://www.bsi.bund.de/SharedDocs/Downloads/EN/BSI/CloudComputing/ComplianceControlsCatalogue/ComplianceControlsCatalogue_tables_editable.xlsx?__blob=publicationFile&v=8";

    LOGGER.info("Fetching BSI C5 from {}...", url);

    var controls = new ArrayList<Control>();

    try (var workbook = new XSSFWorkbook(new URL(url).openStream())) {
      var sheet = workbook.getSheetAt(1);

      int max = 115;

      // starts at row 2
      for (int i = 1; i < max; i++) {
        var control = new Control();

        control.setDomain(new Domain(sheet.getRow(i).getCell(0).toString()));
        control.setControlId(sheet.getRow(i).getCell(1).toString().trim());
        control.setName(sheet.getRow(i).getCell(2).toString());
        control.setDescription(sheet.getRow(i).getCell(3).toString());

        controls.add(control);
      }
    } catch (IOException e) {
      LOGGER.error("Could not parse BSI C5 from xlsx: {}", e.getMessage());
    }

    Certification certification = new Certification();
    certification.setId(this.getName());
    certification.setPublisher("BSI");
    certification.setDescription(
        "The Cloud Computing Compliance Controls Catalogue (abbreviated \"C5\") is intended primarily for professional cloud service providers, their auditors and customers of the cloud service providers. It is defined which requirements (also referred to as controls in this context) the cloud providers have to comply with or which minimum requirements the cloud providers should be obliged to meet.\n");
    certification.setWebsite(
        "https://www.bsi.bund.de/EN/Topics/CloudComputing/Compliance_Controls_Catalogue/Compliance_Controls_Catalogue_node.html");

    certification.setControls(controls);

    return certification;
  }

  @Override
  public String getName() {
    return "BSI C5";
  }
}
