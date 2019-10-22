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

import java.util.Collection;
import java.util.Map;
import java.util.Objects;

public class InExpression extends Expression {

  private Scope scope = Scope.ANY;
  private String field;
  private Expression expression;

  public void setExpression(Expression expression) {
    this.expression = expression;
  }

  public void setField(String field) {
    this.field = field;
  }

  public void setScope(Scope scope) {
    this.scope = scope;
  }

  @Override
  public boolean evaluate(Map properties) {
    var fieldValue = getValueFromField(properties, this.field);

    Collection<Map> list;

    if (fieldValue instanceof Collection) {
      list = (Collection) fieldValue;
    } else if (fieldValue instanceof Map) {
      list = ((Map) fieldValue).values();
    } else {
      // TODO: or throw exception?
      return false;
    }

    if (scope == Scope.ALL) {
      return list.stream().allMatch(x -> this.expression.evaluate(x));
    } else {
      return list.stream().anyMatch(x -> this.expression.evaluate(x));
    }
  }

  public enum Scope {
    ANY,
    ALL;

    public static Scope of(String text) {
      return Objects.equals("all", text.toLowerCase()) ? ALL : ANY;
    }
  }
}
