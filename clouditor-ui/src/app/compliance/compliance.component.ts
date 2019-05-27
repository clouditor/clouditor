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
