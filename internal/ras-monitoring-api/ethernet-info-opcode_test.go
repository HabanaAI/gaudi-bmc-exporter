package rasmonitoringapi

import (
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeLinkSpeed(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "speed 25",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric

				metric = append(metric, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_field_name", rasmonitoring.PrefixEthernetInfo),
					MetricValue: 25,
				})
				return metric
			},
		},
		{
			name:  "speed 56",
			value: "1",
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric

				metric = append(metric, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_field_name", rasmonitoring.PrefixEthernetInfo),
					MetricValue: 56,
				})
				return metric
			},
		},
		{
			name:  "speed 112",
			value: "2",
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric

				metric = append(metric, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_field_name", rasmonitoring.PrefixEthernetInfo),
					MetricValue: 112,
				})
				return metric
			},
		},
		{
			name:  "unsupported speed",
			value: "3",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := decodeLinkSpeed(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixEthernetInfo)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodePortMaxSpeed(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "max speed 50",
			value: "1",
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric

				metric = append(metric, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_field_name", rasmonitoring.PrefixEthernetInfo),
					MetricValue: 50,
				})
				return metric
			},
		},
		{
			name:  "max speed 100",
			value: "2",
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric

				metric = append(metric, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_field_name", rasmonitoring.PrefixEthernetInfo),
					MetricValue: 100,
				})
				return metric
			},
		},
		{
			name:  "max speed 200",
			value: "3",
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric

				metric = append(metric, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_field_name", rasmonitoring.PrefixEthernetInfo),
					MetricValue: 200,
				})
				return metric
			},
		},
		{
			name:  "max speed 400",
			value: "4",
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric

				metric = append(metric, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_field_name", rasmonitoring.PrefixEthernetInfo),
					MetricValue: 400,
				})
				return metric
			},
		},
		{
			name:  "unsupported max speed",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := decodePortMaxSpeed(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixEthernetInfo)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodeANLTStatus(t *testing.T) {

	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "OK",
			value: "10",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetInfo, "field_name"),
					CustomLabels: map[string]string{
						"state": "Enabled",
						"port":  "0",
					},
				})

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetInfo, "field_name"),
					CustomLabels: map[string]string{
						"state": "Disabled",
						"port":  "1",
					},
				})
				return metrics
			},
		},
		{
			name:  "unsupported port state",
			value: "2",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodeANLTStatus(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixEthernetInfo)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodeSerDesAvailability(t *testing.T) {

	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "Available",
			value: "1",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetInfo, "field_name"),
					CustomLabels: map[string]string{
						"state": "Available",
					},
				})

				return metrics
			},
		},
		{
			name:  "Unavailable",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetInfo, "field_name"),
					CustomLabels: map[string]string{
						"state": "Unavailable",
					},
				})
				return metrics
			},
		},
		{
			name:  "unsupported state",
			value: "2",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodeSerDesAvailability(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixEthernetInfo)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}
