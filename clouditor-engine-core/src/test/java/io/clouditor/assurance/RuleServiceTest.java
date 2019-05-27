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

package io.clouditor.assurance;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

import io.clouditor.Engine;
import io.clouditor.assurance.ccl.BinaryComparison;
import io.clouditor.assurance.ccl.BinaryComparison.Operator;
import io.clouditor.util.FileSystemManager;
import java.io.IOException;
import org.junit.jupiter.api.Test;

class RuleServiceTest {

  private Engine engine = new Engine();

  @Test
  void testParseMarkdown() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .parseFromMarkDown(
                FileSystemManager.getInstance().getPathForResource("rules/test/test.md"));

    assertNotNull(rule);
    assertEquals("Title", rule.getName());
    assertNotNull(rule.getDescription());

    assertEquals(1, rule.getConditions().size());

    var condition = rule.getConditions().get(0);

    assertEquals("Asset", rule.getAssetType());

    var expression = condition.getExpression();

    assertTrue(expression instanceof BinaryComparison);
    assertEquals(Operator.EQUALS, ((BinaryComparison) expression).getOperator());
    assertEquals("value", ((BinaryComparison) expression).getValue().getValue());
    assertEquals("property", ((BinaryComparison) expression).getField());

    assertEquals(2, rule.getControls().size());
  }
}
