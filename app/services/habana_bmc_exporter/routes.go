package main

import (
	"context"
	"habana_bmc_exporter/pkg/web"
	"habana_bmc_exporter/pkg/web/checkgrp"
	"time"
)

// routes initialize the http mux and register mux path.
func (app application) routes(ctx context.Context, port string) *web.Server {

	return web.NewServer(ctx, web.ServerOpts{
		Port: port,
		Middleware: []web.Middleware{
			web.PanicMid,
			web.TimeSinceMid,
			app.combineContextMiddleware,
		},
		Log: app.log,
	}, checkgrp.Handlers{
		Log:   app.log,
		Build: app.build,
	}, []web.WebHandler{
		{
			Route:   "/metrics",
			Handler: app.metrics,
		},
		{
			Route:   "/probe",
			Handler: app.probe,
		},
		{
			Route:   "/temperature",
			Handler: app.temperature,
		},
		{
			Route:   "/ethernet-info",
			Handler: app.ethernetInfo,
		},
		{
			Route:   "/ethernet-status",
			Handler: app.ethernetStatus,
		},
		{
			Route:   "/ethernet-status-counters",
			Handler: app.ethernetStatusCounters,
		},
		{
			Route:   "/pcie-info",
			Handler: app.pcieInfo,
		},
		{
			Route:   "/status",
			Handler: app.status,
		},
		{
			Route:   "/frequency",
			Handler: app.frequency,
		},
		{
			Route:   "/power",
			Handler: app.power,
		},
		{
			Route:   "/sensor-temperature",
			Handler: app.sensorTemperature,
		},
		{
			Route:   "/sensor-voltage",
			Handler: app.sensorVoltage,
		},
		{
			Route:   "/sensor-current",
			Handler: app.sensorCurrent,
		},
		{
			Route:   "/info",
			Handler: app.info,
		},
		{
			Route:   "/security",
			Handler: app.security,
		},
		{
			Route:   "/direct",
			Handler: app.direct,
		},
		{
			Route:   "/ctemperature",
			Handler: app.ctemperature,
		},
		{
			Route:   "/sensor-voltage-monitor",
			Handler: app.sensorVoltageMonitor,
		},
		{
			Route:   "/hbm",
			Handler: app.hbm,
		},
		{
			Route:   "/bmc-state",
			Handler: app.bmcState,
		},
		{
			Route:   "/alerts",
			Handler: app.alerts,
		},
		{
			Route:   "/lane-info",
			Handler: app.laneInfo,
		},
		{
			Route:   "/debug/status",
			Handler: app.serversStatus,
		},
		{
			Route:   "/exporter-info",
			Handler: app.exporterInfo,
		},
	}, web.WithWriteTimeout(time.Minute*5))

}
