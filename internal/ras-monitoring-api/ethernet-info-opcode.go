package rasmonitoringapi

import (
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
)

var (
	ethernetInfoOpcodes = map[string]Opcode{
		rasmonitoring.EthernetInfoMetricSerDesAvailability: {
			OpcodeNumber: 5,
			Offset:       0,
			Length:       1,
			ExpectedType: BinArray,
			Decoder:      decodeSerDesAvailability,
		},
		rasmonitoring.EthernetInfoMetricPortMaxSpeed: {
			OpcodeNumber: 5,
			Offset:       1,
			Length:       1,
			ExpectedType: Int,
			Decoder:      decodePortMaxSpeed,
		},
		rasmonitoring.EthernetInfoMetricANLTStatus: {
			OpcodeNumber: 5,
			Offset:       2,
			Length:       16,
			ExpectedType: BinArray,
			Decoder:      decodeANLTStatus,
		},
		rasmonitoring.EthernetInfoMetricNumberOfLanes: {
			OpcodeNumber: 5,
			Offset:       18,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.EthernetInfoMetricNumberOfLinks: {
			OpcodeNumber: 5,
			Offset:       19,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.EthernetInfoMetricLinkSpeed: {
			OpcodeNumber: 5,
			Offset:       20,
			Length:       2,
			ExpectedType: Int,
			Decoder:      decodeLinkSpeed,
		},
	}
)

func decodeLinkSpeed(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var value int
	switch val {
	case "0":
		value = 25
	case "1":
		value = 56
	case "2":
		value = 112
	default:
		return nil, fmt.Errorf("unsupported link speed: %s", val)
	}
	return []rasmonitoring.Metric{
		{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
		},
	}, nil
}

func decodePortMaxSpeed(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var value int

	switch val {
	case "1":
		value = 50
	case "2":
		value = 100
	case "3":
		value = 200
	case "4":
		value = 400
	default:
		return nil, fmt.Errorf("unsupported max speed: %s", val)
	}

	return []rasmonitoring.Metric{
		{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
		},
	}, nil
}

// decodeANLTStatus will return information per port.
func decodeANLTStatus(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {

	var metrics []rasmonitoring.Metric
	for port, portState := range val {

		if portState != '0' && portState != '1' {
			return nil, fmt.Errorf("unsupported port state: %c", portState)
		}

		value := 1
		state := "Enabled"
		if portState == '0' {
			value = 0
			state = "Disabled"
		}

		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				"state": state,
				"port":  fmt.Sprintf("%d", port),
			},
		})
	}

	return metrics, nil
}

func decodeSerDesAvailability(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {

	if val != "0" && val != "1" {
		return nil, fmt.Errorf("unsupported state: %s", val)
	}

	state := "Available"
	value := 1
	if val == "0" {
		value = 0
		state = "Unavailable"
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
