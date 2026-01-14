package rasmonitoringapi

import rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

var (
	powerOpcodes = map[string]Opcode{
		rasmonitoring.PowerMetricCurrentPowerConsumption: {
			OpcodeNumber: 4,
			Offset:       0,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PowerMetricPeakPowerConsumption: {
			OpcodeNumber: 4,
			Offset:       4,
			Length:       4,
			ExpectedType: Int,
		},
	}
)
