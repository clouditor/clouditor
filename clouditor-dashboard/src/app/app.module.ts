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
import { AnalysisDetailComponent } from './analysis-detail/analysis-detail.component';
import { AnalysisComponent } from './analysis/analysis.component';
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
import { GaugeChartComponent } from './gauge-chart/gauge-chart.component';
import { LoginComponent } from './login/login.component';
import { MapComponent } from './map/map.component';
import { ScanBadgesComponent } from './scan-badges/scan-badges.component';
import { ScoreComponent } from './score/score.component';
import { ServiceDescriptionService } from './service-description.service';
import { StatisticService } from './statistic.service';
import { TruncateMiddlePipe } from './truncate-middle.pipe';

@NgModule({
  declarations: [
    AboutComponent,
    AnalysisComponent,
    AnalysisDetailComponent,
    AppComponent,
    DashboardComponent,
    DiscoveryDetailComponent,
    DiscoveryComponent,
    LoginComponent,
    GaugeChartComponent,
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
