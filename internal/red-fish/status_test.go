package redfish

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
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
					StatusResponse{
						Response: StatusResp{
							BootStage:                    "Preboot",
							EmergencyPowerReduction:      "Normal",
							ClockThrottling:              "Power",
							LastClockThrottlingDuration:  100,
							PowerState:                   "Reduced by 1/16",
							TotalClockThrottlingDuration: 100,
							GlobalTimeFromReset:          1662833000,
							ChipStatus:                   "Idle",
							DeviceActivity:               "Device In-use",
							DeviceActivityCounter:        10,
							DevicePowerReduction:         "Normal Power",
							LastPowerReductionDuration:   1500,
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

				prefix := bmcmonitoring.PrefixStatus
				// BootStage:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricBootStage),
					MetricValue: 2,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricBootStage: "Preboot",
					},
				})
				// EmergencyPowerReduction:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricEmergencyPowerReduction),
					MetricValue: 0,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricEmergencyPowerReduction: "Normal",
					},
				})
				// ClockThrottling:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricClockThrottling),
					MetricValue: 1,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricClockThrottling: "Power",
					},
				})

				// LastClockThrottlingDuration:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricLastClockThrottlingDuration),
					MetricValue: 100,
				})
				// PowerState:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricPowerState),
					MetricValue: 1,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricPowerState: "Reduced by 1/16",
					},
				})
				// TotalClockThrottlingDuration:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricTotalClockThrottlingDuration),
					MetricValue: 100,
				})
				// GlobalTimeFromReset:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricGlobalTimeFromReset),
					MetricValue: 1662833000,
				})
				// ChipStatus:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricChipStatus),
					MetricValue: 1,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricChipStatus: "Idle",
					},
				})
				// DeviceActivity:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricDeviceActivity),
					MetricValue: 1,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricDeviceActivity: "Device In-use",
					},
				})

				// DeviceActivityCounter:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricDeviceActivityCounter),
					MetricValue: 10,
				})
				// DevicePowerReduction:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricDevicePowerReduction),
					MetricValue: 3,
					CustomLabels: map[string]string{
						"state": "Normal Power",
					},
				})

				// LastPowerReductionDuration:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricLastPowerReductionDuration),
					MetricValue: 1500,
				})
				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifyStatus(metric)
					require.NoError(t, err)
				}
			},
		},
		{
			name: "Red fish returns invalid values of BootStage,EmergencyPowerReduction,ClockThrottling,ChipStatus,DeviceActivity,DevicePowerReduction",
			sendResponse: func() (*http.Response, error) {
				// send a valid response

				body, err := json.Marshal(
					StatusResponse{
						Response: StatusResp{
							BootStage:                    "unknown",
							EmergencyPowerReduction:      "unknown",
							ClockThrottling:              "unknown",
							LastClockThrottlingDuration:  100,
							PowerState:                   "Reduced by 1/16",
							TotalClockThrottlingDuration: 100,
							GlobalTimeFromReset:          1662833000,
							ChipStatus:                   "unknown",
							DeviceActivity:               "unknown",
							DeviceActivityCounter:        10,
							DevicePowerReduction:         "unknown",
							LastPowerReductionDuration:   1500,
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

				prefix := bmcmonitoring.PrefixStatus
				// BootStage:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricBootStage),
					MetricValue: -1,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricBootStage: "unknown",
					},
				})
				// EmergencyPowerReduction:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricEmergencyPowerReduction),
					MetricValue: -1,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricEmergencyPowerReduction: "unknown",
					},
				})
				// ClockThrottling:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricClockThrottling),
					MetricValue: -1,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricClockThrottling: "unknown",
					},
				})

				// LastClockThrottlingDuration:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricLastClockThrottlingDuration),
					MetricValue: 100,
				})
				// PowerState:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricPowerState),
					MetricValue: 1,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricPowerState: "Reduced by 1/16",
					},
				})
				// TotalClockThrottlingDuration:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricTotalClockThrottlingDuration),
					MetricValue: 100,
				})
				// GlobalTimeFromReset:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricGlobalTimeFromReset),
					MetricValue: 1662833000,
				})
				// ChipStatus:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricChipStatus),
					MetricValue: -1,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricChipStatus: "unknown",
					},
				})
				// DeviceActivity:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricDeviceActivity),
					MetricValue: -1,
					CustomLabels: map[string]string{
						bmcmonitoring.StatusMetricDeviceActivity: "unknown",
					},
				})

				// DeviceActivityCounter:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricDeviceActivityCounter),
					MetricValue: 10,
				})
				// DevicePowerReduction:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricDevicePowerReduction),
					MetricValue: -1,
					CustomLabels: map[string]string{
						"state": "unknown",
					},
				})

				// LastPowerReductionDuration:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricLastPowerReductionDuration),
					MetricValue: 1500,
				})
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/Status", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.Status(context.Background(), log)
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
