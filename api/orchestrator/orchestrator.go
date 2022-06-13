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

package orchestrator

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"

	"k8s.io/utils/strings/slices"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/internal/util"
)

type CloudServiceHookFunc func(ctx context.Context, cld *CloudService, err error)

var (
	ErrRequestIsNil      = errors.New("request is empty")
	ErrCertificateIsNil  = errors.New("certificate is empty")
	ErrServiceIsNil      = errors.New("service is empty")
	ErrNameIsMissing     = errors.New("service name is empty")
	ErrIDIsMissing       = errors.New("service ID is empty")
	ErrCertIDIsMissing   = errors.New("certificate ID is empty")
	ErrInvalidColumnName = errors.New("column name is invalid")
)

// Value implements https://pkg.go.dev/database/sql/driver#Valuer to indicate
// how this struct will be saved into an SQL database field.
func (c *CloudService_Requirements) Value() (driver.Value, error) {
	if c == nil || c.RequirementIds == nil {
		return nil, nil
	} else {
		return driver.Value(strings.Join(c.RequirementIds, ",")), nil
	}
}

// Scan implements https://pkg.go.dev/database/sql#Scanner to indicate how
// this struct can be loaded from an SQL database field.
func (c *CloudService_Requirements) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		(*c).RequirementIds = strings.Split(v, ",")
	default:
		return errors.New("unsupported type")
	}

	return nil
}

// GormDataType implements GormDataTypeInterface to give an indication how
// this struct will be serialized into a database using GORM.
func (*CloudService_Requirements) GormDataType() string {
	return "string"
}

// Validate validates a ListCertificatesRequest
func (req *ListCertificatesRequest) Validate() (err error) {
	// req must be non-nil
	if req == nil {
		err = ErrRequestIsNil
		return
	}

	// Avoid DB injections by whitelisting the valid orderBy statements
	whitelist, err := util.GetFieldNames(Certificate{})
	// Add empty string indicating no explicit ordering
	whitelist = append(whitelist, "")
	if !slices.Contains(whitelist, req.OrderBy) {
		err = ErrInvalidColumnName
		return
	}

	return
}

// Validate validates a ListCloudServicesRequest
func (req *ListCloudServicesRequest) Validate() (err error) {
	// req must be non-nil
	if req == nil {
		err = ErrRequestIsNil
		return
	}

	// Avoid DB injections by whitelisting the valid orderBy statements
	whitelist, err := util.GetFieldNames(CloudService{})
	// Add empty string indicating no explicit ordering
	whitelist = append(whitelist, "")
	if !slices.Contains(whitelist, req.OrderBy) {
		err = ErrInvalidColumnName
		return
	}

	return
}

// Validate validates a ListMetricsRequest
func (req *ListMetricsRequest) Validate() (err error) {
	// req must be non-nil
	if req == nil {
		err = ErrRequestIsNil
		return
	}

	// Avoid DB injections by whitelisting the valid orderBy statements
	whitelist, err := util.GetFieldNames(assessment.Metric{})
	// Add empty string indicating no explicit ordering
	whitelist = append(whitelist, "")
	if !slices.Contains(whitelist, req.OrderBy) {
		err = ErrInvalidColumnName
		return
	}

	return
}
