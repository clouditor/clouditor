// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/evidence/evidence_store.proto

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
var _evidence_store_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on StoreEvidenceRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StoreEvidenceRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StoreEvidenceRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StoreEvidenceRequestMultiError, or nil if none found.
func (m *StoreEvidenceRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *StoreEvidenceRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetEvidence() == nil {
		err := StoreEvidenceRequestValidationError{
			field:  "Evidence",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetEvidence()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, StoreEvidenceRequestValidationError{
					field:  "Evidence",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, StoreEvidenceRequestValidationError{
					field:  "Evidence",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetEvidence()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return StoreEvidenceRequestValidationError{
				field:  "Evidence",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return StoreEvidenceRequestMultiError(errors)
	}

	return nil
}

// StoreEvidenceRequestMultiError is an error wrapping multiple validation
// errors returned by StoreEvidenceRequest.ValidateAll() if the designated
// constraints aren't met.
type StoreEvidenceRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StoreEvidenceRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StoreEvidenceRequestMultiError) AllErrors() []error { return m }

// StoreEvidenceRequestValidationError is the validation error returned by
// StoreEvidenceRequest.Validate if the designated constraints aren't met.
type StoreEvidenceRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StoreEvidenceRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StoreEvidenceRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StoreEvidenceRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StoreEvidenceRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StoreEvidenceRequestValidationError) ErrorName() string {
	return "StoreEvidenceRequestValidationError"
}

// Error satisfies the builtin error interface
func (e StoreEvidenceRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStoreEvidenceRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StoreEvidenceRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StoreEvidenceRequestValidationError{}

// Validate checks the field values on StoreEvidenceResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StoreEvidenceResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StoreEvidenceResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StoreEvidenceResponseMultiError, or nil if none found.
func (m *StoreEvidenceResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *StoreEvidenceResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return StoreEvidenceResponseMultiError(errors)
	}

	return nil
}

// StoreEvidenceResponseMultiError is an error wrapping multiple validation
// errors returned by StoreEvidenceResponse.ValidateAll() if the designated
// constraints aren't met.
type StoreEvidenceResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StoreEvidenceResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StoreEvidenceResponseMultiError) AllErrors() []error { return m }

// StoreEvidenceResponseValidationError is the validation error returned by
// StoreEvidenceResponse.Validate if the designated constraints aren't met.
type StoreEvidenceResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StoreEvidenceResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StoreEvidenceResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StoreEvidenceResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StoreEvidenceResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StoreEvidenceResponseValidationError) ErrorName() string {
	return "StoreEvidenceResponseValidationError"
}

// Error satisfies the builtin error interface
func (e StoreEvidenceResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStoreEvidenceResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StoreEvidenceResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StoreEvidenceResponseValidationError{}

// Validate checks the field values on StoreEvidencesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *StoreEvidencesResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on StoreEvidencesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// StoreEvidencesResponseMultiError, or nil if none found.
func (m *StoreEvidencesResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *StoreEvidencesResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Status

	// no validation rules for StatusMessage

	if len(errors) > 0 {
		return StoreEvidencesResponseMultiError(errors)
	}

	return nil
}

// StoreEvidencesResponseMultiError is an error wrapping multiple validation
// errors returned by StoreEvidencesResponse.ValidateAll() if the designated
// constraints aren't met.
type StoreEvidencesResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m StoreEvidencesResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m StoreEvidencesResponseMultiError) AllErrors() []error { return m }

// StoreEvidencesResponseValidationError is the validation error returned by
// StoreEvidencesResponse.Validate if the designated constraints aren't met.
type StoreEvidencesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e StoreEvidencesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e StoreEvidencesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e StoreEvidencesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e StoreEvidencesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e StoreEvidencesResponseValidationError) ErrorName() string {
	return "StoreEvidencesResponseValidationError"
}

// Error satisfies the builtin error interface
func (e StoreEvidencesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sStoreEvidencesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = StoreEvidencesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = StoreEvidencesResponseValidationError{}

// Validate checks the field values on ListEvidencesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListEvidencesRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListEvidencesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListEvidencesRequestMultiError, or nil if none found.
func (m *ListEvidencesRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListEvidencesRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for PageSize

	// no validation rules for PageToken

	// no validation rules for OrderBy

	// no validation rules for Asc

	if m.Filter != nil {

		if all {
			switch v := interface{}(m.GetFilter()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListEvidencesRequestValidationError{
						field:  "Filter",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListEvidencesRequestValidationError{
						field:  "Filter",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetFilter()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListEvidencesRequestValidationError{
					field:  "Filter",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListEvidencesRequestMultiError(errors)
	}

	return nil
}

// ListEvidencesRequestMultiError is an error wrapping multiple validation
// errors returned by ListEvidencesRequest.ValidateAll() if the designated
// constraints aren't met.
type ListEvidencesRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListEvidencesRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListEvidencesRequestMultiError) AllErrors() []error { return m }

// ListEvidencesRequestValidationError is the validation error returned by
// ListEvidencesRequest.Validate if the designated constraints aren't met.
type ListEvidencesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListEvidencesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListEvidencesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListEvidencesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListEvidencesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListEvidencesRequestValidationError) ErrorName() string {
	return "ListEvidencesRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListEvidencesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListEvidencesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListEvidencesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListEvidencesRequestValidationError{}

// Validate checks the field values on Filter with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Filter) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Filter with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in FilterMultiError, or nil if none found.
func (m *Filter) ValidateAll() error {
	return m.validate(true)
}

func (m *Filter) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.CloudServiceId != nil {

		if err := m._validateUuid(m.GetCloudServiceId()); err != nil {
			err = FilterValidationError{
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

	if m.ToolId != nil {

		if err := m._validateUuid(m.GetToolId()); err != nil {
			err = FilterValidationError{
				field:  "ToolId",
				reason: "value must be a valid UUID",
				cause:  err,
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return FilterMultiError(errors)
	}

	return nil
}

func (m *Filter) _validateUuid(uuid string) error {
	if matched := _evidence_store_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// FilterMultiError is an error wrapping multiple validation errors returned by
// Filter.ValidateAll() if the designated constraints aren't met.
type FilterMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FilterMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FilterMultiError) AllErrors() []error { return m }

// FilterValidationError is the validation error returned by Filter.Validate if
// the designated constraints aren't met.
type FilterValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FilterValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FilterValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FilterValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FilterValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FilterValidationError) ErrorName() string { return "FilterValidationError" }

// Error satisfies the builtin error interface
func (e FilterValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFilter.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FilterValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FilterValidationError{}

// Validate checks the field values on ListEvidencesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListEvidencesResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListEvidencesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListEvidencesResponseMultiError, or nil if none found.
func (m *ListEvidencesResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListEvidencesResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetEvidences() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListEvidencesResponseValidationError{
						field:  fmt.Sprintf("Evidences[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListEvidencesResponseValidationError{
						field:  fmt.Sprintf("Evidences[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListEvidencesResponseValidationError{
					field:  fmt.Sprintf("Evidences[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for NextPageToken

	if len(errors) > 0 {
		return ListEvidencesResponseMultiError(errors)
	}

	return nil
}

// ListEvidencesResponseMultiError is an error wrapping multiple validation
// errors returned by ListEvidencesResponse.ValidateAll() if the designated
// constraints aren't met.
type ListEvidencesResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListEvidencesResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListEvidencesResponseMultiError) AllErrors() []error { return m }

// ListEvidencesResponseValidationError is the validation error returned by
// ListEvidencesResponse.Validate if the designated constraints aren't met.
type ListEvidencesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListEvidencesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListEvidencesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListEvidencesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListEvidencesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListEvidencesResponseValidationError) ErrorName() string {
	return "ListEvidencesResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListEvidencesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListEvidencesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListEvidencesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListEvidencesResponseValidationError{}

// Validate checks the field values on GetEvidenceRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetEvidenceRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetEvidenceRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetEvidenceRequestMultiError, or nil if none found.
func (m *GetEvidenceRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetEvidenceRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetEvidenceId()); err != nil {
		err = GetEvidenceRequestValidationError{
			field:  "EvidenceId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetEvidenceRequestMultiError(errors)
	}

	return nil
}

func (m *GetEvidenceRequest) _validateUuid(uuid string) error {
	if matched := _evidence_store_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// GetEvidenceRequestMultiError is an error wrapping multiple validation errors
// returned by GetEvidenceRequest.ValidateAll() if the designated constraints
// aren't met.
type GetEvidenceRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetEvidenceRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetEvidenceRequestMultiError) AllErrors() []error { return m }

// GetEvidenceRequestValidationError is the validation error returned by
// GetEvidenceRequest.Validate if the designated constraints aren't met.
type GetEvidenceRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetEvidenceRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetEvidenceRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetEvidenceRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetEvidenceRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetEvidenceRequestValidationError) ErrorName() string {
	return "GetEvidenceRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetEvidenceRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetEvidenceRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetEvidenceRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetEvidenceRequestValidationError{}
