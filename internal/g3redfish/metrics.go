package g3redfish

import (
	"context"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) Alerts(ctx context.Context, alertsData bmcmonitoring.AlertsData, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) BmcState(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	metrics := c.getMetrics(ctx, log, "")
	return metrics
}

func (c *Client) CTemperature(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) Direct(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) EthernetInfo(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) EthernetStatusCounters(context.Context, *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) EthernetStatus(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) Frequency(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) Hbm(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) Info(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) LaneInfo(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) PcieInfo(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) Power(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) Security(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) SensorCurrent(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) SensorTemperature(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) SensorVoltageMonitor(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) SensorVoltage(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) Status(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}

func (c *Client) Temperature(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	return metrics
}
