package assessment

import (
	"encoding/json"
	"errors"
)

var (
	//ErrMetricIdMissing   = errors.New("metric id is missing")
	ErrMetricNameMissing = errors.New("metric name is missing")
	ErrMetricEmpty       = errors.New("metric is missing or empty")
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

// MetricValidationOption is a function-style option to fine-tune metric validation.
type MetricValidationOption func(*Metric) error

// WithMetricRequiresId is a validation option that specifies that Id must not be empty.
func WithMetricRequiresId() MetricValidationOption {
	return func(m *Metric) error {
		if m.Id == "" {
			return ErrMetricIdMissing
		}

		return nil
	}
}

// Validate validates the metric according to several required fields.
func (m *Metric) Validate(opts ...MetricValidationOption) (err error) {
	if m == nil {
		return ErrMetricEmpty
	}

	// Check for extra validation options
	for _, o := range opts {
		err = o(m)
		if err != nil {
			return err
		}
	}

	if m.Name == "" {
		return ErrMetricNameMissing
	}

	return nil
}
