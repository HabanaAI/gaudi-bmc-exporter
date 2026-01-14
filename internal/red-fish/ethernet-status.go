package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) EthernetStatus(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.ethernetStatus(ctx, ll, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get ethernet status")
				ch <- nil
				return
			}

			ch <- c.ethernetStatusMetrics(ll, resp, oam)
		}(oam)

	}

	for oam := 0; oam < c.Oams; oam++ {
		alertMetrics := <-ch
		if alertMetrics != nil {
			metrics = append(metrics, alertMetrics...)
		}
	}
	return metrics
}

func (c *Client) ethernetStatus(ctx context.Context, log *logrus.Entry, oam int) ([]EthernetStatusResp, error) {

	var result []EthernetStatusResp
	for port := 0; port < c.Ports; port++ {
		var res EthernetStatusResponse
		url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/EthernetStatus/port/%d", c.Hostname, oam, port)
		err := c.request(ctx, log, url, &res, nil)
		if err != nil {
			log.WithError(err).Errorf("failed to get information on port %d", port)
			continue
		}

		// add the port number
		res.Response.PortNum = port

		result = append(result, res.Response)
	}

	return result, nil

}

func (c *Client) ethernetStatusMetrics(log *logrus.Entry, resps []EthernetStatusResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixEthernetStatus

	for _, resp := range resps {
		// Port Mapping
		portMapVal := -1
		switch strings.TrimSpace(strings.ToLower(resp.PortMap)) {
		case "internal port":
			portMapVal = 0
		case "external port":
			portMapVal = 1
		default:
			log.WithField("port", resp.PortNum).WithError(fmt.Errorf("unexpected port map value %s", resp.PortMap)).Error()
		}

		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: portMapVal,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricPortMapping),
			CustomLabels: map[string]string{
				"type": resp.PortMap,
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})

		// Toggle Count

		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.ToggleCount,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricStateTogglingCounter),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})

		// ExternalLinkStatus
		externalLinkStatusVal := -1
		switch strings.TrimSpace(strings.ToLower(resp.ExternalLinkStatus)) {
		case "non-active port", "non-activeport":
			externalLinkStatusVal = 0
		case "active port":
			externalLinkStatusVal = 1
		default:
			log.WithField("port", resp.PortNum).WithError(fmt.Errorf("unexpected external link status %s", resp.ExternalLinkStatus)).Error()
		}
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: externalLinkStatusVal,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricExternalLinkStatus),
			CustomLabels: map[string]string{
				"state": resp.ExternalLinkStatus,
				"link":  fmt.Sprintf("%d", resp.PortNum),
			},
		})

		// LinkStatus
		linkStatusVal := -1
		switch strings.TrimSpace(strings.ToLower(resp.LinkStatus)) {
		case "not connected":
			linkStatusVal = 0
		case "connected":
			linkStatusVal = 1
		default:
			log.WithField("port", resp.PortNum).WithError(fmt.Errorf("unexpected link status %s", resp.LinkStatus)).Error()
		}
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: linkStatusVal,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricLinkStatus),
			CustomLabels: map[string]string{
				"state": resp.LinkStatus,
				"link":  fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// PHYStatus

		phyStatusVal := -1
		switch strings.TrimSpace(strings.ToLower(resp.PHYStatus)) {
		case "not ready":
			phyStatusVal = 0
		case "ready":
			phyStatusVal = 1
		default:
			log.WithField("port", resp.PortNum).WithError(fmt.Errorf("unexpected phy status %s", resp.PHYStatus)).Error()
		}

		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: phyStatusVal,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricPHYStatus),
			CustomLabels: map[string]string{
				"state": resp.PHYStatus,
				"phy":   fmt.Sprintf("%d", resp.PortNum),
			},
		})

		// BERCorrectable
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.BERCorrectable,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricBERCorrectable),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// BERUncorrectable
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.BERUncorrectable,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricBERUncorrectable),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// Nack
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.Nack,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricNack),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})

		// CRC
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.CRC,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricCRC),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// RetransmissionTimeout
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.RetransmissionTimeout,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricRetransmissionTimeout),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// LinkRetrainingDueToBER
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.LinkRetrainingDueToBER,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricLinkRetrainingDueToBER),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// MACRemote
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.MACRemote,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricMACRemote),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// Retransmission
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.Retransmission,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricRetransmission),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// Retraining
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.Retraining,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricRetraining),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// SERPreFEC
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.SERPreFEC,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricSERPreFEC),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// SERPostFEC
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.SERPostFEC,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricSERPostFEC),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// Latency
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.Latency,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricLatency),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
		// Throughput
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: resp.Throughput,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetStatusMetricThroughput),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", resp.PortNum),
			},
		})
	}

	return metrics
}
