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

import { Component, OnInit, OnDestroy } from '@angular/core';
import { DiscoveryService } from '../discovery.service';
import { ActivatedRoute } from '@angular/router';
import { Scan } from '../scan';
import { Rule } from '../rule';
import { RuleService } from '../rule.service';
import { AssetService } from '../asset.service';
import { timer } from 'rxjs';
import { componentDestroyed } from '@w11k/ngx-componentdestroyed';
import { takeUntil } from 'rxjs/operators';
import { ConfigService } from '../config.service';

@Component({
  selector: 'clouditor-scan-detail',
  templateUrl: './discovery-detail.component.html',
  styleUrls: ['./discovery-detail.component.css']
})
export class DiscoveryDetailComponent implements OnInit, OnDestroy {
  scan: Scan;
  rules: Rule[];
  assets: any[];

  isExpanded: Map<string, boolean> = new Map();

  processing: Map<string, boolean> = new Map();

  constructor(private discoveryService: DiscoveryService,
    private ruleService: RuleService,
    private assetService: AssetService,
    private route: ActivatedRoute,
    private config: ConfigService) {
    this.route.params.subscribe(params => {
      timer(0, 10000)
        .pipe(
          takeUntil(componentDestroyed(this)),
        )
        .subscribe(x => {
          this.updateScan(params['id']);
        });
    });
  }

  ngOnInit(): void {
  }

  ngOnDestroy(): void {
  }

  onEnable(scan: Scan) {
    this.processing[scan._id] = true;

    this.discoveryService.enableScan(scan).subscribe(() => {
      this.processing[scan._id] = false;

      // would actually be enough to just update this particular scan
      this.updateScan(scan._id);
    });
  }

  onDisable(scan: Scan) {
    this.processing[scan._id] = true;

    this.discoveryService.disableScan(scan).subscribe(() => {
      this.processing[scan._id] = false;

      // would actually be enough to just update this particular scan
      this.updateScan(scan._id);
    });
  }

  updateScan(id: string): any {
    this.discoveryService.getScan(id).subscribe(scan => {
      this.scan = scan;

      /*const source = new EventSource(this.config.get().apiUrl + '/assets/' + scan.assetType + '/subscribe');
      source.addEventListener('message', (e) => {
        console.log(e);
      });*/

      this.assetService.getAssetsWithType(this.scan.assetType).subscribe(assets => this.assets = assets);
    });
  }

}
