package bmcmonitoring

import (
	"fmt"
	"strings"
)

// Here we verify the metric values.

func basicVerification(metric Metric) error {

	if metric.Hostname == "" {
		return fmt.Errorf("hostname must not be empty")
	}

	if metric.MetricName == "" {
		return fmt.Errorf("metric name must not be empty")
	}

	return nil
}

func VerifyTemperature(metric Metric) error {

	err := basicVerification(metric)
	if err != nil {
		return err
	}

	switch {
	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

// validateCustomLabels will validate custom labels.
func validateCustomLabels(metric Metric, labels []string) error {
	if metric.CustomLabels == nil {
		return fmt.Errorf("custom labels should not be empty")
	}

	// ensure we got the correct number of labels.
	if len(metric.CustomLabels) != len(labels) {
		return fmt.Errorf("%s expecting %d labels, got %d", metric.MetricName, len(labels), len(metric.CustomLabels))
	}

	for _, label := range labels {
		if _, ok := metric.CustomLabels[label]; !ok {
			return fmt.Errorf("label %s is missing", label)
		}
	}

	return nil
}

func VerifyDirect(metric Metric) error {

	err := basicVerification(metric)
	if err != nil {
		return err
	}

	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixDirect)) {
	case DirectMetricPCIeVendorID:
		return validateCustomLabels(metric, []string{DirectMetricPCIeVendorID})

	case DirectMetricASICSerialNumber:
		return validateCustomLabels(metric, []string{DirectMetricASICSerialNumber})

	case DirectMetricIBAccessState, DirectMetricOOBAccessState:
		return validateCustomLabels(metric, []string{"state"})

	case DirectMetricOSStage:
		return validateCustomLabels(metric, []string{"stage"})
	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifyCTemperature(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}

	switch {
	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifyInfo(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}

	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixInfo)) {
	case InfoMetricUUID:
		return validateCustomLabels(metric, []string{InfoMetricUUID})

	case InfoMetricDeviceID:
		return validateCustomLabels(metric, []string{InfoMetricDeviceID})

	case InfoMetricSubsystemDeviceID:
		return validateCustomLabels(metric, []string{InfoMetricSubsystemDeviceID})

	case InfoMetricSubsystemVendorID:
		return validateCustomLabels(metric, []string{InfoMetricSubsystemVendorID})

	case InfoMetricASICSerialNumber:
		return validateCustomLabels(metric, []string{InfoMetricASICSerialNumber})

	case InfoMetricBoardSerialNumber:
		return validateCustomLabels(metric, []string{InfoMetricBoardSerialNumber})

	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifyStatus(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}

	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixStatus)) {
	case StatusMetricChipStatus:
		return validateCustomLabels(metric, []string{StatusMetricChipStatus})

	case StatusMetricBootStage:
		return validateCustomLabels(metric, []string{StatusMetricBootStage})

	case StatusMetricClockThrottling:
		return validateCustomLabels(metric, []string{StatusMetricClockThrottling})

	case StatusMetricPowerState:
		return validateCustomLabels(metric, []string{StatusMetricPowerState})

	case StatusMetricDeviceActivity:
		return validateCustomLabels(metric, []string{StatusMetricDeviceActivity})

	case StatusMetricEmergencyPowerReduction:
		return validateCustomLabels(metric, []string{StatusMetricEmergencyPowerReduction})

	case StatusMetricDevicePowerReduction:
		return validateCustomLabels(metric, []string{"state"})

	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifyFrequency(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}

	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixFrequency)) {
	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifyPower(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}

	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixPower)) {
	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifyEthernetInfo(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}

	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixEthernetInfo)) {
	case EthernetInfoMetricPHYStatus:
		return validateCustomLabels(metric, []string{"state", "phy"})
	case EthernetInfoMetricPortMapping:
		return validateCustomLabels(metric, []string{"port", "type"})
	case EthernetInfoMetricSerDesAvailability:
		return validateCustomLabels(metric, []string{"state"})
	case EthernetInfoMetricExternalLinkStatus, EthernetInfoMetricLinkStatus:
		return validateCustomLabels(metric, []string{"state", "link"})
	case EthernetInfoMetricANLTStatus:
		return validateCustomLabels(metric, []string{"state", "port"})
	case EthernetInfoMetricStateTogglingCounter:
		return validateCustomLabels(metric, []string{"port"})
	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifyEthernetStatus(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}

	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixEthernetStatus)) {
	case EthernetStatusMetricMACRemote, EthernetStatusMetricRetransmission, EthernetStatusMetricRetraining,
		EthernetStatusMetricSERPreFEC, EthernetStatusMetricLinkRetrainingDueToBER, EthernetStatusMetricSERPostFEC, EthernetStatusMetricThroughput, EthernetStatusMetricBERCorrectable,
		EthernetStatusMetricBERUncorrectable, EthernetStatusMetricNack, EthernetStatusMetricCRC, EthernetStatusMetricRetransmissionTimeout,
		EthernetStatusMetricLatency, EthernetStatusMetricStateTogglingCounter:
		return validateCustomLabels(metric, []string{"port"})
	case EthernetStatusMetricExternalLinkStatus, EthernetStatusMetricLinkStatus:
		return validateCustomLabels(metric, []string{"state", "link"})
	case EthernetStatusMetricPHYStatus:
		return validateCustomLabels(metric, []string{"state", "phy"})
	case EthernetStatusMetricPortMapping:
		return validateCustomLabels(metric, []string{"port", "type"})
	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifyPcieInfo(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}
	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixPcieInfo)) {
	case PcieInfoMetricMaxPCIeLinkSpeed:
		return validateCustomLabels(metric, []string{PcieInfoMetricMaxPCIeLinkSpeed})
	case PcieInfoMetricCurrentPCIeLinkSpeed:
		return validateCustomLabels(metric, []string{PcieInfoMetricCurrentPCIeLinkSpeed})
	case PcieInfoMetricPCIeSubsystemVendorID:
		return validateCustomLabels(metric, []string{PcieInfoMetricPCIeSubsystemVendorID})
	case PcieInfoMetricPCIeBusAndDevice:
		return validateCustomLabels(metric, []string{PcieInfoMetricPCIeBusAndDevice})
	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifySensorTemperature(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}
	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixSensorTemperature)) {

	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifySensorVoltage(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}
	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixSensorVoltage)) {

	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifySensorVoltageMonitor(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}
	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixSensorVoltageMonitor)) {

	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifySensorCurrent(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}
	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixSensorCurrent)) {

	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifySecurity(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}
	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixSecurity)) {
	case SecurityMetricTPMPCRPPBOOT:
		return validateCustomLabels(metric, []string{SecurityMetricTPMPCRPPBOOT})

	case SecurityMetricTPMPCRLINUX:
		return validateCustomLabels(metric, []string{SecurityMetricTPMPCRLINUX})
	case SecurityMetricTPMPCRBOOT:
		return validateCustomLabels(metric, []string{SecurityMetricTPMPCRBOOT})
	case SecurityMetricTPMPCRPREBOOT:
		return validateCustomLabels(metric, []string{SecurityMetricTPMPCRPREBOOT})

	case SecurityMetricCPLDVersion:
		return validateCustomLabels(metric, []string{SecurityMetricCPLDVersion})
	case SecurityMetricMinimalSVNindex:
		return validateCustomLabels(metric, []string{SecurityMetricMinimalSVNindex})
	case SecurityMetricCPLDVersionTimestamp:
		return validateCustomLabels(metric, []string{SecurityMetricCPLDVersionTimestamp})
	case SecurityMetricFWImageSource:
		return validateCustomLabels(metric, []string{SecurityMetricFWImageSource})
	case SecurityMetricKey0Revocation:
		return validateCustomLabels(metric, []string{SecurityMetricKey0Revocation})
	case SecurityMetricKey1Revocation:
		return validateCustomLabels(metric, []string{SecurityMetricKey1Revocation})
	case SecurityMetricKey2Revocation:
		return validateCustomLabels(metric, []string{SecurityMetricKey2Revocation})
	case SecurityMetricKey3Revocation:
		return validateCustomLabels(metric, []string{SecurityMetricKey3Revocation})
	case SecurityMetricKey4Revocation:
		return validateCustomLabels(metric, []string{SecurityMetricKey4Revocation})

	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func VerifyHbm(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}
	switch strings.TrimPrefix(metric.MetricName, fmt.Sprintf("%s_", PrefixHbm)) {
	case HBMMetricMbistRepair:
		return validateCustomLabels(metric, []string{HBMMetricMbistRepairLabelState, HBMMetricMbistRepairLabelIndex})
	case HBMMetricGlobalECC:
		return validateCustomLabels(metric, []string{HBMMetricGlobalECCLabelIndex})
	case HBMMetricRepairedLanes:
		return validateCustomLabels(metric, []string{HBMMetricRepairedLanesLabelHBMIndex, HBMMetricRepairedLanesLabelMCChannel})
	case HBMMetricReplaceRows:
		return validateCustomLabels(metric, []string{HBMMetricReplaceRowsLabelHBMIndex, HBMMetricReplaceRowsLabelPCIndex,
			HBMMetricReplaceRowsLabelStackID, HBMMetricReplaceRowsLabelBankIndex, HBMMetricReplaceRowsLabelCause, HBMMetricReplaceRowsLabelRowAddress})

	default:
		if metric.CustomLabels != nil && len(metric.CustomLabels) > 0 {
			return fmt.Errorf("metric %s shouldn't have any custom labels", metric.MetricName)
		}
	}

	return nil
}

func verifyBmcState(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}

	return nil
}

func VerifyAlerts(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}
	return validateCustomLabels(metric, []string{"description", "error_type", "time", "severity"})

}

func VerifyLaneInfo(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}
	return validateCustomLabels(metric, []string{"lane"})

}

func verifyEthernetStatusCounters(metric Metric) error {
	err := basicVerification(metric)
	if err != nil {
		return err
	}

	return validateCustomLabels(metric, []string{"port"})

}
