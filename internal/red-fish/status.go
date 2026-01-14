package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) Status(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {
		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.status(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Errorf("failed to get status")
				ch <- nil
				return
			}

			ch <- c.statusMetrics(ll, resp, oam)
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

func (c *Client) status(ctx context.Context, log *logrus.Entry, oam int) (StatusResp, error) {
	var res StatusResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/Status", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return StatusResp{}, err
	}

	return res.Response, nil

}

func (c *Client) statusMetrics(log *logrus.Entry, resp StatusResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixStatus
	// Boot stage.

	bootStageVal := -1
	switch strings.TrimSpace(strings.ToLower(resp.BootStage)) {
	case "linux":
		bootStageVal = 0
	case "uboot":
		bootStageVal = 1
	case "preboot":
		bootStageVal = 2
	case "zephyr":
		bootStageVal = 3
	default:
		log.WithError(fmt.Errorf("unexpected boot stage %s", resp.BootStage)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricBootStage),
		MetricValue: bootStageVal,
		CustomLabels: map[string]string{
			bmcmonitoring.StatusMetricBootStage: resp.BootStage,
		},
	})

	// Emergency Power Reduction

	emergencyPowerReductionVal := -1
	switch strings.ToLower(resp.EmergencyPowerReduction) {
	case "normal":
		emergencyPowerReductionVal = 0
	case "reduced":
		emergencyPowerReductionVal = 1
	default:
		log.WithError(fmt.Errorf("unexpected emergency Power Reduction %s", resp.EmergencyPowerReduction)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricEmergencyPowerReduction),
		MetricValue: emergencyPowerReductionVal,
		CustomLabels: map[string]string{
			bmcmonitoring.StatusMetricEmergencyPowerReduction: resp.EmergencyPowerReduction,
		},
	})

	// Clock Throttling

	clockThrottlingVal := -1
	switch strings.ToLower(resp.ClockThrottling) {
	case "none":
		clockThrottlingVal = 0
	case "power":
		clockThrottlingVal = 1
	case "thermal":
		clockThrottlingVal = 2
	default:
		log.WithError(fmt.Errorf("unexpected Clock Throttling %s", resp.ClockThrottling)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricClockThrottling),
		MetricValue: clockThrottlingVal,
		CustomLabels: map[string]string{
			bmcmonitoring.StatusMetricClockThrottling: resp.ClockThrottling,
		},
	})

	// Last Clock Throttling Duration

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricLastClockThrottlingDuration),
		MetricValue: resp.LastClockThrottlingDuration,
	})

	// Power State

	// default value full performance
	powerStateVal := 0
	var err error
	// "Reduced by <number>/16"
	if strings.Contains(resp.PowerState, "Reduced") {
		v := strings.Fields(resp.PowerState)

		// check the format, expecting: "Reduced by 1/16"
		if len(v) == 3 {
			val := strings.Split(v[2], "/")[0]

			// convert the string into an int
			powerStateVal, err = strconv.Atoi(val)
			if err != nil {
				log.WithError(fmt.Errorf("failed in conversion of power state value %s", val)).Error()
				powerStateVal = -1
			}

		} else {
			log.WithError(fmt.Errorf("power state reduced in an unexpected format, wanted len 3, got %d", len(v))).Error()
			powerStateVal = -1
		}
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricPowerState),
		MetricValue: powerStateVal,
		CustomLabels: map[string]string{
			bmcmonitoring.StatusMetricPowerState: resp.PowerState,
		},
	})

	// Total Clock Throttling Duration
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricTotalClockThrottlingDuration),
		MetricValue: resp.TotalClockThrottlingDuration,
	})

	// GlobalTimeFromReset
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricGlobalTimeFromReset),
		MetricValue: resp.GlobalTimeFromReset,
	})

	// ChipStatus

	chipStatusVal := -1
	switch strings.ToLower(resp.ChipStatus) {
	case "processing":
		chipStatusVal = 0
	case "idle":
		chipStatusVal = 1
	default:
		log.WithError(fmt.Errorf("unexpected Chip Status %s", resp.ChipStatus)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricChipStatus),
		MetricValue: chipStatusVal,
		CustomLabels: map[string]string{
			bmcmonitoring.StatusMetricChipStatus: resp.ChipStatus,
		},
	})

	// DeviceActivity

	deviceActivityVal := -1
	switch strings.ToLower(resp.DeviceActivity) {
	case "device not in use", "device not in-use":
		deviceActivityVal = 0
	case "device in-use":
		deviceActivityVal = 1
	default:
		log.WithError(fmt.Errorf("unexpected Device Activity %s", resp.DeviceActivity)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricDeviceActivity),
		MetricValue: deviceActivityVal,
		CustomLabels: map[string]string{
			bmcmonitoring.StatusMetricDeviceActivity: resp.DeviceActivity,
		},
	})

	// DeviceActivityCounter
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricDeviceActivityCounter),
		MetricValue: resp.DeviceActivityCounter,
	})

	// DevicePowerReduction

	devicePowerReductionValue := -1
	switch strings.ToLower(resp.DevicePowerReduction) {
	case "max reduction":
		devicePowerReductionValue = 0
	case "2nd reduction":
		devicePowerReductionValue = 1
	case "1st reduction":
		devicePowerReductionValue = 2
	case "normal power":
		devicePowerReductionValue = 3
	default:
		log.WithError(fmt.Errorf("unexpected Device Power Reduction %s", resp.DevicePowerReduction)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricDevicePowerReduction),
		MetricValue: devicePowerReductionValue,
		CustomLabels: map[string]string{
			"state": resp.DevicePowerReduction,
		},
	})

	// LastPowerReductionDuration
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.StatusMetricLastPowerReductionDuration),
		MetricValue: resp.LastPowerReductionDuration,
	})
	return metrics
}
