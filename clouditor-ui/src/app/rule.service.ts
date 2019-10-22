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
import { HttpClient } from '@angular/common/http';

import { ConfigService } from './config.service';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { Rule } from './rule';
import { RuleEvaluation } from './rule-evaluation';

@Injectable({ providedIn: 'root' })
export class RuleService {
  constructor(private http: HttpClient,
    private config: ConfigService) { }

  getRules(type: string): Observable<Rule[]> {
    return this.http.get<any>(this.config.get().apiUrl + '/rules/assets/' + type).pipe(map(entries => {
      return entries.map(entry => Object.assign(new Rule(), entry));
    }));
  }

  getRuleEvaluation(id: string): Observable<RuleEvaluation> {
    return this.http.get<RuleEvaluation>(this.config.get().apiUrl + '/rules/' + id).pipe(map(data => {
      const evaluation = Object.assign(new RuleEvaluation(), data);
      evaluation.rule = Object.assign(new Rule(), data.rule);

      return evaluation;
    }));
  }
}
