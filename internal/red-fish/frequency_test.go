package redfish

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFrequency(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name         string
		sendResponse func() (*http.Response, error)
		expected     func() []bmcmonitoring.Metric

		// check that the metric are with the expected labels
		verifyMetrics func(metrics []bmcmonitoring.Metric)
	}

	tests := []testCase{
		{
			name: "Red fish returns valid metrics",
			sendResponse: func() (*http.Response, error) {
				// send a valid response

				body, err := json.Marshal(
					FrequencyResponse{
						Response: FrequencyResp{
							HBMFrequency:          1600,
							MaxTPCFrequency:       1800,
							MaxMMEFrequency:       1650,
							MaxDMAFrequency:       1750,
							MaxMediaFrequency:     750,
							MaxPCIeFrequency:      500,
							MaxARMFrequency:       500,
							MaxNICFrequency:       976,
							MaxNoCFrequency:       1750,
							CurrentTPCFrequency:   1800,
							CurrentMMEFrequency:   1650,
							CurrentDMAFrequency:   1750,
							CurrentMediaFrequency: 750,
							CurrentPCIeFrequency:  500,
							CurrentARMFrequency:   500,
							CurrentNICFrequency:   976,
							CurrentNoCFrequency:   1750,
							CurrentSRAMFrequency:  350,
							CurrentMSSFrequency:   1200,
						},
					},
				)
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				prefix := bmcmonitoring.PrefixFrequency

				// HBMFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1600,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricHBMFrequency),
				})

				// MaxTPCFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1800,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxTPCFrequency),
				})

				// MaxMMEFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1650,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxMMEFrequency),
				})

				// MaxDMAFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1750,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxDMAFrequency),
				})

				// MaxMediaFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 750,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxMediaFrequency),
				})

				// MaxPCIeFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 500,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxPCIeFrequency),
				})

				// MaxARMFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 500,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxARMFrequency),
				})

				// MaxNICFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 976,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxNICFrequency),
				})

				// MaxNoCFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1750,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxNoCFrequency),
				})

				// CurrentTPCFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1800,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentTPCFrequency),
				})

				// CurrentMMEFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1650,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentMMEFrequency),
				})

				// CurrentDMAFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1750,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentDMAFrequency),
				})

				// CurrentMediaFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 750,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentMediaFrequency),
				})

				// CurrentPCIeFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 500,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentPCIeFrequency),
				})

				// CurrentARMFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 500,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentARMFrequency),
				})

				// CurrentNICFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 976,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentNICFrequency),
				})

				// CurrentNoCFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1750,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentNoCFrequency),
				})

				// CurrentSRAMFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 350,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentSRAMFrequency),
				})

				// CurrentMSSFrequency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1200,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentMSSFrequency),
				})

				return metrics

			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifyFrequency(metric)
					require.NoError(t, err)
				}
			},
		},
		{
			name: "Red fish returns error",
			sendResponse: func() (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
				}, errors.New("connection error")
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			transport := &transport{}

			client := &Client{
				log: log,
				Client: &bmcmonitoring.Client{
					ClientOpts: bmcmonitoring.ClientOpts{
						Oams:     1,
						Hostname: "hostname",
					},
					Client: http.Client{
						Transport: transport,
					},
				},
			}

			// check request and send response
			transport.RoundTripFunc = func(r *http.Request) (*http.Response, error) {

				checkAuth(t, r)
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/Frequency", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.Frequency(context.Background(), log)
			}()

			// wait until we finish to check the request
			got := <-ch
			require.Equal(t, test.expected(), got)

			if test.verifyMetrics != nil {
				test.verifyMetrics(got)
			}
		})
	}
}
