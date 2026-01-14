package rasmonitoringapi

import (
	"context"
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	directCodes = map[string]Opcode{
		rasmonitoring.DirectMetricPCIeVendorID: {
			Offset:       9,
			Length:       2,
			ExpectedType: String,
		},
		rasmonitoring.DirectMetricASICSerialNumber: {
			Offset:       11,
			Length:       8,
			ExpectedType: AsciiString,
		},
		rasmonitoring.DirectMetricFWVersionMajor: {
			Offset:       106,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.DirectMetricFWVersionMinor: {
			Offset:       107,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.DirectMetricFWVersionPatch: {
			Offset:       108,
			Length:       1,
			ExpectedType: Int,
		},
		rasmonitoring.DirectMetricCoreVDD: {
			Offset:       133,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.DirectMetricHBMVDDq: {
			Offset:       115,
			Length:       2,
			ExpectedType: Int,
		},
		rasmonitoring.DirectMetricV12: {
			Offset:       117,
			Length:       2,
			ExpectedType: Int,
		},
	}
)

var (
	apiVersionOpcode = 240
	osStateOpcode    = 241
)

func (c *Client) osState(ctx context.Context, log *logrus.Entry, opcode int, hostname, prefix string) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric

	for oam := 0; oam < c.Oams; oam++ {

		val, err := c.directCommand(ctx, "control_command", oam, opcode)
		if err != nil {
			log.WithError(err).Error("failed getting os state data")
			continue
		}

		// turn to hex.
		hexStr := strconv.FormatInt(int64(val), 16)
		labels := []string{rasmonitoring.DirectMetricOSStage, rasmonitoring.DirectMetricIBAccessState, rasmonitoring.DirectMetricOOBAccessState}

		if len(hexStr) != 3 {
			log.WithError(fmt.Errorf("%s os state result in wrong format, expected len %d, got %d", prefix, 3, len(hexStr))).Error("os state unexpected data")
			continue
		}

		for index, v := range hexStr {
			fieldName := labels[index]

			switch fieldName {
			case rasmonitoring.DirectMetricIBAccessState, rasmonitoring.DirectMetricOOBAccessState:
				state := "Full Access"

				var value int
				if v == '1' {
					value = 1
					state = "Restricted Access"
				}

				metrics = append(metrics, rasmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
					Hostname:    hostname,
					Oam:         fmt.Sprintf("%d", oam),
					MetricValue: value,
					CustomLabels: map[string]string{
						"state": state,
					},
				})
			case rasmonitoring.DirectMetricOSStage:
				var state string

				var value int
				switch v {
				case '0':
					state = "LINUX"
				case '1':
					state = "UBOOT"
					value = 1
				case '2':
					state = "PREBOOT"
					value = 2
				case '3':
					state = "ZEPHYR"
					value = 3
				}
				metrics = append(metrics, rasmonitoring.Metric{
					MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
					Hostname:    hostname,
					Oam:         fmt.Sprintf("%d", oam),
					MetricValue: value,
					CustomLabels: map[string]string{
						"stage": state,
					},
				})
			}
		}
	}

	return metrics
}

// apiVersion will return the major, minor, patch.
func (c *Client) apiVersion(ctx context.Context, log *logrus.Entry, opcode int, hostname, prefix string) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric

	for oam := 0; oam < c.Oams; oam++ {

		val, err := c.directCommand(ctx, "control_command", oam, opcode)
		if err != nil {
			log.WithError(err).Error("failed getting api version data")
			continue
		}

		// turn to hex.
		hexStr := strconv.FormatInt(int64(val), 16)

		v, err := strconv.Atoi(hexStr)
		if err != nil {
			log.WithError(err).Error()
			continue
		}
		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", prefix, rasmonitoring.DirectMetricApiVersion),
			Hostname:    hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: v,
		})

	}

	return metrics
}
