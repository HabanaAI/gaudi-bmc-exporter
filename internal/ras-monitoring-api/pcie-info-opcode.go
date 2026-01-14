package rasmonitoringapi

import (
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strconv"

	"golang.org/x/exp/slices"
)

var (
	pciInfoOpcodes = map[string]Opcode{
		rasmonitoring.PcieInfoMetricMaxPCIeLinkSpeed: {
			OpcodeNumber: 6,
			Offset:       0,
			Length:       1,
			ExpectedType: Int,
			Decoder:      decodePcieLinkSpeed,
		},
		rasmonitoring.PcieInfoMetricCurrentPCIeLinkSpeed: {
			OpcodeNumber: 6,
			Offset:       1,
			Length:       1,
			ExpectedType: Int,
			Decoder:      decodePcieLinkSpeed,
		},
		rasmonitoring.PcieInfoMetricMaxPCIeLinkWidth: {
			OpcodeNumber: 6,
			Offset:       2,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricCurrentPCIeLinkWidth: {
			OpcodeNumber: 6,
			Offset:       4,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricPCIeDeviceID: {
			OpcodeNumber: 6,
			Offset:       6,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricPCIeSubsystemID: {
			OpcodeNumber: 6,
			Offset:       8,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricPCIeSubsystemVendorID: {
			OpcodeNumber: 6,
			Offset:       10,
			Length:       2,
			ExpectedType: ReverseString,
		},
		rasmonitoring.PcieInfoMetricPCIeBusAndDevice: {
			OpcodeNumber: 6,
			Offset:       12,
			Length:       2,
			ExpectedType: ReverseString,
		},
		rasmonitoring.PcieInfoMetricCorrectedInternalErrorStatus: {
			OpcodeNumber: 6,
			Offset:       22,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricReplayBufferNumRolloverError: {
			OpcodeNumber: 6,
			Offset:       26,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricReplayTimerTimeoutError: {
			OpcodeNumber: 6,
			Offset:       30,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricBadTLPCounter: {
			OpcodeNumber: 6,
			Offset:       34,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricBadDLLPCounter: {
			OpcodeNumber: 6,
			Offset:       38,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricReceiverErrorCounter: {
			OpcodeNumber: 6,
			Offset:       42,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricLCRCErrorCounter: {
			OpcodeNumber: 6,
			Offset:       46,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricECRCErrorCounter: {
			OpcodeNumber: 6,
			Offset:       50,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricCompletionTimeoutIndication: {
			OpcodeNumber: 6,
			Offset:       54,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricUncorrectableInternalErrorIndication: {
			OpcodeNumber: 6,
			Offset:       55,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricReceiverOverflowIndication: {
			OpcodeNumber: 6,
			Offset:       56,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricFlowControlProtocolErrorIndication: {
			OpcodeNumber: 6,
			Offset:       57,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricSurpriseLinkDownIndication: {
			OpcodeNumber: 6,
			Offset:       58,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricMalfunctionTLPErrorIndication: {
			OpcodeNumber: 6,
			Offset:       59,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricDLLPProtocolErrorIndication: {
			OpcodeNumber: 6,
			Offset:       60,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricRXNakDLLPCounter: {
			OpcodeNumber: 6,
			Offset:       61,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricTxNakDLLPCounter: {
			OpcodeNumber: 6,
			Offset:       65,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricRetryTLPcounter: {
			OpcodeNumber: 6,
			Offset:       69,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricPWRBRKindication: {
			OpcodeNumber: 6,
			Offset:       89,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricPCIeRXMemoryWriteCounter: {
			OpcodeNumber: 6,
			Offset:       90,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricPCIeRXMemoryReadCounter: {
			OpcodeNumber: 6,
			Offset:       94,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricPCIeTXMemoryWriteCounter: {
			OpcodeNumber: 6,
			Offset:       98,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricPCIeTXMemoryReadCounter: {
			OpcodeNumber: 6,
			Offset:       102,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricAERCapabilityControlOffset: {
			OpcodeNumber: 6,
			Offset:       122,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricAERerrorlog: {
			OpcodeNumber: 6,
			Offset:       126,
			Length:       16,
			ExpectedType: Int,
		},
		rasmonitoring.PcieInfoMetricPCIeFWversion: {
			OpcodeNumber: 6,
			Offset:       195,
			Length:       4,
			ExpectedType: Int,
		},
	}
)

// we receive number (1-6) representing the gen and we need to return string.
func decodePcieLinkSpeed(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {

	if !slices.Contains([]string{"1", "2", "3", "4", "5", "6"}, val) {
		return nil, fmt.Errorf("unexpected value: %s", val)
	}

	gen := fmt.Sprintf("Gen%s", val)

	value, err := strconv.Atoi(val)
	if err != nil {
		return nil, err
	}

	return []rasmonitoring.Metric{
		{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				fieldName: gen,
			},
		},
	}, nil
}
