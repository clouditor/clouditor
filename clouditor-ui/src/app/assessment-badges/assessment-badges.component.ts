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

import { Component, OnInit, Input } from '@angular/core';

import { Scan } from '../scan';
import { AssetService } from '../asset.service';


@Component({
  selector: 'clouditor-assessment-badges',
  templateUrl: './assessment-badges.component.html',
  styleUrls: ['./assessment-badges.component.css']
})
export class AssessmentBadgesComponent implements OnInit {

  @Input() scan: Scan;
  @Input() assetType: string;
  compliantAssets: number;
  nonCompliantAssets: number;

  constructor(private assetService: AssetService) { }

  ngOnInit() {
    this.assetService.getAssetsWithType(this.assetType).subscribe(assets => {
      const assetArr: any[] = assets;
      let compliantAssets = 0;
      let nonCompliantAssets = 0;
      assetArr.forEach(asset => {
        if (!asset.compliant) {
          nonCompliantAssets += 1;
        } else {
          compliantAssets += 1;
        }
        this.compliantAssets = compliantAssets;
        this.nonCompliantAssets = nonCompliantAssets;
      });
    });
  }

}
