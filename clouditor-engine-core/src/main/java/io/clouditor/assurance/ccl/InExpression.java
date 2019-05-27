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
