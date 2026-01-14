package redfish

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEthernetInfo(t *testing.T) {
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
					EthernetInfoResponse{
						Response: EthernetInfoResp{
							SerdesAvailability: "Available",
							PortMaxSpeed:       400,
							ANLTStatus: []string{
								"Enabled",
								"Disabled",
							},
							NumberOfLanes: 48,
							NumberOfLinks: 48,
							LinkSpeed:     56,
							PortMap: []string{
								"Internal Port",
								"External Port",
							},
							ToggleCount: []int{
								16,
							},
							ExternalLinkStatus: []string{
								"Non-active Port",
								"Active Port",
							},
							LinkStatus: []string{
								"Not connected",
								"Connected",
							},
							PHYStatus: []string{
								"Not ready",
								"Ready",
							},
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
				prefix := rasmonitoring.PrefixEthernetInfo

				// SerdesAvailability
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricSerDesAvailability),
					CustomLabels: map[string]string{
						"state": "Available",
					},
				})

				// PortMaxSpeed
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 400,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricPortMaxSpeed),
				})

				// ANLTStatus
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricANLTStatus),
					CustomLabels: map[string]string{
						"state": "Enabled",
						"port":  "0",
					},
				})

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricANLTStatus),
					CustomLabels: map[string]string{
						"state": "Disabled",
						"port":  "1",
					},
				})

				// NumberOfLanes
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 48,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricNumberOfLanes),
				})

				// NumberOfLinks
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 48,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricNumberOfLinks),
				})

				// LinkSpeed
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 56,
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricLinkSpeed),
				})

				// port map
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricPortMapping),
					CustomLabels: map[string]string{
						"type": "Internal Port",
						"port": "0",
					},
				})

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricPortMapping),
					CustomLabels: map[string]string{
						"type": "External Port",
						"port": "1",
					},
				})

				// Toggle count
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 16,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricStateTogglingCounter),
					CustomLabels: map[string]string{
						"port": "0",
					},
				})

				// ExternalLinkStatus
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricExternalLinkStatus),
					CustomLabels: map[string]string{
						"state": "Non-active Port",
						"link":  "0",
					},
				})

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricExternalLinkStatus),
					CustomLabels: map[string]string{
						"state": "Active Port",
						"link":  "1",
					},
				})

				// LinkStatus
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricLinkStatus),
					CustomLabels: map[string]string{
						"state": "Not connected",
						"link":  "0",
					},
				})

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricLinkStatus),
					CustomLabels: map[string]string{
						"state": "Connected",
						"link":  "1",
					},
				})

				// PHYStatus
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricPHYStatus),
					CustomLabels: map[string]string{
						"state": "Not ready",
						"phy":   "0",
					},
				})

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricPHYStatus),
					CustomLabels: map[string]string{
						"state": "Ready",
						"phy":   "1",
					},
				})

				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifyEthernetInfo(metric)
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/EthernetInfo", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.EthernetInfo(context.Background(), log)
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
