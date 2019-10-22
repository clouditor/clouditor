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

import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { APP_INITIALIZER, NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { MomentModule } from 'ngx-moment';
import { ToastrModule } from 'ngx-toastr';
import { AboutComponent } from './about/about.component';
import { AccountsService } from './accounts.service';
import { RuleDetailComponent } from './rule-detail/rule-detail.component';
import { RulesComponent } from './rules/rules.component';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { AssessmentBadgesComponent } from './assessment-badges/assessment-badges.component';
import { AssetTypeService } from './asset-type.service';
import { AssetService } from './asset.service';
import { AssetsComponent } from './assets/assets.component';
import { AuthGuard } from './auth.guard';
import { AuthInterceptor } from './auth.interceptor';
import { AuthService } from './auth.service';
import { CertificationService } from './certification.service';
import { ComplianceDetailComponent } from './compliance-detail/compliance-detail.component';
import { ComplianceComponent } from './compliance/compliance.component';
import { ConfigService } from './config.service';
import { ConfigureAccountComponent } from './configure-account/configure-account.component';
import { ControlDetailComponent } from './control-detail/control-detail.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { DiscoveryDetailComponent } from './discovery-detail/discovery-detail.component';
import { DiscoveryService } from './discovery.service';
import { DiscoveryComponent } from './discovery/discovery.component';
import { ErrorService } from './error.service';
import { LoginComponent } from './login/login.component';
import { MapComponent } from './map/map.component';
import { ScanBadgesComponent } from './scan-badges/scan-badges.component';
import { ScoreComponent } from './score/score.component';
import { ServiceDescriptionService } from './service-description.service';
import { StatisticService } from './statistic.service';
import { TruncateMiddlePipe } from './truncate-middle.pipe';
import { UsersComponent } from './users/users.component';
import { EditUserComponent } from './edit-user/edit-user.component';

@NgModule({
  declarations: [
    AboutComponent,
    RulesComponent,
    RuleDetailComponent,
    AppComponent,
    DashboardComponent,
    DiscoveryDetailComponent,
    DiscoveryComponent,
    LoginComponent,
    ScoreComponent,
    MapComponent,
    AssetsComponent,
    ComplianceComponent,
    TruncateMiddlePipe,
    ComplianceDetailComponent,
    ControlDetailComponent,
    ScanBadgesComponent,
    ConfigureAccountComponent,
    AssessmentBadgesComponent,
    UsersComponent,
    EditUserComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpClientModule,
    AppRoutingModule,
    NgbModule,
    MomentModule,
    ToastrModule.forRoot({ positionClass: 'toast-bottom-left', timeOut: 0, extendedTimeOut: 0 })
  ],
  providers: [
    AccountsService,
    AuthGuard,
    AuthService,
    CertificationService,
    ConfigService,
    {
      provide: APP_INITIALIZER,
      useFactory: (config: ConfigService) => () => config.load(),
      deps: [ConfigService], multi: true
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: AuthInterceptor,
      multi: true
    },
    ServiceDescriptionService,
    DiscoveryService,
    AssetService,
    StatisticService,
    AssetTypeService,
    ErrorService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
