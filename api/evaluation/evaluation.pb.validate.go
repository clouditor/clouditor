// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/evaluation/evaluation.proto

package evaluation

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
var _evaluation_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on ListEvaluationResultsRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListEvaluationResultsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListEvaluationResultsRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListEvaluationResultsRequestMultiError, or nil if none found.
func (m *ListEvaluationResultsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListEvaluationResultsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PageSize

	// no validation rules for PageToken

	// no validation rules for OrderBy

	// no validation rules for Asc

	if m.FilteredCloudServiceId != nil {

		if err := m._validateUuid(m.GetFilteredCloudServiceId()); err != nil {
			err = ListEvaluationResultsRequestValidationError{
				field:  "FilteredCloudServiceId",
				reason: "value must be a valid UUID",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.FilteredControlId != nil {

		if utf8.RuneCountInString(m.GetFilteredControlId()) < 1 {
			err := ListEvaluationResultsRequestValidationError{
				field:  "FilteredControlId",
				reason: "value length must be at least 1 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.FilteredSubControls != nil {

		if utf8.RuneCountInString(m.GetFilteredSubControls()) < 1 {
			err := ListEvaluationResultsRequestValidationError{
				field:  "FilteredSubControls",
				reason: "value length must be at least 1 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.LatestByResourceId != nil {
		// no validation rules for LatestByResourceId
	}

	if len(errors) > 0 {
		return ListEvaluationResultsRequestMultiError(errors)
	}

	return nil
}

func (m *ListEvaluationResultsRequest) _validateUuid(uuid string) error {
	if matched := _evaluation_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// ListEvaluationResultsRequestMultiError is an error wrapping multiple
// validation errors returned by ListEvaluationResultsRequest.ValidateAll() if
// the designated constraints aren't met.
type ListEvaluationResultsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListEvaluationResultsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListEvaluationResultsRequestMultiError) AllErrors() []error { return m }

// ListEvaluationResultsRequestValidationError is the validation error returned
// by ListEvaluationResultsRequest.Validate if the designated constraints
// aren't met.
type ListEvaluationResultsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListEvaluationResultsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListEvaluationResultsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListEvaluationResultsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListEvaluationResultsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListEvaluationResultsRequestValidationError) ErrorName() string {
	return "ListEvaluationResultsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListEvaluationResultsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListEvaluationResultsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListEvaluationResultsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListEvaluationResultsRequestValidationError{}

// Validate checks the field values on ListEvaluationResultsResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListEvaluationResultsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListEvaluationResultsResponse with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// ListEvaluationResultsResponseMultiError, or nil if none found.
func (m *ListEvaluationResultsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListEvaluationResultsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetResults() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListEvaluationResultsResponseValidationError{
						field:  fmt.Sprintf("Results[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListEvaluationResultsResponseValidationError{
						field:  fmt.Sprintf("Results[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListEvaluationResultsResponseValidationError{
					field:  fmt.Sprintf("Results[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for NextPageToken

	if len(errors) > 0 {
		return ListEvaluationResultsResponseMultiError(errors)
	}

	return nil
}

// ListEvaluationResultsResponseMultiError is an error wrapping multiple
// validation errors returned by ListEvaluationResultsResponse.ValidateAll()
// if the designated constraints aren't met.
type ListEvaluationResultsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListEvaluationResultsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListEvaluationResultsResponseMultiError) AllErrors() []error { return m }

// ListEvaluationResultsResponseValidationError is the validation error
// returned by ListEvaluationResultsResponse.Validate if the designated
// constraints aren't met.
type ListEvaluationResultsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListEvaluationResultsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListEvaluationResultsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListEvaluationResultsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListEvaluationResultsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListEvaluationResultsResponseValidationError) ErrorName() string {
	return "ListEvaluationResultsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListEvaluationResultsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListEvaluationResultsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListEvaluationResultsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListEvaluationResultsResponseValidationError{}

// Validate checks the field values on StartEvaluationRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StartEvaluationRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StartEvaluationRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StartEvaluationRequestMultiError, or nil if none found.
func (m *StartEvaluationRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *StartEvaluationRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetCloudServiceId()); err != nil {
		err = StartEvaluationRequestValidationError{
			field:  "CloudServiceId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetCatalogId()) < 1 {
		err := StartEvaluationRequestValidationError{
			field:  "CatalogId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.Interval != nil {

		if m.GetInterval() <= 0 {
			err := StartEvaluationRequestValidationError{
				field:  "Interval",
				reason: "value must be greater than 0",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return StartEvaluationRequestMultiError(errors)
	}

	return nil
}

func (m *StartEvaluationRequest) _validateUuid(uuid string) error {
	if matched := _evaluation_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// StartEvaluationRequestMultiError is an error wrapping multiple validation
// errors returned by StartEvaluationRequest.ValidateAll() if the designated
// constraints aren't met.
type StartEvaluationRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StartEvaluationRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StartEvaluationRequestMultiError) AllErrors() []error { return m }

// StartEvaluationRequestValidationError is the validation error returned by
// StartEvaluationRequest.Validate if the designated constraints aren't met.
type StartEvaluationRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StartEvaluationRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StartEvaluationRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StartEvaluationRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StartEvaluationRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StartEvaluationRequestValidationError) ErrorName() string {
	return "StartEvaluationRequestValidationError"
}

// Error satisfies the builtin error interface
func (e StartEvaluationRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStartEvaluationRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StartEvaluationRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StartEvaluationRequestValidationError{}

// Validate checks the field values on StartEvaluationResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StartEvaluationResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StartEvaluationResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StartEvaluationResponseMultiError, or nil if none found.
func (m *StartEvaluationResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *StartEvaluationResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Status

	// no validation rules for StatusMessage

	if len(errors) > 0 {
		return StartEvaluationResponseMultiError(errors)
	}

	return nil
}

// StartEvaluationResponseMultiError is an error wrapping multiple validation
// errors returned by StartEvaluationResponse.ValidateAll() if the designated
// constraints aren't met.
type StartEvaluationResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StartEvaluationResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StartEvaluationResponseMultiError) AllErrors() []error { return m }

// StartEvaluationResponseValidationError is the validation error returned by
// StartEvaluationResponse.Validate if the designated constraints aren't met.
type StartEvaluationResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StartEvaluationResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StartEvaluationResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StartEvaluationResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StartEvaluationResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StartEvaluationResponseValidationError) ErrorName() string {
	return "StartEvaluationResponseValidationError"
}

// Error satisfies the builtin error interface
func (e StartEvaluationResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStartEvaluationResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StartEvaluationResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StartEvaluationResponseValidationError{}

// Validate checks the field values on StopEvaluationRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StopEvaluationRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StopEvaluationRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StopEvaluationRequestMultiError, or nil if none found.
func (m *StopEvaluationRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *StopEvaluationRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetCloudServiceId()); err != nil {
		err = StopEvaluationRequestValidationError{
			field:  "CloudServiceId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetCatalogId()) < 1 {
		err := StopEvaluationRequestValidationError{
			field:  "CatalogId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return StopEvaluationRequestMultiError(errors)
	}

	return nil
}

func (m *StopEvaluationRequest) _validateUuid(uuid string) error {
	if matched := _evaluation_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// StopEvaluationRequestMultiError is an error wrapping multiple validation
// errors returned by StopEvaluationRequest.ValidateAll() if the designated
// constraints aren't met.
type StopEvaluationRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StopEvaluationRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StopEvaluationRequestMultiError) AllErrors() []error { return m }

// StopEvaluationRequestValidationError is the validation error returned by
// StopEvaluationRequest.Validate if the designated constraints aren't met.
type StopEvaluationRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StopEvaluationRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StopEvaluationRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StopEvaluationRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StopEvaluationRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StopEvaluationRequestValidationError) ErrorName() string {
	return "StopEvaluationRequestValidationError"
}

// Error satisfies the builtin error interface
func (e StopEvaluationRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStopEvaluationRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StopEvaluationRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StopEvaluationRequestValidationError{}

// Validate checks the field values on StopEvaluationResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StopEvaluationResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StopEvaluationResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StopEvaluationResponseMultiError, or nil if none found.
func (m *StopEvaluationResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *StopEvaluationResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return StopEvaluationResponseMultiError(errors)
	}

	return nil
}

// StopEvaluationResponseMultiError is an error wrapping multiple validation
// errors returned by StopEvaluationResponse.ValidateAll() if the designated
// constraints aren't met.
type StopEvaluationResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StopEvaluationResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StopEvaluationResponseMultiError) AllErrors() []error { return m }

// StopEvaluationResponseValidationError is the validation error returned by
// StopEvaluationResponse.Validate if the designated constraints aren't met.
type StopEvaluationResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StopEvaluationResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StopEvaluationResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StopEvaluationResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StopEvaluationResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StopEvaluationResponseValidationError) ErrorName() string {
	return "StopEvaluationResponseValidationError"
}

// Error satisfies the builtin error interface
func (e StopEvaluationResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStopEvaluationResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StopEvaluationResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StopEvaluationResponseValidationError{}

// Validate checks the field values on EvaluationResult with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *EvaluationResult) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on EvaluationResult with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// EvaluationResultMultiError, or nil if none found.
func (m *EvaluationResult) ValidateAll() error {
	return m.validate(true)
}

func (m *EvaluationResult) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetId() != "" {

		if err := m._validateUuid(m.GetId()); err != nil {
			err = EvaluationResultValidationError{
				field:  "Id",
				reason: "value must be a valid UUID",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.GetCloudServiceId() != "" {

		if err := m._validateUuid(m.GetCloudServiceId()); err != nil {
			err = EvaluationResultValidationError{
				field:  "CloudServiceId",
				reason: "value must be a valid UUID",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if utf8.RuneCountInString(m.GetControlId()) < 1 {
		err := EvaluationResultValidationError{
			field:  "ControlId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetCategoryName()) < 1 {
		err := EvaluationResultValidationError{
			field:  "CategoryName",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetCatalogId()) < 1 {
		err := EvaluationResultValidationError{
			field:  "CatalogId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetResourceId()) < 1 {
		err := EvaluationResultValidationError{
			field:  "ResourceId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := EvaluationResult_EvaluationStatus_name[int32(m.GetStatus())]; !ok {
		err := EvaluationResultValidationError{
			field:  "Status",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetTimestamp()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, EvaluationResultValidationError{
					field:  "Timestamp",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, EvaluationResultValidationError{
					field:  "Timestamp",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTimestamp()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return EvaluationResultValidationError{
				field:  "Timestamp",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return EvaluationResultMultiError(errors)
	}

	return nil
}

func (m *EvaluationResult) _validateUuid(uuid string) error {
	if matched := _evaluation_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// EvaluationResultMultiError is an error wrapping multiple validation errors
// returned by EvaluationResult.ValidateAll() if the designated constraints
// aren't met.
type EvaluationResultMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m EvaluationResultMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m EvaluationResultMultiError) AllErrors() []error { return m }

// EvaluationResultValidationError is the validation error returned by
// EvaluationResult.Validate if the designated constraints aren't met.
type EvaluationResultValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e EvaluationResultValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e EvaluationResultValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e EvaluationResultValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e EvaluationResultValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e EvaluationResultValidationError) ErrorName() string { return "EvaluationResultValidationError" }

// Error satisfies the builtin error interface
func (e EvaluationResultValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sEvaluationResult.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = EvaluationResultValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = EvaluationResultValidationError{}
