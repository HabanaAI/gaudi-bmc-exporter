package rasmonitoringapi

import (
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"log/slog"
	"strconv"
)

var (
	statusOpcodes = map[string]Opcode{
		rasmonitoring.StatusMetricBootStage: {
			OpcodeNumber: 1,
			Offset:       0,
			Length:       1,
			ExpectedType: Int,
			Decoder:      decodeBootStage,
		},
		rasmonitoring.StatusMetricEmergencyPowerReduction: {
			OpcodeNumber: 1,
			Offset:       1,
			Length:       1,
			ExpectedType: Int,
			Decoder:      decodeEmergencyPowerReduction,
		},
		rasmonitoring.StatusMetricClockThrottling: {
			OpcodeNumber: 1,
			Offset:       10,
			Length:       1,
			ExpectedType: Int,
			Decoder:      decodeClockThrottling,
		},
		rasmonitoring.StatusMetricLastClockThrottlingDuration: {
			OpcodeNumber: 1,
			Offset:       11,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.StatusMetricPowerState: {
			OpcodeNumber: 1,
			Offset:       15,
			Length:       2,
			ExpectedType: Int,
			Decoder:      decodePowerState,
		},
		rasmonitoring.StatusMetricTotalClockThrottlingDuration: {
			OpcodeNumber: 1,
			Offset:       18,
			Length:       8,
			ExpectedType: Int,
		},
		rasmonitoring.StatusMetricGlobalTimeFromReset: {
			OpcodeNumber: 1,
			Offset:       26,
			Length:       8,
			ExpectedType: Int,
		},
		rasmonitoring.StatusMetricChipStatus: {
			OpcodeNumber: 1,
			Offset:       44,
			Length:       1,
			ExpectedType: Int,
			Decoder:      decodeChipStatus,
		},
		rasmonitoring.StatusMetricDeviceActivity: {
			OpcodeNumber: 1,
			Offset:       45,
			Length:       1,
			ExpectedType: Int,
			Decoder:      decodeDeviceActivity,
		},
		rasmonitoring.StatusMetricDeviceActivityCounter: {
			OpcodeNumber: 1,
			Offset:       46,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.StatusMetricDevicePowerReduction: {
			OpcodeNumber: 1,
			Offset:       50,
			Length:       1,
			ExpectedType: Int,
			Decoder:      decodeDevicePowerReduction,
		},
		rasmonitoring.StatusMetricLastPowerReductionDuration: {
			OpcodeNumber: 1,
			Offset:       51,
			Length:       4,
			ExpectedType: Int,
		},
	}
)

func decodeDevicePowerReduction(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var state string
	switch val {
	case "0":
		state = "Max Reduction"
	case "1":
		state = "2nd Reduction"
	case "2":
		state = "1st Reduction"
	case "3":
		state = "Normal Power"
	default:
		return nil, fmt.Errorf("unexpected device power reduction: %s", val)
	}

	value, err := strconv.Atoi(val)
	if err != nil {
		slog.Error("failed to convert image source to int", "error", err, "value", val)
	}

	return []rasmonitoring.Metric{
		{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				"state": state,
			},
		},
	}, nil
}
func decodeBootStage(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var stage string
	switch val {
	case "0":
		stage = "Linux"
	case "1":
		stage = "Uboot"
	case "2":
		stage = "Preboot"
	case "3":
		stage = "Zephyr"
	default:
		return nil, fmt.Errorf("unsupported stage: %s", val)
	}

	value, err := strconv.Atoi(val)
	if err != nil {
		slog.Error("failed to convert image source to int", "error", err, "value", val)
	}

	return []rasmonitoring.Metric{
		{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				fieldName: stage,
			},
		},
	}, nil
}

func decodeEmergencyPowerReduction(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var stage string
	switch val {
	case "0":
		stage = "Normal"
	case "1":
		stage = "Reduced"
	default:
		return nil, fmt.Errorf("unsupported stage: %s", val)
	}

	value, err := strconv.Atoi(val)
	if err != nil {
		slog.Error("failed to convert image source to int", "error", err, "value", val)
	}

	return []rasmonitoring.Metric{
		{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				fieldName: stage,
			},
		},
	}, nil
}

func decodeClockThrottling(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var stage string
	switch val {
	case "0":
		stage = "None"
	case "1":
		stage = "Power"
	case "2":
		stage = "Thermal"
	default:
		return nil, fmt.Errorf("unsupported stage: %s", val)
	}

	value, err := strconv.Atoi(val)
	if err != nil {
		slog.Error("failed to convert image source to int", "error", err, "value", val)
	}

	return []rasmonitoring.Metric{
		{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				fieldName: stage,
			},
		},
	}, nil
}

func decodePowerState(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var stage string
	switch val {
	case "0":
		stage = "Full performance"
	case "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15":
		stage = fmt.Sprintf("Reduced by %s/16", val)
	default:
		return nil, fmt.Errorf("unsupported stage: %s", val)
	}

	value, err := strconv.Atoi(val)
	if err != nil {
		slog.Error("failed to convert image source to int", "error", err, "value", val)
	}

	return []rasmonitoring.Metric{
		{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				fieldName: stage,
			},
		},
	}, nil
}

func decodeChipStatus(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var stage string
	switch val {
	case "0":
		stage = "Processing"
	case "1":
		stage = "Idle"
	default:
		return nil, fmt.Errorf("unsupported stage: %s", val)
	}

	value, err := strconv.Atoi(val)
	if err != nil {
		slog.Error("failed to convert image source to int", "error", err, "value", val)
	}

	return []rasmonitoring.Metric{
		{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				fieldName: stage,
			},
		},
	}, nil
}

func decodeDeviceActivity(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var state string
	switch val {
	case "0":
		state = "Device not in use"
	case "1":
		state = "Device in-use"
	default:
		return nil, fmt.Errorf("unsupported state: %s", val)
	}

	value, err := strconv.Atoi(val)
	if err != nil {
		slog.Error("failed to convert image source to int", "error", err, "value", val)
	}

	return []rasmonitoring.Metric{
		{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				fieldName: state,
			},
		},
	}, nil
}
