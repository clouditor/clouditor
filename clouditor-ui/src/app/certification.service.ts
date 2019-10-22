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
