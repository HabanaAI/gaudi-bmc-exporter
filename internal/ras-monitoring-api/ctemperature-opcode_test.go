package rasmonitoringapi

import (
	"bytes"
	"context"
	"encoding/json"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCTemperature(t *testing.T) {
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
			name:    "Ok",
			ctxFunc: context.Background,
			checkRequest: func(c *Client, r *http.Request) {
				checkCTemperatureRequest(t, c, r)
			},
			sendResponse: func() (r *http.Response, err error) {

				body, err := json.Marshal(CTemperatureResponse{Temperature: 101})
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			expected: func() []rasmonitoring.Metric {
				return []rasmonitoring.Metric{
					{
						Hostname:    "hostname",
						Oam:         "0",
						MetricValue: 101,
						MetricName:  rasmonitoring.PrefixCTemperature,
					},
				}
			},
		},
		{
			name: "Error getting ctemperature",
			checkRequest: func(c *Client, r *http.Request) {
				checkCTemperatureRequest(t, c, r)
			},
			sendResponse: func() (*http.Response, error) {
				body, err := json.Marshal(
					RasError{Err: "HLMI call failed, oam 0, ret = 3", Code: 1002},
				)
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			ctxFunc: context.Background,
			expected: func() []rasmonitoring.Metric {
				var metric []rasmonitoring.Metric

				return metric
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
				ch <- client.CTemperature(test.ctxFunc(), log)
			}()

			// wait until we finish to check the request
			got := <-ch
			require.Equal(t, test.expected(), got)
		})
	}
}

func verifyHeaders(t *testing.T, c *Client, r *http.Request) {

	require.Equal(t, c.Auth.CSRFToken, r.Header.Get(tokenHeader))
	require.Equal(t, r.Header.Get("Content-Type"), "application/json")

	cookie, err := r.Cookie("QSESSIONID")
	require.NoError(t, err)
	require.Equal(t, c.Auth.Cookie, cookie)

}

func checkCTemperatureRequest(t *testing.T, c *Client, r *http.Request) {
	require.Equal(t, "https://hostname/ext/ras/direct/temperature?oam=0", r.URL.String())

	verifyHeaders(t, c, r)
}
