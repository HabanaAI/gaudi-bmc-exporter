package rasmonitoringapi

import (
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodePHYStatus(t *testing.T) {

	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:    "OK",
			value:   "10",
			wantErr: false,
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, "field_name"),
					CustomLabels: map[string]string{
						"state": "Ready",
						"phy":   "0",
					},
				})

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, "field_name"),
					CustomLabels: map[string]string{
						"state": "Not ready",
						"phy":   "1",
					},
				})
				return metrics
			},
		},
		{
			name:  "unsupported link state",
			value: "2",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got, err := decodePHYStatus(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixEthernetStatus)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodeLinkStatus(t *testing.T) {

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
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, "field_name"),
					CustomLabels: map[string]string{
						"state": "Connected",
						"link":  "0",
					},
				})

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, "field_name"),
					CustomLabels: map[string]string{
						"state": "Not connected",
						"link":  "1",
					},
				})
				return metrics
			},
		},
		{
			name:  "unsupported link state",
			value: "2",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodeLinkStatus(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixEthernetStatus)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodeExternalLinkStatus(t *testing.T) {
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
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, "field_name"),
					CustomLabels: map[string]string{
						"state": "Active Port",
						"link":  "0",
					},
				})

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, "field_name"),
					CustomLabels: map[string]string{
						"state": "Non-active port",
						"link":  "1",
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
			got, err := decodeExternalLinkStatus(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixEthernetStatus)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}

func TestDecodePortMapping(t *testing.T) {
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
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, "field_name"),
					CustomLabels: map[string]string{
						"type": "External",
						"port": "0",
					},
				})

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, "field_name"),
					CustomLabels: map[string]string{
						"type": "Internal",
						"port": "1",
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
			got, err := decodePortMapping(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixEthernetStatus)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}
