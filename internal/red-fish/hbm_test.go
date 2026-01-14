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

func TestHbm(t *testing.T) {
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
					HbmResponse{
						Response: HbmResp{
							EccErrors: 2,
							RepairedLanes: []RepairedLane{
								{
									HBMIndex:  0,
									MCChannel: 1,
								},
								{
									HBMIndex:  2,
									MCChannel: 1,
								},
							},
							ReplacedRows: []ReplacedRow{
								{
									HBMIndex:   0,
									PCIndex:    1,
									StackID:    1,
									BankIndex:  1,
									Cause:      "DoubleECC",
									RowAddress: 1,
								},
							},
							HBMRepairStatus: []RepairStatus{
								{
									HBMIndex:    1,
									MBISTRepair: "Flow did not run",
									GlobalECC:   4,
								},
							},
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

				prefix := bmcmonitoring.PrefixHbm

				// EccErrors
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 2,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricEccErrors),
				})

				// RepairedLanes
				metrics = append(metrics, bmcmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricNumOfRepairedLanes),
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 2,
				})

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricRepairedLanes),
					CustomLabels: map[string]string{
						bmcmonitoring.HBMMetricRepairedLanesLabelHBMIndex:  "0",
						bmcmonitoring.HBMMetricRepairedLanesLabelMCChannel: "1",
					},
				})

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricRepairedLanes),
					CustomLabels: map[string]string{
						bmcmonitoring.HBMMetricRepairedLanesLabelHBMIndex:  "2",
						bmcmonitoring.HBMMetricRepairedLanesLabelMCChannel: "1",
					},
				})

				// ReplacedRow

				metrics = append(metrics, bmcmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricNumOfReplacedRows),
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
				})

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricReplaceRows),
					CustomLabels: map[string]string{
						bmcmonitoring.HBMMetricReplaceRowsLabelHBMIndex:   "0",
						bmcmonitoring.HBMMetricReplaceRowsLabelPCIndex:    "1",
						bmcmonitoring.HBMMetricReplaceRowsLabelStackID:    "1",
						bmcmonitoring.HBMMetricReplaceRowsLabelBankIndex:  "1",
						bmcmonitoring.HBMMetricReplaceRowsLabelCause:      "DoubleECC",
						bmcmonitoring.HBMMetricReplaceRowsLabelRowAddress: "1",
					},
				})

				// MbistRepair
				metrics = append(metrics, bmcmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricMbistRepair),
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					CustomLabels: map[string]string{
						bmcmonitoring.HBMMetricMbistRepairLabelState: "Flow did not run",
						bmcmonitoring.HBMMetricMbistRepairLabelIndex: "1",
					},
				})

				// Global ECC
				metrics = append(metrics, bmcmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricGlobalECC),
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 4,
					CustomLabels: map[string]string{
						bmcmonitoring.HBMMetricGlobalECCLabelIndex: "1",
					},
				})

				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifyHbm(metric)
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/HBM", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.Hbm(context.Background(), log)
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
