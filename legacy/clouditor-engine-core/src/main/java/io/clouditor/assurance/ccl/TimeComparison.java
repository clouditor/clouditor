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
