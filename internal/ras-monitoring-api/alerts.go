package rasmonitoringapi

import (
	"context"
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/ipmi"
	"strconv"

	"github.com/sirupsen/logrus"
)

func (c *Client) Alerts(ctx context.Context, alertsData rasmonitoring.AlertsData, log *logrus.Entry) []rasmonitoring.Metric {

	// Get the sel alerts
	selInfo, err := c.ipmiClient.SelList()
	if err != nil {
		log.WithField("metric", "alerts").WithError(err).Error("failed getting sel elist")
		return nil
	}

	// filter only the habana alerts
	var RasAlertsRecordIDs []string
	var metrics []rasmonitoring.Metric
	for _, a := range selInfo {

		if a.SensorType == "Add-in Card" {

			oam, err := oamFromSensorNum(a)
			if err != nil {
				log.WithError(err).Error("failed to get oam from sensor number")
				continue
			}

			alertType, err := alertTypeFromEventData(a)
			if err != nil {
				log.WithError(err).Error("failed to get alert type from event data")
				continue
			}

			// skip POWER_THRESHOLD_CROSSING_ERR, NIC_MULTIPLE_NACK_ERR,NIC_MULTIPLE_LINK_RE_TRANSMISSIONS_ERR alerts
			if alertType == rasmonitoring.AlertTypePowerThresholdCrossingErr || alertType == rasmonitoring.AlertTypeNicMultipleLinkReTransmissionsErr ||
				alertType == rasmonitoring.AlertTypeNicMultipleNackErr || alertType == rasmonitoring.AlertTypeSecurityPermissionAccessErr ||
				alertType == rasmonitoring.AlertTypeFatalSramEccErr || alertType == rasmonitoring.AlertTypeHbmMbistRepairFailureErr {
				continue
			}

			RasAlertsRecordIDs = append(RasAlertsRecordIDs, a.RecordID)

			// Add the severity and description to the alert if exists in the file
			description := "unknown description"
			severity := "unknown severity"

			if alertInfo, ok := alertsData.AlertsInformation[alertType]; ok {
				description = alertInfo.Description
				severity = alertInfo.Severity
			} else {
				log.WithError(fmt.Errorf("alert %s is not in the alerts description file", alertType)).Error()
			}

			metrics = append(metrics, rasmonitoring.Metric{
				Hostname:    c.Hostname,
				Oam:         fmt.Sprintf("%d", oam),
				MetricName:  rasmonitoring.PrefixAlerts,
				MetricValue: 1,
				CustomLabels: map[string]string{
					"error_type":  alertType,
					"description": description,
					"time":        a.Timestamp,
					"severity":    severity,
				},
			})

		}
	}

	// delete only RAS alerts from the sel entries.
	for _, id := range RasAlertsRecordIDs {
		_, err := c.ipmiClient.ExecCommand(fmt.Sprintf("sel delete 0x%s", id))
		if err != nil {
			log.WithError(err).Error()
		}
	}

	return metrics
}

// oamFromSensorNum will return the oam number from the sensor number.
// 14 - oam0, 15 - oam1, etc.
func oamFromSensorNum(a ipmi.SelInfo) (int, error) {
	sn, err := strconv.ParseInt(a.SensorNumber, 16, 64)
	if err != nil {
		return 0, err
	}

	// 20 - 14 in hex
	return int(sn) - 20, nil
}

const (
	DOUBLE_ECC_BLOCK_ERR = iota /* 0 */
	SINGLE_ECC_BLOCK_ERR
	HBM_COMMAND_ADDRESS_ERR
	HBM_DATA_WRITE_ERR
	HBM_DATA_READ_ERR
	DOUBLE_ECC_ERR /* 5 */
	SINGLE_ECC_ERR
	SRAM_UNREPAIRABLE_RINGS_ERR
	SRAM_REPAIRED_RINGS_ERR
	HBM_CATASTROPHIC_TEMP_ERR
	PCIE_UNCORRECTABLE_FATAL_ERR /* 10 */
	PCIE_RECEIVER_OVERFLOW_ERR
	PCIE_FLOW_CONTROL_PROTOCOL_ERR
	PCIE_SURPRISE_LINK_DOWN_ERR
	PCIE_MALFUNCTION_TLP_ERR
	PCIE_DLLP_PROTOCOL_ERR /* 15 */
	PCIE_ECRC_ERR
	PCIE_COMPLETION_TIMEOUT_ERR
	PCIE_CORRECTED_INTERNAL_ERR
	PCIE_REPLAY_BUFFER_ROLLOVER_ERR
	PCIE_REPLAY_TIMER_TIMEOUT_ERR /* 20 */
	PCIE_BAD_DLLP_ERR
	PCIE_BAD_TLP_ERR
	PCIE_RECEIVER_ERR
	PCIE_LCRD_ERR
	PCIE_RX_NAK_DLLP_ERR /* 25 */
	PCIE_TX_NAK_DLLP_ERR
	PCIE_RETRY_TLP_ERR
	NIC_HIGH_BER_ERR
	NIC_MULTIPLE_NACK_ERR
	NIC_MULTIPLE_CRC_ERR /* 30 */
	NIC_MULTIPLE_LINK_RE_INITIALIZATIONS_ERR
	NIC_MULTIPLE_LINK_RE_TRANSMISSIONS_ERR
	THERMAL_THRESHOLD_CROSSING_ERR
	THERMAL_TEMP_RISE_TIME_VIOLATION_ERR
	VOLTAGE_MONITOR_THRESHOLD_1_ERR /* 35 */
	VOLTAGE_MONITOR_THRESHOLD_2_ERR
	SECURITY_SPI_REVOKED_PUB_KEY_ERR
	SECURITY_SPI_SVN_ERR
	SECURITY_SPI_SIGNATURE_ERR
	SECURITY_AGENT_REVOKED_PUB_KEY_ERR /* 40 */
	SECURITY_AGENT_SVN_ERR
	SECURITY_AGENT_SIGNATURE_ERR
	SECURITY_AGENT_UBOOT_ERR
	SECURITY_AGENT_LINUX_ERR
	SECURITY_AGENT_ZEPHYR_ERR /* 45 */
	SECURITY_REVOKE_LAST_KEY_ERR
	SECURITY_PERMISSION_ACCESS_ERR
	WD_HEARTBEAT_ERR
	WD_TDR_ERR
	WD_HW_ERR /* 50 */
	CHECKSUM_PPBOOT_PRIMARY_ERR
	CHECKSUM_PPBOOT_SECONDARY_ERR
	CHECKSUM_PREBOOT_PRIMARY_ERR
	CHECKSUM_PREBOOT_SECONDARY_ERR
	CHECKSUM_CPLD_ERR /* 55 */
	CHECKSUM_HBM_RMA_PERITY_ERR
	CHECKSUM_HBM_RMA_ECC_ERR
	CHECKSUM_HBM_RMA_SRAM_ECC_ERR
	CHECKSUM_HBM_RMA_ROW_REPLACEMENT_ERR
	CHECKSUM_HBM_RMA_CAUSE_ERR /* 60 */
	FIT_INVALID_IMAGE_FORMAT_ERR
	FIT_TARGET_ADDREES_ERR
	FIT_NOT_RUNABLE_ERR
	POWER_THRESHOLD_CROSSING_ERR
	FATAL_SRAM_ECC_ERR
	FATAL_HBM_ECC_ERR
	FATAL_HBM_PARITY_ERR
	FATAL_HBM_INIT_ERR
	FATAL_ROW_REPLACEMENT_ERR
	HBM_MBIST_REPAIR_FAILURE_ERR
	VOLTAGE_MONITOR_THRESHOLD
	CHECKSUM_EEPROM_DEVICE_ACCESS_ERR
	CHECKSUM_EEPROM_CRC_ERR
	CHECKSUM_EEPROM_UNKNOWN_ERR
	DYNAMIC_ROW_REPLACEMENT_ERR
	PLL_LOCK_ERR
	DOUBLE_ECC_COUNTER_ERR
	SINGLE_ECC_COUNTER_ERR
	DOUBLE_HBM_COUNTER_ERR
	SINGLE_HBM_COUNTER_ERR
)

var alertTypesMapping = map[int64]string{
	DOUBLE_ECC_BLOCK_ERR:                     rasmonitoring.AlertTypeDoubleECCBlockErr,
	SINGLE_ECC_BLOCK_ERR:                     rasmonitoring.AlertTypeSingleECCBlockErr,
	HBM_COMMAND_ADDRESS_ERR:                  rasmonitoring.AlertTypeHBMCommandAddressErr,
	HBM_DATA_WRITE_ERR:                       rasmonitoring.AlertTypeHBMDataWriteErr,
	HBM_DATA_READ_ERR:                        rasmonitoring.AlertTypeHBMDataReadErr,
	DOUBLE_ECC_ERR:                           rasmonitoring.AlertTypeDoubleECCErr,
	SINGLE_ECC_ERR:                           rasmonitoring.AlertTypeSingleECCErr,
	SRAM_UNREPAIRABLE_RINGS_ERR:              rasmonitoring.AlertTypeSramUnrepairableRingsErr,
	SRAM_REPAIRED_RINGS_ERR:                  rasmonitoring.AlertTypeSramRepairedRingsErr,
	HBM_CATASTROPHIC_TEMP_ERR:                rasmonitoring.AlertTypeHbmCatastrophicTempErr,
	PCIE_UNCORRECTABLE_FATAL_ERR:             rasmonitoring.AlertTypePcieUncorrectableFatalErr,
	PCIE_RECEIVER_OVERFLOW_ERR:               rasmonitoring.AlertTypePcieReceiverOverflowErr,
	PCIE_FLOW_CONTROL_PROTOCOL_ERR:           rasmonitoring.AlertTypePcieFlowControlProtocolErr,
	PCIE_SURPRISE_LINK_DOWN_ERR:              rasmonitoring.AlertTypePcieSurpriseLinkDownErr,
	PCIE_MALFUNCTION_TLP_ERR:                 rasmonitoring.AlertTypePcieMalfunctionTlpErr,
	PCIE_DLLP_PROTOCOL_ERR:                   rasmonitoring.AlertTypePcieDllpProtocolErr,
	PCIE_ECRC_ERR:                            rasmonitoring.AlertTypePcieEcrcErr,
	PCIE_COMPLETION_TIMEOUT_ERR:              rasmonitoring.AlertTypePcieCompletionTimeoutErr,
	PCIE_CORRECTED_INTERNAL_ERR:              rasmonitoring.AlertTypePcieCorrectedInternalErr,
	PCIE_REPLAY_BUFFER_ROLLOVER_ERR:          rasmonitoring.AlertTypePcieReplayBufferRolloverErr,
	PCIE_REPLAY_TIMER_TIMEOUT_ERR:            rasmonitoring.AlertTypePcieReplayTimerTimeoutErr,
	PCIE_BAD_DLLP_ERR:                        rasmonitoring.AlertTypePcieBadDllpErr,
	PCIE_BAD_TLP_ERR:                         rasmonitoring.AlertTypePcieBadTlpErr,
	PCIE_RECEIVER_ERR:                        rasmonitoring.AlertTypePcieReceiverErr,
	PCIE_LCRD_ERR:                            rasmonitoring.AlertTypePcieLcrdErr,
	PCIE_RX_NAK_DLLP_ERR:                     rasmonitoring.AlertTypePcieRxNakDllpErr,
	PCIE_TX_NAK_DLLP_ERR:                     rasmonitoring.AlertTypePcieTxNakDllpErr,
	PCIE_RETRY_TLP_ERR:                       rasmonitoring.AlertTypePcieRetryTlpErr,
	NIC_HIGH_BER_ERR:                         rasmonitoring.AlertTypeNicHighBerErr,
	NIC_MULTIPLE_NACK_ERR:                    rasmonitoring.AlertTypeNicMultipleNackErr,
	NIC_MULTIPLE_CRC_ERR:                     rasmonitoring.AlertTypeNicMultipleCrcErr,
	NIC_MULTIPLE_LINK_RE_INITIALIZATIONS_ERR: rasmonitoring.AlertTypeNicMultipleLinkReInitializationsErr,
	NIC_MULTIPLE_LINK_RE_TRANSMISSIONS_ERR:   rasmonitoring.AlertTypeNicMultipleLinkReTransmissionsErr,
	THERMAL_THRESHOLD_CROSSING_ERR:           rasmonitoring.AlertTypeThermalThresholdCrossingErr,
	THERMAL_TEMP_RISE_TIME_VIOLATION_ERR:     rasmonitoring.AlertTypeThermalTempRiseTimeViolationErr,
	VOLTAGE_MONITOR_THRESHOLD_1_ERR:          rasmonitoring.AlertTypeVoltageMonitorThreshold1Err,
	VOLTAGE_MONITOR_THRESHOLD_2_ERR:          rasmonitoring.AlertTypeVoltageMonitorThreshold2Err,
	SECURITY_SPI_REVOKED_PUB_KEY_ERR:         rasmonitoring.AlertTypeSecuritySpiRevokedPubKeyErr,
	SECURITY_SPI_SVN_ERR:                     rasmonitoring.AlertTypeSecuritySpiSvnErr,
	SECURITY_SPI_SIGNATURE_ERR:               rasmonitoring.AlertTypeSecuritySpiSignatureErr,
	SECURITY_AGENT_REVOKED_PUB_KEY_ERR:       rasmonitoring.AlertTypeSecurityAgentRevokedPubKeyErr,
	SECURITY_AGENT_SVN_ERR:                   rasmonitoring.AlertTypeSecurityAgentSvnErr,
	SECURITY_AGENT_SIGNATURE_ERR:             rasmonitoring.AlertTypeSecurityAgentSignatureErr,
	SECURITY_AGENT_UBOOT_ERR:                 rasmonitoring.AlertTypeSecurityAgentUbootErr,
	SECURITY_AGENT_LINUX_ERR:                 rasmonitoring.AlertTypeSecurityAgentLinuxErr,
	SECURITY_AGENT_ZEPHYR_ERR:                rasmonitoring.AlertTypeSecurityAgentZephyrErr,
	SECURITY_REVOKE_LAST_KEY_ERR:             rasmonitoring.AlertTypeSecurityRevokeLastKeyErr,
	SECURITY_PERMISSION_ACCESS_ERR:           rasmonitoring.AlertTypeSecurityPermissionAccessErr,
	WD_HEARTBEAT_ERR:                         rasmonitoring.AlertTypeWdHeartbeatErr,
	WD_TDR_ERR:                               rasmonitoring.AlertTypeWdTdrErr,
	WD_HW_ERR:                                rasmonitoring.AlertTypeWdHwErr,
	CHECKSUM_PPBOOT_PRIMARY_ERR:              rasmonitoring.AlertTypeChecksumPpbootPrimaryErr,
	CHECKSUM_PPBOOT_SECONDARY_ERR:            rasmonitoring.AlertTypeChecksumPpbootSecondaryErr,
	CHECKSUM_PREBOOT_PRIMARY_ERR:             rasmonitoring.AlertTypeChecksumPrebootPrimaryErr,
	CHECKSUM_PREBOOT_SECONDARY_ERR:           rasmonitoring.AlertTypeChecksumPrebootSecondaryErr,
	CHECKSUM_CPLD_ERR:                        rasmonitoring.AlertTypeChecksumCpldErr,
	CHECKSUM_HBM_RMA_PERITY_ERR:              rasmonitoring.AlertTypeChecksumHbmRmaPerityErr,
	CHECKSUM_HBM_RMA_ECC_ERR:                 rasmonitoring.AlertTypeChecksumHbmRmaEccErr,
	CHECKSUM_HBM_RMA_SRAM_ECC_ERR:            rasmonitoring.AlertTypeChecksumHbmRmaSramEccErr,
	CHECKSUM_HBM_RMA_ROW_REPLACEMENT_ERR:     rasmonitoring.AlertTypeChecksumHbmRmaRowReplacementErr,
	CHECKSUM_HBM_RMA_CAUSE_ERR:               rasmonitoring.AlertTypeChecksumHbmRmaCauseErr,
	FIT_INVALID_IMAGE_FORMAT_ERR:             rasmonitoring.AlertTypeFitInvalidImageFormatErr,
	FIT_TARGET_ADDREES_ERR:                   rasmonitoring.AlertTypeFitTargetAddreesErr,
	FIT_NOT_RUNABLE_ERR:                      rasmonitoring.AlertTypeFitNotRunableErr,
	POWER_THRESHOLD_CROSSING_ERR:             rasmonitoring.AlertTypePowerThresholdCrossingErr,
	FATAL_SRAM_ECC_ERR:                       rasmonitoring.AlertTypeFatalSramEccErr,
	FATAL_HBM_ECC_ERR:                        rasmonitoring.AlertTypeFatalHbmEccErr,
	FATAL_HBM_PARITY_ERR:                     rasmonitoring.AlertTypeFatalHbmParityErr,
	FATAL_HBM_INIT_ERR:                       rasmonitoring.AlertTypeFatalHbmInitErr,
	FATAL_ROW_REPLACEMENT_ERR:                rasmonitoring.AlertTypeFatalRowReplacementErr,
	HBM_MBIST_REPAIR_FAILURE_ERR:             rasmonitoring.AlertTypeHbmMbistRepairFailureErr,
	VOLTAGE_MONITOR_THRESHOLD:                rasmonitoring.AlertTypeVoltageMonitorThreshold,
	CHECKSUM_EEPROM_DEVICE_ACCESS_ERR:        rasmonitoring.AlertTypeChecksumEepromDeviceAccessErr,
	CHECKSUM_EEPROM_CRC_ERR:                  rasmonitoring.AlertTypeChecksumEepromCrcErr,
	CHECKSUM_EEPROM_UNKNOWN_ERR:              rasmonitoring.AlertTypeChecksumEepromUnknownErr,
	DYNAMIC_ROW_REPLACEMENT_ERR:              rasmonitoring.AlertTypeDynamicRowReplacementErr,
	PLL_LOCK_ERR:                             rasmonitoring.AlertTypePllLockErr,
	DOUBLE_ECC_COUNTER_ERR:                   rasmonitoring.AlertTypeDoubleEccCounterErr,
	SINGLE_ECC_COUNTER_ERR:                   rasmonitoring.AlertTypeSingleEccCounterErr,
	DOUBLE_HBM_COUNTER_ERR:                   rasmonitoring.AlertTypeDoubleHbmCounterErr,
	SINGLE_HBM_COUNTER_ERR:                   rasmonitoring.AlertTypeSingleHbmCounterErr,
}

// alertTypeFromEventData will return the alert type from the event data.
func alertTypeFromEventData(a ipmi.SelInfo) (string, error) {

	if len(a.EventData) < 2 {
		return "", fmt.Errorf("unexpected event data length %d", len(a.EventData))
	}

	// only address the first 2 bytes.
	data := a.EventData[:2]

	alertCode, err := strconv.ParseInt(string(data), 16, 64)
	if err != nil {
		return "", err
	}

	if _, ok := alertTypesMapping[alertCode]; !ok {
		return fmt.Sprintf("unknown error, number: %d", alertCode), nil

	}

	return alertTypesMapping[alertCode], nil
}
