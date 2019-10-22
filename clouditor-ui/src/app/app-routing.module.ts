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
import { EditUserComponent } from './edit-user/edit-user.component';

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
  { path: 'users/:id', component: EditUserComponent },
  { path: 'about', component: AboutComponent },
  { path: '', redirectTo: '/accounts', pathMatch: 'full', resolve: { config: ConfigService } },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
