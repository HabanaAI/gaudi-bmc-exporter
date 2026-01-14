package bmcmonitoring

import (
	"time"
)

type ClientOpts struct {
	Username string
	Password string

	Hostname string

	Oams  int
	Ports int
	Lanes int

	ConnectionTimeout time.Duration

	// timeout for the creation of the client
	CreationTimeout time.Duration
}

type Metric struct {

	// CustomLabels allows to add more labels to the metric.
	CustomLabels map[string]string
	Hostname     string `label:"hostname"`
	KVMName      string `label:"kvm"`
	Oam          string `label:"oam"`
	MetricName   string `label:"metric_name"`

	// ras/red-fish
	Exporter string `label:"backend"`

	MetricValue int `label:"metric_value"`
}

func (m *Metric) AddLabel(key, value string) *Metric {
	if key == "oam" {
		m.Oam = value
	} else {
		if m.CustomLabels == nil {
			m.CustomLabels = make(map[string]string)
		}
		m.CustomLabels[key] = value
	}
	return m
}

type AlertsData struct {
	AlertsInformation map[string]AlertInfo `json:"rules"`
}
type AlertInfo struct {
	Severity    string `json:"severity"`
	Description string `json:"description"`
}
