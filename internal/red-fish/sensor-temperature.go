package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) SensorTemperature(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.sensorTemperature(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Error("failed to get sensor temperature")
				ch <- nil
				return
			}

			ch <- c.sensorTemperatureMetrics(ll, resp, oam)

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

func (c *Client) sensorTemperature(ctx context.Context, log *logrus.Entry, oam int) (SensorTemperatureResp, error) {
	var res SensorTemperatureResponse

	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/SensorsTemperature", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return SensorTemperatureResp{}, err
	}

	return res.Response, nil

}

func (c *Client) sensorTemperatureMetrics(log *logrus.Entry, resp SensorTemperatureResp, oam int) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric

	prefix := rasmonitoring.PrefixSensorTemperature

	// OnDie0
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.OnDie0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnDie0),
	})

	// OnDie1
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.OnDie1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnDie1),
	})

	// OnDie2
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.OnDie2,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnDie2),
	})

	// OnDie3
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.OnDie3,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnDie3),
	})

	// HBM0
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.HBM0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM0),
	})

	// HBM1
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.HBM1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM1),
	})

	// HBM2
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.HBM2,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM2),
	})

	// HBM3
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.HBM3,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM3),
	})

	// HBM4
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.HBM4,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM4),
	})

	// HBM5
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.HBM5,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricHBM5),
	})

	// CPLDLocal
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CPLDLocal,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLDLocal),
	})

	// CPLD0
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CPLD0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLD0),
	})

	// CPLD1
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CPLD1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLD1),
	})

	// CPLD2
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CPLD2,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLD2),
	})

	// CPLD3
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CPLD3,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLD3),
	})

	// OnBoard0
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.OnBoard0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnboard0),
	})

	// OnBoard1
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.OnBoard1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnboard1),
	})

	// OnBoard2
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.OnBoard2,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnboard2),
	})

	// OnBoard3
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.OnBoard3,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricOnboard3),
	})

	// CPLDTemp
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CPLDTemp,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricCPLDTemp),
	})

	// PSUStage1
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.PSUStage1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricPSUStage1),
	})

	// PSUStage2
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.PSUStage2,
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.SensorTemperatureMetricPSUStage2),
	})

	return metrics
}
