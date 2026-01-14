package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) Alerts(ctx context.Context, alertsData bmcmonitoring.AlertsData, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	ch := make(chan []bmcmonitoring.Metric)
	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)
			resp, err := c.alerts(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get alerts")
				ch <- nil
				return
			}
			ch <- c.alertsMetrics(ll, alertsData, resp, oam)

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

func (c *Client) alerts(ctx context.Context, log *logrus.Entry, oam int) (AlertsResp, error) {
	var res AlertsResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/Alert", c.Hostname, oam)
	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return AlertsResp{}, err
	}

	return res.Response, nil

}

func (c *Client) alertsMetrics(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp AlertsResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	// sram double ecc
	metrics = append(metrics, sramDoubleEccAlerts(log, alertsData, resp.SRAMDoubleECC, c.Hostname, oam)...)

	// sram correctable
	metrics = append(metrics, sramCorrectableAlerts(log, alertsData, resp.SRAMCorrectable, c.Hostname, oam)...)

	// pcie
	metrics = append(metrics, pcieAlerts(log, alertsData, resp.PCIe, c.Hostname, oam)...)

	// hbm
	metrics = append(metrics, hbmAlerts(log, alertsData, resp.HBM, c.Hostname, oam)...)

	// wd cause
	metrics = append(metrics, wdCauseAlerts(log, alertsData, resp.WD, c.Hostname, oam)...)

	// temperature
	metrics = append(metrics, temperatureAlerts(log, alertsData, resp.Temperature, c.Hostname, oam)...)

	// security
	metrics = append(metrics, securityAlerts(log, alertsData, resp.Security, c.Hostname, oam)...)

	// rma
	metrics = append(metrics, rmaAlerts(log, alertsData, resp.RMA, c.Hostname, oam)...)

	// voltage monitor
	metrics = append(metrics, voltageMonitorAlerts(log, alertsData, resp.VoltageMonitor, c.Hostname, oam)...)

	// nic
	metrics = append(metrics, nicAlerts(log, alertsData, resp.NIC, c.Hostname, oam)...)

	// checksum
	metrics = append(metrics, checksumAlerts(log, alertsData, resp.Checksum, c.Hostname, oam)...)

	// pll
	metrics = append(metrics, pllAlerts(log, alertsData, resp.PLL, c.Hostname, oam)...)

	// fit
	metrics = append(metrics, fitAlerts(log, alertsData, resp.FIT, c.Hostname, oam)...)

	return metrics
}

func sramDoubleEccAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *SRAMDoubleECCAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}
	sramAlertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeDoubleECCBlockErr]
	for _, sramAlert := range *resp {

		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", sramAlertData.Description, sramAlert.Description()),
				"error_type":  bmcmonitoring.AlertTypeDoubleECCBlockErr,
				"time":        convertAlertTime(sramAlert.Timestamp),
				"severity":    sramAlertData.Severity,
			},
		})
	}

	return metrics
}

func sramCorrectableAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *SRAMCorrectableAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}

	// SINGLE_ECC_BLOCK_ERR
	for _, singleEcc := range resp.SingleECC {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSingleECCBlockErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, singleEcc.Description()),
				"error_type":  bmcmonitoring.AlertTypeSingleECCBlockErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(singleEcc.Timestamp),
				"severity": alertData.Severity,
			},
		})

	}

	// FATAL_SRAM_ECC_ERR
	for _, multipleEcc := range resp.MultipleECC {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFatalSramEccErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, multipleEcc.Description()),
				"error_type":  bmcmonitoring.AlertTypeFatalSramEccErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(multipleEcc.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	return metrics
}

func pcieAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *PCIeAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}

	// Once we have fatal we send all this alerts
	if strings.ToLower(resp.Fatal) == "true" {

		// PCIE_UNCORRECTABLE_FATAL_ERR
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieUncorrectableFatalErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieUncorrectableFatalErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_RECEIVER_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieReceiverErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieReceiverErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_RECEIVER_OVERFLOW_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieReceiverOverflowErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieReceiverOverflowErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_FLOW_CONTROL_PROTOCOL_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieFlowControlProtocolErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieFlowControlProtocolErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_SURPRISE_LINK_DOWN_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieSurpriseLinkDownErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieSurpriseLinkDownErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_MALFUNCTION_TLP_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieMalfunctionTlpErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieMalfunctionTlpErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_DLLP_PROTOCOL_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieDllpProtocolErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieDllpProtocolErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	if strings.ToLower(resp.NonFatal) == "true" {

		// PCIE_ECRC_ERR
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieEcrcErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieEcrcErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_COMPLETION_TIMEOUT_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieCompletionTimeoutErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieCompletionTimeoutErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}
	if strings.ToLower(resp.Correctable) == "true" {

		// PCIE_CORRECTED_INTERNAL_ERR
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieCorrectedInternalErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieCorrectedInternalErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_REPLAY_BUFFER_ROLLOVER_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieReplayBufferRolloverErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieReplayBufferRolloverErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_REPLAY_TIMER_TIMEOUT_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieReplayTimerTimeoutErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieReplayTimerTimeoutErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_BAD_DLLP_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieBadDllpErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieBadDllpErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_BAD_TLP_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieBadTlpErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieBadTlpErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_RECEIVER_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieReceiverErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieReceiverErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_LCRD_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieLcrdErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieLcrdErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_RX_NAK_DLLP_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieRxNakDllpErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieRxNakDllpErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_TX_NAK_DLLP_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieTxNakDllpErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieTxNakDllpErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})

		// PCIE_RETRY_TLP_ERR
		alertData = alertsData.AlertsInformation[bmcmonitoring.AlertTypePcieRetryTlpErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypePcieRetryTlpErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	return metrics
}
func hbmAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *HBMAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}

	// DOUBLE_ECC_ERR
	for _, derr := range resp.DERR {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeDoubleECCErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, derr.Description()),
				"error_type":  bmcmonitoring.AlertTypeDoubleECCErr,

				"time":     convertAlertTime(derr.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// SINGLE_ECC_ERR
	for _, multipleSerr := range resp.MultiSERR {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSingleECCErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, multipleSerr.Description()),
				"error_type":  bmcmonitoring.AlertTypeSingleECCErr,

				"time":     convertAlertTime(multipleSerr.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// HBM_CATASTROPHIC_TEMP_ERR
	if strings.ToLower(resp.CATTRIP) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeHbmCatastrophicTempErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeHbmCatastrophicTempErr,

				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// FATAL_HBM_PARITY_ERR
	for _, parity := range resp.Parity {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFatalHbmParityErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, parity.Description()),
				"error_type":  bmcmonitoring.AlertTypeFatalHbmParityErr,

				"time":     convertAlertTime(parity.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// FATAL_HBM_ECC_ERR
	for _, serr := range resp.SameAddressMultiSERR {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFatalHbmEccErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, serr.Description()),
				"error_type":  bmcmonitoring.AlertTypeFatalHbmEccErr,

				"time":     convertAlertTime(serr.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	return metrics
}

func wdCauseAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *WDAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	if resp == nil {
		return metrics
	}

	// check which alerts do we have

	// TODO: all those fields are string, i asked that they will be replaced to boolean
	if strings.ToLower(resp.Heartbeat) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeWdHeartbeatErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeWdHeartbeatErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	if strings.ToLower(resp.TDR) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeWdTdrErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeWdTdrErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	if strings.ToLower(resp.HW) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeWdHwErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeWdHwErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	return metrics
}

func temperatureAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *TemperatureAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}

	// THERMAL_THRESHOLD_CROSSING_ERR

	if resp.ThermalTemperatureCrossing != nil && strings.ToLower(resp.ThermalTemperatureCrossing.Crossed) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeThermalThresholdCrossingErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, resp.ThermalTemperatureCrossing.Description()),
				"error_type":  bmcmonitoring.AlertTypeThermalThresholdCrossingErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	//THERMAL_TEMP_RISE_TIME_VIOLATION_ERR
	if strings.ToLower(resp.RiseTimeViolation) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeThermalTempRiseTimeViolationErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeThermalTempRiseTimeViolationErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	return metrics
}

func securityAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *SecurityAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}
	if resp.SPIFailure != nil {

		// SECURITY_SPI_REVOKED_PUB_KEY_ERR
		if strings.ToLower(resp.SPIFailure.RevokedPublicKey) == "true" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSecuritySpiRevokedPubKeyErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeSecuritySpiRevokedPubKeyErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// SECURITY_SPI_SVN_ERR
		if strings.ToLower(resp.SPIFailure.SVN) == "true" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSecuritySpiSvnErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeSecuritySpiSvnErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// SECURITY_SPI_SIGNATURE_ERR
		if strings.ToLower(resp.SPIFailure.SVN) == "true" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSecuritySpiSignatureErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeSecuritySpiSignatureErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}
	}

	if resp.AgentFailure != nil {

		// SECURITY_AGENT_REVOKED_PUB_KEY_ERR
		if strings.ToLower(resp.AgentFailure.RevokedPublicKey) == "true" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSecurityAgentRevokedPubKeyErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeSecurityAgentRevokedPubKeyErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// SECURITY_AGENT_SVN_ERR
		if strings.ToLower(resp.AgentFailure.SVN) == "true" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSecurityAgentSvnErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeSecurityAgentSvnErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// SECURITY_AGENT_SIGNATURE_ERR
		if strings.ToLower(resp.AgentFailure.Signature) == "true" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSecurityAgentSignatureErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeSecurityAgentSignatureErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// SECURITY_AGENT_UBOOT_ERR
		if strings.ToLower(resp.AgentFailure.UBoot) == "true" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSecurityAgentUbootErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeSecurityAgentUbootErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// SECURITY_AGENT_LINUX_ERR
		if strings.ToLower(resp.AgentFailure.Linux) == "true" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSecurityAgentLinuxErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeSecurityAgentLinuxErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// SECURITY_AGENT_ZEPHYR_ERR
		if strings.ToLower(resp.AgentFailure.Zephyr) == "true" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSecurityAgentZephyrErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeSecurityAgentZephyrErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}
	}

	// SECURITY_PERMISSION_ACCESS_ERR
	if strings.ToLower(resp.PermissionAccessError) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSecurityPermissionAccessErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeSecurityPermissionAccessErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}
	return metrics
}

func rmaAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *RMAAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}

	// FATAL_SRAM_ECC_ERR
	for _, fatal := range resp.FatalSRAMDoubleECC {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFatalSramEccErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, fatal.Description()),
				"error_type":  bmcmonitoring.AlertTypeFatalSramEccErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(fatal.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// FATAL_HBM_ECC_ERR
	for _, hbm := range resp.HBMECC {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFatalHbmEccErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, hbm.Description()),
				"error_type":  bmcmonitoring.AlertTypeFatalHbmEccErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(hbm.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// FATAL_HBM_PARITY_ERR
	for _, parity := range resp.HBMParity {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFatalHbmParityErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, parity.Description()),
				"error_type":  bmcmonitoring.AlertTypeFatalHbmParityErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(parity.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// FATAL_ROW_REPLACEMENT_ERR
	if resp.HBMRowReplacement != "" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFatalRowReplacementErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, resp.HBMRowReplacement),
				"error_type":  bmcmonitoring.AlertTypeFatalRowReplacementErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// FATAL_HBM_INIT_ERR
	if resp.HBMInitialization != "No failure" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFatalHbmInitErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, resp.HBMInitialization),
				"error_type":  bmcmonitoring.AlertTypeFatalHbmInitErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	if resp.SRAMRepairFailure != nil {

		// SRAM_UNREPAIRABLE_RINGS_ERR
		for _, nonRepairable := range resp.SRAMRepairFailure.NonRepairable {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSramUnrepairableRingsErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": fmt.Sprintf("%s. %s", alertData.Description, nonRepairable.Description()),
					"error_type":  bmcmonitoring.AlertTypeSramUnrepairableRingsErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// SRAM_REPAIRED_RINGS_ERR
		for _, repairable := range resp.SRAMRepairFailure.Repairable {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeSramRepairedRingsErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": fmt.Sprintf("%s. %s", alertData.Description, repairable.Description()),
					"error_type":  bmcmonitoring.AlertTypeSramRepairedRingsErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}
	}

	if resp.HBMMBISTUnrepairable != nil && resp.HBMMBISTUnrepairable.InFieldFailure != "" {
		// HBM_MBIST_REPAIR_FAILURE_ERR
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeHbmMbistRepairFailureErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": fmt.Sprintf("%s. %s", alertData.Description, resp.HBMMBISTUnrepairable.Description()),
				"error_type":  bmcmonitoring.AlertTypeHbmMbistRepairFailureErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	return metrics
}

func voltageMonitorAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *VoltageMonitorAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}

	// VOLTAGE_MONITOR_THRESHOLD_1_ERR
	if strings.ToLower(resp.Threshold1Crossing) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeVoltageMonitorThreshold1Err]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeVoltageMonitorThreshold1Err,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// VOLTAGE_MONITOR_THRESHOLD_2_ERR
	if strings.ToLower(resp.Threshold2Crossing) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeVoltageMonitorThreshold2Err]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeVoltageMonitorThreshold2Err,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	return metrics
}

func nicAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *NICAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}

	// NIC_HIGH_BER_ERR
	if strings.ToLower(resp.HighBER) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeNicHighBerErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeNicHighBerErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// NIC_MULTIPLE_NACK_ERR
	if strings.ToLower(resp.MultipleNAC) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeNicMultipleNackErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeNicMultipleNackErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// NIC_MULTIPLE_CRC_ERR
	if strings.ToLower(resp.MultipleCRCError) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeNicMultipleCrcErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeNicMultipleCrcErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// NIC_MULTIPLE_LINK_RE_INITIALIZATIONS_ERR
	if strings.ToLower(resp.MultipleLinkReinitializations) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeNicMultipleLinkReInitializationsErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeNicMultipleLinkReInitializationsErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// NIC_MULTIPLE_LINK_RE_TRANSMISSIONS_ERR
	if strings.ToLower(resp.MultipleLinkRetransmissions) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeNicMultipleLinkReTransmissionsErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeNicMultipleLinkReTransmissionsErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	return metrics
}

func checksumAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *ChecksumAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}

	if resp.RMATables != nil {

		// CHECKSUM_HBM_RMA_PERITY_ERR
		if strings.ToLower(resp.RMATables.HBMParity) == "fail" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumHbmRmaPerityErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaPerityErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// CHECKSUM_HBM_RMA_ECC_ERR
		if strings.ToLower(resp.RMATables.HBMECC) == "fail" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumHbmRmaEccErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaEccErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// CHECKSUM_HBM_RMA_SRAM_ECC_ERR
		if strings.ToLower(resp.RMATables.SRAMECC) == "fail" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumHbmRmaSramEccErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaSramEccErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// CHECKSUM_HBM_RMA_ROW_REPLACEMENT_ERR
		if strings.ToLower(resp.RMATables.RowReplacement) == "fail" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumHbmRmaRowReplacementErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaRowReplacementErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// CHECKSUM_HBM_RMA_CAUSE_ERR
		if strings.ToLower(resp.RMATables.RMACauseFailure) == "fail" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumHbmRmaCauseErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaCauseErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// CHECKSUM_HBM_RMA_SPARE_ROW_AVAILABILITY_ERR
		if strings.ToLower(resp.RMATables.SpareRowAvailability) == "fail" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumHbmRmaSpareRowAvailabilityErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaSpareRowAvailabilityErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}
	}

	if resp.SPIChecksum != nil {

		// CHECKSUM_PPBOOT_PRIMARY_ERR
		if strings.ToLower(resp.SPIChecksum.PpBootPrimary) == "fail" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumPpbootPrimaryErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeChecksumPpbootPrimaryErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// CHECKSUM_PPBOOT_SECONDARY_ERR
		if strings.ToLower(resp.SPIChecksum.PpBootSecondary) == "fail" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumPpbootSecondaryErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeChecksumPpbootSecondaryErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// CHECKSUM_PREBOOT_PRIMARY_ERR
		if strings.ToLower(resp.SPIChecksum.PreBootPrimary) == "fail" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumPrebootPrimaryErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeChecksumPrebootPrimaryErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}

		// CHECKSUM_PREBOOT_SECONDARY_ERR
		if strings.ToLower(resp.SPIChecksum.PreBootSecondary) == "fail" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumPrebootSecondaryErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": alertData.Description,
					"error_type":  bmcmonitoring.AlertTypeChecksumPrebootSecondaryErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     convertAlertTime(resp.Timestamp),
					"severity": alertData.Severity,
				},
			})
		}
	}

	// CHECKSUM_CPLD_ERR
	if strings.ToLower(resp.CPLDStatus) == "fail" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumCpldErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeChecksumCpldErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	switch strings.ToLower(resp.EEPROM) {
	case "i/o fail":

		// AlertTypeChecksumEepromDeviceAccessErr
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumEepromDeviceAccessErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeChecksumEepromDeviceAccessErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	case "checksum fail":

		// CHECKSUM_EEPROM_CRC_ERR
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumEepromCrcErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeChecksumEepromCrcErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	case "pass":
	default:

		// CHECKSUM_EEPROM_UNKNOWN_ERR
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeChecksumEepromUnknownErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeChecksumEepromUnknownErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}
	return metrics
}

func pllAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *PLLAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}

	// if one of the values is not locked then we send an alert
	for name, value := range *resp {

		// PLL_LOCK_ERR
		if strings.ToLower(value) == "not locked" {
			alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypePllLockErr]
			metrics = append(metrics, bmcmonitoring.Metric{
				Hostname:   hostname,
				Oam:        fmt.Sprintf("%d", oam),
				MetricName: bmcmonitoring.PrefixAlerts,
				CustomLabels: map[string]string{
					"description": fmt.Sprintf("%s. %s is not locked", alertData.Description, name),
					"error_type":  bmcmonitoring.AlertTypePllLockErr,

					// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
					"time":     "",
					"severity": alertData.Severity,
				},
			})

		}
	}

	return metrics
}

func fitAlerts(log *logrus.Entry, alertsData bmcmonitoring.AlertsData, resp *FITAlert, hostname string, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	if resp == nil {
		return metrics
	}

	// FIT_INVALID_IMAGE_FORMAT_ERR
	if strings.ToLower(resp.InvalidImageFormat) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFitInvalidImageFormatErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeFitInvalidImageFormatErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// FIT_TARGET_ADDREES_ERR
	if strings.ToLower(resp.FWUpgradeTargetAddressViolation) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFitTargetAddreesErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeFitTargetAddreesErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	// FIT_NOT_RUNABLE_ERR
	if strings.ToLower(resp.FITNotRunnable) == "true" {
		alertData := alertsData.AlertsInformation[bmcmonitoring.AlertTypeFitNotRunableErr]
		metrics = append(metrics, bmcmonitoring.Metric{
			Hostname:   hostname,
			Oam:        fmt.Sprintf("%d", oam),
			MetricName: bmcmonitoring.PrefixAlerts,
			CustomLabels: map[string]string{
				"description": alertData.Description,
				"error_type":  bmcmonitoring.AlertTypeFitNotRunableErr,

				// TODO: to what the timestamp here refers to? what happens if we have multiple alerts?
				"time":     convertAlertTime(resp.Timestamp),
				"severity": alertData.Severity,
			},
		})
	}

	return metrics
}
