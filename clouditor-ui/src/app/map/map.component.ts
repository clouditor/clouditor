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
