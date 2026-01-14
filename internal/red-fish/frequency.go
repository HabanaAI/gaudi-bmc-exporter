package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) Frequency(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)
			resp, err := c.frequency(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get frequency")
				ch <- nil
				return
			}

			ch <- c.frequencyMetrics(ll, resp, oam)
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

func (c *Client) frequency(ctx context.Context, log *logrus.Entry, oam int) (FrequencyResp, error) {
	var res FrequencyResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/Frequency", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return FrequencyResp{}, err
	}

	return res.Response, nil

}

func (c *Client) frequencyMetrics(log *logrus.Entry, resp FrequencyResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixFrequency

	// HBMFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.HBMFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricHBMFrequency),
	})

	// MaxTPCFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.MaxTPCFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxTPCFrequency),
	})

	// MaxMMEFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.MaxMMEFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxMMEFrequency),
	})

	// MaxDMAFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.MaxDMAFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxDMAFrequency),
	})

	// MaxMediaFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.MaxMediaFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxMediaFrequency),
	})

	// MaxPCIeFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.MaxPCIeFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxPCIeFrequency),
	})

	// MaxARMFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.MaxARMFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxARMFrequency),
	})

	// MaxNICFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.MaxNICFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxNICFrequency),
	})

	// MaxNoCFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.MaxNoCFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricMaxNoCFrequency),
	})

	// CurrentTPCFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentTPCFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentTPCFrequency),
	})

	// CurrentMMEFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentMMEFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentMMEFrequency),
	})

	// CurrentDMAFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentDMAFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentDMAFrequency),
	})

	// CurrentMediaFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentMediaFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentMediaFrequency),
	})

	// CurrentPCIeFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentPCIeFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentPCIeFrequency),
	})

	// CurrentARMFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentARMFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentARMFrequency),
	})

	// CurrentNICFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentNICFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentNICFrequency),
	})

	// CurrentNoCFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentNoCFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentNoCFrequency),
	})

	// CurrentSRAMFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentSRAMFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentSRAMFrequency),
	})

	// CurrentMSSFrequency
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentMSSFrequency,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.FrequencyMetricCurrentMSSFrequency),
	})

	return metrics
}
