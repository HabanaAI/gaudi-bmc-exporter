package main

import (
	"context"
	"flag"
	"fmt"
	g3redfish "habana_bmc_exporter/internal/g3redfish"
	rasmonitoringapi "habana_bmc_exporter/internal/ras-monitoring-api"
	redfish "habana_bmc_exporter/internal/red-fish"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"habana_bmc_exporter/pkg/web"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"go.uber.org/automaxprocs/maxprocs"
)

var build = "development"

func main() {
	log := logger.New().WithField("service", "bmc-exporter")

	if err := run(log); err != nil {
		log.WithError(err).Error()
		os.Exit(1)
	}
}

type application struct {
	log           *logrus.Entry
	store         *rasmonitoring.ExporterStore
	failedServers *map[string]serverInfo
	alertsData    rasmonitoring.AlertsData
	build         string
	exporter      string
	confPath      string
}

func run(log *logrus.Entry) error {

	daemonSet := flag.Bool("daemon-set", true, "the exporter will run as a daemon set, meaning every node will monitor it's own Bmc")
	exporter := flag.String("exporter", RedFish, "set exporter to raw RAS or Red fish")
	confPath := flag.String("config", "./config.json", "path to json config")
	tlsCert := flag.String("tls_cert", "", "path to tls certificate")
	tlsKey := flag.String("tls_key", "", "path to tls key")
	alertsInfoPath := flag.String("alerts-description", "/etc/alert_description.json", "path to json file containing the information about the alerts")
	flag.Parse()

	if (*tlsCert != "" && *tlsKey == "") || (*tlsKey != "" && *tlsCert == "") {
		return fmt.Errorf("tls_cert and tls_key flags must be provided together")
	}

	// ============================================================
	// GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either bt the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}

	// Load alerts description
	alertsData, err := rasmonitoring.InitAlertsDescription(*alertsInfoPath)
	if err != nil {
		return err
	}

	// Load config.
	conf, err := initConfig(*confPath)
	if err != nil {
		return err
	}

	log.Info("Config Loaded")

	exporters, failedServers, err := getExporters(*exporter, conf.Servers, *daemonSet, log)
	if err != nil {
		return err
	}

	log.Infof("got %d servers", len(exporters))
	log.Infof("got %d failed servers", len(failedServers))

	store := rasmonitoring.NewExporterStore(exporters, log)

	// When this context is cancel all the http call using this context will be revoked.
	ctx, cancelFunc := context.WithCancel(context.Background())

	go func() {
		for {

			time.Sleep(3 * time.Minute)

			err := store.DeleteAllTokens(ctx)
			if err != nil {
				log.WithError(err).Error("failed to delete all tokens")
			}

			// refresh token in background
			go store.RefreshToken(ctx)

			// Initialize failed servers.

			log.Info("trying to create failed clients")
			var exs []rasmonitoring.Exporter
			for _, server := range failedServers {
				log.Infof("trying to add failed server %s", server.Hostname)

				var c rasmonitoring.Exporter
				var err error

				switch *exporter {
				case G3RedFish:
					c, err = g3redfish.NewClient(rasmonitoring.ClientOpts{
						Hostname:          server.Hostname,
						Username:          server.Username,
						Password:          server.Password,
						ConnectionTimeout: 10 * time.Second,
					}, log)

				case RedFish:
					c, err = redfish.NewClient(rasmonitoring.ClientOpts{
						Hostname:          server.Hostname,
						Username:          server.Username,
						Password:          server.Password,
						ConnectionTimeout: 10 * time.Second,
					}, log)

				case Ras:
					c, err = rasmonitoringapi.NewClient(rasmonitoring.ClientOpts{
						Hostname:          server.Hostname,
						Username:          server.Username,
						Password:          server.Password,
						ConnectionTimeout: 10 * time.Second,
					}, log)
				}

				if err != nil {
					log.WithField("hostname", server.Hostname).WithError(err).Error("failed to create exporter")
					continue
				}

				log.Infof("created client %s", server.Hostname)

				// add the client to the list of exports.
				exs = append(exs, c)

				// remove server that we added successfully from the failed servers.
				delete(failedServers, server.Hostname)
			}

			// add the new exports to our store.
			store.Append(exs)

		}

	}()

	log.Infof("started service %s, port %s", build, conf.Port)

	app := application{
		store:         store,
		log:           log,
		build:         build,
		exporter:      *exporter,
		failedServers: &failedServers,
		alertsData:    alertsData,
		confPath:      *confPath,
	}

	srv := app.routes(ctx, conf.Port)

	if *tlsCert != "" && *tlsKey != "" {
		return srv.Start(app.log, cancelFunc, store, web.TlsCerts{
			CertFile: *tlsCert,
			KeyFile:  *tlsKey,
		})

	}

	return srv.Start(app.log, cancelFunc, store)

}
