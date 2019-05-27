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
import { Title } from '@angular/platform-browser';

import { StatisticService } from '../statistic.service';
import { AccountsService } from '../accounts.service';
import { Account } from '../account';
import { catchError } from 'rxjs/operators';
import { empty, throwError, EMPTY } from 'rxjs';

@Component({
  selector: 'clouditor-ui',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  statistics: any;
  accounts: Map<string, Account> = new Map();

  constructor(titleService: Title,
    private statisticService: StatisticService, private accountsService: AccountsService) {
    titleService.setTitle('Accounts');
  }

  async ngOnInit() {
    this.statistics = await
      this.statisticService
        .getStatistics()
        .toPromise();

    this.accountsService.getAccounts().subscribe(accounts => this.accounts = accounts);
  }
}
