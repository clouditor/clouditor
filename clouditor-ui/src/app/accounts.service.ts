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
