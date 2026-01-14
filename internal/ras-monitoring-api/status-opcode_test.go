package rasmonitoringapi

import (
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeBootStage(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "stage Linux",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Linux",
					},
				})

				return metrics
			},
		},
		{
			name:  "stage Uboot",
			value: "1",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Uboot",
					},
				})

				return metrics
			},
		},
		{
			name:  "stage Preboot",
			value: "2",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 2,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Preboot",
					},
				})

				return metrics
			},
		},
		{
			name:  "stage Zephyr",
			value: "3",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 3,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Zephyr",
					},
				})

				return metrics
			},
		},
		{
			name:  "unsupported stage",
			value: "4",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodeBootStage(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixStatus)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodeEmergencyPowerReduction(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "stage Normal",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Normal",
					},
				})

				return metrics
			},
		},
		{
			name:  "stage Reduced",
			value: "1",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Reduced",
					},
				})

				return metrics
			},
		},
		{
			name:  "unsupported stage",
			value: "4",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodeEmergencyPowerReduction(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixStatus)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodeClockThrottling(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "stage None",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "None",
					},
				})

				return metrics
			},
		},
		{
			name:  "stage Power",
			value: "1",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Power",
					},
				})

				return metrics
			},
		},
		{
			name:  "stage Thermal",
			value: "2",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 2,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Thermal",
					},
				})

				return metrics
			},
		},
		{
			name:  "unsupported stage",
			value: "4",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodeClockThrottling(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixStatus)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodePowerState(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "stage Full performance",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Full performance",
					},
				})

				return metrics
			},
		},
		{
			name:  "stage Reduced 2/16",
			value: "2",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 2,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Reduced by 2/16",
					},
				})

				return metrics
			},
		},
		{
			name:  "unsupported stage",
			value: "16",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodePowerState(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixStatus)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodeChipStatus(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "stage Processing",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Processing",
					},
				})

				return metrics
			},
		},
		{
			name:  "stage Idle",
			value: "1",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Idle",
					},
				})

				return metrics
			},
		},
		{
			name:  "unsupported stage",
			value: "2",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodeChipStatus(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixStatus)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodeDeviceActivity(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "Device not in use",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Device not in use",
					},
				})

				return metrics
			},
		},
		{
			name:  "Device in-use",
			value: "1",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixStatus, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Device in-use",
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
			got, err := decodeDeviceActivity(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixStatus)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}
