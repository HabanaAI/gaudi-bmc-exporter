package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) Hbm(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {
		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.hbm(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get hbm information")
				ch <- nil
				return
			}
			ch <- c.hbmMetrics(ll, resp, oam)

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

func (c *Client) hbm(ctx context.Context, log *logrus.Entry, oam int) (HbmResp, error) {
	var res HbmResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/HBM", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return HbmResp{}, err
	}

	return res.Response, nil

}

func (c *Client) hbmMetrics(log *logrus.Entry, resp HbmResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixHbm

	// EccErrors
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.EccErrors,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricEccErrors),
	})

	// RepairedLanes

	metrics = append(metrics, bmcmonitoring.Metric{
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricNumOfRepairedLanes),
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: len(resp.RepairedLanes),
	})

	for _, lane := range resp.RepairedLanes {
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: 0,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricRepairedLanes),
			CustomLabels: map[string]string{
				bmcmonitoring.HBMMetricRepairedLanesLabelHBMIndex:  fmt.Sprintf("%d", lane.HBMIndex),
				bmcmonitoring.HBMMetricRepairedLanesLabelMCChannel: fmt.Sprintf("%d", lane.MCChannel),
			},
		})
	}

	// ReplacedRow
	metrics = append(metrics, bmcmonitoring.Metric{
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricNumOfReplacedRows),
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: len(resp.ReplacedRows),
	})

	for _, row := range resp.ReplacedRows {
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: 0,
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricReplaceRows),
			CustomLabels: map[string]string{
				bmcmonitoring.HBMMetricReplaceRowsLabelHBMIndex:   fmt.Sprintf("%d", row.HBMIndex),
				bmcmonitoring.HBMMetricReplaceRowsLabelPCIndex:    fmt.Sprintf("%d", row.PCIndex),
				bmcmonitoring.HBMMetricReplaceRowsLabelStackID:    fmt.Sprintf("%d", row.StackID),
				bmcmonitoring.HBMMetricReplaceRowsLabelBankIndex:  fmt.Sprintf("%d", row.BankIndex),
				bmcmonitoring.HBMMetricReplaceRowsLabelCause:      row.Cause,
				bmcmonitoring.HBMMetricReplaceRowsLabelRowAddress: fmt.Sprintf("%d", row.RowAddress),
			},
		})
	}

	for _, repairStatus := range resp.HBMRepairStatus {
		// MbistRepair

		mbistRepairVal := -1

		switch strings.TrimSpace(strings.ToLower(repairStatus.MBISTRepair)) {
		case "flow did not run":
			mbistRepairVal = 0
		case "flow ran":
			mbistRepairVal = 1
		default:
			log.WithError(fmt.Errorf("unexpected mbist repair %s", repairStatus.MBISTRepair)).Error()
		}
		metrics = append(metrics, bmcmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricMbistRepair),
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: mbistRepairVal,
			CustomLabels: map[string]string{
				bmcmonitoring.HBMMetricMbistRepairLabelState: repairStatus.MBISTRepair,
				bmcmonitoring.HBMMetricMbistRepairLabelIndex: fmt.Sprintf("%d", repairStatus.HBMIndex),
			},
		})

		// Global ECC
		metrics = append(metrics, bmcmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.HBMMetricGlobalECC),
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: repairStatus.GlobalECC,
			CustomLabels: map[string]string{
				bmcmonitoring.HBMMetricGlobalECCLabelIndex: fmt.Sprintf("%d", repairStatus.HBMIndex),
			},
		})
	}

	return metrics
}
