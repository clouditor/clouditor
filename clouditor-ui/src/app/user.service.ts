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
import { User } from './user';
import { Observable } from 'rxjs';
import { ConfigService } from './config.service';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  constructor(private http: HttpClient,
    private config: ConfigService) { }

  getUsers(): Observable<User[]> {
    return this.http.get<any>(this.config.get().apiUrl + '/users').pipe(map(users => {
      return users.map(user => Object.assign(new User(), user));
    }));
  }

  getUser(id: string): Observable<User> {
    return this.http.get<User>(this.config.get().apiUrl + '/users/' + id).pipe(map(user => {
      return Object.assign(new User(), user);
    }));
  }

  updateUser(id: string, user: User): Observable<any> {
    return this.http.put(this.config.get().apiUrl + '/users/' + id, user);
  }

  createUser(user: User): Observable<any> {
    return this.http.post(this.config.get().apiUrl + '/users/', user);
  }

  deleteUser(id: string): Observable<any> {
    return this.http.delete(this.config.get().apiUrl + '/users/' + id);
  }
}
