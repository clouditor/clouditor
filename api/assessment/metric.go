package assessment

import (
	"encoding/json"
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
