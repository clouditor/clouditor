// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/assessment/metric.proto

package assessment

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
var _metric_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on Metric with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Metric) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Metric with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in MetricMultiError, or nil if none found.
func (m *Metric) ValidateAll() error {
	return m.validate(true)
}

func (m *Metric) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetId()) < 1 {
		err := MetricValidationError{
			field:  "Id",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetName()) < 1 {
		err := MetricValidationError{
			field:  "Name",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Description

	// no validation rules for Category

	if _, ok := Metric_Scale_name[int32(m.GetScale())]; !ok {
		err := MetricValidationError{
			field:  "Scale",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetRange() == nil {
		err := MetricValidationError{
			field:  "Range",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetRange()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MetricValidationError{
					field:  "Range",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MetricValidationError{
					field:  "Range",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetRange()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MetricValidationError{
				field:  "Range",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Interval

	// no validation rules for Deprecated

	if m.Implementation != nil {

		if all {
			switch v := interface{}(m.GetImplementation()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, MetricValidationError{
						field:  "Implementation",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, MetricValidationError{
						field:  "Implementation",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetImplementation()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return MetricValidationError{
					field:  "Implementation",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return MetricMultiError(errors)
	}

	return nil
}

// MetricMultiError is an error wrapping multiple validation errors returned by
// Metric.ValidateAll() if the designated constraints aren't met.
type MetricMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MetricMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MetricMultiError) AllErrors() []error { return m }

// MetricValidationError is the validation error returned by Metric.Validate if
// the designated constraints aren't met.
type MetricValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetricValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetricValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetricValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetricValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetricValidationError) ErrorName() string { return "MetricValidationError" }

// Error satisfies the builtin error interface
func (e MetricValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetric.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetricValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetricValidationError{}

// Validate checks the field values on Range with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Range) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Range with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in RangeMultiError, or nil if none found.
func (m *Range) ValidateAll() error {
	return m.validate(true)
}

func (m *Range) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	switch v := m.Range.(type) {
	case *Range_AllowedValues:
		if v == nil {
			err := RangeValidationError{
				field:  "Range",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetAllowedValues()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, RangeValidationError{
						field:  "AllowedValues",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, RangeValidationError{
						field:  "AllowedValues",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetAllowedValues()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RangeValidationError{
					field:  "AllowedValues",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Range_Order:
		if v == nil {
			err := RangeValidationError{
				field:  "Range",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetOrder()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, RangeValidationError{
						field:  "Order",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, RangeValidationError{
						field:  "Order",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetOrder()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RangeValidationError{
					field:  "Order",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	case *Range_MinMax:
		if v == nil {
			err := RangeValidationError{
				field:  "Range",
				reason: "oneof value cannot be a typed-nil",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if all {
			switch v := interface{}(m.GetMinMax()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, RangeValidationError{
						field:  "MinMax",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, RangeValidationError{
						field:  "MinMax",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetMinMax()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RangeValidationError{
					field:  "MinMax",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	default:
		_ = v // ensures v is used
	}

	if len(errors) > 0 {
		return RangeMultiError(errors)
	}

	return nil
}

// RangeMultiError is an error wrapping multiple validation errors returned by
// Range.ValidateAll() if the designated constraints aren't met.
type RangeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RangeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RangeMultiError) AllErrors() []error { return m }

// RangeValidationError is the validation error returned by Range.Validate if
// the designated constraints aren't met.
type RangeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RangeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RangeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RangeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RangeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RangeValidationError) ErrorName() string { return "RangeValidationError" }

// Error satisfies the builtin error interface
func (e RangeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRange.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RangeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RangeValidationError{}

// Validate checks the field values on MinMax with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *MinMax) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MinMax with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in MinMaxMultiError, or nil if none found.
func (m *MinMax) ValidateAll() error {
	return m.validate(true)
}

func (m *MinMax) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Min

	// no validation rules for Max

	if len(errors) > 0 {
		return MinMaxMultiError(errors)
	}

	return nil
}

// MinMaxMultiError is an error wrapping multiple validation errors returned by
// MinMax.ValidateAll() if the designated constraints aren't met.
type MinMaxMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MinMaxMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MinMaxMultiError) AllErrors() []error { return m }

// MinMaxValidationError is the validation error returned by MinMax.Validate if
// the designated constraints aren't met.
type MinMaxValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MinMaxValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MinMaxValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MinMaxValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MinMaxValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MinMaxValidationError) ErrorName() string { return "MinMaxValidationError" }

// Error satisfies the builtin error interface
func (e MinMaxValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMinMax.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MinMaxValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MinMaxValidationError{}

// Validate checks the field values on AllowedValues with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *AllowedValues) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AllowedValues with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in AllowedValuesMultiError, or
// nil if none found.
func (m *AllowedValues) ValidateAll() error {
	return m.validate(true)
}

func (m *AllowedValues) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetValues() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, AllowedValuesValidationError{
						field:  fmt.Sprintf("Values[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, AllowedValuesValidationError{
						field:  fmt.Sprintf("Values[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return AllowedValuesValidationError{
					field:  fmt.Sprintf("Values[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return AllowedValuesMultiError(errors)
	}

	return nil
}

// AllowedValuesMultiError is an error wrapping multiple validation errors
// returned by AllowedValues.ValidateAll() if the designated constraints
// aren't met.
type AllowedValuesMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AllowedValuesMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AllowedValuesMultiError) AllErrors() []error { return m }

// AllowedValuesValidationError is the validation error returned by
// AllowedValues.Validate if the designated constraints aren't met.
type AllowedValuesValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AllowedValuesValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AllowedValuesValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AllowedValuesValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AllowedValuesValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AllowedValuesValidationError) ErrorName() string { return "AllowedValuesValidationError" }

// Error satisfies the builtin error interface
func (e AllowedValuesValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAllowedValues.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AllowedValuesValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AllowedValuesValidationError{}

// Validate checks the field values on Order with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Order) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Order with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in OrderMultiError, or nil if none found.
func (m *Order) ValidateAll() error {
	return m.validate(true)
}

func (m *Order) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetValues() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, OrderValidationError{
						field:  fmt.Sprintf("Values[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, OrderValidationError{
						field:  fmt.Sprintf("Values[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return OrderValidationError{
					field:  fmt.Sprintf("Values[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return OrderMultiError(errors)
	}

	return nil
}

// OrderMultiError is an error wrapping multiple validation errors returned by
// Order.ValidateAll() if the designated constraints aren't met.
type OrderMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrderMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrderMultiError) AllErrors() []error { return m }

// OrderValidationError is the validation error returned by Order.Validate if
// the designated constraints aren't met.
type OrderValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrderValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrderValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrderValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrderValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrderValidationError) ErrorName() string { return "OrderValidationError" }

// Error satisfies the builtin error interface
func (e OrderValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrder.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrderValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrderValidationError{}

// Validate checks the field values on MetricConfiguration with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *MetricConfiguration) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MetricConfiguration with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// MetricConfigurationMultiError, or nil if none found.
func (m *MetricConfiguration) ValidateAll() error {
	return m.validate(true)
}

func (m *MetricConfiguration) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if !_MetricConfiguration_Operator_Pattern.MatchString(m.GetOperator()) {
		err := MetricConfigurationValidationError{
			field:  "Operator",
			reason: "value does not match regex pattern \"^(|<|>|<=|>=|==|isIn|allIn)$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetTargetValue() == nil {
		err := MetricConfigurationValidationError{
			field:  "TargetValue",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetTargetValue()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MetricConfigurationValidationError{
					field:  "TargetValue",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MetricConfigurationValidationError{
					field:  "TargetValue",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTargetValue()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MetricConfigurationValidationError{
				field:  "TargetValue",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for IsDefault

	if all {
		switch v := interface{}(m.GetUpdatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MetricConfigurationValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MetricConfigurationValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MetricConfigurationValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if utf8.RuneCountInString(m.GetMetricId()) < 1 {
		err := MetricConfigurationValidationError{
			field:  "MetricId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateUuid(m.GetCloudServiceId()); err != nil {
		err = MetricConfigurationValidationError{
			field:  "CloudServiceId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return MetricConfigurationMultiError(errors)
	}

	return nil
}

func (m *MetricConfiguration) _validateUuid(uuid string) error {
	if matched := _metric_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// MetricConfigurationMultiError is an error wrapping multiple validation
// errors returned by MetricConfiguration.ValidateAll() if the designated
// constraints aren't met.
type MetricConfigurationMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MetricConfigurationMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MetricConfigurationMultiError) AllErrors() []error { return m }

// MetricConfigurationValidationError is the validation error returned by
// MetricConfiguration.Validate if the designated constraints aren't met.
type MetricConfigurationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetricConfigurationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetricConfigurationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetricConfigurationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetricConfigurationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetricConfigurationValidationError) ErrorName() string {
	return "MetricConfigurationValidationError"
}

// Error satisfies the builtin error interface
func (e MetricConfigurationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetricConfiguration.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetricConfigurationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetricConfigurationValidationError{}

var _MetricConfiguration_Operator_Pattern = regexp.MustCompile("^(|<|>|<=|>=|==|isIn|allIn)$")

// Validate checks the field values on MetricImplementation with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *MetricImplementation) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on MetricImplementation with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// MetricImplementationMultiError, or nil if none found.
func (m *MetricImplementation) ValidateAll() error {
	return m.validate(true)
}

func (m *MetricImplementation) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetMetricId()) < 1 {
		err := MetricImplementationValidationError{
			field:  "MetricId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if _, ok := MetricImplementation_Language_name[int32(m.GetLang())]; !ok {
		err := MetricImplementationValidationError{
			field:  "Lang",
			reason: "value must be one of the defined enum values",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetCode()) < 1 {
		err := MetricImplementationValidationError{
			field:  "Code",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetUpdatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, MetricImplementationValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, MetricImplementationValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return MetricImplementationValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return MetricImplementationMultiError(errors)
	}

	return nil
}

// MetricImplementationMultiError is an error wrapping multiple validation
// errors returned by MetricImplementation.ValidateAll() if the designated
// constraints aren't met.
type MetricImplementationMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m MetricImplementationMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m MetricImplementationMultiError) AllErrors() []error { return m }

// MetricImplementationValidationError is the validation error returned by
// MetricImplementation.Validate if the designated constraints aren't met.
type MetricImplementationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e MetricImplementationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e MetricImplementationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e MetricImplementationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e MetricImplementationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e MetricImplementationValidationError) ErrorName() string {
	return "MetricImplementationValidationError"
}

// Error satisfies the builtin error interface
func (e MetricImplementationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sMetricImplementation.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = MetricImplementationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = MetricImplementationValidationError{}
