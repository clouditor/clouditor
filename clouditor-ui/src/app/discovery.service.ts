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
