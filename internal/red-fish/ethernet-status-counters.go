package redfish

import (
	"context"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

// TODO: we only split it in RAS to 2 different metrics because the Ras is slow
func (c *Client) EthernetStatusCounters(context.Context, *logrus.Entry) []bmcmonitoring.Metric {
	return nil
}
