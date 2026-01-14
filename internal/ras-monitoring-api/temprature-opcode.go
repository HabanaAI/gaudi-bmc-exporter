package rasmonitoringapi

import rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

var (
	temperatureOpcodes = map[string]Opcode{
		rasmonitoring.TemperatureMetricCurrentBoardTemp: {
			OpcodeNumber: 2,
			Offset:       0,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricCurrentVRMTemp: {
			OpcodeNumber: 2,
			Offset:       4,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricCurrentDRAMTemp: {
			OpcodeNumber: 2,
			Offset:       8,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricCurrentOnDieTemp: {
			OpcodeNumber: 2,
			Offset:       12,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricHistoricalBoardTemp: {
			OpcodeNumber: 2,
			Offset:       16,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricHistoricalVRMTemp: {
			OpcodeNumber: 2,
			Offset:       20,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricHistoricalDRAMTemp: {
			OpcodeNumber: 2,
			Offset:       24,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricHistoricalOnDieTemp: {
			OpcodeNumber: 2,
			Offset:       28,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricMaxTempRiseTime: {
			OpcodeNumber: 2,
			Offset:       32,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricMaxSocTempErrorThreshold: {
			OpcodeNumber: 2,
			Offset:       36,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricMaxSocTempWarmingThreshold: {
			OpcodeNumber: 2,
			Offset:       40,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricMaxHbmTempThreshold: {
			OpcodeNumber: 2,
			Offset:       44,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricCurrentSocTempErrorThreshold: {
			OpcodeNumber: 2,
			Offset:       48,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricCurrentSocTempWarningThreshold: {
			OpcodeNumber: 2,
			Offset:       52,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.TemperatureMetricCurrentHbmTempThreshold: {
			OpcodeNumber: 2,
			Offset:       56,
			Length:       4,
			ExpectedType: Int,
		},
	}
)
