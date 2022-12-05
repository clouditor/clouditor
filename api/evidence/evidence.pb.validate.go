// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/evidence/evidence.proto

package evidence

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// define the regex for a UUID once up-front
var _evidence_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on Evidence with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Evidence) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Evidence with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in EvidenceMultiError, or nil
// if none found.
func (m *Evidence) ValidateAll() error {
	return m.validate(true)
}

func (m *Evidence) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetId()); err != nil {
		err = EvidenceValidationError{
			field:  "Id",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetTimestamp() == nil {
		err := EvidenceValidationError{
			field:  "Timestamp",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateUuid(m.GetCloudServiceId()); err != nil {
		err = EvidenceValidationError{
			field:  "CloudServiceId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetToolId()) < 1 {
		err := EvidenceValidationError{
			field:  "ToolId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Raw

	if m.GetResource() == nil {
		err := EvidenceValidationError{
			field:  "Resource",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetResource()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, EvidenceValidationError{
					field:  "Resource",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, EvidenceValidationError{
					field:  "Resource",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetResource()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return EvidenceValidationError{
				field:  "Resource",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return EvidenceMultiError(errors)
	}

	return nil
}

func (m *Evidence) _validateUuid(uuid string) error {
	if matched := _evidence_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// EvidenceMultiError is an error wrapping multiple validation errors returned
// by Evidence.ValidateAll() if the designated constraints aren't met.
type EvidenceMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EvidenceMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EvidenceMultiError) AllErrors() []error { return m }

// EvidenceValidationError is the validation error returned by
// Evidence.Validate if the designated constraints aren't met.
type EvidenceValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EvidenceValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EvidenceValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EvidenceValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EvidenceValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EvidenceValidationError) ErrorName() string { return "EvidenceValidationError" }

// Error satisfies the builtin error interface
func (e EvidenceValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEvidence.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EvidenceValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EvidenceValidationError{}
