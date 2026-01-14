package rasmonitoringapi

import (
	"context"
	"encoding/json"
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"net/http"
	"net/url"
	"sync"

	"github.com/sirupsen/logrus"
)

type voltageMonitorInfo struct {
	Index  int
	Sensor int
}

var (
	voltageMonitorData = map[string]voltageMonitorInfo{
		rasmonitoring.SensorVoltageMonitorSwVm: {
			Index:  0,
			Sensor: 0,
		},
		rasmonitoring.SensorVoltageMonitorSwCpeEu_0: {
			Index:  0,
			Sensor: 1,
		},
		rasmonitoring.SensorVoltageMonitorSwCpeEu_1: {
			Index:  0,
			Sensor: 2,
		},
		rasmonitoring.SensorVoltageMonitorSwCpeEu_2: {
			Index:  0,
			Sensor: 3,
		},
		rasmonitoring.SensorVoltageMonitorSwcpeEu_3: {
			Index:  0,
			Sensor: 4,
		},
		rasmonitoring.SensorVoltageMonitorSwCpeHbm: {
			Index:  0,
			Sensor: 5,
		},
		rasmonitoring.SensorVoltageMonitorSwTpcSb: {
			Index:  0,
			Sensor: 6,
		},
		rasmonitoring.SensorVoltageMonitorSwTft: {
			Index:  0,
			Sensor: 7,
		},
		rasmonitoring.SensorVoltageMonitorSwCnt: {
			Index:  0,
			Sensor: 8,
		},
		rasmonitoring.SensorVoltageMonitorSwMcHbm0_0: {
			Index:  0,
			Sensor: 9,
		},
		rasmonitoring.SensorVoltageMonitorSwMcHbm0_1: {
			Index:  0,
			Sensor: 10,
		},
		rasmonitoring.SensorVoltageMonitorSwHconHbm0: {
			Index:  0,
			Sensor: 11,
		},
		rasmonitoring.SensorVoltageMonitorSwPpw_0: {
			Index:  0,
			Sensor: 12,
		},
		rasmonitoring.SensorVoltageMonitorSwPpw_1: {
			Index:  0,
			Sensor: 13,
		},
		rasmonitoring.SensorVoltageMonitorSwL2cMacro: {
			Index:  0,
			Sensor: 14,
		},
		rasmonitoring.SensorVoltageMonitorSwVcdMacro: {
			Index:  0,
			Sensor: 15,
		},
		rasmonitoring.SensorVoltageMonitorSeVm: {
			Index:  1,
			Sensor: 0,
		},
		rasmonitoring.SensorVoltageMonitorSePe0_0: {
			Index:  1,
			Sensor: 1,
		},
		rasmonitoring.SensorVoltageMonitorSePe0_1: {
			Index:  1,
			Sensor: 2,
		},
		rasmonitoring.SensorVoltageMonitorSeEuCore: {
			Index:  1,
			Sensor: 3,
		},
		rasmonitoring.SensorVoltageMonitorSeBpyramid: {
			Index:  1,
			Sensor: 4,
		},
		rasmonitoring.SensorVoltageMonitorSeRx: {
			Index:  1,
			Sensor: 5,
		},
		rasmonitoring.SensorVoltageMonitorSeMx: {
			Index:  1,
			Sensor: 6,
		},
		rasmonitoring.SensorVoltageMonitorSeHconHbm1: {
			Index:  1,
			Sensor: 7,
		},
		rasmonitoring.SensorVoltageMonitorSeMcHbm1: {
			Index:  1,
			Sensor: 8,
		},
		rasmonitoring.SensorVoltageMonitorSeGasket: {
			Index:  1,
			Sensor: 9,
		},
		rasmonitoring.SensorVoltageMonitorSeMmeCtrl: {
			Index:  1,
			Sensor: 10,
		},
		rasmonitoring.SensorVoltageMonitorSeMmeQman: {
			Index:  1,
			Sensor: 11,
		},
		rasmonitoring.SensorVoltageMonitorSeSbte: {
			Index:  1,
			Sensor: 12,
		},
		rasmonitoring.SensorVoltageMonitorSeRtrDn: {
			Index:  1,
			Sensor: 13,
		},
		rasmonitoring.SensorVoltageMonitorSeRtrUp: {
			Index:  1,
			Sensor: 14,
		},
		rasmonitoring.SensorVoltageMonitorSeSram: {
			Index:  1,
			Sensor: 15,
		},
		rasmonitoring.SensorVoltageMonitorNwVm: {
			Index:  2,
			Sensor: 0,
		},
		rasmonitoring.SensorVoltageMonitorNwPe0_0: {
			Index:  2,
			Sensor: 1,
		},
		rasmonitoring.SensorVoltageMonitorNwPe0_1: {
			Index:  2,
			Sensor: 2,
		},
		rasmonitoring.SensorVoltageMonitorNwPe0_2: {
			Index:  2,
			Sensor: 3,
		},
		rasmonitoring.SensorVoltageMonitorNwPe0_3: {
			Index:  2,
			Sensor: 4,
		},
		rasmonitoring.SensorVoltageMonitorNwEuCore: {
			Index:  2,
			Sensor: 5,
		},
		rasmonitoring.SensorVoltageMonitorNwBpyramid: {
			Index:  2,
			Sensor: 6,
		},
		rasmonitoring.SensorVoltageMonitorNwMcHbm4_0: {
			Index:  2,
			Sensor: 7,
		},
		rasmonitoring.SensorVoltageMonitorNwMcHbm4_1: {
			Index:  2,
			Sensor: 8,
		},
		rasmonitoring.SensorVoltageMonitorNwHconHbm4: {
			Index:  2,
			Sensor: 9,
		},
		rasmonitoring.SensorVoltageMonitorNwAcc: {
			Index:  2,
			Sensor: 10,
		},
		rasmonitoring.SensorVoltageMonitorNwWap: {
			Index:  2,
			Sensor: 11,
		},
		rasmonitoring.SensorVoltageMonitorNwTif: {
			Index:  2,
			Sensor: 12,
		},
		rasmonitoring.SensorVoltageMonitorNwRtrDn: {
			Index:  2,
			Sensor: 13,
		},
		rasmonitoring.SensorVoltageMonitorNwRtrUp: {
			Index:  2,
			Sensor: 14,
		},
		rasmonitoring.SensorVoltageMonitorNwSram: {
			Index:  2,
			Sensor: 15,
		},
		rasmonitoring.SensorVoltageMonitorNeVM: {
			Index:  3,
			Sensor: 0,
		},
		rasmonitoring.SensorVoltageMonitorNeTx: {
			Index:  3,
			Sensor: 1,
		},
		rasmonitoring.SensorVoltageMonitorNeMx_0: {
			Index:  3,
			Sensor: 2,
		},
		rasmonitoring.SensorVoltageMonitorNeMx_1: {
			Index:  3,
			Sensor: 3,
		},
		rasmonitoring.SensorVoltageMonitorNeHconHbm5: {
			Index:  3,
			Sensor: 4,
		},
		rasmonitoring.SensorVoltageMonitorNeMcHbm5: {
			Index:  3,
			Sensor: 5,
		},
		rasmonitoring.SensorVoltageMonitorNeSob: {
			Index:  3,
			Sensor: 6,
		},
		rasmonitoring.SensorVoltageMonitorNeQnt: {
			Index:  3,
			Sensor: 7,
		},
		rasmonitoring.SensorVoltageMonitorNeCnt_0: {
			Index:  3,
			Sensor: 8,
		},
		rasmonitoring.SensorVoltageMonitorNeCnt_1: {
			Index:  3,
			Sensor: 9,
		},
		rasmonitoring.SensorVoltageMonitorNeCpeEu: {
			Index:  3,
			Sensor: 10,
		},
		rasmonitoring.SensorVoltageMonitorNeCntHbm: {
			Index:  3,
			Sensor: 11,
		},
		rasmonitoring.SensorVoltageMonitorNeSbte_0: {
			Index:  3,
			Sensor: 12,
		},
		rasmonitoring.SensorVoltageMonitorNeSbte_1: {
			Index:  3,
			Sensor: 13,
		},
		rasmonitoring.SensorVoltageMonitorNeMif: {
			Index:  3,
			Sensor: 14,
		},
		rasmonitoring.SensorVoltageMonitorNeRtrDn: {
			Index:  3,
			Sensor: 15,
		},
	}
)

func (c *Client) SensorVoltageMonitor(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {

	var metrics []rasmonitoring.Metric

	ch := make(chan result)
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for oam := 0; oam < c.Oams; oam++ {
			wg.Add(1)

			go func(oam int) {
				defer wg.Done()

				for fieldName, v := range voltageMonitorData {

					val, err := c.getVoltageMonitor(ctx, v, oam)
					if err != nil {
						ch <- result{err: fmt.Errorf("failed getting voltage monitoring data: %w", err)}
						continue
					}

					ch <- result{
						metric: rasmonitoring.Metric{
							MetricName:  fmt.Sprintf("%s_%s", rasmonitoring.PrefixSensorVoltageMonitor, fieldName),
							Hostname:    c.Hostname,
							Oam:         fmt.Sprintf("%d", oam),
							MetricValue: val,
						},
					}
				}
			}(oam)

		}
	}()

	go func() {
		wg.Wait()
		done <- struct{}{}
		close(ch)
	}()

	for {
		select {
		case <-done:
			return metrics
		case data := <-ch:
			if data.err != nil {
				log.WithError(data.err).Error()
				continue
			}
			metrics = append(metrics, data.metric)

		}
	}

}

type monitorVoltageResponse struct {
	Value int `json:"value"`
}

func (c *Client) getVoltageMonitor(ctx context.Context, v voltageMonitorInfo, oam int) (int, error) {

	u, err := url.Parse(fmt.Sprintf("https://%s/ext/ras/indirect/sensor_voltage_mon", c.Hostname))
	if err != nil {
		return 0, err
	}

	q := u.Query()
	q.Add("oam", fmt.Sprintf("%d", oam))
	q.Add("sensor", fmt.Sprintf("%d", v.Sensor))
	q.Add("index", fmt.Sprintf("%d", v.Index))

	u.RawQuery = q.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return 0, err
	}

	var respBody monitorVoltageResponse
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return 0, err
	}

	return respBody.Value, nil
}
