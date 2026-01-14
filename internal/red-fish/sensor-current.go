package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) SensorCurrent(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.sensorCurrent(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Error("failed to get sensor current")
				ch <- nil
				return
			}

			ch <- c.sensorCurrentMetrics(ll, resp, oam)
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

func (c *Client) sensorCurrent(ctx context.Context, log *logrus.Entry, oam int) (SensorCurrentResp, error) {
	var res SensorCurrentResponse

	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/SensorsCurrent", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return SensorCurrentResp{}, err
	}

	return res.Response, nil

}

func (c *Client) sensorCurrentMetrics(log *logrus.Entry, resp SensorCurrentResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixSensorCurrent

	// VIn54
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.VIn54,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorCurrentVin54),
	})

	// P1VIn12
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.P1VIn12,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorCurrentP1Vin12),
	})

	// Stage154VIn
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Stage154VIn,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorCurrentStage154Vin),
	})

	// Stage113P5VOut
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Stage113P5VOut,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorCurrentStage113P5VOut),
	})

	// Stage213P5VIn
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Stage213P5VIn,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorCurrentStage213P5Vin),
	})

	// Stage2CoreOut
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Stage2CoreOut,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorCurrentStage2CoreOut),
	})

	// Stage2HBMOut
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Stage2HBMOut,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorCurrentStage2HBMout),
	})
	return metrics
}
