package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) BmcState(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var res BmcStateResponse

	url := fmt.Sprintf("https://%s/redfish/v1/Chassis/1", c.Hostname)
	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		log.WithError(err).Error()
		return nil
	}

	bmcStateVal := 0
	if strings.TrimSpace(strings.ToLower(res.PowerState)) == "on" {
		bmcStateVal = 1
	}

	metric := bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		MetricValue: bmcStateVal,
		MetricName:  bmcmonitoring.PrefixBMCState,
	}

	return []bmcmonitoring.Metric{metric}
}
