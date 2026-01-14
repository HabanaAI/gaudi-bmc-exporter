package rasmonitoringapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMetrics(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name         string
		checkRequest func(c *Client, r *http.Request)
		sendResponse func() (*http.Response, error)
		ctxFunc      func() context.Context
		expected     func() []rasmonitoring.Metric
		opcodes      map[string]Opcode
		method       string
	}

	tests := []testCase{
		{
			name:    "Ok string without decoder, with direct method",
			ctxFunc: context.Background,
			opcodes: map[string]Opcode{
				"field": {
					OpcodeNumber: 1,
					Offset:       1,
					Length:       1,
					ExpectedType: String,
				}},
			method: methodDirect,
			checkRequest: func(c *Client, r *http.Request) {
				require.Equal(t, fmt.Sprintf("https://hostname/ext/ras/%s/read?length=1&oam=0&offset=1", methodDirect), r.URL.String())
				verifyHeaders(t, c, r)
			},
			sendResponse: func() (*http.Response, error) {
				var data OpcodeData
				data.Data = []byte{29, 163}

				body, err := json.Marshal(data)
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  "prefix_field",
					CustomLabels: map[string]string{
						"field": "1da3",
					},
				})
				return metrics
			},
		},
		{
			name:    "Ok int without decoder, with indirect method",
			ctxFunc: context.Background,
			opcodes: map[string]Opcode{
				"field": {
					OpcodeNumber: 1,
					Offset:       1,
					Length:       1,
					ExpectedType: Int,
				}},
			method: methodIndirect,
			checkRequest: func(c *Client, r *http.Request) {
				require.Equal(t, fmt.Sprintf("https://hostname/ext/ras/%s/read?length=1&oam=0&offset=1&opcode=1", methodIndirect), r.URL.String())
				verifyHeaders(t, c, r)
			},
			sendResponse: func() (*http.Response, error) {
				var data OpcodeData
				data.Data = []byte{67, 121, 235, 141}

				body, err := json.Marshal(data)
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 2381019459,
					MetricName:  "prefix_field",
				})
				return metrics
			},
		},
		{
			name: "no data sent",
			checkRequest: func(c *Client, r *http.Request) {
				require.Equal(t, fmt.Sprintf("https://hostname/ext/ras/%s/read?length=1&oam=0&offset=1&opcode=1", methodIndirect), r.URL.String())
				verifyHeaders(t, c, r)
			},
			sendResponse: func() (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
				}, nil
			},
			ctxFunc: context.Background,
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric
				return metric
			},
			opcodes: map[string]Opcode{
				"field": {
					OpcodeNumber: 1,
					Offset:       1,
					Length:       1,
					ExpectedType: Int,
				}},
			method: methodIndirect,
		},
		{
			name: "data sent empty",
			checkRequest: func(c *Client, r *http.Request) {
				require.Equal(t, fmt.Sprintf("https://hostname/ext/ras/%s/read?length=1&oam=0&offset=1&opcode=1", methodIndirect), r.URL.String())
				verifyHeaders(t, c, r)
			},
			sendResponse: func() (*http.Response, error) {
				var data OpcodeData
				body, err := json.Marshal(data)
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			ctxFunc: context.Background,
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric
				return metric
			},
			opcodes: map[string]Opcode{
				"field": {
					OpcodeNumber: 1,
					Offset:       1,
					Length:       1,
					ExpectedType: Int,
				}},
			method: methodIndirect,
		},
		{
			name: "data sent empty",
			checkRequest: func(c *Client, r *http.Request) {
				require.Equal(t, fmt.Sprintf("https://hostname/ext/ras/%s/read?length=1&oam=0&offset=1&opcode=1", methodIndirect), r.URL.String())
				verifyHeaders(t, c, r)
			},
			sendResponse: func() (*http.Response, error) {
				var data OpcodeData
				data.Data = []byte{67, 121, 235, 141}

				body, err := json.Marshal(data)
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			ctxFunc: context.Background,
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric
				return metric
			},
			opcodes: map[string]Opcode{
				"field": {
					OpcodeNumber: 1,
					Offset:       1,
					Length:       1,
					ExpectedType: "unsupported",
				}},
			method: methodIndirect,
		},
		{
			name:    "Ok with custom decoder",
			ctxFunc: context.Background,
			opcodes: map[string]Opcode{
				"field": {
					OpcodeNumber: 1,
					Offset:       1,
					Length:       1,
					ExpectedType: Int,
					Decoder: func(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
						var metrics []rasmonitoring.Metric

						value, err := strconv.Atoi(val)
						if err != nil {
							return nil, err
						}
						metrics = append(metrics, rasmonitoring.Metric{
							Hostname:    hostname,
							Oam:         fmt.Sprintf("%d", oam),
							MetricName:  fmt.Sprintf("prefix_%s", fieldName),
							MetricValue: value,
						})
						return metrics, nil
					},
				}},
			method: methodIndirect,
			checkRequest: func(c *Client, r *http.Request) {
				require.Equal(t, fmt.Sprintf("https://hostname/ext/ras/%s/read?length=1&oam=0&offset=1&opcode=1", methodIndirect), r.URL.String())
				verifyHeaders(t, c, r)
			},
			sendResponse: func() (*http.Response, error) {
				var data OpcodeData
				data.Data = []byte{67, 121, 235, 141}

				body, err := json.Marshal(data)
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 2381019459,
					MetricName:  "prefix_field",
				})
				return metrics
			},
		},
		{
			name:    "custom decoder error",
			ctxFunc: context.Background,
			opcodes: map[string]Opcode{
				"field": {
					OpcodeNumber: 1,
					Offset:       1,
					Length:       1,
					ExpectedType: Int,
					Decoder: func(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
						return nil, errors.New("failed to decode")
					},
				}},
			method: methodIndirect,
			checkRequest: func(c *Client, r *http.Request) {
				require.Equal(t, fmt.Sprintf("https://hostname/ext/ras/%s/read?length=1&oam=0&offset=1&opcode=1", methodIndirect), r.URL.String())
				verifyHeaders(t, c, r)
			},
			sendResponse: func() (*http.Response, error) {
				var data OpcodeData
				data.Data = []byte{67, 121, 235, 141}

				body, err := json.Marshal(data)
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
				ch <- client.getMetrics(test.ctxFunc(), log, test.opcodes, "prefix", test.method)
			}()

			// wait until we finish to check the request
			got := <-ch
			require.Equal(t, test.expected(), got)
		})
	}
}
