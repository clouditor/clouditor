// Copyright 2021 Fraunhofer AISEC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"fmt"

	"github.com/logrusorgru/aurora/v3"
	"github.com/sirupsen/logrus"
)

type GRPCFormatter struct {
	logrus.TextFormatter
}

func (f *GRPCFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if _, ok := entry.Data["grpc.service"]; ok {
		entry.Message = fmt.Sprintf("gRPC call %s | %s | %s",
			f.getService(entry.Data),
			f.getMethod(entry.Data),
			f.getCode(entry.Data))
	}

	return f.TextFormatter.Format(entry)
}

func (GRPCFormatter) getCode(data logrus.Fields) aurora.Value {
	var (
		ok   bool
		code string
	)

	if code, ok = data["grpc.code"].(string); !ok {
		code = "UNKNOWN"
	}

	paddedCode := aurora.BrightWhite(" " + code + " ")

	switch code {
	case "Unauthenticated":
		return aurora.BgRed(paddedCode)
	case "OK":
		return aurora.BgGreen(paddedCode)
	default:
		return aurora.BgGray(128, paddedCode)
	}
}

func (GRPCFormatter) getService(data logrus.Fields) aurora.Value {
	var (
		ok      bool
		service string
	)

	if service, ok = data["grpc.service"].(string); !ok {
		service = "UNKNOWN"
	}

	padded := aurora.BrightWhite(" " + service + " ")

	switch service {
	case "UNKNOWN":
		return aurora.BgGray(128, padded)
	default:
		return aurora.BgCyan(padded)
	}
}

func (GRPCFormatter) getMethod(data logrus.Fields) aurora.Value {
	var (
		ok     bool
		method string
	)

	if method, ok = data["grpc.method"].(string); !ok {
		method = "UNKNOWN"
	}

	padded := aurora.BrightWhite(" " + method + " ")

	switch method {
	case "UNKNOWN":
		return aurora.BgGray(128, padded)
	default:
		return aurora.BgGreen(padded)
	}
}
