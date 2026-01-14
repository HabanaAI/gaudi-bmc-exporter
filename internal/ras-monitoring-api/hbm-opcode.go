package rasmonitoringapi

import (
	"context"
	"encoding/json"
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"net/http"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	numOfHbmIndex = 6
)

func (c *Client) Hbm(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	metrics := c.getEccErrors(ctx, log)

	metrics = append(metrics, c.getRepairedLanes(ctx, log)...)

	metrics = append(metrics, c.getReplacedRows(ctx, log)...)

	metrics = append(metrics, c.getHbmRepairStatusArray(ctx, log, numOfHbmIndex)...)

	return metrics
}

var (
	hbmWriteOp = Opcode{
		OpcodeNumber: 19,
		Offset:       104,
		Length:       1,
	}

	hbmMbistRepairOpcode = Opcode{
		OpcodeNumber: 19,
		Offset:       105,
		Length:       1,
	}
	hbmcGlobalEccOpcode = Opcode{
		OpcodeNumber: 19,
		Offset:       106,
		Length:       1,
	}

	numOfLanesOpcode = Opcode{
		OpcodeNumber: 19,
		Offset:       0,
		Length:       2,
		ExpectedType: Int,
	}

	numOfReplacedRowsOp = Opcode{
		OpcodeNumber: 19,
		Offset:       16,
		Length:       2,
		ExpectedType: Int,
	}
)

func (c *Client) getHbmRepairStatusArray(ctx context.Context, log *logrus.Entry, numOfHbmIndex int) []rasmonitoring.Metric {

	var metrics []rasmonitoring.Metric
	for oam := 0; oam < c.Oams; oam++ {

		for hbmIndex := 0; hbmIndex < numOfHbmIndex; hbmIndex++ {

			// first write so we can read.
			err := c.write(ctx, oam, hbmWriteOp, hbmIndex)
			if err != nil {
				log.WithError(err).Error()
				continue
			}

			mbistRepairData, err := c.opcodeData(ctx, hbmMbistRepairOpcode, oam, methodIndirect)
			if err != nil {
				log.WithError(err).Error()
				continue
			}

			if len(mbistRepairData.Data) == 0 {
				log.WithError(fmt.Errorf("mbistRepairData data is empty")).Error()
				continue
			}

			state := "Flow did not run"
			if mbistRepairData.Data[0] == 1 {
				state = "Flow ran"
			}

			metrics = append(metrics, rasmonitoring.Metric{
				MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricMbistRepair),
				Hostname:    c.Hostname,
				Oam:         fmt.Sprintf("%d", oam),
				MetricValue: int(mbistRepairData.Data[0]),
				CustomLabels: map[string]string{
					rasmonitoring.HBMMetricMbistRepairLabelState: state,
					rasmonitoring.HBMMetricMbistRepairLabelIndex: fmt.Sprintf("%d", hbmIndex),
				},
			})

			globalEccData, err := c.opcodeData(ctx, hbmcGlobalEccOpcode, oam, methodIndirect)
			if err != nil {
				log.WithError(err).Error()
				return nil
			}

			if len(globalEccData.Data) == 0 {
				log.WithError(fmt.Errorf("globalEccData data is empty")).Error()
				continue
			}

			metrics = append(metrics, rasmonitoring.Metric{
				MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricGlobalECC),
				Hostname:    c.Hostname,
				Oam:         fmt.Sprintf("%d", oam),
				MetricValue: int(globalEccData.Data[0]),
				CustomLabels: map[string]string{
					rasmonitoring.HBMMetricGlobalECCLabelIndex: fmt.Sprintf("%d", hbmIndex),
				},
			})
		}
	}

	return metrics

}

func (c *Client) getRepairedLanes(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric
	for oam := 0; oam < c.Oams; oam++ {

		val, err := c.decodeOpcode(ctx, numOfLanesOpcode, oam, "hbm_repaired_lanes", methodIndirect)
		if err != nil {
			log.WithError(err).Error("failed getting repaired lanes data")
			continue
		}

		value, err := strconv.Atoi(val)
		if err != nil {
			log.WithError(err).Error()
			continue
		}

		// ADD the number of repaired Lanes.
		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricNumOfRepairedLanes),
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
		})

		numLanes, err := strconv.Atoi(val)
		if err != nil {
			log.WithError(err).Error("repaired lanes failed parsing num of lanes to int")
			continue
		}

		// read all the repaired lanes.
		for index := 0; index < numLanes; index++ {
			m, err := c.getRepairLane(ctx, oam, index)
			if err != nil {
				log.WithField("metric", rasmonitoring.PrefixHbm).WithError(err).Error(fmt.Sprintf("failed getting repair lane, index %d", index))
				continue
			}
			metrics = append(metrics, m)
		}
	}

	return metrics
}

type RepairLaneResponse struct {
	LaneHBMIndex int `json:"lane_HBM_index"`
	LaneCHIndex  int `json:"lane_ch_index"`
}

// getRepairLane will return 1 instance in the array of repaired lanes.
func (c *Client) getRepairLane(ctx context.Context, oam, index int) (rasmonitoring.Metric, error) {

	u, err := url.Parse(fmt.Sprintf("https://%s/ext/ras/indirect/dram_repaired_lane", c.Hostname))
	if err != nil {
		return rasmonitoring.Metric{}, err
	}

	q := u.Query()
	q.Add("oam", fmt.Sprintf("%d", oam))
	q.Add("index", fmt.Sprintf("%d", index))

	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return rasmonitoring.Metric{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return rasmonitoring.Metric{}, err
	}

	var respBody RepairLaneResponse

	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return rasmonitoring.Metric{}, err
	}

	return rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: 0,
		MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricRepairedLanes),
		CustomLabels: map[string]string{
			rasmonitoring.HBMMetricRepairedLanesLabelHBMIndex:  fmt.Sprintf("%d", respBody.LaneHBMIndex),
			rasmonitoring.HBMMetricRepairedLanesLabelMCChannel: fmt.Sprintf("%d", respBody.LaneCHIndex),
		},
	}, nil
}

func (c *Client) getReplacedRows(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric
	for oam := 0; oam < c.Oams; oam++ {

		val, err := c.decodeOpcode(ctx, numOfReplacedRowsOp, oam, "num_replaced_rows", methodIndirect)
		if err != nil {
			log.WithError(err).Error("failed getting replaced rows data")
			continue
		}

		value, err := strconv.Atoi(val)
		if err != nil {
			log.WithError(err).Error()
			continue
		}

		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricNumOfReplacedRows),
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
		})

		numRows, err := strconv.Atoi(val)
		if err != nil {
			log.WithError(err).Error("replaced rows failed parsing data to int")
			continue
		}

		// read all the replaced rows.
		for index := 0; index < numRows; index++ {
			m, err := c.getReplacedRow(ctx, oam, index)
			if err != nil {
				log.WithField("metric", "hbm").WithError(err).Error(fmt.Sprintf("failed getting replace row, index %d", index))
				continue
			}
			metrics = append(metrics, m)
		}
	}

	return metrics
}

type replacedRowResponse struct {
	HBMIndex           int `json:"hbm_index"`
	SidIndex           int `json:"sid_index"`
	PCIndex            int `json:"pc_index"`
	BankIndex          int `json:"bank_index"`
	Cause              int `json:"cause"`
	ReplacedRowAddress int `json:"replaced_row_address"`
}

func (c *Client) getReplacedRow(ctx context.Context, oam, index int) (rasmonitoring.Metric, error) {

	u, err := url.Parse(fmt.Sprintf("https://%s/ext/ras/indirect/dram_replaced_row", c.Hostname))
	if err != nil {
		return rasmonitoring.Metric{}, err
	}

	q := u.Query()
	q.Add("oam", fmt.Sprintf("%d", oam))
	q.Add("index", fmt.Sprintf("%d", index))

	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return rasmonitoring.Metric{}, err
	}
	body, err := c.doRequest(req)
	if err != nil {
		return rasmonitoring.Metric{}, err
	}

	var respBody replacedRowResponse

	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return rasmonitoring.Metric{}, err
	}

	return rasmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: 0,
		MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricReplaceRows),
		CustomLabels: map[string]string{
			rasmonitoring.HBMMetricReplaceRowsLabelHBMIndex:   fmt.Sprintf("%d", respBody.HBMIndex),
			rasmonitoring.HBMMetricReplaceRowsLabelPCIndex:    fmt.Sprintf("%d", respBody.SidIndex),
			rasmonitoring.HBMMetricReplaceRowsLabelStackID:    fmt.Sprintf("%d", respBody.PCIndex),
			rasmonitoring.HBMMetricReplaceRowsLabelBankIndex:  fmt.Sprintf("%d", respBody.BankIndex),
			rasmonitoring.HBMMetricReplaceRowsLabelCause:      fmt.Sprintf("%d", respBody.Cause),
			rasmonitoring.HBMMetricReplaceRowsLabelRowAddress: fmt.Sprintf("%d", respBody.ReplacedRowAddress),
		},
	}, nil
}

func (c *Client) getEccErrors(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric

	for oam := 0; oam < c.Oams; oam++ {

		eccErrorsOpcode := Opcode{
			OpcodeNumber: 19,
			Offset:       10,
			Length:       2,
			ExpectedType: Int,
		}

		val, err := c.decodeOpcode(ctx, eccErrorsOpcode, oam, rasmonitoring.HBMMetricEccErrors, methodIndirect)
		if err != nil {
			log.WithError(err).Error("failed getting ecc errors")
			return nil
		}

		value, err := strconv.Atoi(val)
		if err != nil {
			log.WithError(err).Error()
			return nil
		}
		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixHbm, rasmonitoring.HBMMetricEccErrors),
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: value,
		})
	}

	return metrics
}
