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
  nonCompliantAssets: number;

  constructor(private assetService: AssetService) { }

  ngOnInit() {
    this.assetService.getAssetsWithType(this.assetType).subscribe(assets => {
      const assetArr: any[] = assets;
      let nonCompliantAssets = 0;
      assetArr.forEach(asset => {
        for (const result of asset.evaluationResults) {
          if (result.failedConditions.length > 0) {
            nonCompliantAssets += 1;
            break;
          }
        }
        this.nonCompliantAssets = nonCompliantAssets;
      });
    });
  }

}
