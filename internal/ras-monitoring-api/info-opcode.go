package rasmonitoringapi

import rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

var (
	infoOpcodes = map[string]Opcode{
		rasmonitoring.InfoMetricDeviceID: {
			OpcodeNumber: 0,
			Offset:       0,
			Length:       2,
			ExpectedType: ReverseString,
		},
		rasmonitoring.InfoMetricSubsystemDeviceID: {
			OpcodeNumber: 0,
			Offset:       4,
			Length:       2,
			ExpectedType: ReverseString,
		},
		rasmonitoring.InfoMetricSubsystemVendorID: {
			OpcodeNumber: 0,
			Offset:       2,
			Length:       2,
			ExpectedType: ReverseString,
		},
		rasmonitoring.InfoMetricASICSerialNumber: {
			OpcodeNumber: 0,
			Offset:       6,
			Length:       8,
			ExpectedType: AsciiString,
		},
		rasmonitoring.InfoMetricBoardSerialNumber: {
			OpcodeNumber: 0,
			Offset:       14,
			Length:       16,
			ExpectedType: AsciiString,
		},
		rasmonitoring.InfoMetricSRAMSize: {
			OpcodeNumber: 0,
			Offset:       30,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.InfoMetricHBMSize: {
			OpcodeNumber: 0,
			Offset:       32,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.InfoMetricUUID: {
			OpcodeNumber: 0,
			Offset:       34,
			Length:       14,
			ExpectedType: AsciiString,
		},
	}
)
