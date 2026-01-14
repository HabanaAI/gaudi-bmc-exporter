package redfish

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSensorTemperature(t *testing.T) {
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
					SensorTemperatureResponse{
						Response: SensorTemperatureResp{
							OnDie0:    48435,
							OnDie1:    48180,
							OnDie2:    48945,
							OnDie3:    47924,
							HBM0:      46000,
							HBM1:      46000,
							HBM2:      47000,
							HBM3:      49000,
							HBM4:      46000,
							HBM5:      46000,
							CPLDLocal: 47125,
							CPLD0:     47312,
							CPLD1:     46937,
							CPLD2:     47000,
							CPLD3:     47437,
							OnBoard0:  46500,
							OnBoard1:  47000,
							OnBoard2:  47000,
							OnBoard3:  47000,
							CPLDTemp:  49250,
							PSUStage1: 48000,
							PSUStage2: 55000,
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

				prefix := rasmonitoring.PrefixSensorTemperature

				// OnDie0
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 48435,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnDie0),
				})

				// OnDie1
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 48180,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnDie1),
				})

				// OnDie2
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 48945,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnDie2),
				})

				// OnDie3
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 47924,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnDie3),
				})

				// HBM0
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 46000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM0),
				})

				// HBM1
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 46000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM1),
				})

				// HBM2
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 47000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM2),
				})

				// HBM3
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 49000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM3),
				})

				// HBM4
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 46000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM4),
				})

				// HBM5
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 46000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM5),
				})

				// CPLDLocal
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 47125,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLDLocal),
				})

				// CPLD0
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 47312,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLD0),
				})

				// CPLD1
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 46937,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLD1),
				})

				// CPLD2
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 47000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLD2),
				})

				// CPLD3
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 47437,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLD3),
				})

				// OnBoard0
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 46500,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnboard0),
				})

				// OnBoard1
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 47000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnboard1),
				})

				// OnBoard2
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 47000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnboard2),
				})

				// OnBoard3
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 47000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnboard3),
				})

				// CPLDTemp
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 49250,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLDTemp),
				})

				// PSUStage1
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 48000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricPSUStage1),
				})

				// PSUStage2
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 55000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricPSUStage2),
				})

				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifySensorTemperature(metric)
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/SensorsTemperature", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.SensorTemperature(context.Background(), log)
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
