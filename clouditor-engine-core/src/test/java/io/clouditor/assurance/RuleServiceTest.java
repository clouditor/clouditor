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

  private final Engine engine = new Engine();

  @Test
  void testParseMarkdown() throws IOException {
    var rule =
        this.engine
            .getService(RuleService.class)
            .loadRule(FileSystemManager.getInstance().getPathForResource("rules/test/test.md"));

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
