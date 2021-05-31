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

type notEqualsOperator struct{}

func (e notEqualsOperator) CompareString(lhs string, rhs string) (bool, error) {
	return lhs != rhs, nil
}
func (e notEqualsOperator) CompareInt(lhs int64, rhs int64) (bool, error) { return lhs != rhs, nil }
func (e notEqualsOperator) CompareFloat(lhs float64, rhs float64) (bool, error) {
	return lhs != rhs, nil
}
func (e notEqualsOperator) CompareBool(lhs bool, rhs bool) (bool, error) { return lhs != rhs, nil }

type lessOperator struct{}

func (e lessOperator) CompareString(lhs string, rhs string) (bool, error)  { return lhs < rhs, nil }
func (e lessOperator) CompareInt(lhs int64, rhs int64) (bool, error)       { return lhs < rhs, nil }
func (e lessOperator) CompareFloat(lhs float64, rhs float64) (bool, error) { return lhs < rhs, nil }
func (e lessOperator) CompareBool(lhs bool, rhs bool) (bool, error) {
	return false, ErrOperationNotSupported
}

type lessOrEqualsOperator struct{}

func (e lessOrEqualsOperator) CompareString(lhs string, rhs string) (bool, error) {
	return lhs <= rhs, nil
}
func (e lessOrEqualsOperator) CompareInt(lhs int64, rhs int64) (bool, error) { return lhs <= rhs, nil }
func (e lessOrEqualsOperator) CompareFloat(lhs float64, rhs float64) (bool, error) {
	return lhs <= rhs, nil
}
func (e lessOrEqualsOperator) CompareBool(lhs bool, rhs bool) (bool, error) {
	return false, ErrOperationNotSupported
}

type greaterOperator struct{}

func (e greaterOperator) CompareString(lhs string, rhs string) (bool, error)  { return lhs > rhs, nil }
func (e greaterOperator) CompareInt(lhs int64, rhs int64) (bool, error)       { return lhs > rhs, nil }
func (e greaterOperator) CompareFloat(lhs float64, rhs float64) (bool, error) { return lhs > rhs, nil }
func (e greaterOperator) CompareBool(lhs bool, rhs bool) (bool, error) {
	return false, ErrOperationNotSupported
}

type greaterOrEqualsOperator struct{}

func (e greaterOrEqualsOperator) CompareString(lhs string, rhs string) (bool, error) {
	return lhs >= rhs, nil
}
func (e greaterOrEqualsOperator) CompareInt(lhs int64, rhs int64) (bool, error) {
	return lhs >= rhs, nil
}
func (e greaterOrEqualsOperator) CompareFloat(lhs float64, rhs float64) (bool, error) {
	return lhs >= rhs, nil
}
func (e greaterOrEqualsOperator) CompareBool(lhs bool, rhs bool) (bool, error) {
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
