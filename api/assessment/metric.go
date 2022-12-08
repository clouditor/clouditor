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

package assessment

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"clouditor.io/clouditor/persistence"
	"github.com/google/uuid"
)

var (
	ErrMetricNameMissing       = errors.New("metric name is missing")
	ErrMetricEmpty             = errors.New("metric is missing or empty")
	ErrCloudServiceIDIsMissing = errors.New("cloud service id is missing")
	ErrCloudServiceIDIsInvalid = errors.New("cloud service id is invalid")
)

func (r *Range) UnmarshalJSON(b []byte) (err error) {
	// Check for the different range types
	var (
		r1 Range_AllowedValues
		r2 Range_Order
		r3 Range_MinMax
	)

	if err = json.Unmarshal(b, &r1); err == nil && r1.AllowedValues != nil {
		r.Range = &r1
		return
	}

	if err = json.Unmarshal(b, &r2); err == nil && r2.Order != nil {
		r.Range = &r2
		return
	}

	if err = json.Unmarshal(b, &r3); err == nil && r3.MinMax != nil {
		r.Range = &r3
		return
	}

	return
}

// MarshalJSON is a custom implementation of JSON marshalling to correctly
// serialize the Range type because the inner types, such as Range_AllowedValues
// are missing json struct tags. This is needed if the Range type is marshalled
// on its own (for example) as a single field in a database. In gRPC messages,
// the protojson.Marshal function takes care of this.
func (r *Range) MarshalJSON() (b []byte, err error) {
	switch v := r.Range.(type) {
	case *Range_AllowedValues:
		return json.Marshal(&struct {
			AllowedValues *AllowedValues `json:"allowedValues"`
		}{
			AllowedValues: v.AllowedValues,
		})
	case *Range_Order:
		return json.Marshal(&struct {
			Order *Order `json:"order"`
		}{
			Order: v.Order,
		})
	case *Range_MinMax:
		return json.Marshal(&struct {
			MinMax *MinMax `json:"minMax"`
		}{
			MinMax: v.MinMax,
		})
	default:
		return nil, persistence.ErrUnsupportedType
	}
}

// Value implements https://pkg.go.dev/database/sql/driver#Valuer to indicate
// how this struct will be saved into an SQL database field.
func (r *Range) Value() (val driver.Value, err error) {
	if r == nil {
		return
	} else {
		val, err = json.Marshal(r)
		if err != nil {
			err = fmt.Errorf("could not marshal JSON: %w", err)
		}

		return
	}
}

// Scan implements https://pkg.go.dev/database/sql#Scanner to indicate how
// this struct can be loaded from an SQL database field.
func (r *Range) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case []byte:
		err = json.Unmarshal(v, r)
		if err != nil {
			err = fmt.Errorf("could not unmarshal JSON: %w", err)
		}
	default:
		err = persistence.ErrUnsupportedType
	}

	return
}

// GormDataType implements GormDataTypeInterface to give an indication how
// this struct will be serialized into a database using GORM.
func (*Range) GormDataType() string {
	return "jsonb"
}

// Hash provides a simple string based hash for this metric configuration. It can be used
// to provide a key for a map or a cache.
func (x *MetricConfiguration) Hash() string {
	return base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%v-%v", x.Operator, x.TargetValue)))
}

// CheckCloudServiceID checks if serviceID is available and in the valid UUID format.
func CheckCloudServiceID(serviceID string) error {
	// Check if ServiceId is missing
	if serviceID == "" {
		return ErrCloudServiceIDIsMissing
	}

	// Check if ServiceId is valid
	if _, err := uuid.Parse(serviceID); err != nil {
		return ErrCloudServiceIDIsInvalid
	}

	return nil
}
