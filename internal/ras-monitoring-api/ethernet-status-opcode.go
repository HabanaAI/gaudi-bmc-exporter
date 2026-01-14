package rasmonitoringapi

import (
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
)

var (
	ethernetStatusOpcodes = map[string]Opcode{
		rasmonitoring.EthernetStatusMetricExternalLinkStatus: {
			OpcodeNumber: 5,
			Offset:       86,
			Length:       16,
			ExpectedType: BinArray,
			Decoder:      decodeExternalLinkStatus,
		},
		rasmonitoring.EthernetStatusMetricLinkStatus: {
			OpcodeNumber: 5,
			Offset:       102,
			Length:       16,
			ExpectedType: BinArray,
			Decoder:      decodeLinkStatus,
		},
		rasmonitoring.EthernetStatusMetricPHYStatus: {
			OpcodeNumber: 5,
			Offset:       118,
			Length:       16,
			ExpectedType: BinArray,
			Decoder:      decodePHYStatus,
		},
		rasmonitoring.EthernetStatusMetricPortMapping: {
			OpcodeNumber: 5,
			Offset:       54,
			Length:       16,
			ExpectedType: BinArray,
			Decoder:      decodePortMapping,
		},
	}
)

func decodePortMapping(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var metrics []rasmonitoring.Metric
	for port, portState := range val {

		if portState != '0' && portState != '1' {
			return nil, fmt.Errorf("unsupported port state: %c", portState)
		}

		value := 1
		portType := "External"
		if portState == '0' {
			value = 0
			portType = "Internal"
		}

		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				"type": portType,
				"port": fmt.Sprintf("%d", port),
			},
		})
	}

	return metrics, nil
}

func decodeExternalLinkStatus(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var metrics []rasmonitoring.Metric
	for link, linkState := range val {

		if linkState != '0' && linkState != '1' {
			return nil, fmt.Errorf("unsupported port state: %c", linkState)
		}

		state := "Active Port"
		value := 1
		if linkState == '0' {
			value = 0
			state = "Non-active port"
		}

		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				"state": state,
				"link":  fmt.Sprintf("%d", link),
			},
		})
	}

	return metrics, nil
}

func decodeLinkStatus(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var metrics []rasmonitoring.Metric
	for link, linkState := range val {

		if linkState != '0' && linkState != '1' {
			return nil, fmt.Errorf("unsupported link state: %c", linkState)
		}

		state := "Connected"
		value := 1
		if linkState == '0' {
			value = 0
			state = "Not connected"
		}

		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				"state": state,
				"link":  fmt.Sprintf("%d", link),
			},
		})
	}

	return metrics, nil
}

func decodePHYStatus(val string, oam int, hostname string, fieldName string, prefix string) ([]rasmonitoring.Metric, error) {
	var metrics []rasmonitoring.Metric
	for link, linkState := range val {

		if linkState != '0' && linkState != '1' {
			return nil, fmt.Errorf("unsupported link state: %c", linkState)
		}

		state := "Ready"
		value := 1
		if linkState == '0' {
			value = 0
			state = "Not ready"
		}

		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
			CustomLabels: map[string]string{
				"state": state,
				"phy":   fmt.Sprintf("%d", link),
			},
		})
	}

	return metrics, nil
}
