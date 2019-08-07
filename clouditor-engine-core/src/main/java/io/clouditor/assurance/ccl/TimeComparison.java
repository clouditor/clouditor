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

package io.clouditor.assurance.ccl;

import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.Map;

public class TimeComparison extends Comparison {

  private TimeOperator timeOperator = TimeOperator.BEFORE;

  private int relativeValue = 0;
  private ChronoUnit timeUnit = ChronoUnit.DAYS;

  @Override
  public boolean evaluate(Map properties) {
    var fieldValue = getValueFromField(properties, this.field);

    // we support two different formats of time fields
    // a) serialized Instants
    // b) unix timestamps in seconds

    Instant instant;

    if (fieldValue == null) {
      return false;
    }

    if (fieldValue instanceof Long) {
      instant = Instant.ofEpochSecond((Long) fieldValue);
    } else if ((fieldValue instanceof Map) && ((Map) fieldValue).get("epochSecond") != null) {
      // field values are not really instants but serialized Instant in the form of epochSecond
      // and nano, we need to re-create the Instant

      instant =
          Instant.ofEpochSecond(
              ((Number) ((Map) fieldValue).get("epochSecond")).longValue(),
              ((Number) ((Map) fieldValue).get("nano")).longValue());
    } else {
      try {
        instant = Instant.parse((String) fieldValue);
      } catch (ClassCastException e) {
        return false;
      }
    }

    Instant value;
    if (this.timeOperator == TimeOperator.BEFORE || this.timeOperator == TimeOperator.AFTER) {
      value = Instant.now().plus(relativeValue, timeUnit);
    } else {
      value = Instant.now().minus(relativeValue, timeUnit);
    }

    if (this.timeOperator == TimeOperator.YOUNGER || this.timeOperator == TimeOperator.AFTER) {
      return instant.isAfter(value);
    } else {
      return instant.isBefore(value);
    }
  }

  public int getRelativeValue() {
    return relativeValue;
  }

  public void setRelativeValue(int relativeValue) {
    this.relativeValue = relativeValue;
  }

  public ChronoUnit getTimeUnit() {
    return timeUnit;
  }

  public void setTimeUnit(ChronoUnit timeUnit) {
    this.timeUnit = timeUnit;
  }

  public void setTimeOperator(TimeOperator timeOperator) {
    this.timeOperator = timeOperator;
  }

  public enum TimeOperator {
    BEFORE,
    AFTER,
    YOUNGER,
    OLDER
  }
}
