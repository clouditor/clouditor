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

import { Certification, Control } from './certification';

import { ConfigService } from './config.service';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable()
export class CertificationService {
  constructor(private http: HttpClient,
    private config: ConfigService) { }

  getCertification(certificationId: string): Observable<Certification> {
    return this.http.get<Certification>(this.config.get().apiUrl + '/certification/' + certificationId).pipe(map(data => {
      return {
        ...new Certification(),
        ...data,
        controls: data.controls.map(x => Object.assign(new Control(), x))
      };
    }));
  }

  getControl(certificationId: string, controlId: string): Observable<Control> {
    return this.http.get<Control>(this.config.get().apiUrl + '/certification/' + certificationId + '/' + controlId).pipe(map(data => {
      return Object.assign(new Control, data);
    }));
  }

  getNonCompliantAssets(certificationId: string, controlId: string): Observable<Map<string, any>> {
    return this.http.get<Map<string, any>>(this.config.get().apiUrl +
      '/certification/' +
      certificationId +
      '/' + controlId +
      '/assets/warning').pipe(map(data => {
        return new Map(Object.entries(data));
      }));
  }

  getCompliantAssets(certificationId: string, controlId: string): Observable<Map<string, any>> {
    return this.http.get<Map<string, any>>(this.config.get().apiUrl +
      '/certification/' +
      certificationId +
      '/' + controlId +
      '/assets/good').pipe(map(data => {
        return new Map(Object.entries(data));
      }));
  }

  getCertifications(): Observable<Map<String, Certification>> {
    return this.http.get<Map<String, Certification>>(this.config.get().apiUrl + '/certification/').pipe(map(data => {
      return new Map(Object.entries(data));
    }));
  }

  modifyControlStatus(certificationId: string, controlId: string, status: boolean): any {
    return this.http.post(this.config.get().apiUrl
      + '/certification/'
      + certificationId
      + '/'
      + controlId + '/status', { 'status': status });
  }

  import(certificationId: string): any {
    return this.http.post(this.config.get().apiUrl + '/certification/import/' + certificationId, null);
  }

  getImporters(): Observable<any[]> {
    return this.http.get<any[]>(this.config.get().apiUrl + '/certification/importers');
  }

}
