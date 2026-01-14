package rasmonitoringapi

import (
	"context"
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	laneWriteOpcode = Opcode{
		OpcodeNumber: 6,
		Offset:       207,
		Length:       4,
	}

	laneInfoOpcodes = map[string]Opcode{
		rasmonitoring.LaneInfoMetricEBUFOverflow: {
			OpcodeNumber: 6,
			Offset:       211,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.LaneInfoMetricEBUFUnderRun: {
			OpcodeNumber: 6,
			Offset:       215,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.LaneInfoMetricDecodeError: {
			OpcodeNumber: 6,
			Offset:       219,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.LaneInfoMetricRunningDisplayError: {
			OpcodeNumber: 6,
			Offset:       223,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.LaneInfoMetricSKPOSParityError: {
			OpcodeNumber: 6,
			Offset:       227,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.LaneInfoMetricSYNCHeaderError: {
			OpcodeNumber: 6,
			Offset:       231,
			Length:       4,
			ExpectedType: Int,
		},
		rasmonitoring.LaneInfoMetricDeskewError: {
			OpcodeNumber: 6,
			Offset:       235,
			Length:       4,
			ExpectedType: Int,
		},
	}
)

func (c *Client) LaneInfo(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {

	var metrics []rasmonitoring.Metric

	type resChan struct {
		metrics []rasmonitoring.Metric
	}

	ch := make(chan resChan)
	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {

			var oamMetrics []rasmonitoring.Metric

			for lane := 0; lane < c.Lanes; lane++ {

				ll := log.WithFields(logrus.Fields{
					"oam":  oam,
					"lane": lane,
				})

				// write the lane number so we can read information about it
				err := c.write(ctx, oam, laneWriteOpcode, lane)
				if err != nil {
					ll.WithError(err).Error("writing lane info")
					continue
				}

				// Read the lane information for the provided lane
				for fieldName, op := range laneInfoOpcodes {

					val, err := c.decodeOpcode(ctx, op, oam, fieldName, methodIndirect)
					if err != nil {
						ll.WithError(err).Errorf("getting %s value", fieldName)
						continue
					}

					value, err := strconv.Atoi(val)
					if err != nil {
						ll.WithError(err).Errorf("converting %s value %s", fieldName, val)
						continue
					}

					oamMetrics = append(oamMetrics, rasmonitoring.Metric{
						MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixLaneInfo, fieldName),
						Hostname:    c.Hostname,
						Oam:         fmt.Sprintf("%d", oam),
						MetricValue: value,
						CustomLabels: map[string]string{
							"lane": fmt.Sprintf("%d", lane),
						},
					})

				}
			}

			ch <- resChan{
				metrics: oamMetrics,
			}

		}(oam)
	}

	for oam := 0; oam < c.Oams; oam++ {
		results := <-ch

		metrics = append(metrics, results.metrics...)
	}

	close(ch)

	return metrics
}
