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

import { Component, OnInit, ViewChild, OnDestroy } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { take, takeUntil } from 'rxjs/operators';
import { timer } from 'rxjs';

import { DiscoveryService } from '../discovery.service';
import { Scan } from '../scan';
import { NgForm } from '@angular/forms';
import { componentDestroyed } from '@w11k/ngx-componentdestroyed';
import { AccountsService } from '../accounts.service';

@Component({
  selector: 'clouditor-discovery',
  templateUrl: './discovery.component.html',
  styleUrls: ['./discovery.component.scss']
})
export class DiscoveryComponent implements OnInit, OnDestroy {
  scans: Scan[] = [];
  groups: string[] = [];

  @ViewChild('searchForm', { static: true }) searchForm: NgForm;

  search: string;
  // to work around the issue, that by default everything should
  // be selected, a value of true means that it is filtered OUT
  deselected: Map<string, boolean[]> = new Map();

  filtered: Scan[] = [];

  processing: Map<string, boolean> = new Map();
  accountsConfigured: Map<string, boolean> = new Map();

  constructor(private discoveryService: DiscoveryService,
    private accountsService: AccountsService,
    private titleService: Title) {
    this.search = localStorage.getItem('search-scans');
    if (this.search === null) {
      this.search = '';
    }

    const deselected = localStorage.getItem('deselected-scans');
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

  isAccountConfigured(provider: string) {
    const configured = this.accountsConfigured.get(provider);

    return configured !== undefined && configured;
  }

  ngOnInit() {
    this.titleService.setTitle('Discovery');

    // quickly fetch the current account status, to see whether we need to redirect the user to the account page
    this.accountsService.getAccounts().subscribe(accounts => {
      for (const entry of accounts) {
        this.accountsConfigured.set(entry[0], true);
      }
    });

    timer(0, 30000)
      .pipe(
        takeUntil(componentDestroyed(this)),
      )
      .subscribe(() => this.updateScans());

    this.searchForm.form.valueChanges.subscribe(changes => {
      console.log(changes);

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
      localStorage.setItem('deselected-scans', JSON.stringify(this.deselected));
      localStorage.setItem('search-scan', this.search);

      this.updateFiltered();
    });
  }

  ngOnDestroy() {

  }

  updateScans() {
    this.discoveryService
      .getScans()
      .pipe(take(1))
      .subscribe(scans => {
        this.scans = Array.from(scans.values());
        this.scans.sort((a, b) => {
          if (!a.enabled && b.enabled) {
            return 1;
          } else if (!b.enabled && a.enabled) {
            return -1;
          } else {
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

    this.discoveryService.enableScan(scan).subscribe(() => {
      this.processing[scan._id] = false;

      // would actually be enough to just update this particular scan
      this.updateScans();
    });
    timer(3000).subscribe(() => this.updateScans());
  }

  onDisable(scan: Scan) {
    this.processing[scan._id] = true;

    this.discoveryService.disableScan(scan).subscribe(() => {
      this.processing[scan._id] = false;

      // would actually be enough to just update this particular scan
      this.updateScans();
    });
  }
}
