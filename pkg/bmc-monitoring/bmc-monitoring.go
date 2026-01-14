package bmcmonitoring

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Exporter interface {
	//TODO: ad comments per group
	Temperature(context.Context, *logrus.Entry) []Metric
	EthernetInfo(context.Context, *logrus.Entry) []Metric
	EthernetStatus(context.Context, *logrus.Entry) []Metric
	PcieInfo(context.Context, *logrus.Entry) []Metric
	Status(context.Context, *logrus.Entry) []Metric
	Frequency(context.Context, *logrus.Entry) []Metric
	Power(context.Context, *logrus.Entry) []Metric
	SensorTemperature(context.Context, *logrus.Entry) []Metric
	SensorVoltage(context.Context, *logrus.Entry) []Metric
	SensorCurrent(context.Context, *logrus.Entry) []Metric
	Info(context.Context, *logrus.Entry) []Metric
	Security(context.Context, *logrus.Entry) []Metric
	Direct(context.Context, *logrus.Entry) []Metric
	CTemperature(context.Context, *logrus.Entry) []Metric
	SensorVoltageMonitor(context.Context, *logrus.Entry) []Metric
	Hbm(context.Context, *logrus.Entry) []Metric
	BmcState(context.Context, *logrus.Entry) []Metric
	EthernetStatusCounters(context.Context, *logrus.Entry) []Metric
	Alerts(context.Context, AlertsData, *logrus.Entry) []Metric
	LaneInfo(context.Context, *logrus.Entry) []Metric
	Logout() error
	RefreshToken(context.Context) error
	// Return the name of the server.
	Name() string

	// Temp solution for RAS to solve the max connections issue.
	DeleteAllTokens(context.Context) error

	// return the reason of why the token is not valid.
	ValidToken(context.Context, *logrus.Entry) (string, bool)
}

type Client struct {
	http.Client
	ClientOpts
}

func NewClient(opts ClientOpts) (*Client, error) {
	c := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client := &Client{
		ClientOpts: opts,
		Client:     c,
	}

	if opts.ConnectionTimeout != 0 {
		client.Timeout = opts.ConnectionTimeout
	}

	return client, nil
}
