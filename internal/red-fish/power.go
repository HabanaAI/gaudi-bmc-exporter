package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) Power(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.power(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get power information")
				ch <- nil
				return
			}

			ch <- c.powerMetrics(ll, resp, oam)
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

func (c *Client) power(ctx context.Context, log *logrus.Entry, oam int) (PowerResp, error) {
	var res PowerResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/Power", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return PowerResp{}, err
	}

	return res.Response, nil

}

func (c *Client) powerMetrics(log *logrus.Entry, resp PowerResp, oam int) []bmcmonitoring.Metric {

	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixPower

	// CurrentPowerConsumption
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentPowerConsumption,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.PowerMetricCurrentPowerConsumption),
	})

	// PeakPowerConsumption

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.PeakPowerConsumption,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.PowerMetricPeakPowerConsumption),
	})

	return metrics
}
