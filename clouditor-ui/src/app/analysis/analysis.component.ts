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

import { Component, OnInit, ViewChild } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { take } from 'rxjs/operators';

import { DiscoveryService } from '../discovery.service';
import { Scan } from '../scan';
import { NgForm } from '@angular/forms';
import { RuleService } from '../rule.service';
import { Rule } from '../rule';

@Component({
  selector: 'clouditor-analysis',
  templateUrl: './analysis.component.html',
  styleUrls: ['./analysis.component.scss']
})
export class AnalysisComponent implements OnInit {
  scans: Scan[] = [];
  groups: string[] = [];
  assets: any[] = [];
  assetRules: Map<string, Rule[]> = new Map();

  selectedScan: Scan = new Scan();

  @ViewChild('searchForm') searchForm: NgForm;

  search: string;
  // to work around the issue, that by default everything should
  // be selected, a value of true means that it is filtered OUT
  deselected: Map<string, boolean[]> = new Map();

  filtered: Scan[] = [];

  processing: Map<string, boolean> = new Map();
  configuring: boolean;

  constructor(private checkService: DiscoveryService,
    private titleService: Title,
    private ruleService: RuleService) {
    this.search = localStorage.getItem('search-analysis');
    if (this.search === null) {
      this.search = '';
    }

    const deselected = localStorage.getItem('deselected-analysis');
    if (deselected !== null) {
      this.deselected = JSON.parse(deselected);
    }
  }

  getClassForGroup(group: string, colored: boolean): string {
    if (group === 'AWS') {
      if (colored) {
        return 'fab fa-aws aws';
      } else {
        return 'fab fa-aws text-muted';
      }
    } else if (group === 'Azure') {
      if (colored) {
        return 'fab fa-windows azure';
      } else {
        return 'fab fa-windows text-muted';
      }
    } else if (group === 'EU-SEC Audit API') {
      if (colored) {
        return 'fas fa-shield-alt eu-sec';
      } else {
        return 'fas fa-shield-alt text-muted';
      }
    } else {
      if (colored) {
        return 'fas fa-cloud';
      } else {
        return 'fas fa-cloud text-muted';
      }
    }
  }

  ngOnInit() {
    this.titleService.setTitle('Analysis');

    this.configuring = false;

    this.updateScans();

    this.searchForm.form.valueChanges.subscribe(changes => {
      const keys = Object.keys(changes);

      // we only need to watch for updates if all keys are there
      // +1 is for the search field
      if (keys.length !== this.groups.length + 1) {
        return;
      }

      // because we are using [ngModel], not [(ngModel)], we need to update our backing map ourselves
      // this is indented, so we can properly observe the changes
      for (const key of keys) {
        if (key === 'search') {
          this.search = changes[key];
          continue;
        }

        // the group keys are called selected:{{group}}
        // TODO: we tried to force the keys to be arrays, but it didn't work, this would be a cleaner solution
        const rr = key.split(':');
        if (rr[0] === 'selected' && rr.length === 2) {
          const group = rr[1];

          // update the backing map
          this.deselected[group] = !changes[key];
          continue;
        }

      }

      // update local storage
      localStorage.setItem('deselected-analysis', JSON.stringify(this.deselected));
      localStorage.setItem('search-scan', this.search);

      this.updateFiltered();

      this.filtered.forEach(scan => {
        this.ruleService.getRules(scan.assetType).subscribe(rules => {
          this.assetRules.set(scan.assetType, rules);
        });
      });
    });
  }

  updateScans() {
    this.checkService
      .getScans()
      .pipe(take(1))
      .subscribe(scans => {
        this.scans = Array.from(scans.values());
        this.scans.sort((a, b) => {
          if (a.group > b.group) {
            return 1;
          } else if (b.group > a.group) {
            return -1;
          } else {
            if (a.service > b.service) {
              return 1;
            } else if (b.service > a.service) {
              return -1;
            } else {
              return 0;
            }
          }
        });

        // update groups (should be unique)
        this.groups = this.scans.map(scan => scan.group).filter((value, index, self) => self.indexOf(value) === index);

        // update filtered scans
        this.updateFiltered();
      });
  }

  updateFiltered() {
    // first, filter according to the group
    let filtered = this.scans.filter(scan => !this.deselected[scan.group]);

    // only show scans that are being discovered
    filtered = filtered.filter(scan => scan.enabled);

    // filter according to the search
    filtered = filtered.filter((scan: Scan) => {
      const search = this.search.toLowerCase();
      return (scan.assetType !== null && scan.assetType.toLowerCase().includes(search)) ||
        (scan.group !== null && scan.group.toLowerCase().includes(search)) ||
        (scan.service !== null && scan.service.toLowerCase().includes(search));
    });

    // set it
    this.filtered = filtered;
  }

  onEnable(scan: Scan) {
    this.processing[scan._id] = true;

    this.checkService.enableScan(scan).subscribe(() => {
      this.processing[scan._id] = false;

      // would actually be enough to just update this particular scan
      this.updateScans();
    });
  }

  onDisable(scan: Scan) {
    this.processing[scan._id] = true;

    this.checkService.disableScan(scan).subscribe(() => {
      this.processing[scan._id] = false;

      // would actually be enough to just update this particular scan
      this.updateScans();
    });
  }

  getRulesForAssetType(assetType: string) {
    return this.assetRules.get(assetType);
  }

  selectScan(scan: Scan) {
    if (this.selectedScan === scan) {
      this.selectedScan = new Scan();
    } else {
      this.selectedScan = scan;
    }
  }
}
