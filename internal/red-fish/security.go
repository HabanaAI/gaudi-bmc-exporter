package redfish

import (
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) Security(ctx context.Context, log *logrus.Entry) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric
	ch := make(chan []bmcmonitoring.Metric)

	for oam := 0; oam < c.Oams; oam++ {

		go func(oam int) {
			ll := log.WithField("oam", oam)

			resp, err := c.security(ctx, log, oam)
			if err != nil {
				ll.WithError(err).Error("failed to get security information")
				ch <- nil
				return
			}

			ch <- c.securityMetrics(ll, resp, oam)
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

func (c *Client) security(ctx context.Context, log *logrus.Entry, oam int) (SecurityResp, error) {
	var res SecurityResponse
	url := fmt.Sprintf("https://%s/redfish/v1/hl/oam/%d/Security", c.Hostname, oam)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return SecurityResp{}, err
	}

	return res.Response, nil

}

func (c *Client) securityMetrics(log *logrus.Entry, resp SecurityResp, oam int) []bmcmonitoring.Metric {
	var metrics []bmcmonitoring.Metric

	prefix := bmcmonitoring.PrefixSecurity

	// CurrentPublicKeyHashIndex
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentPublicKeyHashIndex,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricCurrentPublicKeyHashIndex),
	})

	// CurrentSVNVersion
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: resp.CurrentSVNVersion,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricCurrentSVNversion),
	})

	// Key0Revocation
	Key0RevocationVal := -1

	switch strings.TrimSpace(strings.ToLower(resp.Key0Revocation)) {
	case "not revoked":
		Key0RevocationVal = 0
	case "revoked":
		Key0RevocationVal = 1
	default:
		log.WithError(fmt.Errorf("unexpected key 0 Revocation value %s", resp.Key0Revocation)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: Key0RevocationVal,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricKey0Revocation),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricKey0Revocation: resp.Key0Revocation,
		},
	})

	// Key1Revocation
	Key1RevocationVal := -1

	switch strings.ToLower(resp.Key1Revocation) {
	case "not revoked":
		Key1RevocationVal = 0
	case "revoked":
		Key1RevocationVal = 1
	default:
		log.WithError(fmt.Errorf("unexpected key 1 Revocation value %s", resp.Key1Revocation)).Error()
	}
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: Key1RevocationVal,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricKey1Revocation),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricKey1Revocation: resp.Key1Revocation,
		},
	})

	// Key2Revocation
	Key2RevocationVal := -1

	switch strings.ToLower(resp.Key2Revocation) {
	case "not revoked":
		Key2RevocationVal = 0
	case "revoked":
		Key2RevocationVal = 1
	default:
		log.WithError(fmt.Errorf("unexpected key 2 Revocation value %s", resp.Key2Revocation)).Error()
	}
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: Key2RevocationVal,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricKey2Revocation),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricKey2Revocation: resp.Key2Revocation,
		},
	})

	// Key3Revocation
	Key3RevocationVal := -1

	switch strings.ToLower(resp.Key3Revocation) {
	case "not revoked":
		Key3RevocationVal = 0
	case "revoked":
		Key3RevocationVal = 1
	default:
		log.WithError(fmt.Errorf("unexpected key 3 Revocation value %s", resp.Key3Revocation)).Error()
	}
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: Key3RevocationVal,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricKey3Revocation),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricKey3Revocation: resp.Key3Revocation,
		},
	})

	// Key4Revocation
	Key4RevocationVal := -1

	switch strings.ToLower(resp.Key4Revocation) {
	case "not revoked":
		Key4RevocationVal = 0
	case "revoked":
		Key4RevocationVal = 1
	default:
		log.WithError(fmt.Errorf("unexpected key 4 Revocation value %s", resp.Key4Revocation)).Error()
	}

	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: Key4RevocationVal,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricKey4Revocation),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricKey4Revocation: resp.Key4Revocation,
		},
	})

	// MinimalSVNIndex
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: 0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricMinimalSVNindex),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricMinimalSVNindex: resp.MinimalSVNIndex,
		},
	})

	// FWImageSource

	fWImageSourceVal := -1
	switch strings.ToLower(resp.FWImageSource) {
	case "primary":
		fWImageSourceVal = 0
	case "secondary":
		fWImageSourceVal = 1
	default:
		log.WithError(fmt.Errorf("unexpected fw image source %s", resp.FWImageSource)).Error()
	}
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: fWImageSourceVal,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricFWImageSource),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricFWImageSource: resp.FWImageSource,
		},
	})

	// TpmPcrPpboot
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: 0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricTPMPCRPPBOOT),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricTPMPCRPPBOOT: resp.TpmPcrPpboot,
		},
	})

	// TpmPcrPreboot
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: 0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricTPMPCRPREBOOT),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricTPMPCRPREBOOT: resp.TpmPcrPreboot,
		},
	})

	// TpmPcrUBoot

	// TODO: In RAS it's called TpmPcrBoot, is it the right one?
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: 0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricTPMPCRBOOT),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricTPMPCRBOOT: resp.TpmPcrUboot,
		},
	})

	// TpmPcrLinux
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: 0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricTPMPCRLINUX),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricTPMPCRLINUX: resp.TpmPcrLinux,
		},
	})

	// CPLDVersion
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: 0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricCPLDVersion),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricCPLDVersion: resp.CPLDVersion,
		},
	})

	// CPLDVersionTimestamp
	metrics = append(metrics, bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		Oam:         fmt.Sprintf("%d", oam),
		MetricValue: 0,
		MetricName:  fmt.Sprintf("%s_%s", prefix, bmcmonitoring.SecurityMetricCPLDVersionTimestamp),
		CustomLabels: map[string]string{
			bmcmonitoring.SecurityMetricCPLDVersionTimestamp: fmt.Sprintf("%d", resp.CPLDVersionTimestamp),
		},
	})

	return metrics
}
