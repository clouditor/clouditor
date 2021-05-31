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

// Package ccl contains the the Cloud Compliance Language (CCL). It is a simple domain specific language to model rules
// that should apply to discovered resources.
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
		if v.IsEmptyExpression() != nil {
			return evaluateIsEmptyExpression(v.IsEmptyExpression(), o)
		} else if v.Comparison() != nil {
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

func evaluateIsEmptyExpression(c parser.IIsEmptyExpressionContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.IsEmptyExpressionContext); ok {
		// evalute the field
		value, err := evaluateField(v.Field().GetText(), o)
		if errors.Is(err, ErrFieldNameNotFound) {
			// if the field is not found, it is also empty
			return true, nil
		} else if err != nil {
			return false, fmt.Errorf("could not evaluate field: %w", err)
		}

		// otherwise, check, if the value is an "empty" value
		return value == "" || value == 0 || value == false, nil
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

		var (
			lhs interface{}
			rhs interface{}
			op  Operator
		)

		if lhs, err = evaluateField(identifier, o); err != nil {
			return false, fmt.Errorf("could not evaluate field %s: %w", identifier, err)
		}

		if rhs, err = evaluateLiteral(v.Value().(*parser.ValueContext)); err != nil {
			return false, fmt.Errorf("could not parse comparison value: %w", err)
		}

		if op, err = evaluateOperator(v.Operator()); err != nil {
			return false, fmt.Errorf("could not parse operator: %w", err)
		}

		switch v := lhs.(type) {
		case string:
			return compareString(v, rhs, op)
		case int:
			return compareInt(int64(v), rhs, 64, op)
		case int16:
			return compareInt(int64(v), rhs, 64, op)
		case int32:
			return compareInt(int64(v), rhs, 64, op)
		case int64:
			return compareInt(int64(v), rhs, 64, op)
		case float32:
			return compareFloat(float64(v), rhs, 64, op)
		case float64:
			return compareFloat(float64(v), rhs, 64, op)
		case bool:
			return compareBool(v, rhs, op)
		}

		return false, fmt.Errorf("could not compare %+v and %+v", lhs, rhs)
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

func evaluateOperator(c parser.IOperatorContext) (Operator, error) {
	if v, ok := c.(*parser.OperatorContext); ok {
		if v.EqualsOperator() != nil {
			return &equalsOperator{}, nil
		}

		if v.NotEqualsOperator() != nil {
			return &notOperator{&equalsOperator{}}, nil
		}

		if v.LessThanOperator() != nil {
			return &lessOperator{}, nil
		}

		if v.LessOrEqualsThanOperator() != nil {
			return &lessOrEqualsOperator{}, nil
		}

		if v.MoreThanOperator() != nil {
			return &greaterOperator{}, nil
		}

		if v.MoreOrEqualsThanOperator() != nil {
			return &greaterOrEqualsOperator{}, nil
		}

		if v.ContainsOperator() != nil {
			return &containsOperator{}, nil
		}
	}

	return nil, ErrUnsupportedContext
}

func compareString(v string, rhs interface{}, op Operator) (success bool, err error) {
	var s = fmt.Sprintf("%v", rhs)

	return op.CompareString(v, s)
}

func compareInt(v int64, rhs interface{}, bitSize int, op Operator) (success bool, err error) {
	var i int64
	// try to convert rhs to integer
	if i, err = strconv.ParseInt(fmt.Sprintf("%v", rhs), 10, bitSize); err != nil {
		return false, fmt.Errorf("could not compare %+v and %+v: %w", v, rhs, err)
	}

	return op.CompareInt(v, i)
}

func compareFloat(v float64, rhs interface{}, bitSize int, op Operator) (success bool, err error) {
	var f float64
	// try to convert rhs to integer
	if f, err = strconv.ParseFloat(fmt.Sprintf("%v", rhs), bitSize); err != nil {
		return false, fmt.Errorf("could not compare %+v and %+v: %w", v, rhs, err)
	}

	return op.CompareFloat(v, f)
}

func compareBool(v bool, rhs interface{}, op Operator) (success bool, err error) {
	var b bool
	// try to convert rhs to integer
	if b, err = strconv.ParseBool(fmt.Sprintf("%v", rhs)); err != nil {
		return false, fmt.Errorf("could not compare %+v and %+v: %w", v, rhs, err)
	}

	return op.CompareBool(v, b)
}
