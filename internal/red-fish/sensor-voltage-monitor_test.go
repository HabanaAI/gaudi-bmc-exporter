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

func TestSensorVoltageMonitor(t *testing.T) {
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
					SensorVoltageMonitorResponse{
						Response: SensorVoltageMonitorResp{
							SwVm:       818,
							SwCpeEu0:   818,
							SwCpeEu1:   818,
							SwCpeEu2:   818,
							SwcpeEu3:   23214,
							SwCpeHbm:   818,
							SwTpcSb:    818,
							SwTft:      818,
							SwCnt:      818,
							SwMcHbm00:  807,
							SwMcHbm01:  23214,
							SwHconHbm0: 818,
							SwPpw0:     818,
							SwPpw1:     818,
							SwL2cMacro: 818,
							SwVcdMacro: 818,
							SeVm:       818,
							SePe00:     818,
							SePe01:     26213,
							SeEuCore:   818,
							SeBpyramid: 818,
							SeRx:       818,
							SeMx:       818,
							SeHconHbm1: 818,
							SeMcHbm1:   805,
							SeGasket:   818,
							SeMmeCtrl:  817,
							SeMmeQman:  817,
							SeSbte:     818,
							SeRtrDn:    817,
							SeRtrUp:    818,
							SeSram:     818,
							NwVm:       818,
							NwPe00:     818,
							NwPe01:     818,
							NwPe02:     818,
							NwPe03:     26213,
							NwEuCore:   818,
							NwBpyramid: 818,
							NwMcHbm40:  806,
							NwMcHbm41:  26214,
							NwHconHbm4: 818,
							NwAcc:      818,
							NwWap:      818,
							NwTif:      818,
							NwRtrDn:    818,
							NwRtrUp:    818,
							NwSram:     818,
							NeVM:       818,
							NeTx:       818,
							NeMx0:      818,
							NeMx1:      818,
							NeHconHbm5: 818,
							NeMcHbm5:   807,
							NeSob:      818,
							NeQnt:      818,
							NeCnt0:     818,
							NeCnt1:     26214,
							NeCpeEu:    818,
							NeCntHbm:   819,
							NeSbte0:    818,
							NeSbte1:    26214,
							NeMif:      818,
							NeRtrDn:    818,
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

				prefix := bmcmonitoring.PrefixSensorVoltageMonitor

				// SwVm

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwVm),
				})

				// SwCpeEu0
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwCpeEu_0),
				})

				// SwCpeEu1
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwCpeEu_1),
				})

				// SwCpeEu2
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwCpeEu_2),
				})

				// SwcpeEu3
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 23214,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwcpeEu_3),
				})

				// SwCpeHbm
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwCpeHbm),
				})

				// SwTpcSb
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwTpcSb),
				})

				// SwTft
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwTft),
				})

				// SwCnt
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwCnt),
				})

				// SwMcHbm00
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 807,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwMcHbm0_0),
				})

				// SwMcHbm01
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 23214,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwMcHbm0_1),
				})

				// SwHconHbm0
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwHconHbm0),
				})

				// SwPpw0
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwPpw_0),
				})

				// SwPpw1
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwPpw_1),
				})

				// SwL2cMacro
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwL2cMacro),
				})

				// SwVcdMacro
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwVcdMacro),
				})

				// SeVm
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeVm),
				})

				// SePe00
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSePe0_0),
				})

				// SePe01
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 26213,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSePe0_1),
				})

				// SeEuCore
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeEuCore),
				})

				// SeBpyramid
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeBpyramid),
				})

				// SeRx
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeRx),
				})

				// SeMx
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeMx),
				})

				// SeHconHbm1
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeHconHbm1),
				})

				// SeMcHbm1
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 805,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeMcHbm1),
				})

				// SeGasket
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeGasket),
				})

				// SeMmeCtrl
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 817,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeMmeCtrl),
				})

				// SeMmeQman
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 817,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeMmeQman),
				})

				// SeSbte
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeSbte),
				})

				// SeRtrDn
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 817,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeRtrDn),
				})

				// SeRtrUp
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeRtrUp),
				})

				// SeSram
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeSram),
				})

				// NwVm
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwVm),
				})

				// NwPe00
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwPe0_0),
				})

				// NwPe01
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwPe0_1),
				})

				// NwPe02
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwPe0_2),
				})

				// NwPe03
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 26213,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwPe0_3),
				})

				// NwEuCore
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwEuCore),
				})

				// NwBpyramid
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwBpyramid),
				})

				// NwMcHbm40
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 806,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwMcHbm4_0),
				})

				// NwMcHbm41
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 26214,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwMcHbm4_1),
				})

				// NwHconHbm4
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwHconHbm4),
				})

				// NwAcc
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwAcc),
				})

				// NwWap
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwWap),
				})

				// NwTif
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwTif),
				})

				// NwRtrDn
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwRtrDn),
				})

				// NwRtrUp
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwRtrUp),
				})

				// NwSram
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwSram),
				})

				// NeVM
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeVM),
				})

				// NeTx
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeTx),
				})

				// NeMx0
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeMx_0),
				})

				// NeMx1
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeMx_1),
				})

				// NeHconHbm5
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeHconHbm5),
				})

				// NeMcHbm5
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 807,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeMcHbm5),
				})

				// NeSob
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeSob),
				})

				// NeQnt
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeQnt),
				})

				// NeCnt0
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeCnt_0),
				})

				// NeCnt1
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 26214,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeCnt_1),
				})

				// NeCpeEu
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeCpeEu),
				})

				// NeCntHbm
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 819,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeCntHbm),
				})

				// NeSbte0
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeSbte_0),
				})

				// NeSbte1
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 26214,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeSbte_1),
				})

				// NeMif
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeMif),
				})

				// NeRtrDn
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeRtrDn),
				})

				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifySensorVoltageMonitor(metric)
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/SensorsVoltageMonitor", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.SensorVoltageMonitor(context.Background(), log)
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
