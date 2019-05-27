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
