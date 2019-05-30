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
import { ActivatedRoute } from '@angular/router';
import { Rule } from '../rule';
import { RuleService } from '../rule.service';
import { AssetService } from '../asset.service';
import { RuleEvaluation } from '../rule-evaluation';
import { Asset } from '../asset';
import { DiscoveryService } from '../discovery.service';
import { Scan } from '../scan';

@Component({
  selector: 'clouditor-rule-detail',
  templateUrl: './rule-detail.component.html',
  styleUrls: ['./rule-detail.component.css']
})
export class RuleDetailComponent implements OnInit {
  evaluation: RuleEvaluation;
  rule: Rule;
  scan: Scan;
  resources: Asset[] = [];

  isExpanded: Map<string, boolean> = new Map();

  constructor(private ruleService: RuleService,
    private assetService: AssetService,
    private discoveryService: DiscoveryService,
    private route: ActivatedRoute) {
    this.route.params.subscribe(params => {
      const ruleId = params['id'];

      this.ruleService.getRuleEvaluation(ruleId).subscribe(evaluation => {
        this.evaluation = evaluation;
        this.rule = evaluation.rule;

        // fetch associated assets
        this.assetService.getAssetsWithType(this.rule.assetType).subscribe(resources => {
          this.resources = resources;
        });

        // fetch associated scan
        this.discoveryService.getScan(this.rule.assetType).subscribe(scan => {
          this.scan = scan;
        });

        /*for (const resourceId of Object.keys(evaluation.compliance)) {
          this.assetService.getAsset(resourceId).subscribe(resource => {
            this.resources.push(resource);
          });
        }*/
      });
    });
  }

  ngOnInit(): void {
  }

  isResourceOk(id: string) {
    return this.evaluation.compliance[id];
  }
}
