package rasmonitoringapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSensorVoltageMonitor(t *testing.T) {
	log := logger.New().WithField("test", "true")
	type testCase struct {
		name         string
		checkRequest func(c *Client, r *http.Request)
		sendResponse func(r *http.Request) (*http.Response, error)
		expected     func() []rasmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			checkRequest: func(c *Client, r *http.Request) {
				checkSensorVoltageMonitorRequest(t, c, r)
			},
			sendResponse: func(r *http.Request) (*http.Response, error) {
				body, err := json.Marshal(monitorVoltageResponse{Value: 100})
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixSensorVoltageMonitor, "sensor_name"),
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 100,
				})
				return metrics
			},
		},
		{
			name: "Ras return an error when trying to get sensor voltage",
			checkRequest: func(c *Client, r *http.Request) {
				checkSensorVoltageMonitorRequest(t, c, r)
			},
			sendResponse: func(r *http.Request) (*http.Response, error) {
				body, err := json.Marshal(RasError{Err: "some error", Code: 1000})
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				return metrics
			},
		},
		{
			name: "Ras returns an empty body",
			checkRequest: func(c *Client, r *http.Request) {
				checkSensorVoltageMonitorRequest(t, c, r)
			},
			sendResponse: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
				}, nil
			},
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			transport := &transport{}

			client := &Client{
				log: log,
				Client: &rasmonitoring.Client{
					ClientOpts: rasmonitoring.ClientOpts{
						Oams:     1,
						Hostname: "hostname",
					},
					Client: http.Client{
						Transport: transport,
					},
				},
			}
			// check request
			transport.RoundTripFunc = func(r *http.Request) (*http.Response, error) {

				test.checkRequest(client, r)
				return test.sendResponse(r)
			}

			client.Auth = authDetails{
				Cookie: &http.Cookie{Name: "QSESSIONID"},
			}

			// override the var for testing
			voltageMonitorData = map[string]voltageMonitorInfo{
				"sensor_name": {
					Index:  1,
					Sensor: 2,
				},
			}
			ch := make(chan []rasmonitoring.Metric)
			go func() {
				ch <- client.SensorVoltageMonitor(context.Background(), log)
			}()

			// wait until we finish to check the request
			got := <-ch
			require.Equal(t, test.expected(), got)
		})
	}
}

func checkSensorVoltageMonitorRequest(t *testing.T, c *Client, r *http.Request) {
	verifyHeaders(t, c, r)
	require.Equal(t, "https://hostname/ext/ras/indirect/sensor_voltage_mon?index=1&oam=0&sensor=2", r.URL.String())
}
