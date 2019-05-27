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
import { DecimalPipe } from '@angular/common';

import * as c3 from 'c3';

@Component({
  selector: 'clouditor-gauge-chart',
  templateUrl: './gauge-chart.component.html',
  styleUrls: ['./gauge-chart.component.css']
})
export class GaugeChartComponent implements OnInit, OnDestroy, OnChanges {
  @Input() value: number;

  private chart: any;

  constructor(private element: ElementRef) {

  }

  ngOnInit(): void {
    this.chart = c3.generate({
      data: {
        columns: [
          ['data', this.value * 100]
        ],
        type: 'gauge',
      },
      bindto: this.element.nativeElement,
      gauge: {
        label: {
          format: function(value, ratio) {
            return new DecimalPipe('en').transform(value, '1.1-1') + ' %';
          }
        },
      },
      color: {
        pattern: ['#FF0000', '#F97600', '#F6C600', '#60B044'], // the three color levels for the percentage values.
        threshold: {
          values: [30, 60, 90, 100]
        }
      }
    });
  }

  loadData(): void {
    this.chart.load({
      columns: [
        ['data', this.value * 100]
      ]
    });
  }

  ngOnChanges(changes): void {
    if (this.chart) {
      this.loadData();
    }
  }

  ngOnDestroy(): void {

  }
}
