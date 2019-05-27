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
