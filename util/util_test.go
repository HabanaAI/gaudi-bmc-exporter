package util

import "testing"

type ExampleMetric struct {
	MetricName  string            `label:"metric_name"`
	MetricValue int               `label:"metric_value"`
	Labels      map[string]string `label:"something"`
}

func TestMetric(t *testing.T) {
	tests := []struct {
		name    string
		obj     any
		want    string
		wantErr bool
	}{
		{
			name: "valid metric",
			obj: &ExampleMetric{
				MetricName:  "test_metric",
				MetricValue: 42,
				Labels:      map[string]string{"something": "else", "another": "label", "123number": "value"},
			},
			want:    `test_metric{123number="value", another="label", something="else"} 42`,
			wantErr: false,
		},
		{
			name:    "nil pointer",
			obj:     (*ExampleMetric)(nil),
			want:    "",
			wantErr: true,
		},
		{
			name:    "non-pointer",
			obj:     ExampleMetric{},
			want:    "",
			wantErr: true,
		},
		{
			name: "missing metric_name tag",
			obj: &struct {
				Value int `label:"metric_value"`
			}{Value: 10},
			want:    "",
			wantErr: true,
		},
		{
			name: "missing metric_value tag",
			obj: &struct {
				Name string `label:"metric_name"`
			}{Name: "test"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetricFromObject(tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatProm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("FormatProm() got = %v, want %v", got, tt.want)
			}
		})
	}
}
