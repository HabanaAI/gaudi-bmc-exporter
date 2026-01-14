package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"

	"github.com/sirupsen/logrus"
)

func (c *Client) SensorVoltage(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {
		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.sensorVoltage(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Error("failed to get sensor voltage")
				ch <- nil
				return
			}

			ch <- c.sensorVoltageMetrics(ll, resp, oam)
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

func (c *Client) sensorVoltage(ctx context.Context, log *logrus.Entry, oam int) (SensorVoltageResp, error) {
	var res SensorVoltageResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/SensorsVoltage", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return SensorVoltageResp{}, err
	}

	return res.Response, nil

}

func (c *Client) sensorVoltageMetrics(log *logrus.Entry, resp SensorVoltageResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixSensorVoltage

	// VADC54
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.VADC54,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVADC54),
	})

	// Vrm1In
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Vrm1In,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVRM1in),
	})

	// Vrm1Out
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Vrm1Out,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVRM1out),
	})

	// Vrm2In
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Vrm2In,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVRM2in),
	})

	// Vrm2VddOut
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Vrm2VddOut,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVRM2VDDout),
	})

	// Vrm2HbmOut
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Vrm2HbmOut,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVRM2HBMout),
	})

	// VmonPcieVph1P8V
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.VmonPcieVph1P8V,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMONPCIEVPH1P8V),
	})

	// Vmon1P8HbmVaa
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Vmon1P8HbmVaa,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMON1P8HBMVAA),
	})

	// Vmon2P5
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Vmon2P5,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMON2P5),
	})

	// Vmon48VHimon
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Vmon48VHimon,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMON48VHIMON),
	})

	// VmonP5V
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.VmonP5V,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMONP5V),
	})

	// Vmon12V1
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.Vmon12V1,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMON12V1),
	})

	// VmonHbm
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.VmonHbm,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMONHBM),
	})

	// VmonCore
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.VmonCore,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageVMONCore),
	})

	// CpldHimon1P8NIC
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CpldHimon1P8NIC,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SensorVoltageCPLDHIMON1P8NIC),
	})

	return metrics
}
