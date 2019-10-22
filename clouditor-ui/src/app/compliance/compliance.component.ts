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

import { Component, OnInit } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { take } from 'rxjs/operators';

import { Certification, Fulfillment } from '../certification';
import { CertificationService } from '../certification.service';

@Component({
  selector: 'clouditor-compliance',
  templateUrl: './compliance.component.html',
  styleUrls: ['./compliance.component.scss']
})
export class ComplianceComponent implements OnInit {
  certifications: Certification[];
  importers: any[];

  importing = new Map<string, boolean>();

  constructor(titleService: Title,
    private certificationService: CertificationService) {
    titleService.setTitle('Compliance');
  }

  ngOnInit() {
    this.updateCertification();
  }

  updateCertification() {
    this.certificationService
      .getCertifications()
      .pipe(take(1))
      .subscribe(certifications => this.certifications = Array.from(certifications.values()));

    this.certificationService.getImporters()
      .pipe(take(1)).subscribe(importers => this.importers = importers);
  }

  onImport(certificationId: string) {
    if (this.importing[certificationId]) {
      return;
    }

    this.importing[certificationId] = true;
    this.certificationService.import(certificationId).subscribe(() => {
      this.updateCertification();
      this.importing[certificationId] = false;
    });
  }

  getInactiveControls(certification: Certification) {
    return certification.controls.filter(control => {
      return !control.active;
    });
  }

  getActiveControls(certification: Certification) {
    return certification.controls.filter(control => {
      return control.active;
    });
  }

  getControlsWithWarnings(certification: Certification) {
    return certification.controls.filter(control => {
      return control.active && control.fulfilled === Fulfillment.WARNING;
    });
  }

  getGoodControls(certification: Certification) {
    return certification.controls.filter(control => {
      return control.active && control.fulfilled === Fulfillment.GOOD;
    });
  }
}
