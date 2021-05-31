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
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/sirupsen/logrus"

	"clouditor.io/clouditor/ccl/parser"
)

var log *logrus.Entry

var (
	ErrUnsupportedContext   = errors.New("unsupported context")
	ErrFieldNameNotFound    = errors.New("invalid field name")
	ErrFieldNoMap           = errors.New("field is not a map")
	ErrFieldNoTime          = errors.New("field has no time value")
	ErrFieldNoArray         = errors.New("field is not an array")
	ErrInvalidScope         = errors.New("invalid scope in in-expression")
	ErrUnexpectedExpression = errors.New("unexpected expression")
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
		} else if v.InExpression() != nil {
			return evaluateInExpression(v.InExpression(), o)
		} else {
			return false, ErrUnexpectedExpression
		}
	}

	return false, ErrUnsupportedContext
}

func evaluateSimpleExpression(c parser.ISimpleExpressionContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.SimpleExpressionContext); ok {
		if v.IsEmptyExpression() != nil {
			return evaluateIsEmptyExpression(v.IsEmptyExpression(), o)
		} else if v.WithinExpression() != nil {
			return evaluateWithinExpression(v.WithinExpression(), o)
		} else if v.Comparison() != nil {
			return evaluateComparison(v.Comparison(), o)
		} else if v.Expression() != nil {
			return evaluateExpression(v.Expression(), o)
		} else {
			return false, ErrUnexpectedExpression
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

func evaluateInExpression(c parser.IInExpressionContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.InExpressionContext); ok {
		var value interface{}
		var arrayValue []map[string]interface{}

		if value, err = evaluateField(v.Field().GetText(), o); err != nil {
			return false, fmt.Errorf("could not evaluate in-expression: %w", err)
		}

		if arrayValue, ok = value.([]map[string]interface{}); !ok {
			return false, ErrFieldNoArray
		}

		var scope string
		if v.Scope().GetText() == "all" {
			scope = "all"
		} else if v.Scope().GetText() == "any" {
			scope = "any"
		} else {
			return false, ErrInvalidScope
		}

		var result bool

		// loop through array
		for _, item := range arrayValue {
			// if any matches
			if success, err = evaluateSimpleExpression(v.SimpleExpression(), item); err != nil {
				return false, fmt.Errorf("could not evaluate in-expression: %w", err)
			}

			if success && scope == "any" {
				return true, nil
			}

			result = result && success
		}

		return result, nil
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

// evaluateWithinExpression checks, if a field value is within a specified array
func evaluateWithinExpression(c parser.IWithinExpressionContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.WithinExpressionContext); ok {
		var (
			lhs interface{}
			rhs interface{}
			op  Operator
		)

		op = &equalsOperator{}

		identifier := v.Field().(*parser.FieldContext).Identifier().GetText()

		if lhs, err = evaluateField(identifier, o); err != nil {
			return false, fmt.Errorf("could not evaluate field %s: %w", identifier, err)
		}

		ctxs := v.AllValue()
		for _, ctx := range ctxs {
			if rhs, err = evaluateValue(ctx); err != nil {
				return false, fmt.Errorf("could not evaluate value: %w", err)
			}

			if success, err = compare(lhs, rhs, op); err != nil {
				// directly return error without wrapping, i.e. if the comparison was of a different type
				return false, err
			}

			// return early, if a match was found
			if success {
				return true, nil
			}
		}

		// no match found
		return false, nil
	}

	return false, ErrUnsupportedContext
}

func evaluateComparison(c parser.IComparisonContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.ComparisonContext); ok {
		if v.BinaryComparison() != nil {
			return evaluteBinaryComparison(v.BinaryComparison(), o)
		} else if v.TimeComparison() != nil {
			return evaluateTimeComparison(v.TimeComparison(), o)
		} else {
			return false, ErrUnexpectedExpression
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

		if rhs, err = evaluateValue(v.Value()); err != nil {
			return false, fmt.Errorf("could not evaluate value: %w", err)
		}

		if op, err = evaluateOperator(v.Operator()); err != nil {
			return false, fmt.Errorf("could not parse operator: %w", err)
		}

		return compare(lhs, rhs, op)
	}

	return false, ErrUnsupportedContext
}

func evaluateTimeComparison(c parser.ITimeComparisonContext, o map[string]interface{}) (success bool, err error) {
	if v, ok := c.(*parser.TimeComparisonContext); ok {
		var (
			lhs       interface{}
			timeValue time.Time
			duration  time.Duration
		)

		identifier := v.Field().(*parser.FieldContext).Identifier().GetText()

		if lhs, err = evaluateField(identifier, o); err != nil {
			return false, fmt.Errorf("could not evaluate field %s: %w", identifier, err)
		}

		if timeValue, ok = lhs.(time.Time); !ok {
			return false, ErrFieldNoTime
		}

		if v.Duration() != nil {
			if duration, err = evaluateDuration(v.Duration()); err != nil {
				return false, fmt.Errorf("could not compare times: %w", err)
			}
		} else if v.NowOperator() != nil {
			duration = 0
		}

		// four different operators, but all relative to the current time
		now := time.Now()

		// younger than X (days|minutes|seconds)
		if v.TimeOperator().GetText() == "younger" {
			// substract the duration specified in rhs from now (go back into the past) to get the
			// beginning of the period we consider 'younger'
			younger := now.Add(-duration)

			// check, if timeValue is after or equal to the 'younger' date
			return timeValue.After(younger) || timeValue.Equal(younger), nil
		}

		// older than X (seconds|days|months)
		if v.TimeOperator().GetText() == "older" {
			// basically, the same as younger, just before - not after
			older := now.Add(-duration)

			return timeValue.Before(older), nil
		}

		// before X (seconds|days|months)
		if v.TimeOperator().GetText() == "before" {
			// add the duration specified in rhs to now (going into the future) to get the
			// end of the period we consider 'before'
			before := now.Add(duration)

			// check if timeValue is before or equal to the 'before' date
			return timeValue.Before(before) || timeValue.Equal(before), nil
		}

		// after X (seconds|days|months)
		if v.TimeOperator().GetText() == "after" {
			// basically, the same as before, just after - not before
			after := now.Add(duration)

			// check if timeValue is before or equal to the 'before' date
			return timeValue.After(after), nil
		}
	}

	return false, ErrUnsupportedContext
}

func evaluateDuration(c parser.IDurationContext) (duration time.Duration, err error) {
	if v, ok := c.(*parser.DurationContext); ok {
		var value int64
		if value, err = evaluteIntegerLiteral(v.IntNumber()); err != nil {
			return 0, fmt.Errorf("invalid value: %w", err)
		}

		t := v.Unit().GetText()
		if t == "seconds" {
			return time.Duration(value) * time.Second, nil
		} else if t == "days" {
			return time.Duration(value) * time.Hour * 24, nil
		} else if t == "months" {
			return time.Duration(value) * time.Hour * 24 * 30, nil
		} else {
			return 0, fmt.Errorf("invalid unit duration: %s", t)
		}
	}

	return 0, ErrUnsupportedContext
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

func evaluateValue(c parser.IValueContext) (interface{}, error) {
	if v, ok := c.(*parser.ValueContext); ok {
		if v.StringLiteral() != nil {
			return evaluateStringLiteral(v.StringLiteral()), nil
		} else if v.IntNumber() != nil {
			return evaluteIntegerLiteral(v.IntNumber())
		} else if v.FloatNumber() != nil {
			return evaluateFloatLiteral(v.FloatNumber())
		} else if v.BooleanLiteral() != nil {
			return evaluateBoolLiteral(v.BooleanLiteral())
		}

		return nil, fmt.Errorf("could not evalute literal: %v", v.GetText())
	}

	return nil, ErrUnsupportedContext
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
			return &notOperator{greaterOperator{}}, nil
		}

		if v.MoreThanOperator() != nil {
			return &greaterOperator{}, nil
		}

		if v.MoreOrEqualsThanOperator() != nil {
			return &notOperator{lessOperator{}}, nil
		}

		if v.ContainsOperator() != nil {
			return &containsOperator{}, nil
		}
	}

	return nil, ErrUnsupportedContext
}

func compare(lhs interface{}, rhs interface{}, op Operator) (success bool, err error) {
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
