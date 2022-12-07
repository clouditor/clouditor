// Copyright 2022 Fraunhofer AISEC
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

package evaluation

import (
	"errors"
)

var (
	ErrRequestIsEmpty              = errors.New("request is empty")
	ErrControlIDIsMissing          = errors.New("control id is missing")
	ErrCategoryNameIsMissing       = errors.New("category name is missing")
	ErrEvalControlIsMissing        = errors.New("controlId and category name mapping is missing")
	ErrTargetOfEvaluationIsInvalid = errors.New("target of evaluation is not valid")
)

// Validate validates the evaluate request
func (r *StartEvaluationRequest) Validate() (err error) {
	if r == nil {
		return ErrRequestIsEmpty
	}

	if err = r.TargetOfEvaluation.Validate(); err != nil {
		return err
	}

	return
}

// Validate validates the evaluate request
// TODO(anatheka): Update if everything is ready
func (r *StopEvaluationRequest) Validate() (err error) {
	if r == nil {
		return ErrRequestIsEmpty
	}

	if err = r.TargetOfEvaluation.Validate(); err != nil {
		return err
	}

	return
}
