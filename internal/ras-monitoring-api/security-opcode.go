package rasmonitoringapi

import (
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"log/slog"
	"strconv"
)

var (
	securityOpcodes = map[string]Opcode{

		rasmonitoring.SecurityMetricCurrentPublicKeyHashIndex: {
			OpcodeNumber: 17,
			Offset:       0,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.SecurityMetricCurrentSVNversion: {
			OpcodeNumber: 17,
			Offset:       1,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.SecurityMetricKeyRevocation: {
			OpcodeNumber: 17,
			Offset:       2,
			Length:       1,
			ExpectedType: Bin,
			Decoder:      decodeKeyRevocation,
		},
		rasmonitoring.SecurityMetricMinimalSVNindex: {
			OpcodeNumber: 17,
			Offset:       3,
			Length:       1,
			ExpectedType: ReverseString,
		},
		rasmonitoring.SecurityMetricFWImageSource: {
			OpcodeNumber: 17,
			Offset:       4,
			Length:       1,
			ExpectedType: Int,
			Decoder:      decodeFWImageSource,
		},
		rasmonitoring.SecurityMetricTPMPCRPPBOOT: {
			OpcodeNumber: 17,
			Offset:       13,
			Length:       48,
			ExpectedType: ReverseString,
		},
		rasmonitoring.SecurityMetricTPMPCRPREBOOT: {
			OpcodeNumber: 17,
			Offset:       61,
			Length:       48,
			ExpectedType: ReverseString,
		},
		rasmonitoring.SecurityMetricTPMPCRBOOT: {
			OpcodeNumber: 17,
			Offset:       109,
			Length:       48,
			ExpectedType: ReverseString,
		},
		rasmonitoring.SecurityMetricTPMPCRLINUX: {
			OpcodeNumber: 17,
			Offset:       157,
			Length:       48,
			ExpectedType: ReverseString,
		},
		rasmonitoring.SecurityMetricCPLDVersion: {
			OpcodeNumber: 17,
			Offset:       213,
			Length:       1,
			ExpectedType: ReverseString,
		},
		rasmonitoring.SecurityMetricCPLDVersionTimestamp: {
			OpcodeNumber: 17,
			Offset:       214,
			Length:       4,
			ExpectedType: ReverseString,
		},
	}
)

func decodeFWImageSource(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var s string
	switch val {
	case "0":
		s = "Primary"
	case "1":
		s = "Secondary"
	default:
		return nil, fmt.Errorf("unsupported image source: %s", val)
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
				fieldName: s,
			},
		},
	}, nil
}

// decodeKeyRevocation we get a 8 bytes bit number and send metric for each.
func decodeKeyRevocation(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {

	var metrics []rasmonitoring.Metric

	keys := []string{
		rasmonitoring.SecurityMetricKey0Revocation,
		rasmonitoring.SecurityMetricKey1Revocation,
		rasmonitoring.SecurityMetricKey2Revocation,
		rasmonitoring.SecurityMetricKey3Revocation,
		rasmonitoring.SecurityMetricKey4Revocation,
	}

	for i, v := range val {

		if v != '1' && v != '0' {
			return nil, fmt.Errorf("unsupported state: %c", v)
		}

		// there is 8 values but we only address keys 0-4
		if i == len(keys) {
			break
		}

		state := "Not Revoked"
		value := 0
		if v == '1' {
			value = 1
			state = "Revoked"
		}

		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", prefix, keys[i]),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				keys[i]: state,
			},
		})

	}
	return metrics, nil
}
