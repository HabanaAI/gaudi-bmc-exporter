package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) Direct(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)
			resp, err := c.direct(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get direct information")
				ch <- nil
				return
			}

			ch <- c.directMetrics(ll, resp, oam)
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

// direct will do the request to red fish to get the response.
func (c *Client) direct(ctx context.Context, log *logrus.Entry, oam int) (DirectResp, error) {

	var res DirectResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/Direct", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return DirectResp{}, err
	}

	return res.Response, nil
}

// directMetrics will format the direct response into the expected metrics
func (c *Client) directMetrics(log *logrus.Entry, directResp DirectResp, oam int) []bmcmonitoring.Metric {

	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixDirect

	// PCIeVendorID
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:   c.Hostname,
		Oam:        fmt.Sprintf("%d", oam),
		MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricPCIeVendorID),
		CustomLabels: map[string]string{
			bmcmonitoring.DirectMetricPCIeVendorID: directResp.PCIeVendorID,
		},
	})

	// AsicSerialNumber
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:   c.Hostname,
		Oam:        fmt.Sprintf("%d", oam),
		MetricName: fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricASICSerialNumber),
		CustomLabels: map[string]string{
			bmcmonitoring.DirectMetricASICSerialNumber: directResp.AsicSerialNumber,
		},
	})

	// FWVersionMajor
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricFWVersionMajor),
		MetricValue: directResp.FWVersionMajor,
	})

	// FWVersionMinor
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricFWVersionMinor),
		MetricValue: directResp.FWVersionMinor,
	})

	// FWVersionPatch
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricFWVersionPatch),
		MetricValue: directResp.FWVersionPatch,
	})

	// CoreVDD
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricCoreVDD),
		MetricValue: directResp.CoreVDD,
	})

	// HBMVddq
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricHBMVDDq),
		MetricValue: directResp.HBMVddq,
	})

	// v12
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricV12),
		MetricValue: directResp.V12,
	})

	// APIversion
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricApiVersion),
		MetricValue: directResp.APIVersion,
	})

	// OSStage
	osStageValue := -1

	switch strings.TrimSpace(strings.ToLower(directResp.OSStage)) {
	case "linux":
		osStageValue = 0
	case "uboot":
		osStageValue = 1
	case "preboot":
		osStageValue = 2
	case "zephyr":
		osStageValue = 3
	default:
		log.WithError(fmt.Errorf("unexpected os stage %s", directResp.OSStage)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricOSStage),
		MetricValue: osStageValue,
		CustomLabels: map[string]string{
			"stage": directResp.OSStage,
		},
	})

	// IBAccessState
	iBAccessStateValue := -1

	switch strings.TrimSpace(strings.ToLower(directResp.IBAccessState)) {
	case "full access":
		iBAccessStateValue = 0
	case "restricted access":
		iBAccessStateValue = 1
	default:
		log.WithError(fmt.Errorf("unexpected IB Access State %s", directResp.IBAccessState)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricIBAccessState),
		MetricValue: iBAccessStateValue,
		CustomLabels: map[string]string{
			"state": directResp.IBAccessState,
		},
	})

	// OOBAccessState
	oOBAccessStateValue := -1

	switch strings.TrimSpace(strings.ToLower(directResp.OOBAccessState)) {
	case "full access":
		oOBAccessStateValue = 0
	case "restricted access":
		oOBAccessStateValue = 1
	default:
		log.WithError(fmt.Errorf("unexpected OOB Access State %s", directResp.OOBAccessState)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.DirectMetricOOBAccessState),
		MetricValue: oOBAccessStateValue,
		CustomLabels: map[string]string{
			"state": directResp.OOBAccessState,
		},
	})

	return metrics
}
