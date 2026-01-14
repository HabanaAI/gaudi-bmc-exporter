package rasmonitoringapi

import rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

var (
	sensorVoltageOpcodes = map[string]Opcode{
		rasmonitoring.SensorVoltageVADC54: {
			OpcodeNumber: 9,
			Offset:       108,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVRM1in: {
			OpcodeNumber: 9,
			Offset:       112,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVRM1out: {
			OpcodeNumber: 9,
			Offset:       116,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVRM2in: {
			OpcodeNumber: 9,
			Offset:       120,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVRM2VDDout: {
			OpcodeNumber: 9,
			Offset:       124,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVRM2HBMout: {
			OpcodeNumber: 9,
			Offset:       128,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVMONPCIEVPH1P8V: {
			OpcodeNumber: 9,
			Offset:       132,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVMON1P8HBMVAA: {
			OpcodeNumber: 9,
			Offset:       136,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVMON2P5: {
			OpcodeNumber: 9,
			Offset:       140,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVMON48VHIMON: {
			OpcodeNumber: 9,
			Offset:       144,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVMONP5V: {
			OpcodeNumber: 9,
			Offset:       148,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVMON12V1: {
			OpcodeNumber: 9,
			Offset:       152,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVMONHBM: {
			OpcodeNumber: 9,
			Offset:       156,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageVMONCore: {
			OpcodeNumber: 9,
			Offset:       160,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.SensorVoltageCPLDHIMON1P8NIC: {
			OpcodeNumber: 9,
			Offset:       164,
			Length:       4,
			ExpectedType: Int,
		},
	}
)
