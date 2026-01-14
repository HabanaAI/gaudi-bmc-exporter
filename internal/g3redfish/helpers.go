package g3redfish

import (
	"context"
	"fmt"
	"habana_bmc_exporter/internal/g3redfish/config"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *Client) createMetric(_ *logrus.Entry, metricName string, metricValue int) bmcmonitoring.Metric {
	metric := bmcmonitoring.Metric{
		Hostname:    c.Hostname,
		MetricName:  metricName,
		MetricValue: metricValue,
	}

	return metric
}

func resolveMapping(value string, mapping map[string]string) string {
	key := strings.TrimSpace(strings.ToLower(value))
	if mappedValue, exists := mapping[key]; exists {
		return mappedValue
	}
	return value
}

func (c *Client) extractMetric(log *logrus.Entry, value interface{}, metricDef config.MetricDefinition, endpointExtraLabels map[string]string) (*bmcmonitoring.Metric, error) {
	value, err := getValueByFieldName(value, metricDef.JsonPath)
	if err != nil {
		return nil, err
	}

	finalValue := resolveMapping(fmt.Sprintf("%v", value), metricDef.ValueMappings)
	metricVal, err := strconv.ParseFloat(finalValue, 64)
	valueLabels := map[string]string{}
	if err != nil {
		metricVal = 0
		valueLabels["value"] = finalValue
	}

	metric := c.createMetric(log, metricDef.Name, int(metricVal))
	for key, value := range endpointExtraLabels {
		metric.AddLabel(key, value)
	}
	for key, value := range metricDef.ExtraLabels {
		metric.AddLabel(key, value)
	}
	for key, value := range valueLabels {
		metric.AddLabel(key, value)
	}
	return &metric, nil
}

func (c *Client) getMetricsFromEndpoint(ctx context.Context, log *logrus.Entry, endpoint config.Endpoint) []bmcmonitoring.Metric {
	var res map[string]interface{}
	url := c.createUrl(endpoint.ApiPath)
	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		log.WithError(err).Error()
		return nil
	}

	var metrics []bmcmonitoring.Metric
	for _, metricDefinition := range endpoint.Metrics {
		metric, err := c.extractMetric(log, res, metricDefinition, endpoint.ExtraParams)
		if err != nil {
			log.WithField("metricName", metricDefinition.Name).
				WithField("jsonPath", metricDefinition.JsonPath).
				WithField("apiPath", endpoint.ApiPath).
				WithError(err).Error("Failed to extract metric.")
			continue
		}
		metrics = append(metrics, *metric)
	}

	if endpoint.ForEach != nil {
		rawMembers, err := getValueByFieldName(res, endpoint.ForEach.JsonPath)
		if err != nil {
			log.WithError(err).Error("Failed to get value by field name")
			return metrics
		}
		members := getMembers(rawMembers)

		for member_id := range members {
			member := members[member_id]
			if shouldSkipMember(endpoint, member) {
				continue
			}
			for _, subEndpoint := range endpoint.ForEach.Endpoints {
				resolvedEndpoint := subEndpoint.ResolveMembers(member)
				memberParts := strings.Split(member, "/")
				extraParamValue := extractByRegex(memberParts[len(memberParts)-1], endpoint.ForEach.ExtraParam.Regex, 1)
				if extraParamValue == "" {
					log.Errorf("failed to extract extra param value from %s", member)
					continue
				}
				resolvedEndpoint.AddExtraParam(endpoint.ForEach.ExtraParam.Name, extraParamValue)
				subMetrics := c.getMetricsFromEndpoint(ctx, log, resolvedEndpoint)
				if subMetrics != nil {
					metrics = append(metrics, subMetrics...)
				}
			}
		}
	}

	return metrics
}

func (c *Client) getMetrics(ctx context.Context, log *logrus.Entry, category string) []bmcmonitoring.Metric {
	endpoints := c.config.GetEndpointByCategory(category)
	if len(endpoints) == 0 {
		log.Errorf("no endpoints found for category %s", category)
		return nil
	}
	dataCh := make(chan []bmcmonitoring.Metric)

	for _, endpoint := range endpoints {
		go func(endpoint config.Endpoint) {
			metrics := c.getMetricsFromEndpoint(ctx, log, endpoint)
			dataCh <- metrics
		}(endpoint)
	}

	var metrics []bmcmonitoring.Metric
	for i := 0; i < len(endpoints); i++ {
		newMetrics := <-dataCh
		if newMetrics != nil {
			metrics = append(metrics, newMetrics...)
		}
	}
	return metrics
}

func getValueByFieldName(obj interface{}, path string) (interface{}, error) {
	v := reflect.ValueOf(obj)

	if v.Kind() == reflect.Map && v.Len() == 0 {
		return nil, fmt.Errorf("failed to get %s", path)
	}

	// Change keys to lowercase
	newMap := make(map[string]interface{})
	for _, key := range v.MapKeys() {
		newMap[strings.ToLower(key.String())] = v.MapIndex(key).Interface()
	}

	// Extract the part of the path
	parts := strings.Split(strings.ToLower(path), ".")
	value, ok := newMap[parts[0]]
	if !ok {
		return nil, fmt.Errorf("failed to get %s", path)
	}

	if len(parts) == 1 {
		return value, nil
	}

	nextPath := strings.Join(parts[1:], ".")
	return getValueByFieldName(value, nextPath)
}

func shouldSkipMember(endpoint config.Endpoint, member string) bool {
	if endpoint.ForEach != nil && endpoint.ForEach.Skip != nil {
		for _, skip := range endpoint.ForEach.Skip {
			if skip == member {
				return true
			}
		}
	}
	return false
}

func extractByRegex(text string, regex string, group int) string {
	r := regexp.MustCompile(regex)
	results := r.FindStringSubmatch(text)
	if len(results) > group {
		return results[group]
	}
	return ""
}

func getMembers(data interface{}) []string {
	var ids []string
	switch data.(type) {
	case []interface{}:
		// {"1": [{"@odata.id": "/redfish/v1/Systems/1"}]}
		for _, member := range data.([]interface{}) {
			memberMap := member.(map[string]interface{})
			for key := range memberMap {
				ids = append(ids, memberMap[key].(string))
			}
		}
	case map[string]interface{}:
		// {"1": {"@odata.id": "/redfish/v1/Systems/1"}}
		memberMap := data.(map[string]interface{})
		for key := range memberMap {
			ids = append(ids, memberMap[key].(string))
		}
	}
	return ids
}
