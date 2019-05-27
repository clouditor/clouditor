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
import { AccountsService } from '../accounts.service';
import { Account } from '../account';
import { ActivatedRoute, Router } from '@angular/router';
import { catchError } from 'rxjs/operators';
import { of, throwError, empty, EMPTY } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';
import { ThrowStmt } from '@angular/compiler';
import { ToastrService } from 'ngx-toastr';

@Component({
  selector: 'clouditor-configure-account',
  templateUrl: './configure-account.component.html',
  styleUrls: ['./configure-account.component.css']
})
export class ConfigureAccountComponent implements OnInit {

  discoveredAccount: Account;
  account: Account;
  provider: string;
  discoveryComplete = false;

  constructor(private toastr: ToastrService,
    private route: ActivatedRoute,
    private router: Router,
    private accountsService: AccountsService) { }

  ngOnInit() {
    this.route.params.subscribe(params => {
      this.provider = params['provider'];

      this.accountsService.getAccount(this.provider)
        .pipe(catchError((err: HttpErrorResponse) => {
          if (err.status === 404) {
            const account = new Account(this.provider);

            if (this.provider === 'Azure') {
              account.authFile = '~/.azure/clouditor.azureauth';
            }

            // create a "new" account
            return of(account);
          }

          throwError(err);
        }))
        .subscribe(account => {
          this.account = account;

          this.discover();
        });
    });
  }

  save() {
    this.accountsService.putAccount(this.provider, this.account)
      .pipe(catchError((err: HttpErrorResponse) => {
        this.toastr.error('Could not add ' + this.provider + ' account: ' + err.error);

        return EMPTY;
      }))
      .subscribe(_ => {
        this.router.navigateByUrl('/');
      });
  }

  discover() {
    this.accountsService.discover(this.provider)
      .pipe(catchError((err: HttpErrorResponse) => {
        if (err.status === 404) {
          // cloud not discover the cloud account, set auto-discovery to false and deactivate it
          this.account.autoDiscovered = false;
          this.discoveryComplete = true;

          return EMPTY;
        }

        throwError(err);
      }))
      .subscribe(account => {
        this.discoveredAccount = account;
        this.discoveryComplete = true;
      });
  }
}
