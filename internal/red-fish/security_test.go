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

func TestSecurity(t *testing.T) {
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
					SecurityResponse{
						Response: SecurityResp{
							CurrentPublicKeyHashIndex: 0,
							CurrentSVNVersion:         1,
							Key0Revocation:            "Not Revoked",
							Key1Revocation:            "Revoked",
							Key2Revocation:            "Not Revoked",
							Key3Revocation:            "Revoked",
							Key4Revocation:            "Not Revoked",
							MinimalSVNIndex:           "min",
							FWImageSource:             "Secondary",
							TpmPcrPpboot:              "N/A",
							TpmPcrPreboot:             "N/A",
							TpmPcrUboot:               "N/A",
							TpmPcrLinux:               "N/A",
							CPLDVersion:               "N/A",
							CPLDVersionTimestamp:      -1606322589,
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

				prefix := bmcmonitoring.PrefixSecurity

				// 				CurrentPublicKeyHashIndex:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricCurrentPublicKeyHashIndex),
				})
				// CurrentSVNVersion:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricCurrentSVNversion),
				})
				// Key0Revocation:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricKey0Revocation),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricKey0Revocation: "Not Revoked",
					},
				})
				// Key1Revocation:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricKey1Revocation),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricKey1Revocation: "Revoked",
					},
				})
				// Key2Revocation:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricKey2Revocation),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricKey2Revocation: "Not Revoked",
					},
				})
				// Key3Revocation:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricKey3Revocation),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricKey3Revocation: "Revoked",
					},
				})
				// Key4Revocation:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricKey4Revocation),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricKey4Revocation: "Not Revoked",
					},
				})
				// MinimalSVNIndex:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricMinimalSVNindex),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricMinimalSVNindex: "min",
					},
				})
				// FWImageSource:
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 1,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricFWImageSource),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricFWImageSource: "Secondary",
					},
				})
				// TpmPcrPpboot

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricTPMPCRPPBOOT),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricTPMPCRPPBOOT: "N/A",
					},
				})

				// TpmPcrPreboot
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricTPMPCRPREBOOT),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricTPMPCRPREBOOT: "N/A",
					},
				})

				// TpmPcrUBoot

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricTPMPCRBOOT),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricTPMPCRBOOT: "N/A",
					},
				})

				// TpmPcrLinux
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricTPMPCRLINUX),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricTPMPCRLINUX: "N/A",
					},
				})

				// CPLDVersion
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricCPLDVersion),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricCPLDVersion: "N/A",
					},
				})

				// CPLDVersionTimestamp
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:    "hostname",
					Oam:         "0",
					MetricValue: 0,
					MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricCPLDVersionTimestamp),
					CustomLabels: map[string]string{
						bmcmonitoring.SecurityMetricCPLDVersionTimestamp: "-1606322589",
					},
				})

				return metrics
			},
			verifyMetrics: func(metrics []bmcmonitoring.Metric) {
				for _, metric := range metrics {
					err := bmcmonitoring.VerifySecurity(metric)
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
				require.Equal(t, "https://hostname/redfish/v1/hl/oam/0/Security", r.URL.String())
				return test.sendResponse()
			}

			ch := make(chan []bmcmonitoring.Metric)
			go func() {
				ch <- client.Security(context.Background(), log)
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
