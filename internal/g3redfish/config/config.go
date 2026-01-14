package config

import (
	"strings"
)

type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	ApiVersion  string             `yaml:"apiVersion,omitempty" default:"v1"`
	ApiPath     string             `yaml:"apiPath"`
	Category    string             `yaml:"category,omitempty" default:""`
	ForEach     *ForEach           `yaml:"forEach,omitempty"`
	Metrics     []MetricDefinition `yaml:"metrics"`
	ExtraParams map[string]string
}

type ForEach struct {
	JsonPath   string     `yaml:"jsonPath"`
	Endpoints  []Endpoint `yaml:"endpoints"`
	Skip       []string   `yaml:"skip,omitempty"`
	ExtraParam ExtraParam `yaml:"extraParam,omitempty"`
}

type ExtraParam struct {
	Name  string `yaml:"name"`
	Regex string `yaml:"regex"`
}

type MetricDefinition struct {
	Name          string            `yaml:"name,omitempty"`
	Type          string            `yaml:"type,omitempty" default:"gauge"`
	JsonPath      string            `yaml:"jsonPath"`
	ValueMappings map[string]string `yaml:"valueMappings,omitempty"`
	ExtraLabels   map[string]string `yaml:"extraLabels,omitempty"`
}

func (c *Config) GetEndpointByCategory(category string) []Endpoint {
	if category == "" {
		return c.Endpoints
	}
	var endpoints []Endpoint
	for _, endpoint := range c.Endpoints {
		if endpoint.Category == category {
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}

func (e Endpoint) ResolveMembers(member string) Endpoint {
	return Endpoint{
		ApiVersion:  e.ApiVersion,
		ApiPath:     strings.ReplaceAll(e.ApiPath, "{{member}}", member),
		Category:    e.Category,
		ForEach:     e.ForEach,
		Metrics:     e.Metrics,
		ExtraParams: e.ExtraParams,
	}
}

func (e *Endpoint) AddExtraParam(key, value string) *Endpoint {
	if e.ExtraParams == nil {
		e.ExtraParams = map[string]string{}
	}
	e.ExtraParams[key] = value
	return e
}
