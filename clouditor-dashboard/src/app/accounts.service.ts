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

import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { ConfigService } from './config.service';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Account } from './account';

@Injectable()
export class AccountsService {
  constructor(private http: HttpClient,
    private config: ConfigService) { }

  getAccounts(): Observable<Map<string, Account>> {
    return this.http.get(this.config.get().apiUrl + '/accounts').pipe(map(data => {
      return new Map(Object.entries(data).map(entry => {
        // for some reason we need to explicitly type this, otherwise all things brake loose
        const p: [string, Account] = [entry[0], { ...new Account(entry[1].provider), ...entry[1] }];

        return p;
      }));
    }));
  }

  getAccount(provider: string): Observable<Account> {
    return this.http.get<Account>(this.config.get().apiUrl + '/accounts/' + provider).pipe(map(data => {
      return Object.assign(new Account(provider), data);
    }));
  }

  putAccount(provider: string, account: Account) {
    return this.http.put<Account>(this.config.get().apiUrl + '/accounts/' + provider, account);
  }

  discover(provider: string): Observable<Account> {
    return this.http.post<Account>(this.config.get().apiUrl + '/accounts/discover/' + provider, null).pipe(map(data => {
      return Object.assign(new Account(provider), data);
    }));
  }
}
