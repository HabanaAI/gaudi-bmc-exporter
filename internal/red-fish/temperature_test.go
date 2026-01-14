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

func TestTemperature(t *testing.T) {
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
					TemperatureResponse{
						Response: TemperatureResp{
							CurrentBoardTemperature:               51000,
							CurrentVRMTemperature:                 59000,
							CurrentDRAMTemperature:                55000,
							CurrentOndieTemperature:               55000,
							HistoricalBoardTemperature:            51000,
							HistoricalVRMTemperature:              59000,
							HistoricalDRAMTemperature:             60000,
							HistoricalOndieTemperature:            55000,
							MaxTemperatureRiseTime:                0,
							MaxSOCTemperatureErrorThreshold:       93000,
							MaxSOCTemperatureWarningThreshold:     83000,
							MaxHBMTemperatureThreshold:            125000,
							CurrentSOCTemperatureErrorThreshold:   93000,
							CurrentSOCTemperatureWarningThreshold: 83000,
							CurrentHBMTemperatureThreshold:        125000,
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
				prefix := bmcmonitoring.PrefixTemperature

				// Current Board Temperature
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentBoardTemp),
					MetricValue: 51000,
				})
				// Current VRM Temperature
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentVRMTemp),
					MetricValue: 59000,
				})
				// CurrentDRAMTemperature
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentDRAMTemp),
					MetricValue: 55000,
				})
				// Current Ondie Temperature
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentOnDieTemp),
					MetricValue: 55000,
				})
				// HistoricalBoardTemperature
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricHistoricalBoardTemp),
					MetricValue: 51000,
				})
				// HistoricalVRMTemperature
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricHistoricalVRMTemp),
					MetricValue: 59000,
				})
				// HistoricalDRAMTemperature
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricHistoricalDRAMTemp),
					MetricValue: 60000,
				})
				// HistoricalOndieTemperature
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricHistoricalOnDieTemp),
					MetricValue: 55000,
				})
				// MaxTemperatureRiseTime
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricMaxTempRiseTime),
					MetricValue: 0,
				})
				// MaxSOCTemperatureErrorThreshold
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricMaxSocTempErrorThreshold),
					MetricValue: 93000,
				})
				// MaxSOCTemperatureWarningThreshold
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricMaxSocTempWarmingThreshold),
					MetricValue: 83000,
				})
				// MaxHBMTemperatureThreshold
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricMaxHbmTempThreshold),
					MetricValue: 125000,
				})
				// CurrentSOCTemperatureErrorThreshold
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentSocTempErrorThreshold),
					MetricValue: 93000,
				})
				// CurrentSOCTemperatureWarningThreshold
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentSocTempWarningThreshold),
					MetricValue: 83000,
				})
				// CurrentHBMTemperatureThreshold
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentHbmTempThreshold),
					MetricValue: 125000,
				})
				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifyTemperature(metric)
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/Temperature", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.Temperature(context.Background(), log)
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
