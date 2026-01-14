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

func TestPcieInfo(t *testing.T) {
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
					PcieInfoResponse{
						Response: PcieInfoResp{
							MaxPCIeLinkSpeed:                     "Gen4",
							CurrentPCIeLinkSpeed:                 "Gen2",
							MaxPCIeLinkWidth:                     16,
							CurrentPCIeLinkWidth:                 16,
							PCIeDeviceID:                         1020,
							PCIeSubsystemID:                      1020,
							PCIeSubsystemVendorID:                "1DA3",
							PCIeBusAndDevice:                     "4D:00",
							CorrectedInternalErrorStatus:         0,
							ReplayBufferNumRolloverError:         1,
							ReplayTimerTimeoutError:              2,
							BadTLPCounter:                        3,
							BadDLLPCounter:                       4,
							ReceiverErrorCounter:                 5,
							LCRCErrorCounter:                     6,
							ECRCErrorCounter:                     7,
							CompletionTimeoutIndication:          8,
							UncorrectableInternalErrorIndication: 9,
							ReceiverOverflowIndication:           10,
							FlowControlProtocolErrorIndication:   11,
							SurpriseLinkDownIndication:           12,
							MalfunctionTLPErrorIndication:        13,
							DLLPProtocolErrorIndication:          14,
							RxNakDLLPCounter:                     15,
							TxNakDLLPCounter:                     16,
							RetryTLPCounter:                      17,
							PWRBRKIndication:                     18,
							PCIeRxMemoryWriteCounter:             19,
							PCIeRxMemoryReadCounter:              20,
							PCIeTxMemoryWriteCounter:             21,
							PCIeTxMemoryReadCounter:              22,
							AERCapabilityControlOffset:           23,
							AERErrorLog:                          24,
							PCIeFWVersion:                        25,
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

				prefix := rasmonitoring.PrefixPcieInfo

				// MaxPCIeLinkSpeed

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricMaxPCIeLinkSpeed),
					MetricValue: 4,
					CustomLabels: map[string]string{
						rasmonitoring.PcieInfoMetricMaxPCIeLinkSpeed: "Gen4",
					},
				})

				// CurrentPCIeLinkSpeed

				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricCurrentPCIeLinkSpeed),
					MetricValue: 2,
					CustomLabels: map[string]string{
						rasmonitoring.PcieInfoMetricCurrentPCIeLinkSpeed: "Gen2",
					},
				})

				// MaxPCIeLinkWidth
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricMaxPCIeLinkWidth),
					MetricValue: 16,
				})

				// CurrentPCIeLinkWidth
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricCurrentPCIeLinkWidth),
					MetricValue: 16,
				})

				// PCIeDeviceID
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeDeviceID),
					MetricValue: 1020,
				})

				// PCIeSubsystemID
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeSubsystemID),
					MetricValue: 1020,
				})

				// PCIeSubsystemVendorID
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeSubsystemVendorID),
					MetricValue: 0,
					CustomLabels: map[string]string{
						rasmonitoring.PcieInfoMetricPCIeSubsystemVendorID: "1DA3",
					},
				})

				// PCIeBusAndDevice
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeBusAndDevice),
					MetricValue: 0,
					CustomLabels: map[string]string{
						rasmonitoring.PcieInfoMetricPCIeBusAndDevice: "4D:00",
					},
				})

				// CorrectedInternalErrorStatus
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricCorrectedInternalErrorStatus),
					MetricValue: 0,
				})

				// ReplayBufferNumRolloverError
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricReplayBufferNumRolloverError),
					MetricValue: 1,
				})

				// ReplayTimerTimeoutError
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricReplayTimerTimeoutError),
					MetricValue: 2,
				})

				// BadTLPCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricBadTLPCounter),
					MetricValue: 3,
				})

				// BadDLLPCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricBadDLLPCounter),
					MetricValue: 4,
				})

				// ReceiverErrorCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricReceiverErrorCounter),
					MetricValue: 5,
				})

				// LCRCErrorCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricLCRCErrorCounter),
					MetricValue: 6,
				})

				// ECRCErrorCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricECRCErrorCounter),
					MetricValue: 7,
				})

				// CompletionTimeoutIndication
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricCompletionTimeoutIndication),
					MetricValue: 8,
				})

				// UncorrectableInternalErrorIndication
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricUncorrectableInternalErrorIndication),
					MetricValue: 9,
				})

				// ReceiverOverflowIndication
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricReceiverOverflowIndication),
					MetricValue: 10,
				})

				// FlowControlProtocolErrorIndication
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricFlowControlProtocolErrorIndication),
					MetricValue: 11,
				})

				// SurpriseLinkDownIndication
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricSurpriseLinkDownIndication),
					MetricValue: 12,
				})

				// MalfunctionTLPErrorIndication
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricMalfunctionTLPErrorIndication),
					MetricValue: 13,
				})

				// DLLPProtocolErrorIndication
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricDLLPProtocolErrorIndication),
					MetricValue: 14,
				})

				// RxNakDLLPCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricRXNakDLLPCounter),
					MetricValue: 15,
				})

				// TxNakDLLPCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricTxNakDLLPCounter),
					MetricValue: 16,
				})

				// RetryTLPCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricRetryTLPcounter),
					MetricValue: 17,
				})

				// PWRBRKIndication
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPWRBRKindication),
					MetricValue: 18,
				})

				// PCIeRxMemoryWriteCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeRXMemoryWriteCounter),
					MetricValue: 19,
				})

				// PCIeRxMemoryReadCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeRXMemoryReadCounter),
					MetricValue: 20,
				})

				// PCIeTxMemoryWriteCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeTXMemoryWriteCounter),
					MetricValue: 21,
				})

				// PCIeTxMemoryReadCounter
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeTXMemoryReadCounter),
					MetricValue: 22,
				})

				// AERCapabilityControlOffset
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricAERCapabilityControlOffset),
					MetricValue: 23,
				})

				// AERErrorLog
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricAERerrorlog),
					MetricValue: 24,
				})

				// PCIeFWVersion
				metrics = append(metrics, rasmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeFWversion),
					MetricValue: 25,
				})
				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifyPcieInfo(metric)
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/PCIeInfo", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.PcieInfo(context.Background(), log)
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
