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

import * as flat from 'flat';

import { DiscoveryService } from '../discovery.service';
import { AssetService } from '../asset.service';
import { AssetTypeService } from '../asset-type.service';
import { AssetTypeInfo } from '../asset-type-info';

@Component({
  selector: 'clouditor-assets',
  templateUrl: './assets.component.html',
  styleUrls: ['./assets.component.css']
})
export class AssetsComponent implements OnInit {
  assets: any;
  assetTypes: AssetTypeInfo[];
  selectedAsset: any;

  constructor(private assetService: AssetService,
    private assetTypesService: AssetTypeService,
    private checkService: DiscoveryService) { }

  async ngOnInit() {
    this.assetTypes = await
      this.assetTypesService
        .getAssetTypes()
        .toPromise();

    this.assets = await
      this.assetService
        .getAssets()
        .toPromise();

    if (this.assets.length > 0) {
      this.selectAsset(this.assets[0]);
    }
  }

  async selectAsset(asset) {
    this.selectedAsset = asset;
    this.selectedAsset.checks = [];

    this.selectedAsset.jobnedChecks = 0;
  }

  getFields(data) {
    data = flat.flatten(data);

    return Object.keys(data)
      .map(k => {
        // TODO: define white-list instead of blacklist
        if (!k.startsWith('seenBy') &&
          !k.startsWith('jobnedChecks') &&
          !k.startsWith('object.tags') &&
          !k.startsWith('check')) {
          const json = { name: k, value: data[k] };
          return json;
        } else {
          return null;
        }
      })
      .filter(v => v != null);
  }
}
