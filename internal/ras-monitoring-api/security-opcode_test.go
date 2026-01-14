package rasmonitoringapi

import (
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeFWImageSource(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "Primary",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixSecurity, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Primary",
					},
				})

				return metrics
			},
		},
		{
			name:  "Secondary",
			value: "1",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixSecurity, "field_name"),
					CustomLabels: map[string]string{
						"field_name": "Secondary",
					},
				})

				return metrics
			},
		},
		{
			name:  "unsupported image sourse",
			value: "4",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodeFWImageSource(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixSecurity)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodeKeyRevocation(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "OK",
			value: "11110000",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				keys := []string{
					rasmonitoring.SecurityMetricKey0Revocation,
					rasmonitoring.SecurityMetricKey1Revocation,
					rasmonitoring.SecurityMetricKey2Revocation,
					rasmonitoring.SecurityMetricKey3Revocation,
					rasmonitoring.SecurityMetricKey4Revocation,
				}

				for i := 0; i < 4; i++ {
					metrics = append(metrics, rasmonitoring.Metric{
						Hostname:    "hostname",
						Oam:         "0",
						MetricValue: 1,
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixSecurity, keys[i]),
						CustomLabels: map[string]string{
							keys[i]: "Revoked",
						},
					})
				}

				for i := 4; i < 5; i++ {
					metrics = append(metrics, rasmonitoring.Metric{
						Hostname:    "hostname",
						Oam:         "0",
						MetricValue: 0,
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixSecurity, keys[i]),
						CustomLabels: map[string]string{
							keys[i]: "Not Revoked",
						},
					})
				}

				return metrics
			},
		},
		{
			name:  "unsupported state",
			value: "44444444",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodeKeyRevocation(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixSecurity)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}
