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

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.deser.std.StdDeserializer;
import io.clouditor.assurance.ccl.BinaryComparison.Operator;
import io.clouditor.assurance.ccl.InExpression.Scope;
import io.clouditor.assurance.ccl.TimeComparison.TimeOperator;
import io.clouditor.assurance.grammar.CCLBaseVisitor;
import io.clouditor.assurance.grammar.CCLLexer;
import io.clouditor.assurance.grammar.CCLParser;
import io.clouditor.assurance.grammar.CCLParser.BinaryComparisonContext;
import io.clouditor.assurance.grammar.CCLParser.ConditionContext;
import io.clouditor.assurance.grammar.CCLParser.EmptyExpressionContext;
import io.clouditor.assurance.grammar.CCLParser.InExpressionContext;
import io.clouditor.assurance.grammar.CCLParser.NotExpressionContext;
import io.clouditor.assurance.grammar.CCLParser.SimpleExpressionContext;
import io.clouditor.assurance.grammar.CCLParser.TimeComparisonContext;
import io.clouditor.assurance.grammar.CCLParser.ValueContext;
import io.clouditor.assurance.grammar.CCLParser.WithinExpressionContext;
import java.io.IOException;
import java.time.temporal.ChronoUnit;
import java.util.stream.Collectors;
import org.antlr.v4.runtime.CharStreams;
import org.antlr.v4.runtime.CommonTokenStream;

public class CCLDeserializer extends StdDeserializer<Condition> {

  public CCLDeserializer() {
    super(Condition.class);
  }

  @Override
  public Condition deserialize(JsonParser p, DeserializationContext ctxt) throws IOException {
    JsonNode node = p.getCodec().readTree(p);

    return this.parse(node.asText());
  }

  private static class ConditionVisitor extends CCLBaseVisitor<Condition> {

    @Override
    public Condition visitCondition(ConditionContext ctx) {
      if (ctx.assetType() != null && ctx.expression() != null) {
        var condition = new Condition();

        condition.setAssetType(ctx.assetType().getText());
        condition.setExpression(ctx.expression().accept(new ExpressionListener()));

        return condition;
      }

      return super.visitCondition(ctx);
    }
  }

  private static class ExpressionListener extends CCLBaseVisitor<Expression> {

    @Override
    public Expression visitSimpleExpression(SimpleExpressionContext ctx) {
      if (ctx.comparison() != null) {
        return ctx.comparison().accept(this);
      } else if (ctx.emptyExpression() != null) {
        return ctx.emptyExpression().accept(this);
      } else if (ctx.expression() != null) {
        // just a simple wrapped expression with parenthesis
        return ctx.expression().accept(this);
      }

      return super.visitSimpleExpression(ctx);
    }

    @Override
    public Expression visitBinaryComparison(BinaryComparisonContext ctx) {
      if (ctx.field() != null && ctx.operator() != null && ctx.value() != null) {
        var comparison = new BinaryComparison();
        comparison.setField(ctx.field().getText());
        comparison.setValue(ctx.value().accept(new ValueListener()));
        comparison.setOperator(Operator.of(ctx.operator().getText()));

        return comparison;
      }

      return super.visitBinaryComparison(ctx);
    }

    @Override
    public Expression visitTimeComparison(TimeComparisonContext ctx) {
      if (ctx.field() != null && ctx.timeOperator() != null) {
        var comparison = new TimeComparison();
        comparison.setField(ctx.field().getText());
        comparison.setTimeOperator(
            TimeOperator.valueOf(ctx.timeOperator().getText().toUpperCase()));

        if (ctx.time() != null && ctx.unit() != null) {
          comparison.setRelativeValue(Integer.parseInt(ctx.time().getText()));
          comparison.setTimeUnit(ChronoUnit.valueOf(ctx.unit().getText().toUpperCase()));
        } // its now per default (relative value 0)

        return comparison;
      }

      return super.visitTimeComparison(ctx);
    }

    @Override
    public Expression visitEmptyExpression(EmptyExpressionContext ctx) {
      if (ctx.field() != null) {
        var expression = new EmptyExpression();
        expression.setField(ctx.field().getText());

        return expression;
      }

      return super.visitEmptyExpression(ctx);
    }

    @Override
    public Expression visitNotExpression(NotExpressionContext ctx) {
      if (ctx.expression() != null) {
        var expression = new NotExpression();
        expression.setExpression(ctx.expression().accept(this));

        return expression;
      }

      return super.visitNotExpression(ctx);
    }

    @Override
    public Expression visitInExpression(InExpressionContext ctx) {
      if (ctx.field() != null && ctx.scope() != null && ctx.simpleExpression() != null) {
        var expression = new InExpression();
        expression.setField(ctx.field().getText());
        expression.setExpression(ctx.simpleExpression().accept(this));
        expression.setScope(Scope.of(ctx.scope().getText()));

        return expression;
      }

      return super.visitInExpression(ctx);
    }

    @Override
    public Expression visitWithinExpression(WithinExpressionContext ctx) {
      if (ctx.field() != null && ctx.value() != null) {
        var expression = new WithinExpression();
        expression.setField(ctx.field().getText());
        expression.setValues(
            ctx.value().stream()
                .map(valueContext -> valueContext.accept(new ValueListener()))
                .collect(Collectors.toList()));

        return expression;
      }

      return super.visitWithinExpression(ctx);
    }
  }

  private static class ValueListener extends CCLBaseVisitor<Value> {

    @Override
    public Value visitValue(ValueContext ctx) {
      if (ctx.BooleanLiteral() != null) {
        return new Value<>(Boolean.valueOf(ctx.BooleanLiteral().getText()));
      } else if (ctx.StringLiteral() != null) {
        var text = ctx.StringLiteral().getSymbol().getText();

        text = text.substring(1, text.length() - 1);

        return new Value<>(text);
      } else if (ctx.Number() != null) {
        return new Value<>(Long.valueOf(ctx.Number().getText()));
      }

      return super.visitValue(ctx);
    }
  }

  public Condition parse(String source) {
    var lexer = new CCLLexer(CharStreams.fromString(source));
    var tokens = new CommonTokenStream(lexer);
    var parser = new CCLParser(tokens);

    var visitor = new ConditionVisitor();

    var condition = visitor.visit(parser.condition());
    condition.setSource(source);

    return condition;
  }
}
