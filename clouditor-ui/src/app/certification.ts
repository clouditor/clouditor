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

export const Fulfillment = {
  NOT_EVALUATED: 'NOT_EVALUATED',
  WARNING: 'WARNING',
  GOOD: 'GOOD'
};

export class Certification {
  constructor(public _id?: string,
    public controls?: Control[],
    public description?: string,
    public publisher?: string,
    public website?: string) { }
}

export class Domain {
  constructor(public name: string) { }
}

export class Control {
  constructor(public objectives?: Objective[],
    public controlId?: string,
    public name?: string,
    public domain?: Domain,
    public description?: string,
    public fulfilled?: string,
    public active?: boolean,
    public automated?: boolean,
    public violations?: number) { }

  hasWarning(): boolean {
    return this.active && this.fulfilled === Fulfillment.WARNING;
  }

  isNotEvaluated(): boolean {
    return this.active && this.fulfilled === Fulfillment.NOT_EVALUATED;
  }

  isGood(): boolean {
    return this.active && this.fulfilled === Fulfillment.GOOD;
  }
}

export class Objective {
  constructor(public metricUri: string,
    public condition: string,
    public fulfilled: boolean) { }
}
