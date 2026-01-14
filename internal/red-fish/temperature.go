package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) Temperature(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.temperature(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Error("failed to get temperature")
				ch <- nil
				return
			}

			ch <- c.temperatureMetrics(ll, resp, oam)
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

func (c *Client) temperature(ctx context.Context, log *logrus.Entry, oam int) (TemperatureResp, error) {
	var res TemperatureResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/Temperature", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return TemperatureResp{}, err
	}

	return res.Response, nil

}

func (c *Client) temperatureMetrics(log *logrus.Entry, resp TemperatureResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixTemperature

	// Current Board Temperature
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentBoardTemp),
		MetricValue: resp.CurrentBoardTemperature,
	})

	// Current VRM Temperature
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentVRMTemp),
		MetricValue: resp.CurrentVRMTemperature,
	})

	// CurrentDRAMTemperature
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentDRAMTemp),
		MetricValue: resp.CurrentDRAMTemperature,
	})

	// Current Ondie Temperature
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentOnDieTemp),
		MetricValue: resp.CurrentOndieTemperature,
	})

	// HistoricalBoardTemperature
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricHistoricalBoardTemp),
		MetricValue: resp.HistoricalBoardTemperature,
	})

	// HistoricalVRMTemperature
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricHistoricalVRMTemp),
		MetricValue: resp.HistoricalVRMTemperature,
	})

	// HistoricalDRAMTemperature
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricHistoricalDRAMTemp),
		MetricValue: resp.HistoricalDRAMTemperature,
	})

	// HistoricalOndieTemperature
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricHistoricalOnDieTemp),
		MetricValue: resp.HistoricalOndieTemperature,
	})

	// MaxTemperatureRiseTime
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricMaxTempRiseTime),
		MetricValue: resp.MaxTemperatureRiseTime,
	})

	// MaxSOCTemperatureErrorThreshold
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricMaxSocTempErrorThreshold),
		MetricValue: resp.MaxSOCTemperatureErrorThreshold,
	})

	// MaxSOCTemperatureWarningThreshold
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricMaxSocTempWarmingThreshold),
		MetricValue: resp.MaxSOCTemperatureWarningThreshold,
	})

	// MaxHBMTemperatureThreshold
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricMaxHbmTempThreshold),
		MetricValue: resp.MaxHBMTemperatureThreshold,
	})

	// CurrentSOCTemperatureErrorThreshold
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentSocTempErrorThreshold),
		MetricValue: resp.CurrentSOCTemperatureErrorThreshold,
	})

	// CurrentSOCTemperatureWarningThreshold
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentSocTempWarningThreshold),
		MetricValue: resp.CurrentSOCTemperatureWarningThreshold,
	})

	// CurrentHBMTemperatureThreshold
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.TemperatureMetricCurrentHbmTempThreshold),
		MetricValue: resp.CurrentHBMTemperatureThreshold,
	})

	return metrics
}
