package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) Info(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.info(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get info")
				ch <- nil
				return
			}

			ch <- c.infoMetrics(ll, resp, oam)
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

func (c *Client) info(ctx context.Context, log *logrus.Entry, oam int) (InfoResp, error) {
	var res InfoResponse

	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/Info", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return InfoResp{}, err
	}

	return res.Response, nil

}

func (c *Client) infoMetrics(log *logrus.Entry, resp InfoResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixInfo

	// Device ID.
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:   c.Hostname,
		Oam:        fmt.Sprintf("%d", oam),
		MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricDeviceID),
		CustomLabels: map[string]string{
			bmcmonitoring.InfoMetricDeviceID: resp.DeviceID,
		},
	})

	// SubSystem Device ID.
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:   c.Hostname,
		Oam:        fmt.Sprintf("%d", oam),
		MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricSubsystemDeviceID),
		CustomLabels: map[string]string{
			bmcmonitoring.InfoMetricSubsystemDeviceID: resp.SubSystemDeviceID,
		},
	})

	// SubSystem Vendor ID.
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:   c.Hostname,
		Oam:        fmt.Sprintf("%d", oam),
		MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricSubsystemVendorID),
		CustomLabels: map[string]string{
			bmcmonitoring.InfoMetricSubsystemVendorID: resp.SubSystemVendorID,
		},
	})

	// ASIC Serial Number
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:   c.Hostname,
		Oam:        fmt.Sprintf("%d", oam),
		MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricASICSerialNumber),
		CustomLabels: map[string]string{
			bmcmonitoring.InfoMetricASICSerialNumber: resp.ASICSerialNumber,
		},
	})

	// Board Serial Number
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:   c.Hostname,
		Oam:        fmt.Sprintf("%d", oam),
		MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricBoardSerialNumber),
		CustomLabels: map[string]string{
			bmcmonitoring.InfoMetricBoardSerialNumber: resp.BoardSerialNumber,
		},
	})

	// SRAMSize
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SRAMSize,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricSRAMSize),
	})

	// HBMSize
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.HBMSize,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricHBMSize),
	})

	// UUID
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:   c.Hostname,
		Oam:        fmt.Sprintf("%d", oam),
		MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.InfoMetricUUID),
		CustomLabels: map[string]string{
			bmcmonitoring.InfoMetricUUID: resp.UUID,
		},
	})
	return metrics
}
