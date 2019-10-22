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

import { Component, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { UserService } from '../user.service';
import { User } from '../user';
import { NgControl, NgModel } from '@angular/forms';
import { catchError } from 'rxjs/operators';
import { HttpErrorResponse } from '@angular/common/http';
import { ToastrService } from 'ngx-toastr';
import { EMPTY } from 'rxjs';

@Component({
  selector: 'clouditor-edit-user',
  templateUrl: './edit-user.component.html',
  styleUrls: ['./edit-user.component.css']
})
export class EditUserComponent implements OnInit {

  creating = false;
  changePassword = false;
  user: User;

  allRoles = ['guest', 'user', 'admin'];

  constructor(private userService: UserService,
    private route: ActivatedRoute,
    private router: Router,
    private toastr: ToastrService) {
    this.route.params.subscribe(params => {
      const userId = params['id'];

      if (userId === 'new') {
        this.creating = true;
        this.changePassword = true;
        this.user = new User();
        this.user.roles = ['guest'];
      } else {
        this.userService.getUser(userId).subscribe(user => {
          this.user = user;
        });
      }
    });
  }

  ngOnInit() {
  }

  hasRole(needle: string) {
    for (const role of this.user.roles) {
      if (role === needle) {
        return true;
      }
    }

    return false;
  }

  onSubmit() {
    if (this.creating) {
      this.userService.createUser(this.user)
        .pipe(catchError((err: HttpErrorResponse) => {
          this.toastr.error('Could not create user: ' + err.error);

          return EMPTY;
        }))
        .subscribe(_ => {
          this.router.navigateByUrl('/users');
        });
    } else {
      this.userService.updateUser(this.user.username, this.user)
        .pipe(catchError((err: HttpErrorResponse) => {
          this.toastr.error('Could not update user: ' + err.error);

          return EMPTY;
        }))
        .subscribe(_ => {
          this.router.navigateByUrl('/users');
        });
    }
  }

}
