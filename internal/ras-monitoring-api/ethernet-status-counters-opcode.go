package rasmonitoringapi

import (
	"context"
	"encoding/json"
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
)

func (c *Client) EthernetStatusCounters(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	ch := make(chan result)
	var wg sync.WaitGroup
	done := make(chan struct{})
	var metrics []rasmonitoring.Metric

	wg.Add(1)

	go func() {
		defer wg.Done()

		for oam := 0; oam < c.Oams; oam++ {

			wg.Add(1)

			go func(oam int) {

				defer wg.Done()
				for port := 0; port < c.Ports; port++ {

					wg.Add(1)
					go func(oam, port int) {
						defer wg.Done()

						// Get port toggling state
						counterVal, err := c.stateTogglingCounter(ctx, oam, port)
						if err != nil {
							log.WithError(err).Error("failed to get ethernet state toggling counter")
							return
						}

						value, err := strconv.Atoi(counterVal)
						if err != nil {
							log.WithError(err).Error()
							return
						}

						ch <- result{metric: rasmonitoring.Metric{
							Hostname:    c.Hostname,
							Oam:         fmt.Sprintf("%d", oam),
							MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricStateTogglingCounter),
							MetricValue: value,
							CustomLabels: map[string]string{
								"port": fmt.Sprintf("%d", port),
							},
						}}
					}(oam, port)

					info, err := c.getEthernetCountersInfo(ctx, oam, port)
					if err != nil {
						log.WithError(err).Error("failed to get ethernet counters")
						continue
					}

					// Add all the fields.

					// BERCorrectable.

					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricBERCorrectable),
						MetricValue: info.BitErrRateCorrectable,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// BERUncorrectable.
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricBERUncorrectable),
						MetricValue: info.BitErrRateUncorrectable,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// Nack.
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricNack),
						MetricValue: info.NackCounter,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// CRC.
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricCRC),
						MetricValue: info.CrcErrorCounter,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// Retransmission Timeout.
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricRetransmissionTimeout),
						MetricValue: info.RetransmissionTimeoutCounter,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// Link Retraining Due To BER.
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricLinkRetrainingDueToBER),
						MetricValue: info.LinkRetrainingBarCounter,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// MACRemote.
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricMACRemote),
						MetricValue: info.MacRemoteErrCounter,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// Retransmission.
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricRetransmission),
						MetricValue: info.RetransmissionCounter,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// Retraining.
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricRetraining),
						MetricValue: info.RetrainingCounter,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// SERPreFEC
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricSERPreFEC),
						MetricValue: info.SerPreFec,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// SERPostFEC
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricSERPostFEC),
						MetricValue: info.SerPostFec,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// Latency
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricLatency),
						MetricValue: info.Latency,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}

					// Throughput.
					ch <- result{metric: rasmonitoring.Metric{
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixEthernetStatus, rasmonitoring.EthernetStatusMetricThroughput),
						MetricValue: info.Throughput,
						CustomLabels: map[string]string{
							"port": fmt.Sprintf("%d", port),
						},
					}}
				}
			}(oam)

		}
	}()

	go func() {
		wg.Wait()
		done <- struct{}{}
		close(ch)
	}()

	for {
		select {
		case <-done:
			return metrics
		case data := <-ch:
			if data.err != nil {
				log.WithField("metric", rasmonitoring.PrefixEthernetStatus).WithError(data.err).Error()
				continue
			}
			metrics = append(metrics, data.metric)
		}
	}
}

// stateTogglingCounter will return the counter of the port toggling.
func (c *Client) stateTogglingCounter(ctx context.Context, oam, port int) (string, error) {

	// write so we can read the result.
	err := c.write(ctx, oam, Opcode{
		OpcodeNumber: 5,
		Offset:       234,
		Length:       4,
	}, port)

	if err != nil {
		return "", err
	}

	// Read the results.
	counterVal, err := c.decodeOpcode(ctx, Opcode{
		OpcodeNumber: 5,
		Offset:       238,
		Length:       4,
		ExpectedType: Int,
	}, oam, rasmonitoring.EthernetStatusMetricStateTogglingCounter, methodIndirect)

	if err != nil {
		return "", err
	}

	return counterVal, nil
}

type ethernetCountersInfoResponse struct {
	BitErrRateCorrectable        int `json:"bit_err_rate_correctable"`
	BitErrRateUncorrectable      int `json:"bit_err_rate_uncorrectable"`
	NackCounter                  int `json:"nack_counter"`
	RetransmissionTimeoutCounter int `json:"retransmission_timeout_counter"`
	LinkRetrainingBarCounter     int `json:"link_retrianing_bar_counter"`
	MacRemoteErrCounter          int `json:"mac_remote_err_counter"`
	RetransmissionCounter        int `json:"retransmission_counter"`
	RetrainingCounter            int `json:"retraining_counter"`
	CrcErrorCounter              int `json:"crc_error_counter"`
	SerPreFec                    int `json:"ser_pre_fec_value"`
	SerPostFec                   int `json:"ser_post_fec_value"`
	Latency                      int `json:"latency_counter _value"`
	Throughput                   int `json:"throughput_counter_value"`
}

// getEthernetCountersInfo will return information about the ethernet counters.
func (c *Client) getEthernetCountersInfo(ctx context.Context, oam, port int) (ethernetCountersInfoResponse, error) {

	u, err := url.Parse(fmt.Sprintf("https://%s/ext/ras/indirect/eth_counters", c.Hostname))
	if err != nil {
		return ethernetCountersInfoResponse{}, err
	}

	q := u.Query()
	q.Add("oam", fmt.Sprintf("%d", oam))
	q.Add("port_idx", fmt.Sprintf("%d", port))

	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return ethernetCountersInfoResponse{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return ethernetCountersInfoResponse{}, err
	}

	var res ethernetCountersInfoResponse

	err = json.Unmarshal(body, &res)
	if err != nil {
		return ethernetCountersInfoResponse{}, err
	}

	return res, nil
}
