package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) LaneInfo(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.laneInfo(ctx, ll, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get lane info")
				ch <- nil
				return
			}

			ch <- c.laneInfoMetrics(ll, resp, oam)
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

func (c *Client) laneInfo(ctx context.Context, log *logrus.Entry, oam int) ([]LaneInfoResp, error) {
	var result []LaneInfoResp

	for lane := 0; lane < c.Lanes; lane++ {
		var res LaneInfoResponse
		url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/LaneInfo/%d", c.Hostname, oam, lane)
		err := c.request(ctx, log, url, &res, nil)
		if err != nil {
			log.WithField("lane", lane).WithError(err).Error()
			continue
		}

		// add the lane number
		res.Response.LaneNum = lane
		result = append(result, res.Response)
	}

	return result, nil

}

func (c *Client) laneInfoMetrics(log *logrus.Entry, resps []LaneInfoResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixLaneInfo

	for _, resp := range resps {
		// EBUF Overflow
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricEBUFOverflow),
			MetricValue: resp.EBUFOverflow,
			CustomLabels: map[string]string{
				"lane": fmt.Sprintf("%d", resp.LaneNum),
			},
		})
		// EBUF Under Run
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricEBUFUnderRun),
			MetricValue: resp.EBUFUnderRun,
			CustomLabels: map[string]string{
				"lane": fmt.Sprintf("%d", resp.LaneNum),
			},
		})
		// Decode Error
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricDecodeError),
			MetricValue: resp.DecodeError,
			CustomLabels: map[string]string{
				"lane": fmt.Sprintf("%d", resp.LaneNum),
			},
		})
		// Running Display Error
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricRunningDisplayError),
			MetricValue: resp.RunningDisplayError,
			CustomLabels: map[string]string{
				"lane": fmt.Sprintf("%d", resp.LaneNum),
			},
		})
		// SKPOS Parity Error
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricSKPOSParityError),
			MetricValue: resp.SKPOSParityError,
			CustomLabels: map[string]string{
				"lane": fmt.Sprintf("%d", resp.LaneNum),
			},
		})
		// SYNC Header Error
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricSYNCHeaderError),
			MetricValue: resp.SYNCHeaderError,
			CustomLabels: map[string]string{
				"lane": fmt.Sprintf("%d", resp.LaneNum),
			},
		})
		// Deskew Error
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.LaneInfoMetricDeskewError),
			MetricValue: resp.DeskewError,
			CustomLabels: map[string]string{
				"lane": fmt.Sprintf("%d", resp.LaneNum),
			},
		})
	}

	return metrics
}
