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
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/sirupsen/logrus"

	"clouditor.io/clouditor/ccl/parser"
)

var log *logrus.Entry

var (
	ErrUnsupportedContext = errors.New("unsupported context")
	ErrFieldNameNotFound  = errors.New("invalid field name")
	ErrFieldNoMap         = errors.New("field is not a map")
)

func init() {
	log = logrus.WithField("component", "ccl")
}

func RunRule(data string, object map[string]interface{}) (bool, error) {
	log.Debugf("Evaluating '%s'...", data)

	input := antlr.NewInputStream(data)

	return runRule(input, object)
}

func RunRuleFromFile(file string, object map[string]interface{}) (bool, error) {
	log.Debugf("Evaluating rule %s...", file)

	input, _ := antlr.NewFileStream(file)

	return runRule(input, object)
}

func runRule(input antlr.CharStream, object map[string]interface{}) (bool, error) {
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

	return false, ErrUnsupportedContext
}

func evaluateExpression(c parser.IExpressionContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.ExpressionContext); ok {
		if v.SimpleExpression() != nil {
			return evaluateSimpleExpression(v.SimpleExpression(), o)
		} else if v.NotExpression() != nil {
			return evaluateNotExpression(v.NotExpression(), o)
		}
	}

	return false, ErrUnsupportedContext
}

func evaluateSimpleExpression(c parser.ISimpleExpressionContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.SimpleExpressionContext); ok {
		if v.Comparison() != nil {
			return evaluateComparison(v.Comparison(), o)
		}
	}

	return false, ErrUnsupportedContext
}

func evaluateNotExpression(c parser.INotExpressionContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.NotExpressionContext); ok {
		// evalute the expression and negate it
		success, err = evaluateExpression(v.Expression(), o)
		success = !success

		return
	}

	return false, ErrUnsupportedContext
}

func evaluateComparison(c parser.IComparisonContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.ComparisonContext); ok {
		if v.BinaryComparison() != nil {
			return evaluteBinaryComparison(v.BinaryComparison(), o)
		}
	}

	return false, ErrUnsupportedContext
}

func evaluteBinaryComparison(c parser.IBinaryComparisonContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.BinaryComparisonContext); ok {
		// now the fun begins
		identifier := v.Field().(*parser.FieldContext).Identifier().GetText()

		var fieldValue interface{}

		if fieldValue, err = evaluateField(identifier, o); err != nil {
			return false, fmt.Errorf("could not evaluate field %s: %w", identifier, err)
		}

		var comparisonValue interface{}
		if comparisonValue, err = evaluateLiteral(v.Value().(*parser.ValueContext)); err != nil {
			return false, fmt.Errorf("could not parse comparison value: %w", err)
		}

		return evaluateEquals(fieldValue, comparisonValue)
	}

	return false, ErrUnsupportedContext
}

func evaluateField(field string, o map[string]interface{}) (interface{}, error) {
	var (
		idx   int
		value interface{}
		ok    bool
	)

	if idx = strings.Index(field, "."); idx == -1 {
		// no sub field left, directly access it
		if value, ok = o[field]; !ok {
			return nil, ErrFieldNameNotFound
		}

		return value, nil
	} else {
		// check for the first part; if that does not exist, we do not need to continue
		firstPart := field[:idx]
		if value, ok = o[firstPart]; !ok {
			return nil, ErrFieldNameNotFound
		}

		mapValue, ok := value.(map[string]interface{})
		if !ok {
			return nil, ErrFieldNoMap
		}

		// recursivly check for the second part
		secondPart := field[idx+1:]

		return evaluateField(secondPart, mapValue)
	}
}

func evaluateLiteral(value *parser.ValueContext) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	if value.StringLiteral() != nil {
		return evaluateStringLiteral(value.StringLiteral()), nil
	} else if value.IntNumber() != nil {
		return evaluteIntegerLiteral(value.IntNumber())
	} else if value.FloatNumber() != nil {
		return evaluateFloatLiteral(value.FloatNumber())
	} else if value.BooleanLiteral() != nil {
		return evaluateBoolLiteral(value.BooleanLiteral())
	}

	return nil, fmt.Errorf("could not evalute literal: %v", value.GetText())
}

func evaluateStringLiteral(node antlr.TerminalNode) string {
	s := node.GetText()

	return strings.Trim(s, "\"'")
}

func evaluteIntegerLiteral(node antlr.TerminalNode) (int64, error) {
	return strconv.ParseInt(node.GetText(), 10, 64)
}

func evaluateFloatLiteral(node antlr.TerminalNode) (float64, error) {
	return strconv.ParseFloat(node.GetText(), 64)
}

func evaluateBoolLiteral(node antlr.TerminalNode) (bool, error) {
	return strconv.ParseBool(node.GetText())
}

func evaluateEquals(lhs interface{}, rhs interface{}) (success bool, err error) {
	switch v := lhs.(type) {
	case string:
		return v == fmt.Sprintf("%v", rhs), nil
	case int:
		return compareInt(int64(v), rhs, 64)
	case int16:
		return compareInt(int64(v), rhs, 64)
	case int32:
		return compareInt(int64(v), rhs, 64)
	case int64:
		return compareInt(int64(v), rhs, 64)
	case float32:
		return compareFloat(float64(v), rhs, 64)
	case float64:
		return compareFloat(float64(v), rhs, 64)
	case bool:
		return compareBool(v, rhs)
	}

	return false, fmt.Errorf("could not compare %+v and %+v", lhs, rhs)
}

func compareInt(v int64, rhs interface{}, bitSize int) (success bool, err error) {
	var i int64
	// try to convert rhs to integer
	if i, err = strconv.ParseInt(fmt.Sprintf("%v", rhs), 10, bitSize); err != nil {
		return false, fmt.Errorf("could not compare %+v and %+v: %w", v, rhs, err)
	}

	return v == i, nil
}

func compareFloat(v float64, rhs interface{}, bitSize int) (success bool, err error) {
	var f float64
	// try to convert rhs to integer
	if f, err = strconv.ParseFloat(fmt.Sprintf("%v", rhs), bitSize); err != nil {
		return false, fmt.Errorf("could not compare %+v and %+v: %w", v, rhs, err)
	}

	return v == f, nil
}

func compareBool(v bool, rhs interface{}) (success bool, err error) {
	var b bool
	// try to convert rhs to integer
	if b, err = strconv.ParseBool(fmt.Sprintf("%v", rhs)); err != nil {
		return false, fmt.Errorf("could not compare %+v and %+v: %w", v, rhs, err)
	}

	return v == b, nil
}
