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
