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

func TestLaneInfo(t *testing.T) {
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
					LaneInfoResponse{
						Response: LaneInfoResp{
							EBUFOverflow:        0,
							EBUFUnderRun:        1,
							RunningDisplayError: 2,
							SKPOSParityError:    3,
							DecodeError:         4,
							SYNCHeaderError:     5,
							DeskewError:         6,
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
				prefix := bmcmonitoring.PrefixLaneInfo

				// EBUF Overflow
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricEBUFOverflow),
					MetricValue: 0,
					CustomLabels: map[string]string{
						"lane": "0",
					},
				})
				// EBUF Under Run
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricEBUFUnderRun),
					MetricValue: 1,
					CustomLabels: map[string]string{
						"lane": "0",
					},
				})
				// Decode Error
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricDecodeError),
					MetricValue: 4,
					CustomLabels: map[string]string{
						"lane": "0",
					},
				})
				// Running Display Error
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricRunningDisplayError),
					MetricValue: 2,
					CustomLabels: map[string]string{
						"lane": "0",
					},
				})
				// SKPOS Parity Error
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricSKPOSParityError),
					MetricValue: 3,
					CustomLabels: map[string]string{
						"lane": "0",
					},
				})
				// SYNC Header Error
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricSYNCHeaderError),
					MetricValue: 5,
					CustomLabels: map[string]string{
						"lane": "0",
					},
				})
				// Deskew Error
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricDeskewError),
					MetricValue: 6,
					CustomLabels: map[string]string{
						"lane": "0",
					},
				})
				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifyLaneInfo(metric)
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
						Lanes:    1,
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/LaneInfo/0", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.LaneInfo(context.Background(), log)
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
