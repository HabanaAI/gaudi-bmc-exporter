package main

type ServerStatus struct {
	RunningServers []string `json:"running_servers"`
	FailedServers  []string `json:"failed_servers"`
}

const (
	exporterInfo = "habana_bmc_exporter_info"
)

type ExporterInfo struct {
	MetricName      string `label:"metric_name"`
	ExporterVersion string `label:"exporter_version"`
	MetricValue     int    `label:"metric_value"`
}
