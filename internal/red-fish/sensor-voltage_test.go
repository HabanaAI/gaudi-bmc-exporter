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

func TestSensorVoltage(t *testing.T) {
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
					SensorVoltageResponse{
						Response: SensorVoltageResp{
							VADC54:          54372,
							Vrm1In:          54000,
							Vrm1Out:         13548,
							Vrm2In:          13500,
							Vrm2VddOut:      815,
							Vrm2HbmOut:      1200,
							VmonPcieVph1P8V: 1794,
							Vmon1P8HbmVaa:   1804,
							Vmon2P5:         2502,
							Vmon48VHimon:    1800,
							VmonP5V:         1112,
							Vmon12V1:        1462,
							VmonHbm:         1204,
							VmonCore:        818,
							CpldHimon1P8NIC: 4976,
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

				prefix := bmcmonitoring.PrefixSensorVoltage

				// VADC54
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 54372,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVADC54),
				})

				// Vrm1In
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 54000,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVRM1in),
				})

				// Vrm1Out
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 13548,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVRM1out),
				})

				// Vrm2In
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 13500,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVRM2in),
				})

				// Vrm2VddOut
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 815,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVRM2VDDout),
				})

				// Vrm2HbmOut
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1200,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVRM2HBMout),
				})

				// VmonPcieVph1P8V
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1794,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMONPCIEVPH1P8V),
				})

				// Vmon1P8HbmVaa
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1804,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMON1P8HBMVAA),
				})

				// Vmon2P5
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 2502,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMON2P5),
				})

				// Vmon48VHimon
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1800,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMON48VHIMON),
				})

				// VmonP5V
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1112,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMONP5V),
				})

				// Vmon12V1
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1462,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMON12V1),
				})

				// VmonHbm
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1204,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMONHBM),
				})

				// VmonCore
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 818,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMONCore),
				})

				// CpldHimon1P8NIC
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 4976,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageCPLDHIMON1P8NIC),
				})
				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifySensorVoltage(metric)
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/SensorsVoltage", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.SensorVoltage(context.Background(), log)
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
