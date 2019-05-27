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

import io.clouditor.assurance.grammar.CCLParser;
import java.util.Map;
import java.util.Objects;

public class BinaryComparison extends Comparison {

  public enum Operator {
    EQUALS,
    NOT_EQUALS,
    LESS_THAN,
    LESS_THAN_OR_EQUALS,
    GREATER_THAN,
    GREATER_THAN_OR_EQUALS,
    CONTAINS;

    public static Operator of(String operatorCode) {
      if (isEqualOperator(operatorCode, CCLParser.ContainsOperator)) {
        return CONTAINS;
      } else if (isEqualOperator(operatorCode, CCLParser.MoreOrEqualsThanOperator)) {
        return GREATER_THAN_OR_EQUALS;
      } else if (isEqualOperator(operatorCode, CCLParser.MoreThanOperator)) {
        return GREATER_THAN;
      } else if (isEqualOperator(operatorCode, CCLParser.LessOrEqualsThanOperator)) {
        return LESS_THAN_OR_EQUALS;
      } else if (isEqualOperator(operatorCode, CCLParser.LessThanOperator)) {
        return LESS_THAN;
      } else if (isEqualOperator(operatorCode, CCLParser.NotEqualsOperator)) {
        return NOT_EQUALS;
      } else {
        return EQUALS;
      }
    }

    private static boolean isEqualOperator(String operatorCode, int containsOperator) {
      var literal = CCLParser.VOCABULARY.getLiteralName(containsOperator);

      if (literal == null || operatorCode == null) {
        return false;
      }

      return operatorCode.equals(literal.substring(1, literal.length() - 1));
    }
  }

  private Operator operator = Operator.EQUALS;

  private Value value;

  public Value getValue() {
    return this.value;
  }

  @Override
  public boolean evaluate(Map properties) {
    // getScanners the value of the field
    var fieldValue = getValueFromField(properties, this.field);

    // TODO: just converting to long is probably a stupid idea
    switch (this.operator) {
      case CONTAINS:
        if (fieldValue instanceof String && value.getValue() instanceof CharSequence) {
          return ((String) fieldValue).contains((CharSequence) value.getValue());
        }
        // TODO: throw exception here?
        return false;
      case NOT_EQUALS:
        return !Objects.equals(fieldValue, value.getValue());
      case LESS_THAN:
        return longOf(fieldValue) < longOf(value.getValue());
      case LESS_THAN_OR_EQUALS:
        return longOf(fieldValue) <= longOf(value.getValue());
      case GREATER_THAN:
        return longOf(fieldValue) > longOf(value.getValue());
      case GREATER_THAN_OR_EQUALS:
        return longOf(fieldValue) >= longOf(value.getValue());
      case EQUALS:
      default:
        return Objects.equals(fieldValue, value.getValue());
    }
  }

  public static long longOf(Object object) {
    if (object == null) {
      return 0;
    }

    try {
      return Long.valueOf(object.toString());
    } catch (NumberFormatException ex) {
      return 0;
    }
  }

  public void setValue(Value value) {
    this.value = value;
  }

  public Operator getOperator() {
    return operator;
  }

  public void setOperator(Operator operator) {
    this.operator = operator;
  }
}
