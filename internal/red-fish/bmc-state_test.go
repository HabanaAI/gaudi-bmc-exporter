package redfish

import (
	"bytes"
	"context"
	"encoding/json"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBmcState(t *testing.T) {
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
			name: "Power on",
			sendResponse: func() (*http.Response, error) {
				body, err := json.Marshal(
					BmcStateResponse{
						PowerState: "On",
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

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					MetricValue: 1,
					MetricName:  bmcmonitoring.PrefixBMCState,
				})
				return metrics
			},
		},
		{
			name: "Power off",
			sendResponse: func() (*http.Response, error) {
				body, err := json.Marshal(
					BmcStateResponse{
						PowerState: "Off",
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

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					MetricValue: 0,
					MetricName:  bmcmonitoring.PrefixBMCState,
				})
				return metrics
			},
		},
		{
			name: "Error getting bmc state for red fish",
			sendResponse: func() (*http.Response, error) {
				body, err := json.Marshal(
					RedFishErr{
						Err: GeneralErr{
							Code:    "12",
							Message: "Some error",
						},
					},
				)
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			expected: func() []bmcmonitoring.Metric {
				return nil
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
				require.Equal(t, "https://hostname/redfish/v1/Chassis/1", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.BmcState(context.Background(), log)
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
