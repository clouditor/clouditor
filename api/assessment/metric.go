package assessment

import (
	"encoding/json"
	"errors"
)

var (
	//ErrMetricIdMissing   = errors.New("metric id is missing")
	ErrMetricNameMissing = errors.New("metric name is missing")
)

func (r *Range) UnmarshalJSON(b []byte) (err error) {
	// check for the different range types
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

// Validate validates the metric according to several required fields
func (m *Metric) Validate() (err error) {
	if m.Id == "" {
		return ErrMetricIdMissing
	}

	if m.Name == "" {
		return ErrMetricNameMissing
	}

	return nil
}
