package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) SensorVoltageMonitor(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {

			ll := log.WithField("oam", oam)

			resp, err := c.sensorVoltageMonitor(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Error("failed to get sensor voltage monitor")
				ch <- nil
				return
			}

			ch <- c.sensorVoltageMonitorMetrics(ll, resp, oam)
		}(oam)

	}

	for oam := 0; oam < c.Oams; oam++ {
		alertMetrics := <-ch
		if alertMetrics != nil {
			metrics = append(metrics, alertMetrics...)
		}
	}
	return metrics
}

func (c *Client) sensorVoltageMonitor(ctx context.Context, log *logrus.Entry, oam int) (SensorVoltageMonitorResp, error) {
	var res SensorVoltageMonitorResponse

	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/SensorsVoltageMonitor", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return SensorVoltageMonitorResp{}, err
	}

	return res.Response, nil

}

func (c *Client) sensorVoltageMonitorMetrics(log *logrus.Entry, resp SensorVoltageMonitorResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixSensorVoltageMonitor

	// SwVm

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwVm,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwVm),
	})

	// SwCpeEu0
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwCpeEu0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwCpeEu_0),
	})

	// SwCpeEu1
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwCpeEu1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwCpeEu_1),
	})

	// SwCpeEu2
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwCpeEu2,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwCpeEu_2),
	})

	// SwcpeEu3
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwcpeEu3,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwcpeEu_3),
	})

	// SwCpeHbm
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwCpeHbm,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwCpeHbm),
	})

	// SwTpcSb
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwTpcSb,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwTpcSb),
	})

	// SwTft
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwTft,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwTft),
	})

	// SwCnt
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwCnt,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwCnt),
	})

	// SwMcHbm00
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwMcHbm00,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwMcHbm0_0),
	})

	// SwMcHbm01
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwMcHbm01,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwMcHbm0_1),
	})

	// SwHconHbm0
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwHconHbm0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwHconHbm0),
	})

	// SwPpw0
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwPpw0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwPpw_0),
	})

	// SwPpw1
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwPpw1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwPpw_1),
	})

	// SwL2cMacro
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwL2cMacro,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwL2cMacro),
	})

	// SwVcdMacro
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SwVcdMacro,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSwVcdMacro),
	})

	// SeVm
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeVm,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeVm),
	})

	// SePe00
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SePe00,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSePe0_0),
	})

	// SePe01
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SePe01,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSePe0_1),
	})

	// SeEuCore
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeEuCore,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeEuCore),
	})

	// SeBpyramid
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeBpyramid,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeBpyramid),
	})

	// SeRx
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeRx,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeRx),
	})

	// SeMx
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeMx,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeMx),
	})

	// SeHconHbm1
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeHconHbm1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeHconHbm1),
	})

	// SeMcHbm1
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeMcHbm1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeMcHbm1),
	})

	// SeGasket
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeGasket,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeGasket),
	})

	// SeMmeCtrl
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeMmeCtrl,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeMmeCtrl),
	})

	// SeMmeQman
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeMmeQman,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeMmeQman),
	})

	// SeSbte
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeSbte,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeSbte),
	})

	// SeRtrDn
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeRtrDn,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeRtrDn),
	})

	// SeRtrUp
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeRtrUp,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeRtrUp),
	})

	// SeSram
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.SeSram,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorSeSram),
	})

	// NwVm
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwVm,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwVm),
	})

	// NwPe00
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwPe00,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwPe0_0),
	})

	// NwPe01
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwPe01,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwPe0_1),
	})

	// NwPe02
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwPe02,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwPe0_2),
	})

	// NwPe03
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwPe03,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwPe0_3),
	})

	// NwEuCore
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwEuCore,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwEuCore),
	})

	// NwBpyramid
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwBpyramid,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwBpyramid),
	})

	// NwMcHbm40
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwMcHbm40,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwMcHbm4_0),
	})

	// NwMcHbm41
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwMcHbm41,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwMcHbm4_1),
	})

	// NwHconHbm4
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwHconHbm4,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwHconHbm4),
	})

	// NwAcc
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwAcc,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwAcc),
	})

	// NwWap
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwWap,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwWap),
	})

	// NwTif
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwTif,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwTif),
	})

	// NwRtrDn
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwRtrDn,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwRtrDn),
	})

	// NwRtrUp
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwRtrUp,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwRtrUp),
	})

	// NwSram
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NwSram,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNwSram),
	})

	// NeVM
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeVM,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeVM),
	})

	// NeTx
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeTx,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeTx),
	})

	// NeMx0
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeMx0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeMx_0),
	})

	// NeMx1
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeMx1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeMx_1),
	})

	// NeHconHbm5
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeHconHbm5,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeHconHbm5),
	})

	// NeMcHbm5
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeMcHbm5,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeMcHbm5),
	})

	// NeSob
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeSob,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeSob),
	})

	// NeQnt
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeQnt,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeQnt),
	})

	// NeCnt0
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeCnt0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeCnt_0),
	})

	// NeCnt1
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeCnt1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeCnt_1),
	})

	// NeCpeEu
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeCpeEu,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeCpeEu),
	})

	// NeCntHbm
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeCntHbm,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeCntHbm),
	})

	// NeSbte0
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeSbte0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeSbte_0),
	})

	// NeSbte1
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeSbte1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeSbte_1),
	})

	// NeMif
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeMif,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeMif),
	})

	// NeRtrDn
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.NeRtrDn,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageMonitorNeRtrDn),
	})

	return metrics
}
