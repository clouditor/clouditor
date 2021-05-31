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

package ccl

import (
	"errors"
	"strings"
)

type Operator interface {
	CompareString(lhs string, rhs string) (bool, error)
	CompareInt(lhs int64, rhs int64) (bool, error)
	CompareFloat(lhs float64, rhs float64) (bool, error)
	CompareBool(lhs bool, rhs bool) (bool, error)
}

var (
	ErrOperationNotSupported = errors.New("the specified operation is not supported between the two value types")
)

type equalsOperator struct{}

func (e equalsOperator) CompareString(lhs string, rhs string) (bool, error)  { return lhs == rhs, nil }
func (e equalsOperator) CompareInt(lhs int64, rhs int64) (bool, error)       { return lhs == rhs, nil }
func (e equalsOperator) CompareFloat(lhs float64, rhs float64) (bool, error) { return lhs == rhs, nil }
func (e equalsOperator) CompareBool(lhs bool, rhs bool) (bool, error)        { return lhs == rhs, nil }

type lessOperator struct{}

func (e lessOperator) CompareString(lhs string, rhs string) (bool, error)  { return lhs < rhs, nil }
func (e lessOperator) CompareInt(lhs int64, rhs int64) (bool, error)       { return lhs < rhs, nil }
func (e lessOperator) CompareFloat(lhs float64, rhs float64) (bool, error) { return lhs < rhs, nil }
func (e lessOperator) CompareBool(lhs bool, rhs bool) (bool, error) {
	return false, ErrOperationNotSupported
}

type greaterOperator struct{}

func (e greaterOperator) CompareString(lhs string, rhs string) (bool, error)  { return lhs > rhs, nil }
func (e greaterOperator) CompareInt(lhs int64, rhs int64) (bool, error)       { return lhs > rhs, nil }
func (e greaterOperator) CompareFloat(lhs float64, rhs float64) (bool, error) { return lhs > rhs, nil }
func (e greaterOperator) CompareBool(lhs bool, rhs bool) (bool, error) {
	return false, ErrOperationNotSupported
}

type containsOperator struct{}

func (e containsOperator) CompareString(lhs string, rhs string) (bool, error) {
	return strings.Contains(lhs, rhs), nil
}
func (e containsOperator) CompareInt(lhs int64, rhs int64) (bool, error) {
	return false, ErrOperationNotSupported
}
func (e containsOperator) CompareFloat(lhs float64, rhs float64) (bool, error) {
	return false, ErrOperationNotSupported
}
func (e containsOperator) CompareBool(lhs bool, rhs bool) (bool, error) {
	return false, ErrOperationNotSupported
}

type notOperator struct {
	op Operator
}

func (e notOperator) CompareString(lhs string, rhs string) (b bool, err error) {
	if b, err = e.op.CompareString(lhs, rhs); err != nil {
		return false, err
	}

	return !b, nil
}
func (e notOperator) CompareInt(lhs int64, rhs int64) (b bool, err error) {
	if b, err = e.op.CompareInt(lhs, rhs); err != nil {
		return false, err
	}

	return !b, nil
}
func (e notOperator) CompareFloat(lhs float64, rhs float64) (b bool, err error) {
	if b, err = e.op.CompareFloat(lhs, rhs); err != nil {
		return false, err
	}

	return !b, nil
}
func (e notOperator) CompareBool(lhs bool, rhs bool) (b bool, err error) {
	if b, err = e.op.CompareBool(lhs, rhs); err != nil {
		return false, err
	}

	return !b, nil
}
