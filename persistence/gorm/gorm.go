// Copyright 2016-2022 Fraunhofer AISEC
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

package gorm

import (
	"errors"
	"fmt"
	"strings"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/discovery"
	"clouditor.io/clouditor/api/evaluation"
	"clouditor.io/clouditor/api/evidence"
	"clouditor.io/clouditor/api/orchestrator"
	"clouditor.io/clouditor/persistence"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var log *logrus.Entry

type DB = gorm.DB

type storage struct {
	db *gorm.DB
	// for options: (set default when not in opts)
	dialector gorm.Dialector
	config    gorm.Config

	// types contains all types that we need to auto-migrate into database tables
	types []any

	// maxConn is the maximum number of connections. 0 means unlimited.
	maxConn int
}

// DefaultTypes contains a list of internal types that need to be migrated by default
var DefaultTypes = []any{
	&assessment.MetricImplementation{},
	&assessment.Metric{},
	&assessment.AssessmentResult{},
	&discovery.Resource{},
	&evidence.Evidence{},
	&orchestrator.CloudService{},
	&orchestrator.Certificate{},
	&orchestrator.State{},
	&orchestrator.Catalog{},
	&orchestrator.Category{},
	&orchestrator.Control{},
	&evaluation.EvaluationResult{},
}

// StorageOption is a functional option type to configure the GORM storage. E.g. WithInMemory or WithPostgres
type StorageOption func(*storage)

// WithInMemory is an option to configure Storage to use an in memory DB
func WithInMemory() StorageOption {
	return func(s *storage) {
		s.dialector = sqlite.Open(":memory:?_pragma=foreign_keys(1)")
	}
}

// WithMaxOpenConns is an option to configure the maximum number of open connections
func WithMaxOpenConns(max int) StorageOption {
	return func(s *storage) {
		s.maxConn = max
	}
}

// WithPostgres is an option to configure Storage to use a Postgres DB
func WithPostgres(host string, port uint16, user string, pw string, db string, sslmode string) StorageOption {
	return func(s *storage) {
		s.dialector = postgres.Open(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", user, pw, host, port, db, sslmode))
	}
}

// WithLogger is an option to configure Storage to use a Logger
func WithLogger(logger logger.Interface) StorageOption {
	return func(s *storage) {
		s.config.Logger = logger
	}
}

// WithAdditionalAutoMigration is an option to add additional types to GORM's auto-migration.
func WithAdditionalAutoMigration(types ...any) StorageOption {
	return func(s *storage) {
		s.types = append(s.types, types...)
	}
}

func init() {
	log = logrus.WithField("component", "storage")
}

// NewStorage creates a new storage using GORM (which DB to use depends on the StorageOption)
func NewStorage(opts ...StorageOption) (s persistence.Storage, err error) {
	log.Println("Creating storage")
	// Create storage with default gorm config
	g := &storage{
		// We ignore Deepsource issue GO-W1004 (SkipDefaultTransaction of config is "false"): skipcq: GO-W1004
		config: gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
		types: DefaultTypes,
	}

	// Add options and/or override default ones
	for _, o := range opts {
		o(g)
	}

	if g.dialector == nil {
		WithInMemory()(g)
	}

	g.db, err = gorm.Open(g.dialector, &g.config)
	if err != nil {
		return nil, err
	}

	if g.maxConn > 0 {
		sql, err := g.db.DB()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve sql.DB: %v", err)
		}

		sql.SetMaxOpenConns(g.maxConn)
	}

	schema.RegisterSerializer("timestamppb", &TimestampSerializer{})
	schema.RegisterSerializer("valuepb", &ValueSerializer{})
	schema.RegisterSerializer("anypb", &AnySerializer{})

	if err = g.db.SetupJoinTable(&orchestrator.CloudService{}, "CatalogsInScope", &orchestrator.TargetOfEvaluation{}); err != nil {
		err = fmt.Errorf("error during join-table: %w", err)
		return
	}

	if err = g.db.SetupJoinTable(orchestrator.CloudService{}, "ConfiguredMetrics", assessment.MetricConfiguration{}); err != nil {
		err = fmt.Errorf("error during join-table: %w", err)
		return
	}

	if err = g.db.SetupJoinTable(orchestrator.TargetOfEvaluation{}, "ControlsInScope", orchestrator.ControlInScope{}); err != nil {
		err = fmt.Errorf("error during join-table: %w", err)
		return
	}

	// After successful DB initialization, migrate the schema
	if err = g.db.AutoMigrate(g.types...); err != nil {
		err = fmt.Errorf("error during auto-migration: %w", err)
		return
	}

	s = g
	return
}

func (s *storage) Create(r any) (err error) {
	err = s.db.Create(r).Error

	if err != nil && (strings.Contains(err.Error(), "constraint failed: UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "duplicate key value violates unique constraint")) {
		return persistence.ErrUniqueConstraintFailed
	}

	if err != nil && strings.Contains(err.Error(), "constraint failed") {
		return persistence.ErrConstraintFailed
	}

	return
}

type preload struct {
	query string
	args  []any
}

type QueryOption interface{}

// WithPreload allows the customization of Gorm's preload feature with the specified query and arguments.
func WithPreload(query string, args ...any) QueryOption {
	return &preload{query: query, args: args}
}

// WithoutPreload disables any kind of preloading of Gorm. This is necessary, if custom join tables are used, otherwise
// Gorm will throws errors.
func WithoutPreload() QueryOption {
	return &preload{query: ""}
}

func (s *storage) Get(r any, conds ...any) (err error) {
	// Preload all associations of r if necessary
	db, conds := applyPreload(s.db, conds...)

	err = db.First(r, conds...).Error

	// if record is not found, use the error message defined in the persistence package
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = persistence.ErrRecordNotFound
	}
	return
}

// applyWhere applies the conditional arguments to db.Where. We now basically distinguish between three cases:
//   - an empty conditions list means no db.Where function is called
//   - one condition specified means that it is takes as the query parameter. This will query for the specified primary key
//   - otherwise, the first condition will be taken as the query parameter and all others will be taken as additional args.
func applyWhere(db *gorm.DB, conds ...any) *gorm.DB {
	if len(conds) == 0 {
		return db
	} else if len(conds) == 1 {
		return db.Where(conds[0])
	} else {
		return db.Where(conds[0], conds[1:]...)
	}
}

// applyPreload checks for any preload options and prepends them to the DB query. If no extra option is specified,
// "clause.Associations" is used as the default preload.
func applyPreload(db *gorm.DB, conds ...any) (*gorm.DB, []any) {
	if len(conds) > 0 {
		if preload, ok := conds[0].(*preload); ok {
			if preload.query != "" {
				return db.Preload(preload.query, preload.args...), conds[1:]
			} else {
				return db, conds[1:]
			}
		}
	}

	return db.Preload(clause.Associations), conds
}

func (s *storage) List(r any, orderBy string, asc bool, offset int, limit int, conds ...any) error {
	var query = s.db
	// Set default direction to "ascending"
	var orderDirection = "asc"

	if limit != -1 {
		query = s.db.Limit(limit)
	}

	// Set direction to "descending"
	if !asc {
		orderDirection = "desc"
	}
	orderStmt := orderBy + " " + orderDirection
	// No explicit ordering
	if orderBy == "" {
		orderStmt = ""
	}

	// Preload all associations of r if necessary
	query, conds = applyPreload(query.Offset(offset), conds...)

	return query.Order(orderStmt).Find(r, conds...).Error
}

func (s *storage) Raw(r any, query string, args ...any) error {
	return s.db.Raw(query, args...).Scan(r).Error
}

func (s *storage) Count(r any, conds ...any) (count int64, err error) {
	db := applyWhere(s.db.Model(r), conds...)

	err = db.Count(&count).Error
	return
}

func (s *storage) Save(r any, conds ...any) error {

	err := s.Update(r, conds...)
	if err == nil {
		return nil
	}

	err = s.Create(r)
	if err != nil {
		return err
	}
	// tx := applyWhere(s.db, conds...).Save(r)
	// err := tx.Error

	// if err != nil && strings.Contains(err.Error(), "constraint failed") {
	// 	return persistence.ErrConstraintFailed
	// }

	return err
}

// Update will update the record with non-zero fields. Note that to get the entire updated record you have to call Get
func (s *storage) Update(r any, conds ...any) error {
	tx := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Model(r)
	tx = applyWhere(tx, conds...).Updates(r)
	if err := tx.Error; err != nil { // db error
		if strings.Contains(err.Error(), "constraint failed") {
			return persistence.ErrConstraintFailed
		} else {
			return err
		}
	}

	// No record with given ID found
	if tx.RowsAffected == 0 {
		return persistence.ErrRecordNotFound
	}

	return nil
}

// Delete deletes record with given id. If no record was found, returns ErrRecordNotFound
func (s *storage) Delete(r any, conds ...any) error {
	// Remove record r with given ID
	tx := s.db.Delete(r, conds...)
	if err := tx.Error; err != nil { // db error
		return err
	}

	// No record with given ID found
	if tx.RowsAffected == 0 {
		return persistence.ErrRecordNotFound
	}

	return nil
}
