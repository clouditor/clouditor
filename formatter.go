package clouditor

import (
	"fmt"

	"github.com/logrusorgru/aurora"
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

func (f *GRPCFormatter) getCode(data logrus.Fields) aurora.Value {
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

func (f *GRPCFormatter) getService(data logrus.Fields) aurora.Value {
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

func (f *GRPCFormatter) getMethod(data logrus.Fields) aurora.Value {
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
