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

import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AboutComponent } from './about/about.component';
import { AssetsComponent } from './assets/assets.component';
import { ComplianceComponent } from './compliance/compliance.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { LoginComponent } from './login/login.component';
import { DiscoveryComponent } from './discovery/discovery.component';
import { ComplianceDetailComponent } from './compliance-detail/compliance-detail.component';
import { ControlDetailComponent } from './control-detail/control-detail.component';
import { AuthGuard } from './auth.guard';
import { ConfigService } from './config.service';
import { DiscoveryDetailComponent } from './discovery-detail/discovery-detail.component';
import { RulesComponent } from './rules/rules.component';
import { RuleDetailComponent } from './rule-detail/rule-detail.component';
import { ConfigureAccountComponent } from './configure-account/configure-account.component';
import { UsersComponent } from './users/users.component';

const routes: Routes = [
  { path: 'accounts', component: DashboardComponent, canActivate: [AuthGuard] },
  { path: 'accounts/configure/:provider', component: ConfigureAccountComponent, canActivate: [AuthGuard] },
  { path: 'compliance', component: ComplianceComponent, canActivate: [AuthGuard] },
  { path: 'compliance/:id', component: ComplianceDetailComponent, canActivate: [AuthGuard] },
  { path: 'compliance/:certificationId/:controlId', component: ControlDetailComponent, canActivate: [AuthGuard] },
  { path: 'assets', component: AssetsComponent, canActivate: [AuthGuard] },
  { path: 'discovery', component: DiscoveryComponent, canActivate: [AuthGuard] },
  { path: 'discovery/:id', component: DiscoveryDetailComponent, canActivate: [AuthGuard] },
  { path: 'rules', component: RulesComponent, canActivate: [AuthGuard] },
  { path: 'rules/:id', component: RuleDetailComponent, canActivate: [AuthGuard] },
  { path: 'login', component: LoginComponent },
  { path: 'users', component: UsersComponent },
  { path: 'about', component: AboutComponent },
  { path: '', redirectTo: '/accounts', pathMatch: 'full', resolve: { config: ConfigService } },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
