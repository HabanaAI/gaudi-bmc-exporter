package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) CTemperature(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {
		go func(oam int) {

			ll := log.WithField("oam", oam)
			resp, err := c.ctemperature(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get cTemperature")
				ch <- nil
				return
			}

			ch <- c.ctemperatureMetrics(ll, resp, oam)
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

func (c *Client) ctemperature(ctx context.Context, log *logrus.Entry, oam int) (CTemperatureResp, error) {
	var res CTemperatureResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/Ctemperature", c.Hostname, oam)
	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return CTemperatureResp{}, err
	}

	return res.Response, nil

}

func (c *Client) ctemperatureMetrics(log *logrus.Entry, resp CTemperatureResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixCTemperature

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  prefix,
		MetricValue: resp.Temperature,
	})

	return metrics
}
