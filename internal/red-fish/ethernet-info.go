package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) EthernetInfo(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.ethernetInfo(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get ethernet info")
				ch <- nil
				return
			}

			ch <- c.ethernetInfoMetrics(ll, resp, oam)
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

func (c *Client) ethernetInfo(ctx context.Context, log *logrus.Entry, oam int) (EthernetInfoResp, error) {
	var res EthernetInfoResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/EthernetInfo", c.Hostname, oam)
	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return EthernetInfoResp{}, err
	}

	return res.Response, nil

}

func (c *Client) ethernetInfoMetrics(log *logrus.Entry, resp EthernetInfoResp, oam int) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric

	prefix := rasmonitoring.PrefixEthernetInfo

	// SerdesAvailability

	serdesAvailabilityVal := -1

	switch strings.TrimSpace(strings.ToLower(resp.SerdesAvailability)) {
	case "unavailable":
		serdesAvailabilityVal = 0
	case "available":
		serdesAvailabilityVal = 1
	default:
		log.WithError(fmt.Errorf("unexpected serdes availability value %s", resp.SerdesAvailability)).Error()
	}

	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: serdesAvailabilityVal,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricSerDesAvailability),
		CustomLabels: map[string]string{
			"state": resp.SerdesAvailability,
		},
	})

	// PortMaxSpeed
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.PortMaxSpeed,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricPortMaxSpeed),
	})

	// ANLTStatus

	for port, portState := range resp.ANLTStatus {

		anltStatusVal := -1
		switch strings.TrimSpace(strings.ToLower(portState)) {
		case "disabled":
			anltStatusVal = 0
		case "enabled":
			anltStatusVal = 1
		default:
			log.WithField("port", port).WithError(fmt.Errorf("unexpected ANLT status %s", portState)).Error()

		}

		metrics = append(metrics, rasmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: anltStatusVal,
			MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricANLTStatus),
			CustomLabels: map[string]string{
				"state": portState,
				"port":  fmt.Sprintf("%d", port),
			},
		})
	}

	// NumberOfLanes
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NumberOfLanes,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricNumberOfLanes),
	})

	// NumberOfLinks
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NumberOfLinks,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricNumberOfLinks),
	})

	// LinkSpeed
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.LinkSpeed,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.EthernetInfoMetricLinkSpeed),
	})

	// port map
	for portNum, portType := range resp.PortMap {
		// Port Mapping
		portMapVal := -1
		switch strings.TrimSpace(strings.ToLower(portType)) {
		case "internal port":
			portMapVal = 0
		case "external port":
			portMapVal = 1
		default:
			log.WithField("port", portNum).WithError(fmt.Errorf("unexpected port map value %s", portType)).Error()
		}

		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: portMapVal,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricPortMapping),
			CustomLabels: map[string]string{
				"type": portType,
				"port": fmt.Sprintf("%d", portNum),
			},
		})
	}

	// Toggle count
	for portNum, toggleCount := range resp.ToggleCount {
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: toggleCount,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricStateTogglingCounter),
			CustomLabels: map[string]string{
				"port": fmt.Sprintf("%d", portNum),
			},
		})
	}

	// ExternalLinkStatus
	for port, linkStatus := range resp.ExternalLinkStatus {
		externalLinkStatusVal := -1
		switch strings.TrimSpace(strings.ToLower(linkStatus)) {
		case "non-active port", "non-activeport":
			externalLinkStatusVal = 0
		case "active port":
			externalLinkStatusVal = 1
		default:
			log.WithField("port", port).WithError(fmt.Errorf("unexpected external link status %s", linkStatus)).Error()
		}

		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: externalLinkStatusVal,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricExternalLinkStatus),
			CustomLabels: map[string]string{
				"state": linkStatus,
				"link":  fmt.Sprintf("%d", port),
			},
		})
	}

	// LinkStatus
	for portNum, linkStatus := range resp.LinkStatus {
		linkStatusVal := -1
		switch strings.TrimSpace(strings.ToLower(linkStatus)) {
		case "not connected":
			linkStatusVal = 0
		case "connected":
			linkStatusVal = 1
		default:
			log.WithField("port", portNum).WithError(fmt.Errorf("unexpected link status %s", linkStatus)).Error()
		}
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: linkStatusVal,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricLinkStatus),
			CustomLabels: map[string]string{
				"state": linkStatus,
				"link":  fmt.Sprintf("%d", portNum),
			},
		})
	}

	// PHYStatus
	for portNum, phyStatus := range resp.PHYStatus {
		phyStatusVal := -1
		switch strings.TrimSpace(strings.ToLower(phyStatus)) {
		case "not ready":
			phyStatusVal = 0
		case "ready":
			phyStatusVal = 1
		default:
			log.WithField("port", portNum).WithError(fmt.Errorf("unexpected phy status %s", phyStatus)).Error()
		}

		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: phyStatusVal,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.EthernetInfoMetricPHYStatus),
			CustomLabels: map[string]string{
				"state": phyStatus,
				"phy":   fmt.Sprintf("%d", portNum),
			},
		})
	}

	return metrics
}
