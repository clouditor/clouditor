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
    private route: ActivatedRoute) {
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
      this.assets = scan.lastResult.discoveredAssets;

      /*const source = new EventSource(this.config.get().apiUrl + '/assets/' + scan.assetType + '/subscribe');
      source.addEventListener('message', (e) => {
        console.log(e);
      });*/
    });
  }

}
