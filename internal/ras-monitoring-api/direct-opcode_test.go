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

func TestOsState(t *testing.T) {
	log := logger.New().WithField("test", "true")
	type testCase struct {
		name         string
		checkRequest func(c *Client, r *http.Request)
		sendResponse func() (*http.Response, error)
		ctxFunc      func() context.Context
		expected     func() []rasmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)
				require.Equal(t, "https://hostname/ext/ras/direct/control_command", r.URL.String())

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				var req directRequest

				err = json.Unmarshal(body, &req)
				require.NoError(t, err)

				require.Equal(t, directRequest{
					Oam:    "0",
					Opcode: fmt.Sprintf("%d", osStateOpcode),
					Cpsp:   "0",
				}, req)

			},
			sendResponse: func() (*http.Response, error) {

				body, err := json.Marshal(directResponse{Cpsr: 769})
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			ctxFunc: context.Background,
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics,
					rasmonitoring.Metric{
						Hostname:    "hostname",
						Oam:         "0",
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixDirect, rasmonitoring.DirectMetricOSStage),
						MetricValue: 3,
						CustomLabels: map[string]string{
							"stage": "ZEPHYR",
						},
					})

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixDirect, rasmonitoring.DirectMetricIBAccessState),
					MetricValue: 0,
					CustomLabels: map[string]string{
						"state": "Full Access",
					},
				})

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixDirect, rasmonitoring.DirectMetricOOBAccessState),
					MetricValue: 1,
					CustomLabels: map[string]string{
						"state": "Restricted Access",
					},
				})
				return metrics
			},
		},
		{
			name: "Ras returns an Unexpected format in the direct opcode",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)
				require.Equal(t, "https://hostname/ext/ras/direct/control_command", r.URL.String())

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				var req directRequest

				err = json.Unmarshal(body, &req)
				require.NoError(t, err)

				require.Equal(t, directRequest{
					Oam:    "0",
					Opcode: fmt.Sprintf("%d", osStateOpcode),
					Cpsp:   "0",
				}, req)

			},
			sendResponse: func() (*http.Response, error) {

				body, err := json.Marshal(directResponse{Cpsr: 100})
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			ctxFunc: context.Background,
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				return metrics
			},
		},
		{
			name: "Ras doesn't send a response in the direct opcode",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)
				require.Equal(t, "https://hostname/ext/ras/direct/control_command", r.URL.String())

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				var req directRequest

				err = json.Unmarshal(body, &req)
				require.NoError(t, err)

				require.Equal(t, directRequest{
					Oam:    "0",
					Opcode: fmt.Sprintf("%d", osStateOpcode),
					Cpsp:   "0",
				}, req)

			},
			sendResponse: func() (*http.Response, error) {

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
				}, nil
			},
			ctxFunc: context.Background,
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
				return test.sendResponse()
			}

			client.Auth = authDetails{
				Cookie: &http.Cookie{Name: "QSESSIONID"},
			}

			ch := make(chan []rasmonitoring.Metric)
			go func() {
				ch <- client.osState(test.ctxFunc(), log, osStateOpcode, "hostname", rasmonitoring.PrefixDirect)
			}()

			// wait until we finish to check the request
			got := <-ch
			require.Equal(t, test.expected(), got)
		})
	}
}

func TestApiVersion(t *testing.T) {
	log := logger.New().WithField("test", "true")
	type testCase struct {
		name         string
		checkRequest func(c *Client, r *http.Request)
		sendResponse func() (*http.Response, error)
		ctxFunc      func() context.Context
		expected     func() []rasmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)
				require.Equal(t, "https://hostname/ext/ras/direct/control_command", r.URL.String())

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				var req directRequest

				err = json.Unmarshal(body, &req)
				require.NoError(t, err)

				require.Equal(t, directRequest{
					Oam:    "0",
					Opcode: fmt.Sprintf("%d", osStateOpcode),
					Cpsp:   "0",
				}, req)

			},
			sendResponse: func() (*http.Response, error) {

				body, err := json.Marshal(directResponse{Cpsr: 769})
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			ctxFunc: context.Background,
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics,
					rasmonitoring.Metric{
						Hostname:    "hostname",
						Oam:         "0",
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixDirect, rasmonitoring.DirectMetricApiVersion),
						MetricValue: 301,
					})

				return metrics
			},
		},
		{
			name: "Ras doesn't send a response in the direct opcode",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)
				require.Equal(t, "https://hostname/ext/ras/direct/control_command", r.URL.String())

				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)

				var req directRequest

				err = json.Unmarshal(body, &req)
				require.NoError(t, err)

				require.Equal(t, directRequest{
					Oam:    "0",
					Opcode: fmt.Sprintf("%d", osStateOpcode),
					Cpsp:   "0",
				}, req)

			},
			sendResponse: func() (*http.Response, error) {

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
				}, nil
			},
			ctxFunc: context.Background,
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
				return test.sendResponse()
			}

			client.Auth = authDetails{
				Cookie: &http.Cookie{Name: "QSESSIONID"},
			}

			ch := make(chan []rasmonitoring.Metric)
			go func() {
				ch <- client.apiVersion(test.ctxFunc(), log, osStateOpcode, "hostname", rasmonitoring.PrefixDirect)
			}()

			// wait until we finish to check the request
			got := <-ch
			require.Equal(t, test.expected(), got)
		})
	}
}
