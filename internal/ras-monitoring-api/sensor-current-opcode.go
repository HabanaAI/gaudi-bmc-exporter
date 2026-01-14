package rasmonitoringapi

import rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

var (
	sensorCurrentOpcodes = map[string]Opcode{
		rasmonitoring.SensorCurrentVin54: {
			OpcodeNumber: 9,
			Offset:       220,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorCurrentP1Vin12: {
			OpcodeNumber: 9,
			Offset:       224,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorCurrentStage154Vin: {
			OpcodeNumber: 9,
			Offset:       228,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorCurrentStage113P5VOut: {
			OpcodeNumber: 9,
			Offset:       232,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorCurrentStage213P5Vin: {
			OpcodeNumber: 9,
			Offset:       236,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorCurrentStage2CoreOut: {
			OpcodeNumber: 9,
			Offset:       240,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorCurrentStage2HBMout: {
			OpcodeNumber: 9,
			Offset:       244,
			Length:       4,
			ExpectedType: Int,
		},
	}
)
