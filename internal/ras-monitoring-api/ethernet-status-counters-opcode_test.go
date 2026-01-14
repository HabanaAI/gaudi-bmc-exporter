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

func TestStateTogglingCounter(t *testing.T) {
	log := logger.New().WithField("test", "true")
	type testCase struct {
		name         string
		checkRequest func(c *Client, r *http.Request)
		sendResponse func(r *http.Request) (*http.Response, error)
		ctxFunc      func() context.Context
		expected     string
		wantErr      bool
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)
				switch r.Method {
				// write
				case http.MethodPut:
					require.Equal(t, "https://hostname/ext/ras/indirect/write", r.URL.String())
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					var req writeRequest
					err = json.Unmarshal(body, &req)
					require.NoError(t, err)
					require.Equal(t, writeRequest{Oam: "0", Opcode: "5", Offset: "234", Length: "4", Data: []int{0, 0, 0, 0}}, req)
				case http.MethodGet:
					require.Equal(t, fmt.Sprintf("https://hostname/ext/ras/%s/read?length=4&oam=0&offset=238&opcode=5", methodIndirect), r.URL.String())
				default:
					t.Error("unexpected method")
				}
			},
			sendResponse: func(r *http.Request) (*http.Response, error) {
				switch r.Method {
				case http.MethodPut:
					return &http.Response{
						StatusCode: http.StatusOK,
					}, nil
				case http.MethodGet:

					body, err := json.Marshal(OpcodeData{Data: []byte{67, 121, 235, 141}})
					require.NoError(t, err)

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(body)),
					}, nil
				default:
					t.Error("unexpected method")
				}

				return nil, nil
			},
			ctxFunc:  context.Background,
			expected: "2381019459",
			wantErr:  false,
		},
		{
			name: "Failed to write",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)
				require.Equal(t, "https://hostname/ext/ras/indirect/write", r.URL.String())
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				var req writeRequest
				err = json.Unmarshal(body, &req)
				require.NoError(t, err)
				require.Equal(t, writeRequest{Oam: "0", Opcode: "5", Offset: "234", Length: "4", Data: []int{0, 0, 0, 0}}, req)
			},
			sendResponse: func(r *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
				}, nil
			},
			ctxFunc:  context.Background,
			wantErr:  true,
			expected: "",
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
						Ports:    1,
					},
					Client: http.Client{
						Transport: transport,
					},
				},
			}
			// check request
			transport.RoundTripFunc = func(r *http.Request) (*http.Response, error) {

				test.checkRequest(client, r)
				resp, err := test.sendResponse(r)
				return resp, err
			}

			client.Auth = authDetails{
				Cookie: &http.Cookie{Name: "QSESSIONID"},
			}

			type res struct {
				err  error
				resp string
			}
			ch := make(chan res)
			go func() {
				ctx := test.ctxFunc()
				result, err := client.stateTogglingCounter(ctx, 0, 0)

				ch <- res{
					err:  err,
					resp: result,
				}
			}()

			// wait until we finish to check the request
			got := <-ch

			if test.wantErr {
				require.Error(t, got.err)
			} else {
				require.NoError(t, got.err)
			}
			require.Equal(t, test.expected, got.resp)
		})
	}
}
func TestGetEthernetCountersInfo(t *testing.T) {
	log := logger.New().WithField("test", "true")
	type testCase struct {
		name         string
		checkRequest func(c *Client, r *http.Request)
		sendResponse func() (*http.Response, error)
		ctxFunc      func() context.Context
		expected     func() ethernetCountersInfoResponse
		wantErr      bool
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)
				require.Equal(t, "https://hostname/ext/ras/indirect/eth_counters?oam=0&port_idx=0", r.URL.String())
			},
			ctxFunc: context.Background,
			expected: func() ethernetCountersInfoResponse {
				return ethernetCountersInfoResponse{
					BitErrRateCorrectable:        1,
					BitErrRateUncorrectable:      2,
					NackCounter:                  3,
					RetransmissionTimeoutCounter: 4,
					LinkRetrainingBarCounter:     5,
					MacRemoteErrCounter:          6,
					RetransmissionCounter:        7,
					RetrainingCounter:            8,
					CrcErrorCounter:              9,
					SerPreFec:                    10,
					SerPostFec:                   11,
					Latency:                      12,
					Throughput:                   13,
				}
			},
			sendResponse: func() (*http.Response, error) {

				body, err := json.Marshal(ethernetCountersInfoResponse{
					BitErrRateCorrectable:        1,
					BitErrRateUncorrectable:      2,
					NackCounter:                  3,
					RetransmissionTimeoutCounter: 4,
					LinkRetrainingBarCounter:     5,
					MacRemoteErrCounter:          6,
					RetransmissionCounter:        7,
					RetrainingCounter:            8,
					CrcErrorCounter:              9,
					SerPreFec:                    10,
					SerPostFec:                   11,
					Latency:                      12,
					Throughput:                   13,
				})

				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			wantErr: false,
		},
		{
			name: "Ras returns an empty body",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)
				require.Equal(t, "https://hostname/ext/ras/indirect/eth_counters?oam=0&port_idx=0", r.URL.String())
			},
			sendResponse: func() (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
				}, nil
			},
			ctxFunc: context.Background,
			expected: func() ethernetCountersInfoResponse {
				return ethernetCountersInfoResponse{}
			},
			wantErr: true,
		},
		{
			name: "Ras return an error",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)
				require.Equal(t, "https://hostname/ext/ras/indirect/eth_counters?oam=0&port_idx=0", r.URL.String())
			},
			sendResponse: func() (*http.Response, error) {
				body, err := json.Marshal(RasError{
					Err:  "some error",
					Code: 1002,
				})

				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			ctxFunc: context.Background,
			wantErr: true,
			expected: func() ethernetCountersInfoResponse {
				return ethernetCountersInfoResponse{}
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
						Ports:    1,
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

			type res struct {
				err  error
				resp ethernetCountersInfoResponse
			}
			ch := make(chan res)
			go func() {
				ctx := test.ctxFunc()
				result, err := client.getEthernetCountersInfo(ctx, 0, 0)

				ch <- res{
					err:  err,
					resp: result,
				}
			}()

			// wait until we finish to check the request
			got := <-ch

			if test.wantErr {
				require.Error(t, got.err)
			} else {
				require.NoError(t, got.err)
			}
			require.Equal(t, test.expected(), got.resp)
		})
	}

}
