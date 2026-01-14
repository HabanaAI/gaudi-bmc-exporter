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

func TestInfo(t *testing.T) {
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
					InfoResponse{
						Response: InfoResp{
							DeviceID:          "1020",
							SubSystemDeviceID: "1DA3",
							SubSystemVendorID: "1020",
							ASICSerialNumber:  "HL2080A0",
							BoardSerialNumber: "AM27043696",
							SRAMSize:          48,
							HBMSize:           96,
							UUID:              "TF8A77-07-06-02",
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

				prefix := bmcmonitoring.PrefixInfo
				// 	DeviceID:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricDeviceID),
					CustomLabels: map[string]string{
						bmcmonitoring.InfoMetricDeviceID: "1020",
					},
				})

				// SubSystemDeviceID:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricSubsystemDeviceID),
					CustomLabels: map[string]string{
						bmcmonitoring.InfoMetricSubsystemDeviceID: "1DA3",
					},
				})
				// SubSystemVendorID:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricSubsystemVendorID),
					CustomLabels: map[string]string{
						bmcmonitoring.InfoMetricSubsystemVendorID: "1020",
					},
				})
				// ASICSerialNumber:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricASICSerialNumber),
					CustomLabels: map[string]string{
						bmcmonitoring.InfoMetricASICSerialNumber: "HL2080A0",
					},
				})
				// BoardSerialNumber:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricBoardSerialNumber),
					CustomLabels: map[string]string{
						bmcmonitoring.InfoMetricBoardSerialNumber: "AM27043696",
					},
				})
				// SRAMSize:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricSRAMSize),
					MetricValue: 48,
				})
				// HBMSize:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricHBMSize),
					MetricValue: 96,
				})
				// UUID:

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricUUID),
					CustomLabels: map[string]string{
						bmcmonitoring.InfoMetricUUID: "TF8A77-07-06-02",
					},
				})
				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifyInfo(metric)
					require.NoError(t, err)
				}
			},
		},
		{
			name: "Red fish returns an error",
			sendResponse: func() (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
				}, fmt.Errorf("connection error")
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/Info", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.Info(context.Background(), log)
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
