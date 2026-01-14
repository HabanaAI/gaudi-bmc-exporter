package rasmonitoringapi

import rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

var (
	frequencyOpcodes = map[string]Opcode{
		rasmonitoring.FrequencyMetricHBMFrequency: {
			OpcodeNumber: 3,
			Offset:       0,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricMaxTPCFrequency: {
			OpcodeNumber: 3,
			Offset:       2,
			Length:       2,
			ExpectedType: Int,
		},

		rasmonitoring.FrequencyMetricMaxMMEFrequency: {
			OpcodeNumber: 3,
			Offset:       4,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricMaxDMAFrequency: {
			OpcodeNumber: 3,
			Offset:       6,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricMaxMediaFrequency: {
			OpcodeNumber: 3,
			Offset:       8,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricMaxPCIeFrequency: {
			OpcodeNumber: 3,
			Offset:       10,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricMaxARMFrequency: {
			OpcodeNumber: 3,
			Offset:       12,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricMaxNICFrequency: {
			OpcodeNumber: 3,
			Offset:       14,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricMaxNoCFrequency: {
			OpcodeNumber: 3,
			Offset:       16,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricCurrentTPCFrequency: {
			OpcodeNumber: 3,
			Offset:       18,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricCurrentMMEFrequency: {
			OpcodeNumber: 3,
			Offset:       20,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricCurrentDMAFrequency: {
			OpcodeNumber: 3,
			Offset:       22,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricCurrentMediaFrequency: {
			OpcodeNumber: 3,
			Offset:       24,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricCurrentPCIeFrequency: {
			OpcodeNumber: 3,
			Offset:       26,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricCurrentARMFrequency: {
			OpcodeNumber: 3,
			Offset:       28,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricCurrentNICFrequency: {
			OpcodeNumber: 3,
			Offset:       30,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricCurrentNoCFrequency: {
			OpcodeNumber: 3,
			Offset:       32,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricCurrentSRAMFrequency: {
			OpcodeNumber: 3,
			Offset:       34,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.FrequencyMetricCurrentMSSFrequency: {
			OpcodeNumber: 3,
			Offset:       36,
			Length:       2,
			ExpectedType: Int,
		},
	}
)
