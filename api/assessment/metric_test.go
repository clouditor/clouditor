package assessment

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetric_Validate(t *testing.T) {
	type fields struct {
		metric *Metric
	}
	type args struct {
		opts []MetricValidationOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Metric is nil",
			fields: fields{
				metric: nil,
			},
			args: args{opts: []MetricValidationOption{
				WithMetricRequiresId(),
			}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, ErrMetricEmpty, err)
			},
		},
		{
			name: "Metric ID is empty",
			fields: fields{
				metric: &Metric{Id: ""},
			},
			args: args{opts: []MetricValidationOption{
				WithMetricRequiresId(),
			}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, ErrMetricIdMissing, err)
			},
		},
		{
			name: "Metric Name is empty",
			fields: fields{
				metric: &Metric{Id: "123"},
			},
			args: args{opts: []MetricValidationOption{
				WithMetricRequiresId(),
			}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, ErrMetricNameMissing, err)
			},
		},
		{
			name: "Successful Validation",
			fields: fields{
				metric: &Metric{
					Id:   "SomeId",
					Name: "SomeName1",
				},
			},
			args: args{opts: []MetricValidationOption{
				WithMetricRequiresId(),
			}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.fields.metric
			tt.wantErr(t, m.Validate(tt.args.opts...), fmt.Sprintf("Validate(%v)", tt.args.opts))
		})
	}
}
