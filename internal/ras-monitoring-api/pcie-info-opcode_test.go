package rasmonitoringapi

import (
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodePcieLinkSpeed(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected func() []rasmonitoring.Metric
		wantErr  bool
	}

	tests := []testCase{
		{
			name:  "Invalid generation",
			value: "0",
			expected: func() []rasmonitoring.Metric {
				return nil
			},
			wantErr: true,
		},
		{
			name:  "Valid generation, returns valid metric",
			value: "1",
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixPcieInfo, "field_name"),
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					CustomLabels: map[string]string{
						"field_name": "Gen1",
					},
				})
				return metrics
			},
			wantErr: false,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			got, err := decodePcieLinkSpeed(test.value, 0, "hostname", "field_name", rasmonitoring.PrefixPcieInfo)

			if test.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expected(), got)
		})
	}
}
