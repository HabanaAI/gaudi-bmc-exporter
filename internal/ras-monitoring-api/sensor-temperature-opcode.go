package rasmonitoringapi

import rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

var (
	sensorTemperatureOpcodes = map[string]Opcode{
		rasmonitoring.SensorTemperatureMetricOnDie0: {
			OpcodeNumber: 9,
			Offset:       4,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricOnDie1: {
			OpcodeNumber: 9,
			Offset:       8,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricOnDie2: {
			OpcodeNumber: 9,
			Offset:       12,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricOnDie3: {
			OpcodeNumber: 9,
			Offset:       16,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricHBM0: {
			OpcodeNumber: 9,
			Offset:       20,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricHBM1: {
			OpcodeNumber: 9,
			Offset:       24,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricHBM2: {
			OpcodeNumber: 9,
			Offset:       28,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricHBM3: {
			OpcodeNumber: 9,
			Offset:       32,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricHBM4: {
			OpcodeNumber: 9,
			Offset:       36,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricHBM5: {
			OpcodeNumber: 9,
			Offset:       40,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricCPLDLocal: {
			OpcodeNumber: 9,
			Offset:       44,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricCPLD0: {
			OpcodeNumber: 9,
			Offset:       48,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricCPLD1: {
			OpcodeNumber: 9,
			Offset:       52,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricCPLD2: {
			OpcodeNumber: 9,
			Offset:       56,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricCPLD3: {
			OpcodeNumber: 9,
			Offset:       60,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricOnboard0: {
			OpcodeNumber: 9,
			Offset:       64,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricOnboard1: {
			OpcodeNumber: 9,
			Offset:       68,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricOnboard2: {
			OpcodeNumber: 9,
			Offset:       72,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricOnboard3: {
			OpcodeNumber: 9,
			Offset:       76,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricCPLDTemp: {
			OpcodeNumber: 9,
			Offset:       80,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricPSUStage1: {
			OpcodeNumber: 9,
			Offset:       84,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorTemperatureMetricPSUStage2: {
			OpcodeNumber: 9,
			Offset:       88,
			Length:       4,
			ExpectedType: Int,
		},
	}
)
