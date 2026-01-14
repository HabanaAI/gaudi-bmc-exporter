package main

import (
	"context"
	"encoding/json"
	rasmonitoringapi "habana_bmc_exporter/internal/ras-monitoring-api"
	redfish "habana_bmc_exporter/internal/red-fish"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/util"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (app application) metrics(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	// Call all the endpoints
	app.alerts(ctx, w, r)
	app.bmcState(ctx, w, r)
	app.ctemperature(ctx, w, r)
	app.direct(ctx, w, r)
	app.ethernetInfo(ctx, w, r)
	app.ethernetStatus(ctx, w, r)
	app.ethernetStatusCounters(ctx, w, r)
	app.hbm(ctx, w, r)
	app.info(ctx, w, r)
	app.laneInfo(ctx, w, r)
	app.pcieInfo(ctx, w, r)
	app.power(ctx, w, r)
	app.security(ctx, w, r)
	app.sensorCurrent(ctx, w, r)
	app.sensorTemperature(ctx, w, r)
	app.sensorVoltage(ctx, w, r)
	app.sensorVoltageMonitor(ctx, w, r)
	app.status(ctx, w, r)
	app.temperature(ctx, w, r)
}

func (app application) probe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// get url params indicating who we need to query

	app.log.Info("got a probe request")

	server := r.URL.Query().Get("target")

	var client bmcmonitoring.Exporter

	// Get connections details for the server
	config, err := initConfig(app.confPath)
	if err != nil {
		app.log.WithError(err).Error()
		return
	}

	password := config.Password
	username := config.Username

	switch app.exporter {
	case Ras:
		client, err = rasmonitoringapi.NewClient(bmcmonitoring.ClientOpts{
			Username:          username,
			Password:          password,
			Hostname:          server,
			ConnectionTimeout: 10 * time.Second,
		}, app.log)
	case RedFish:
		client, err = redfish.NewClient(bmcmonitoring.ClientOpts{
			Username:          username,
			Password:          password,
			Hostname:          server,
			ConnectionTimeout: 10 * time.Second,
		}, app.log)
	default:
		app.log.Error("unknown exporter type")
		return
	}

	if err != nil {
		app.log.WithError(err).Error("probe create client")
		return
	}

	defer func() {
		if client != nil {
			err := client.Logout()
			if err != nil {
				app.log.WithError(err).Error("failed to logout")
			}
		}
	}()

	// call all the endpoints of the client

	var metrics []bmcmonitoring.Metric
	metrics = append(metrics, client.Temperature(ctx, app.log)...)
	metrics = append(metrics, client.EthernetInfo(ctx, app.log)...)
	metrics = append(metrics, client.EthernetStatus(ctx, app.log)...)
	metrics = append(metrics, client.PcieInfo(ctx, app.log)...)
	metrics = append(metrics, client.Frequency(ctx, app.log)...)
	metrics = append(metrics, client.Power(ctx, app.log)...)
	metrics = append(metrics, client.SensorTemperature(ctx, app.log)...)
	metrics = append(metrics, client.SensorVoltage(ctx, app.log)...)
	metrics = append(metrics, client.SensorCurrent(ctx, app.log)...)
	metrics = append(metrics, client.Info(ctx, app.log)...)
	metrics = append(metrics, client.Security(ctx, app.log)...)
	metrics = append(metrics, client.Direct(ctx, app.log)...)
	metrics = append(metrics, client.SensorVoltageMonitor(ctx, app.log)...)
	metrics = append(metrics, client.Hbm(ctx, app.log)...)
	metrics = append(metrics, client.BmcState(ctx, app.log)...)
	metrics = append(metrics, client.EthernetStatusCounters(ctx, app.log)...)
	metrics = append(metrics, client.Alerts(ctx, app.alertsData, app.log)...)
	metrics = append(metrics, client.LaneInfo(ctx, app.log)...)

	// print the metrics
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}
}

func (app application) exporterInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	info := ExporterInfo{
		MetricValue:     0,
		MetricName:      exporterInfo,
		ExporterVersion: app.build,
	}

	err := util.Print(&info, w)
	if err != nil {
		app.log.WithError(err).Error()
	}
}

// serversStatus will report the status of the servers.
func (app application) serversStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var status ServerStatus

	for _, s := range app.store.Exporters {
		status.RunningServers = append(status.RunningServers, s.Name())
	}

	for _, s := range *app.failedServers {
		status.FailedServers = append(status.FailedServers, s.Hostname)
	}
	err := json.NewEncoder(w).Encode(status)
	if err != nil {
		app.log.WithError(err).Error()
	}
}

func (app application) temperature(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.Temperature(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) laneInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.LaneInfo(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) ethernetInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.EthernetInfo(r.Context(), u.String())
	for _, metric := range metrics {

		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) ethernetStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.EthernetStatus(r.Context(), u.String())
	for _, metric := range metrics {

		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) ethernetStatusCounters(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.EthernetStatusCounters(r.Context(), u.String())
	for _, metric := range metrics {

		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) pcieInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.PcieInfo(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) status(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.Status(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) frequency(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.Frequency(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) power(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.Power(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) sensorTemperature(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.SensorTemperature(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) sensorVoltage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.SensorVoltage(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) sensorCurrent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.SensorCurrent(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) info(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.Info(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) security(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.Security(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}

}

func (app application) direct(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.Direct(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}
}

func (app application) ctemperature(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.CTemperature(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}
}

func (app application) sensorVoltageMonitor(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.SensorVoltageMonitor(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}
}

func (app application) hbm(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.Hbm(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}
}

func (app application) bmcState(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.BmcState(r.Context(), u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}
}

func (app application) alerts(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	u := uuid.New()

	metrics := app.store.Alerts(r.Context(), app.alertsData, u.String())
	for _, metric := range metrics {
		metric.Exporter = app.exporter
		err := util.Print(&metric, w)
		if err != nil {
			app.log.WithError(err).Error()
		}
	}
}
