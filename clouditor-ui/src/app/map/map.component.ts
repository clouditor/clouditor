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

import { Component, ElementRef, Input, OnChanges, OnDestroy, OnInit } from '@angular/core';

@Component({
  selector: 'clouditor-map',
  templateUrl: './map.component.html',
  styleUrls: ['./map.component.css']
})
export class MapComponent implements OnInit, OnDestroy, OnChanges {
  @Input() locations: string[];

  url: string;
  url2: string;

  constructor(private element: ElementRef) {

  }

  ngOnInit(): void {
    this.updateUrl();
  }

  ngOnChanges(changes): void {
    this.updateUrl();
  }

  ngOnDestroy(): void {

  }

  updateUrl(): void {
    if (this.locations && this.locations.length > 0) {
      this.url = 'https://maps.googleapis.com/maps/api/staticmap?center='
        + (this.locations[0])
        + '&zoom=2&size=640x400&maptype=roadmap'
        + this.locations.map(location => {
          return '&markers=color:green%7C' + location;
        }).join('');
      this.url2 = 'https://maps.googleapis.com/maps/api/staticmap?center='
        + (this.locations[0])
        + '&zoom=4&size=640x400&maptype=roadmap'
        + this.locations.map(location => {
          return '&markers=color:green%7C' + location;
        }).join('');
    }
  }
}
