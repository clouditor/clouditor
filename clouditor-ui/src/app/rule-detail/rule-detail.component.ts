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
