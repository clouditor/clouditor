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
