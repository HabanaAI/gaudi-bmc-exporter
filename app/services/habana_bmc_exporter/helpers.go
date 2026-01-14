package main

import (
	"encoding/json"
	"fmt"
	"habana_bmc_exporter/internal/g3redfish"
	rasmonitoringapi "habana_bmc_exporter/internal/ras-monitoring-api"
	redfish "habana_bmc_exporter/internal/red-fish"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const defaultPort = "4000"

// Supported exporters
const (
	Ras       = "ras"
	RedFish   = "red-fish"
	G3RedFish = "g3-red-fish"
)

type config struct {
	Username string       `json:"username"`
	Password string       `json:"password"`
	Port     string       `json:"port"`
	Servers  []serverInfo `json:"servers"`
}

type serverInfo struct {

	// Hib bmc name
	Hostname string `json:"hostname"`
	Password string `json:"password"`
	Username string `json:"username"`

	// k8s node name
	Srv string `json:"srv"`
}

// initConfig will read the config.
func initConfig(path string) (config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return config{}, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	var conf config

	err = json.Unmarshal(data, &conf)
	if err != nil {
		return config{}, fmt.Errorf("failed to load the config, please check your config: %w", err)
	}

	if conf.Port == "" {
		conf.Port = defaultPort
	}

	// If user didn't put the username and password for the server use the master username and password.
	for i := range conf.Servers {
		if conf.Servers[i].Username == "" {
			conf.Servers[i].Username = conf.Username
		}

		if conf.Servers[i].Password == "" {
			conf.Servers[i].Password = conf.Password
		}

		// Get from env vars
		if conf.Servers[i].Username == "" {
			conf.Servers[i].Username = os.Getenv("USERNAME")
		}

		if conf.Servers[i].Password == "" {
			conf.Servers[i].Password = os.Getenv("PASSWORD")

		}
	}

	return conf, nil
}

func getExporters(exporter string, servers []serverInfo, daemonSet bool, log *logrus.Entry) ([]rasmonitoring.Exporter, map[string]serverInfo, error) {

	var exporters []rasmonitoring.Exporter

	type exporterCh struct {
		exporter     rasmonitoring.Exporter
		err          error
		failedServer serverInfo
	}

	ch := make(chan exporterCh)

	// failedServers are servers we failed to initiate.
	failedServers := make(map[string]serverInfo)

	// filter servers

	var filteredServers []serverInfo

	for _, server := range servers {

		nodeName := os.Getenv("KUBERNETES_NODENAME")

		// If we are running in kubernetes
		if nodeName != "" {

			// we are running on a demon set mode = every node will only look at his own bmc
			if !strings.Contains(server.Srv, nodeName) && daemonSet {
				continue
			}

		}

		filteredServers = append(filteredServers, server)
	}

	for _, server := range filteredServers {

		switch exporter {
		case Ras:
			go func(server serverInfo) {
				c, err := rasmonitoringapi.NewClient(rasmonitoring.ClientOpts{
					Hostname: server.Hostname,
					Username: server.Username,
					Password: server.Password,

					// TODO: Maybe get it from config
					ConnectionTimeout: 10 * time.Second,
				}, log)

				if err != nil {
					ch <- exporterCh{
						err:          err,
						failedServer: server,
					}
					return
				}

				ch <- exporterCh{
					exporter: c,
				}

			}(server)

		case RedFish:
			go func(server serverInfo) {
				c, err := redfish.NewClient(rasmonitoring.ClientOpts{
					Hostname:          server.Hostname,
					Username:          server.Username,
					Password:          server.Password,
					ConnectionTimeout: 10 * time.Second,
				}, log)

				if err != nil {
					ch <- exporterCh{
						err:          err,
						failedServer: server,
					}
					return
				}

				ch <- exporterCh{
					exporter: c,
				}
			}(server)
		case G3RedFish:
			go func(server serverInfo) {
				c, err := g3redfish.NewClient(rasmonitoring.ClientOpts{
					Hostname:          server.Hostname,
					Username:          server.Username,
					Password:          server.Password,
					ConnectionTimeout: 10 * time.Second,
				}, log)

				if err != nil {
					ch <- exporterCh{
						err:          err,
						failedServer: server,
					}
					return
				}

				ch <- exporterCh{
					exporter: c,
				}
			}(server)
		default:
			return nil, nil, fmt.Errorf("exporter %s is not supported", exporter)
		}
	}

	for i := 0; i < len(filteredServers); i++ {
		m := <-ch

		if m.err != nil {
			failedServers[m.failedServer.Hostname] = m.failedServer
			log.WithField("hostname", m.failedServer.Hostname).WithError(m.err).Error("failed to create exporter")
			continue
		}

		log.Infof("created exporter: %s", m.exporter.Name())
		exporters = append(exporters, m.exporter)
	}

	return exporters, failedServers, nil
}
