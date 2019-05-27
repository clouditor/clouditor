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

  model = new User('', '');
  submitted = false;

  constructor(titleService: Title,
    private router: Router,
    private auth: AuthService,
    private http: HttpClient,
    private error: ErrorService,
    private config: ConfigService) {
    titleService.setTitle('Login');
  }

  ngOnInit() {

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
