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

func TestGetReplacedRows(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name         string
		checkRequest func(c *Client, r *http.Request)
		sendResponse func(r *http.Request) (*http.Response, error)
		expected     func() []rasmonitoring.Metric
	}

	numOfReplacedRowsUrl := fmt.Sprintf("https://hostname/ext/ras/%s/read?length=%d&oam=0&offset=%d&opcode=%d",
		methodIndirect, numOfReplacedRowsOp.Length, numOfReplacedRowsOp.Offset, numOfReplacedRowsOp.OpcodeNumber)

	replacedRowsUrl := "https://hostname/ext/ras/indirect/dram_replaced_row?index=0&oam=0"

	tests := []testCase{
		{
			name: "1 replaced row, returns a valid metric",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)

				require.True(t, r.URL.String() == numOfReplacedRowsUrl || r.URL.String() == replacedRowsUrl)
			},
			sendResponse: func(r *http.Request) (*http.Response, error) {
				switch r.URL.String() {

				case numOfReplacedRowsUrl:
					body, err := json.Marshal(OpcodeData{Data: []byte{1}})
					require.NoError(t, err)

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(body)),
					}, nil

				case replacedRowsUrl:

					body, err := json.Marshal(replacedRowResponse{
						HBMIndex:           1,
						SidIndex:           2,
						PCIndex:            3,
						BankIndex:          4,
						Cause:              5,
						ReplacedRowAddress: 6,
					})
					require.NoError(t, err)

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(body)),
					}, nil
				}

				return nil, nil
			},
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricNumOfReplacedRows),
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
				})

				metrics = append(metrics, rasmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricReplaceRows),
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					CustomLabels: map[string]string{
						rasmonitoring.HBMMetricReplaceRowsLabelHBMIndex:   "1",
						rasmonitoring.HBMMetricReplaceRowsLabelPCIndex:    "2",
						rasmonitoring.HBMMetricReplaceRowsLabelStackID:    "3",
						rasmonitoring.HBMMetricReplaceRowsLabelBankIndex:  "4",
						rasmonitoring.HBMMetricReplaceRowsLabelCause:      "5",
						rasmonitoring.HBMMetricReplaceRowsLabelRowAddress: "6",
					},
				})
				return metrics
			},
		},
		{
			name: "0 replaced rows",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)

				require.True(t, r.URL.String() == numOfReplacedRowsUrl)
			},
			sendResponse: func(r *http.Request) (*http.Response, error) {
				body, err := json.Marshal(OpcodeData{Data: []byte{0}})
				require.NoError(t, err)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricNumOfReplacedRows),
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
				})

				return metrics
			},
		},
	}

	for _, test := range tests {
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

		ch := make(chan []rasmonitoring.Metric)
		go func() {
			ch <- client.getReplacedRows(context.Background(), log)
		}()

		// wait until we finish to check the request
		got := <-ch
		require.Equal(t, test.expected(), got)
	}
}
func TestGetRepairedLanes(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name         string
		checkRequest func(c *Client, r *http.Request)
		sendResponse func(r *http.Request) (*http.Response, error)
		expected     func() []rasmonitoring.Metric
	}

	numOfLanesUrl := fmt.Sprintf("https://hostname/ext/ras/%s/read?length=%d&oam=0&offset=%d&opcode=%d",
		methodIndirect, numOfLanesOpcode.Length, numOfLanesOpcode.Offset, numOfLanesOpcode.OpcodeNumber)

	repairedLaneUrl := "https://hostname/ext/ras/indirect/dram_repaired_lane?index=0&oam=0"
	tests := []testCase{
		{
			name: "Number of lanes is 1, returns valid metrics",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)

				require.True(t, r.URL.String() == repairedLaneUrl || r.URL.String() == numOfLanesUrl)
			},
			sendResponse: func(r *http.Request) (*http.Response, error) {

				switch r.URL.String() {
				case numOfLanesUrl:

					body, err := json.Marshal(OpcodeData{Data: []byte{1}})
					require.NoError(t, err)

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(body)),
					}, nil
				case repairedLaneUrl:
					body, err := json.Marshal(RepairLaneResponse{LaneHBMIndex: 1, LaneCHIndex: 2})
					require.NoError(t, err)

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(body)),
					}, nil
				}

				return nil, nil
			},
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricNumOfRepairedLanes),
				})

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricRepairedLanes),
					CustomLabels: map[string]string{
						rasmonitoring.HBMMetricRepairedLanesLabelHBMIndex:  "1",
						rasmonitoring.HBMMetricRepairedLanesLabelMCChannel: "2",
					},
				})
				return metrics
			},
		},
		{
			name: "Number of lanes is 0",
			checkRequest: func(c *Client, r *http.Request) {
				verifyHeaders(t, c, r)

				require.True(t, r.URL.String() == numOfLanesUrl)
			},
			sendResponse: func(r *http.Request) (*http.Response, error) {
				body, err := json.Marshal(OpcodeData{Data: []byte{0}})
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
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricNumOfRepairedLanes),
				})
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

			ch := make(chan []rasmonitoring.Metric)
			go func() {
				ch <- client.getRepairedLanes(context.Background(), log)
			}()

			// wait until we finish to check the request
			got := <-ch
			require.Equal(t, test.expected(), got)
		})
	}
}

func TestGetHbmRepairStatusArray(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name         string
		checkRequest func(c *Client, r *http.Request)
		sendResponse func(r *http.Request) (*http.Response, error)
		expected     func() []rasmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Ras returns valid metrics",
			checkRequest: func(c *Client, r *http.Request) {

				verifyHeaders(t, c, r)
				switch r.Method {
				// write
				case http.MethodPut:

					checkHbmWriteRequest(t, r)

				case http.MethodGet:

					require.True(t, r.URL.String() == fmt.Sprintf("https://hostname/ext/ras/%s/read?length=%d&oam=0&offset=%d&opcode=%d",
						methodIndirect, hbmMbistRepairOpcode.Length, hbmMbistRepairOpcode.Offset, hbmMbistRepairOpcode.OpcodeNumber) ||

						r.URL.String() == fmt.Sprintf("https://hostname/ext/ras/%s/read?length=%d&oam=0&offset=%d&opcode=%d",
							methodIndirect, hbmcGlobalEccOpcode.Length, hbmcGlobalEccOpcode.Offset, hbmcGlobalEccOpcode.OpcodeNumber))

				default:
					require.Fail(t, "unsupported method")
				}

			},
			sendResponse: func(r *http.Request) (*http.Response, error) {

				switch r.Method {

				case http.MethodPut:
					// write successfully
					return &http.Response{
						StatusCode: http.StatusOK,
					}, nil

				case http.MethodGet:

					MbistRepairUrl := fmt.Sprintf("https://hostname/ext/ras/%s/read?length=%d&oam=0&offset=%d&opcode=%d",
						methodIndirect, hbmMbistRepairOpcode.Length, hbmMbistRepairOpcode.Offset, hbmMbistRepairOpcode.OpcodeNumber)

					globalEccUrl := fmt.Sprintf("https://hostname/ext/ras/%s/read?length=%d&oam=0&offset=%d&opcode=%d",
						methodIndirect, hbmcGlobalEccOpcode.Length, hbmcGlobalEccOpcode.Offset, hbmcGlobalEccOpcode.OpcodeNumber)

					switch r.URL.String() {
					case MbistRepairUrl:

						body, err := json.Marshal(OpcodeData{Data: []byte{1}})
						require.NoError(t, err)

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewReader(body)),
						}, nil
					case globalEccUrl:

						body, err := json.Marshal(OpcodeData{Data: []byte{2}})
						require.NoError(t, err)

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewReader(body)),
						}, nil

					}

				default:
					require.Fail(t, "unsupported method")
				}

				return nil, nil
			},
			expected: func() []rasmonitoring.Metric {
				var metrics []rasmonitoring.Metric

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricMbistRepair),
					CustomLabels: map[string]string{
						rasmonitoring.HBMMetricMbistRepairLabelState: "Flow ran",
						rasmonitoring.HBMMetricMbistRepairLabelIndex: "0",
					},
				})

				metrics = append(metrics, rasmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricGlobalECC),
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 2,
					CustomLabels: map[string]string{
						rasmonitoring.HBMMetricGlobalECCLabelIndex: "0",
					},
				})
				return metrics
			},
		},
		{
			name: "Write operation failed",
			checkRequest: func(c *Client, r *http.Request) {
				checkHbmWriteRequest(t, r)
			},
			sendResponse: func(r *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("failed write request")
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

			ch := make(chan []rasmonitoring.Metric)
			go func() {
				ch <- client.getHbmRepairStatusArray(context.Background(), log, 1)
			}()

			// wait until we finish to check the request
			got := <-ch
			require.Equal(t, test.expected(), got)
		})
	}
}

func checkHbmWriteRequest(t *testing.T, r *http.Request) {
	require.Equal(t, "https://hostname/ext/ras/indirect/write", r.URL.String())
	body, err := io.ReadAll(r.Body)
	require.NoError(t, err)
	var req writeRequest
	err = json.Unmarshal(body, &req)
	require.NoError(t, err)
	require.Equal(t, writeRequest{
		Oam:    "0",
		Opcode: fmt.Sprintf("%d", hbmWriteOp.OpcodeNumber),
		Offset: fmt.Sprintf("%d", hbmWriteOp.Offset),
		Length: fmt.Sprintf("%d", hbmWriteOp.Length),
		Data:   []int{0}}, req)
}
