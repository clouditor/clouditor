// Copyright 2021 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//           $$\                           $$\ $$\   $$\
//           $$ |                          $$ |\__|  $$ |
//  $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
// $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
// $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ | \__|
// $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
// \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
//  \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
//
// This file is part of Clouditor Community Edition.

// Package ccl contains the the Cloud Compliance Language (CCL)
package ccl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/sirupsen/logrus"

	"clouditor.io/clouditor/ccl/parser"
)

var log *logrus.Entry

func init() {
	log = logrus.WithField("component", "ccl")
}

func RunRule(file string, object map[string]interface{}) (bool, error) {
	log.Debugf("Evaluating rule %s...", file)

	input, _ := antlr.NewFileStream(file)
	lexer := parser.NewCCLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewCCLParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	tree := p.Condition()

	return evaluateCondition(tree, object)
}

func evaluateCondition(c parser.IConditionContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.ConditionContext); ok {
		return evaluateExpression(v.Expression(), o)
	}

	return false, errors.New("unsupported context")
}

func evaluateExpression(c parser.IExpressionContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.ExpressionContext); ok {
		if v.SimpleExpression() != nil {
			return evaluateSimpleExpression(v.SimpleExpression(), o)
		}
	}

	return false, errors.New("unsupported context")
}

func evaluateSimpleExpression(c parser.ISimpleExpressionContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.SimpleExpressionContext); ok {
		if v.Comparison() != nil {
			return evaluateComparison(v.Comparison(), o)
		}
	}

	return false, errors.New("unsupported context")
}

func evaluateComparison(c parser.IComparisonContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.ComparisonContext); ok {
		if v.BinaryComparison() != nil {
			return evaluteBinaryComparison(v.BinaryComparison(), o)
		}
	}

	return false, errors.New("unsupported context")
}

func evaluteBinaryComparison(c parser.IBinaryComparisonContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.BinaryComparisonContext); ok {
		// now the fun begins
		fieldIdentifier := v.Field().(*parser.FieldContext).Identifier().GetText()

		fieldValue := o[fieldIdentifier]

		comparisonValue := evaluateStringLiteral(v.Value().(*parser.ValueContext).StringLiteral())

		return evaluateEquals(fieldValue, comparisonValue)
	}

	return false, errors.New("unsupported context")
}

func evaluateStringLiteral(node antlr.TerminalNode) string {
	s := node.GetText()

	return strings.Trim(s, "\"'")
}

func evaluateEquals(lhs interface{}, rhs interface{}) (success bool, err error) {
	switch v := lhs.(type) {
	case string:
		return v == fmt.Sprintf("%v", rhs), nil
	}

	return false, fmt.Errorf("could not compare %+v and %+v", lhs, rhs)
}
