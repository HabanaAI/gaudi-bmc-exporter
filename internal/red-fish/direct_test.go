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

func TestDirect(t *testing.T) {
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
					DirectResponse{
						Response: DirectResp{
							PCIeVendorID:     "1DA3",
							AsicSerialNumber: "HL2080A0",
							FWVersionMajor:   40,
							FWVersionMinor:   0,
							FWVersionPatch:   2,
							CoreVDD:          075,
							HBMVddq:          120,
							V12:              135000,
							APIVersion:       010000,
							OSStage:          "PREBOOT",
							IBAccessState:    "Full Access",
							OOBAccessState:   "Restricted Access",
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

				prefix := bmcmonitoring.PrefixDirect

				// PCIeVendorID:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricPCIeVendorID),
					CustomLabels: map[string]string{
						bmcmonitoring.DirectMetricPCIeVendorID: "1DA3",
					},
				})
				// AsicSerialNumber:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricASICSerialNumber),
					CustomLabels: map[string]string{
						bmcmonitoring.DirectMetricASICSerialNumber: "HL2080A0",
					},
				})
				// FWVersionMajor:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricFWVersionMajor),
					MetricValue: 40,
				})
				// FWVersionMinor:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricFWVersionMinor),
					MetricValue: 0,
				})
				// FWVersionPatch:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricFWVersionPatch),
					MetricValue: 2,
				})
				// CoreVDD:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricCoreVDD),
					MetricValue: 075,
				})
				// HBMVddq:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricHBMVDDq),
					MetricValue: 120,
				})
				// V12:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricV12),
					MetricValue: 135000,
				})
				// APIVersion:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricApiVersion),
					MetricValue: 010000,
				})
				// OSStage:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricOSStage),
					MetricValue: 2,
					CustomLabels: map[string]string{
						"stage": "PREBOOT",
					},
				})
				// IBAccessState:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricIBAccessState),
					MetricValue: 0,
					CustomLabels: map[string]string{
						"state": "Full Access",
					},
				})
				// OOBAccessState:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricOOBAccessState),
					MetricValue: 1,
					CustomLabels: map[string]string{
						"state": "Restricted Access",
					},
				})

				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifyDirect(metric)
					require.NoError(t, err)
				}
			},
		},
		{
			name: "Red fish returns unexpected states for OSStage,IBAccessState and OOBAccessState",
			sendResponse: func() (*http.Response, error) {
				body, err := json.Marshal(
					DirectResponse{
						Response: DirectResp{
							PCIeVendorID:     "1DA3",
							AsicSerialNumber: "HL2080A0",
							FWVersionMajor:   40,
							FWVersionMinor:   0,
							FWVersionPatch:   2,
							CoreVDD:          075,
							HBMVddq:          120,
							V12:              135000,
							APIVersion:       010000,
							OSStage:          "unknown",
							IBAccessState:    "unknown",
							OOBAccessState:   "unknown",
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

				prefix := bmcmonitoring.PrefixDirect

				// PCIeVendorID:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricPCIeVendorID),
					CustomLabels: map[string]string{
						bmcmonitoring.DirectMetricPCIeVendorID: "1DA3",
					},
				})
				// AsicSerialNumber:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricASICSerialNumber),
					CustomLabels: map[string]string{
						bmcmonitoring.DirectMetricASICSerialNumber: "HL2080A0",
					},
				})
				// FWVersionMajor:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricFWVersionMajor),
					MetricValue: 40,
				})
				// FWVersionMinor:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricFWVersionMinor),
					MetricValue: 0,
				})
				// FWVersionPatch:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricFWVersionPatch),
					MetricValue: 2,
				})
				// CoreVDD:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricCoreVDD),
					MetricValue: 075,
				})
				// HBMVddq:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricHBMVDDq),
					MetricValue: 120,
				})
				// V12:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricV12),
					MetricValue: 135000,
				})
				// APIVersion:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricApiVersion),
					MetricValue: 010000,
				})
				// OSStage:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricOSStage),
					MetricValue: -1,
					CustomLabels: map[string]string{
						"stage": "unknown",
					},
				})
				// IBAccessState:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricIBAccessState),
					MetricValue: -1,
					CustomLabels: map[string]string{
						"state": "unknown",
					},
				})
				// OOBAccessState:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricOOBAccessState),
					MetricValue: -1,
					CustomLabels: map[string]string{
						"state": "unknown",
					},
				})

				return metrics

			},
		},
		{
			name: "Red fish returns an empty body",
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/Direct", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.Direct(context.Background(), log)
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
