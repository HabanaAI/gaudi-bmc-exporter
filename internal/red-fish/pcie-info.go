package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) PcieInfo(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.pcieInfo(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get pcie info")
				ch <- nil
				return
			}

			ch <- c.pcieInfoMetrics(ll, resp, oam)
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

func (c *Client) pcieInfo(ctx context.Context, log *logrus.Entry, oam int) (PcieInfoResp, error) {
	var res PcieInfoResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/PCIeInfo", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return PcieInfoResp{}, err
	}

	return res.Response, nil

}

func (c *Client) pcieInfoMetrics(log *logrus.Entry, resp PcieInfoResp, oam int) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric

	prefix := rasmonitoring.PrefixPcieInfo

	// MaxPCIeLinkSpeed

	// max pcie link speed values: Gen1-Gen6
	maxPCIeLinkSpeedVal := -1
	switch strings.TrimSpace(strings.ToLower(resp.MaxPCIeLinkSpeed)) {
	case "gen1":
		maxPCIeLinkSpeedVal = 1
	case "gen2":
		maxPCIeLinkSpeedVal = 2
	case "gen3":
		maxPCIeLinkSpeedVal = 3
	case "gen4":
		maxPCIeLinkSpeedVal = 4
	case "gen5":
		maxPCIeLinkSpeedVal = 5
	case "gen6":
		maxPCIeLinkSpeedVal = 6
	default:
		log.WithError(fmt.Errorf("unknown max pcie link speed %s", resp.MaxPCIeLinkSpeed)).Error()
	}

	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricMaxPCIeLinkSpeed),
		MetricValue: maxPCIeLinkSpeedVal,
		CustomLabels: map[string]string{
			rasmonitoring.PcieInfoMetricMaxPCIeLinkSpeed: resp.MaxPCIeLinkSpeed,
		},
	})

	// CurrentPCIeLinkSpeed

	currentPCIeLinkSpeedVal := -1

	switch strings.TrimSpace(strings.ToLower(resp.CurrentPCIeLinkSpeed)) {
	case "gen1":
		currentPCIeLinkSpeedVal = 1
	case "gen2":
		currentPCIeLinkSpeedVal = 2
	case "gen3":
		currentPCIeLinkSpeedVal = 3
	case "gen4":
		currentPCIeLinkSpeedVal = 4
	case "gen5":
		currentPCIeLinkSpeedVal = 5
	case "gen6":
		currentPCIeLinkSpeedVal = 6
	default:
		log.WithError(fmt.Errorf("unknown current pcie link speed %s", resp.CurrentPCIeLinkSpeed)).Error()
	}
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricCurrentPCIeLinkSpeed),
		MetricValue: currentPCIeLinkSpeedVal,
		CustomLabels: map[string]string{
			rasmonitoring.PcieInfoMetricCurrentPCIeLinkSpeed: resp.CurrentPCIeLinkSpeed,
		},
	})

	// MaxPCIeLinkWidth
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricMaxPCIeLinkWidth),
		MetricValue: resp.MaxPCIeLinkWidth,
	})

	// CurrentPCIeLinkWidth
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricCurrentPCIeLinkWidth),
		MetricValue: resp.CurrentPCIeLinkWidth,
	})

	// PCIeDeviceID
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeDeviceID),
		MetricValue: resp.PCIeDeviceID,
	})

	// PCIeSubsystemID
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeSubsystemID),
		MetricValue: resp.PCIeSubsystemID,
	})

	// PCIeSubsystemVendorID
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeSubsystemVendorID),
		MetricValue: 0,
		CustomLabels: map[string]string{
			rasmonitoring.PcieInfoMetricPCIeSubsystemVendorID: resp.PCIeSubsystemVendorID,
		},
	})

	// PCIeBusAndDevice
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeBusAndDevice),
		MetricValue: 0,
		CustomLabels: map[string]string{
			rasmonitoring.PcieInfoMetricPCIeBusAndDevice: resp.PCIeBusAndDevice,
		},
	})

	// CorrectedInternalErrorStatus
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricCorrectedInternalErrorStatus),
		MetricValue: resp.CorrectedInternalErrorStatus,
	})

	// ReplayBufferNumRolloverError
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricReplayBufferNumRolloverError),
		MetricValue: resp.ReplayBufferNumRolloverError,
	})

	// ReplayTimerTimeoutError
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricReplayTimerTimeoutError),
		MetricValue: resp.ReplayTimerTimeoutError,
	})

	// BadTLPCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricBadTLPCounter),
		MetricValue: resp.BadTLPCounter,
	})

	// BadDLLPCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricBadDLLPCounter),
		MetricValue: resp.BadDLLPCounter,
	})

	// ReceiverErrorCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricReceiverErrorCounter),
		MetricValue: resp.ReceiverErrorCounter,
	})

	// LCRCErrorCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricLCRCErrorCounter),
		MetricValue: resp.LCRCErrorCounter,
	})

	// ECRCErrorCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricECRCErrorCounter),
		MetricValue: resp.ECRCErrorCounter,
	})

	// CompletionTimeoutIndication
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricCompletionTimeoutIndication),
		MetricValue: resp.CompletionTimeoutIndication,
	})

	// UncorrectableInternalErrorIndication
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricUncorrectableInternalErrorIndication),
		MetricValue: resp.UncorrectableInternalErrorIndication,
	})

	// ReceiverOverflowIndication
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricReceiverOverflowIndication),
		MetricValue: resp.ReceiverOverflowIndication,
	})

	// FlowControlProtocolErrorIndication
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricFlowControlProtocolErrorIndication),
		MetricValue: resp.FlowControlProtocolErrorIndication,
	})

	// SurpriseLinkDownIndication
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricSurpriseLinkDownIndication),
		MetricValue: resp.SurpriseLinkDownIndication,
	})

	// MalfunctionTLPErrorIndication
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricMalfunctionTLPErrorIndication),
		MetricValue: resp.MalfunctionTLPErrorIndication,
	})

	// DLLPProtocolErrorIndication
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricDLLPProtocolErrorIndication),
		MetricValue: resp.DLLPProtocolErrorIndication,
	})

	// RxNakDLLPCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricRXNakDLLPCounter),
		MetricValue: resp.RxNakDLLPCounter,
	})

	// TxNakDLLPCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricTxNakDLLPCounter),
		MetricValue: resp.TxNakDLLPCounter,
	})

	// RetryTLPCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricRetryTLPcounter),
		MetricValue: resp.RetryTLPCounter,
	})

	// PWRBRKIndication
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPWRBRKindication),
		MetricValue: resp.PWRBRKIndication,
	})

	// PCIeRxMemoryWriteCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeRXMemoryWriteCounter),
		MetricValue: resp.PCIeRxMemoryWriteCounter,
	})

	// PCIeRxMemoryReadCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeRXMemoryReadCounter),
		MetricValue: resp.PCIeRxMemoryReadCounter,
	})

	// PCIeTxMemoryWriteCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeTXMemoryWriteCounter),
		MetricValue: resp.PCIeTxMemoryWriteCounter,
	})

	// PCIeTxMemoryReadCounter
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeTXMemoryReadCounter),
		MetricValue: resp.PCIeTxMemoryReadCounter,
	})

	// AERCapabilityControlOffset
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricAERCapabilityControlOffset),
		MetricValue: resp.AERCapabilityControlOffset,
	})

	// AERErrorLog
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricAERerrorlog),
		MetricValue: resp.AERErrorLog,
	})

	// PCIeFWVersion
	metrics = append(metrics, rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.PcieInfoMetricPCIeFWversion),
		MetricValue: resp.PCIeFWVersion,
	})

	return metrics
}
