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

import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Title } from '@angular/platform-browser';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { EMPTY } from 'rxjs';
import { catchError } from 'rxjs/operators';

import { AuthService } from '../auth.service';
import { ConfigService } from '../config.service';
import { TokenResponse } from '../token-response';
import { User } from '../user';
import { ErrorService } from '../error.service';

@Component({
  selector: 'clouditor-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  model = new User();
  submitted = false;
  oAuthEnabled = false;

  constructor(titleService: Title,
    private router: Router,
    private auth: AuthService,
    private http: HttpClient,
    private error: ErrorService,
    private config: ConfigService) {
    titleService.setTitle('Login');
  }

  ngOnInit() {
    this.http.get<any>(this.config.get().authUrl + '/profile').subscribe(profile => {
      this.oAuthEnabled = profile.enabled;
    });
  }

  getOAuthLogin() {
    return this.config.get().authUrl + '/login';
  }

  onSubmit() {
    this.submitted = true;

    this.http.post(this.config.get().apiUrl + '/authenticate', this.model)
      .pipe(
        catchError((err: HttpErrorResponse) => {
          if (err.status === 401) {
            // incorrect username or password
            this.error.set('Username or password is incorrect.');
            // clear the password
            this.model.password = '';
          } else if (err.status === 0 || err.status === 404) {
            // connection refused
            this.error.set('Could not connect to Clouditor API at ' + this.config.get().apiUrl + '.');

            console.log(this.error);
          }
          return EMPTY;
        })
      )
      .subscribe((response: TokenResponse) => {
        if (response.token) {
          this.error.set(undefined);
          this.auth.login(response.token);
          this.router.navigate(['/']);
        }
      });
  }

  // TODO: Remove this when we're done
  get diagnostic() { return JSON.stringify(this.model); }

}
