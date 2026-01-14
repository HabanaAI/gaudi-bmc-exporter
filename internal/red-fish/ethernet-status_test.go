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

func TestEthernetStatus(t *testing.T) {
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
					EthernetStatusResponse{
						Response: EthernetStatusResp{
							PortMap:                "Internal Port",
							ToggleCount:            10,
							ExternalLinkStatus:     "Active Port",
							LinkStatus:             "Connected",
							PHYStatus:              "Ready",
							BERCorrectable:         0,
							BERUncorrectable:       1,
							Nack:                   2,
							CRC:                    3,
							RetransmissionTimeout:  4,
							LinkRetrainingDueToBER: 5,
							MACRemote:              6,
							Retransmission:         7,
							Retraining:             8,
							SERPreFEC:              9,
							SERPostFEC:             10,
							Latency:                11,
							Throughput:             12,
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
				prefix := bmcmonitoring.PrefixEthernetStatus

				// Port Mapping
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricPortMapping),
					CustomLabels: map[string]string{
						"type": "Internal Port",
						"port": "0",
					},
				})

				// Toggle Count

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 10,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricStateTogglingCounter),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})

				// ExternalLinkStatus

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricExternalLinkStatus),
					CustomLabels: map[string]string{
						"state": "Active Port",
						"link":  "0",
					},
				})

				// LinkStatus

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricLinkStatus),
					CustomLabels: map[string]string{
						"state": "Connected",
						"link":  "0",
					},
				})

				// PHYStatus
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricPHYStatus),
					CustomLabels: map[string]string{
						"state": "Ready",
						"phy":   "0",
					},
				})

				// BERCorrectable
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricBERCorrectable),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// BERUncorrectable
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricBERUncorrectable),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// Nack
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 2,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricNack),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})

				// CRC
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 3,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricCRC),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// RetransmissionTimeout
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 4,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricRetransmissionTimeout),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// LinkRetrainingDueToBER
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 5,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricLinkRetrainingDueToBER),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// MACRemote
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 6,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricMACRemote),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// Retransmission
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 7,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricRetransmission),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// Retraining
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 8,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricRetraining),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// SERPreFEC
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 9,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricSERPreFEC),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// SERPostFEC
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 10,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricSERPostFEC),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// Latency
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 11,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricLatency),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				// Throughput
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 12,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricThroughput),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})
				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifyEthernetStatus(metric)
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
						Ports:    1,
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/EthernetStatus/port/0", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.EthernetStatus(context.Background(), log)
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
