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

package policies

import (
	"os"
	"strings"
	"sync"

	"clouditor.io/clouditor/api/assessment"
	"clouditor.io/clouditor/api/evidence"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("component", "policies")
)

type metricsResourceTypeCache struct {
	sync.RWMutex
	m map[string][]string
}

type PolicyEval interface {
	Eval(evidence *evidence.Evidence, holder MetricConfigurationSource) (data []*Result, err error)
}

type Result struct {
	Applicable  bool
	Compliant   bool
	TargetValue interface{} `mapstructure:"target_value"`
	Operator    string
	MetricId    string
}

// MetricConfigurationSource can be used to retrieve a metric configuration for a particular metric (and target service)
type MetricConfigurationSource interface {
	MetricConfiguration(metric string) (*assessment.MetricConfiguration, error)
}

func scanBundleDir(baseDir string) ([]os.FileInfo, error) {
	dirname := baseDir + "/policies/bundles"

	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	files, err := f.Readdir(-1)
	_ = f.Close()
	if err != nil {
		return nil, err
	}

	return files, err
}

func createKey(types []string) (key string) {
	key = strings.Join(types, "-")
	key = strings.ReplaceAll(key, " ", "")
	return
}
