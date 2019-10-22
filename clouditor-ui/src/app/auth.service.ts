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
import { JwtHelperService } from '@auth0/angular-jwt';
import { Router, ActivatedRoute } from '@angular/router';
import { HttpParams } from '@angular/common/http';

export const TOKEN_NAME = 'token';

const helper = new JwtHelperService();

@Injectable()
export class AuthService {

  constructor(private router: Router, private route: ActivatedRoute) {
    const params = new HttpParams({ fromString: window.location.hash.replace('#?', '') });

    const token = params.get('token');

    if (token) {
      this.login(token);
      this.router.navigate(['/']);
    }
  }

  isLoggedIn() {
    return !helper.isTokenExpired(this.getToken());
  }

  getToken() {
    return localStorage.getItem(TOKEN_NAME);
  }

  getUser() {
    return helper.decodeToken(this.getToken()).sub;
  }

  logout() {
    localStorage.removeItem(TOKEN_NAME);
  }

  login(token: string) {
    localStorage.setItem(TOKEN_NAME, token);
  }
}
