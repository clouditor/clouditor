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

import { Component, OnInit, ViewChild, OnDestroy } from '@angular/core';
import { Certification, Fulfillment, Control } from '../certification';
import { ActivatedRoute } from '@angular/router';
import { CertificationService } from '../certification.service';
import { NgForm } from '@angular/forms';
import { timer } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { componentDestroyed } from '@w11k/ngx-componentdestroyed';

@Component({
    selector: 'clouditor-compliance-detail',
    templateUrl: './compliance-detail.component.html',
    styleUrls: ['./compliance-detail.component.scss']
})
export class ComplianceDetailComponent implements OnInit, OnDestroy {

    @ViewChild('searchForm', { static: true }) searchForm: NgForm;

    selected = {};
    isCollapsed = false;
    certification: Certification;

    search = ''; // default to an empty search field

    // in our default filter settings, everything is turned on except those not monitored to easier view the monitored ones
    filterOptions = {
        waiting: true, // not enough data
        notMonitored: false,
        passed: true,
        failed: true
    };

    filteredControls: Control[] = [];
    processing: Map<string, boolean> = new Map();

    constructor(private route: ActivatedRoute, private certificationService: CertificationService) {
        this.route.params.subscribe(params => {
            timer(0, 10000)
                .pipe(
                    takeUntil(componentDestroyed(this)),
                    // TODO: it would make sense to handle this globally for all components
                    // catchError(this.onError.bind(this))
                )
                .subscribe(x => {
                    this.updateCertification(params['id']);
                });
        });
    }

    ngOnInit() {
        this.route.queryParams.subscribe(params => {
            if ('filter' in params) {
                // if we see a filter incoming, don't disable all options and only show those specified
                this.filterOptions = {
                    waiting: false,
                    notMonitored: false,
                    passed: false,
                    failed: false
                };

                if ('passed' in params) {
                    this.filterOptions['passed'] = true;
                }

                if ('failed' in params) {
                    this.filterOptions['failed'] = true;
                }

                if ('waiting' in params) {
                    this.filterOptions['waiting'] = true;
                }

                if ('notMonitored' in params) {
                    this.filterOptions['notMonitored'] = true;
                }
            }
        });

        this.searchForm.form.valueChanges.subscribe(params => {
            if (params.search != null &&
                params.filterWaiting != null &&
                params.filterNotMonitored != null &&
                params.filterPassed != null &&
                params.filterFailed != null) {
                this.search = params.search;
                this.filterOptions = {
                    waiting: params.filterWaiting,
                    notMonitored: params.filterNotMonitored,
                    passed: params.filterPassed,
                    failed: params.filterFailed
                };
                this.updateFilteredControls();
            }
        });
    }

    ngOnDestroy(): void {

    }

    onSearchChanged(value) {

    }

    onError(): void {

    }

    updateFilteredControls() {
        if (this.certification === undefined) {
            this.filteredControls = [];
            return;
        }

        if (this.search === undefined || this.search === '') {
            this.filteredControls = this.certification.controls;
        } else {
            this.filteredControls = this.certification.controls.filter((control: Control) => {
                const search = this.search.toLowerCase();
                return (control.controlId !== null && control.controlId.toLowerCase().includes(search)) ||
                    (control.name !== null && control.name.toLowerCase().includes(search)) ||
                    (control.domain.name !== null && control.domain.name.toLowerCase().includes(search)) ||
                    (control.description !== null && control.description.toLowerCase().includes(search));
            });
        }

        this.filteredControls = this.filteredControls.filter((control: Control) => {
            return (this.filterOptions.waiting && control.isNotEvaluated()) ||
                (this.filterOptions.notMonitored && !control.active) ||
                (this.filterOptions.passed && control.isGood()) ||
                (this.filterOptions.failed && control.hasWarning())
                ;
        });
    }

    getMonitoredControls(certification: Certification) {
        return certification.controls.filter(control => {
            return control.active;
        });
    }

    getInactiveControls(certification: Certification) {
        return certification.controls.filter(control => {
            return !control.active;
        });
    }

    getFailedControls(certification: Certification) {
        return certification.controls.filter(control => {
            return control.active && control.fulfilled === Fulfillment.WARNING;
        });
    }

    getPassedControls(certification: Certification) {
        return certification.controls.filter(control => {
            return control.active && control.fulfilled === Fulfillment.GOOD;
        });
    }

    doEnable(controlId: string, status: boolean) {
        const controlIds = Object.keys(this.selected).filter(key => this.selected[key] === true);

        this.processing[controlId] = true;

        this.certificationService.modifyControlStatus(this.certification._id, controlId, status).subscribe(() => {
            this.processing[controlId] = false;

            this.updateCertification(this.certification._id);
        });
    }

    doSelectAll() {
        for (const control of this.certification.controls) {
            if (control.automated) {
                this.selected[control.controlId] = true;
            }
        }
    }

    updateCertification(certificationId: string): any {
        this.certificationService.getCertification(certificationId).subscribe(certification => {
            this.certification = certification;

            this.updateFilteredControls();
        });
    }

}
