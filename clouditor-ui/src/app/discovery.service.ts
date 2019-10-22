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
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { Scan } from './scan';

@Injectable()
export class DiscoveryService {
  constructor(private http: HttpClient,
    private config: ConfigService) { }

  getScans(): Observable<Map<string, Scan>> {
    return this.http.get<any>(this.config.get().apiUrl + '/discovery/').pipe(map(entries => entries.map(entry => Scan.fromJSON(entry))));
  }

  getScan(id: string): Observable<Scan> {
    return this.http.get<Scan>(this.config.get().apiUrl + '/discovery/' + id).pipe(map(entry => Scan.fromJSON(entry)));
  }

  enableScan(scan: Scan): Observable<any> {
    return this.http.post<any>(this.config.get().apiUrl + '/discovery/' + scan._id + '/enable', {});
  }

  disableScan(scan: Scan): Observable<any> {
    return this.http.post<any>(this.config.get().apiUrl + '/discovery/' + scan._id + '/disable', {});
  }
}
