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

import { Component, OnInit, OnDestroy, ViewChild} from '@angular/core';
import { NgForm } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { CertificationService } from '../certification.service';
import { timer } from 'rxjs';
import { takeUntil, flatMap } from 'rxjs/operators';
import { componentDestroyed } from '@w11k/ngx-componentdestroyed';
import { Control, Certification } from '../certification';

@Component({
  selector: 'clouditor-control-detail',
  templateUrl: './control-detail.component.html',
  styleUrls: ['./control-detail.component.css']
})
export class ControlDetailComponent implements OnInit, OnDestroy {
  @ViewChild('filterForm', { static: true }) filterForm: NgForm;
  control: Control;

  viewedOnMobileDevice: Boolean;

  goodAssets: Map<string, any>;
  goodAssetsFiltered: Map<string, any>;

  assetsWithWarnings: Map<string, any>;
  assetsWithWarningsFiltered: Map<string, any>;

  certification: Certification;
  hasPassedAssets: Boolean;
  controlStatistics = {
    passedAssets: 0,
    failedAssets: 0,
    timestamp: '-'
  };

  filterOptions = {
    filterPassed: true,
    filterFailed: true,
    searchTerm: ''
  };

constructor(private route: ActivatedRoute, private certificationService: CertificationService) {
    this.route.params.subscribe(params => {
      this.certificationService.getCertification(params['certificationId']).subscribe(certification => {
        this.certification = certification;
      });
      timer(0, 10000)
        .pipe(
          takeUntil(componentDestroyed(this)),
          flatMap(() => {
            return this.certificationService.getControl(params['certificationId'], params['controlId']);
          }),
        )
        .subscribe(control => {
          this.updateControl(control);
          this.certificationService.getCompliantAssets(params['certificationId'], params['controlId']).subscribe(goodAssets => {
            this.goodAssets = goodAssets;
            this.controlStatistics.passedAssets = goodAssets.size;
            this.hasPassedAssets = (this.goodAssets.size > 0);
            this.filterControlDetails();
          });

          this.certificationService.getNonCompliantAssets(params['certificationId'], params['controlId']).subscribe(assetsWithWarnings => {
            this.assetsWithWarnings = assetsWithWarnings;
            this.controlStatistics.failedAssets = assetsWithWarnings.size;
            this.filterControlDetails();
          });
          this.controlStatistics.timestamp = control.objectives[0]['result']['endTime'];

        });
    });
  }

  ngOnInit() {
    this.filterForm.form.valueChanges.subscribe(params => {
      this.filterOptions = params;
      if (params.searchTerm !== undefined) {
        this.filterControlDetails();
      }
    });
  }

  ngOnDestroy() {

  }

  filterControlDetails() {
      if (this.filterOptions.searchTerm !== null && this.filterOptions.searchTerm !== '') {
        this.goodAssetsFiltered = new Map<string, any>();
        this.assetsWithWarningsFiltered = new Map<string, any>();
        if (this.assetsWithWarnings !== null) {
          for (const key of Array.from( this.assetsWithWarnings.keys())) {
            if (JSON.parse(key)['name'].toLowerCase().includes(this.filterOptions.searchTerm.toLowerCase())) {
              this.assetsWithWarningsFiltered.set(key, this.assetsWithWarnings.get(key));
            } else {
              const details = this.assetsWithWarnings.get(key);
              if (details !== null &&
                  details['message'] !== null &&
                  details.message.toLowerCase().includes(this.filterOptions.searchTerm.toLowerCase())) {
                this.assetsWithWarningsFiltered.set(key, this.assetsWithWarnings.get(key));
              }
            }
          }
        }

        if (this.goodAssets !== null) {
          for (const key of Array.from( this.goodAssets.keys()) ) {
            if (JSON.parse(key)['name'].toLowerCase().includes(this.filterOptions.searchTerm.toLowerCase())) {
              this.goodAssetsFiltered.set(key, this.goodAssets.get(key));
            } else {
              const details = this.goodAssets.get(key);
              if (details !== null &&
                  details['message'] !== null &&
                  details.message.toLowerCase().includes(this.filterOptions.searchTerm.toLowerCase())) {
                this.goodAssetsFiltered.set(key, this.goodAssets.get(key));
              }
            }
          }
        }
    } else {
      this.assetsWithWarningsFiltered = new Map<string, any>(this.assetsWithWarnings);
      this.goodAssetsFiltered = new Map<string, any>(this.goodAssets);
    }
  }

  updateControl(control: Control) {
    this.control = control;
  }

  toAsset(key: string): any {
    return JSON.parse(key);
  }
}
