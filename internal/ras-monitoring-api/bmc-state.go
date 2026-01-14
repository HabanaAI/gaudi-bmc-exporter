package rasmonitoringapi

import (
	"context"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) BmcState(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {

	var metric rasmonitoring.Metric
	up, err := c.ipmiClient.IsUp()
	if err != nil {
		log.WithError(err).Error("failed getting bmc state")
		return nil
	}

	val := 0
	if up {
		val = 1
	}
	metric = rasmonitoring.Metric{
		Hostname:    c.Hostname,
		MetricValue: val,
		MetricName:  rasmonitoring.PrefixBMCState,
	}
	return []rasmonitoring.Metric{metric}
}
