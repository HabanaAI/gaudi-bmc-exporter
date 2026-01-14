package rasmonitoringapi

import (
	"context"
	"encoding/json"
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (c *Client) CTemperature(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	var metrics []rasmonitoring.Metric
	for oam := 0; oam < c.Oams; oam++ {
		temp, err := c.getCTemperature(ctx, oam)
		if err != nil {
			log.WithField("metric", rasmonitoring.PrefixCTemperature).WithError(err).Error("failed getting metric")
			continue
		}

		metrics = append(metrics, rasmonitoring.Metric{
			MetricName:  rasmonitoring.PrefixCTemperature,
			Hostname:    c.Hostname,
			Oam:         fmt.Sprintf("%d", oam),
			MetricValue: temp,
		})
	}

	return metrics
}

type CTemperatureResponse struct {
	Temperature int `json:"temperature"`
}

// getCTemperature will return the temperature.
func (c *Client) getCTemperature(ctx context.Context, oam int) (int, error) {

	url := fmt.Sprintf("https://%s/ext/ras/direct/temperature?oam=%d", c.Hostname, oam)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return 0, err
	}

	var respBody CTemperatureResponse

	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal: %w, %s", err, string(body))
	}

	return respBody.Temperature, nil
}
