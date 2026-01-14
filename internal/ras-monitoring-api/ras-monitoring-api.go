package rasmonitoringapi

import (
	"context"
	"fmt"
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/ipmi"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	log *logrus.Entry
	*rasmonitoring.Client
	ipmiClient *ipmi.Client
	Auth       authDetails
	sync.Mutex
}

type res struct {
	client *Client
	err    error
}

const (
	numOfPorts     = 24
	defaultNumOams = 8
	defaultLanes   = 16
)

func NewClient(opts rasmonitoring.ClientOpts, log *logrus.Entry) (*Client, error) {

	if opts.CreationTimeout == 0 {
		opts.CreationTimeout = 10 * time.Second
	}

	if opts.Oams == 0 {
		opts.Oams = defaultNumOams
	}

	if opts.Ports == 0 {
		opts.Ports = numOfPorts
	}

	if opts.Lanes == 0 {
		opts.Lanes = defaultLanes
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.CreationTimeout)
	defer cancel()

	ch := make(chan res)

	go func() {
		c, err := rasmonitoring.NewClient(opts)
		if err != nil {
			ch <- res{err: err}
			return
		}

		ll := log.WithField("hostname", opts.Hostname)

		client := &Client{
			Client: c,
			log:    ll,
			Mutex:  sync.Mutex{},
			ipmiClient: ipmi.NewClient(ipmi.ClientOpts{
				BMCUsername: opts.Username,
				BMCPassword: opts.Password,
				BMCName:     opts.Hostname,
			}),
		}

		err = client.authenticate()
		if err != nil {
			ch <- res{err: fmt.Errorf("failed to authenticate: %w", err)}
			return
		}

		// verify that the token is valid.
		if reason, valid := client.ValidToken(context.Background(), ll); !valid {
			ch <- res{err: fmt.Errorf("invalid token: %s", reason)}
			return
		}

		ch <- res{client: client}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout exceeded for creating client %s", opts.Hostname)
	case resp := <-ch:
		if resp.err != nil {
			return nil, resp.err
		}
		return resp.client, nil
	}

}

func (c *Client) Direct(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {

	var metrics []rasmonitoring.Metric
	prefix := rasmonitoring.PrefixDirect

	metrics = c.getMetrics(ctx, log, directCodes, prefix, methodDirect)

	m := c.apiVersion(ctx, log, apiVersionOpcode, c.Hostname, prefix)
	metrics = append(metrics, m...)

	m = c.osState(ctx, log, osStateOpcode, c.Hostname, prefix)
	metrics = append(metrics, m...)

	return metrics
}

func (c *Client) Security(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, securityOpcodes, rasmonitoring.PrefixSecurity, methodIndirect)
}

func (c *Client) Info(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, infoOpcodes, rasmonitoring.PrefixInfo, methodIndirect)
}

func (c *Client) SensorCurrent(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, sensorCurrentOpcodes, rasmonitoring.PrefixSensorCurrent, methodIndirect)
}

func (c *Client) SensorVoltage(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, sensorVoltageOpcodes, rasmonitoring.PrefixSensorVoltage, methodIndirect)
}

func (c *Client) SensorTemperature(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, sensorTemperatureOpcodes, rasmonitoring.PrefixSensorTemperature, methodIndirect)
}

func (c *Client) Power(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, powerOpcodes, rasmonitoring.PrefixPower, methodIndirect)
}

func (c *Client) Frequency(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, frequencyOpcodes, rasmonitoring.PrefixFrequency, methodIndirect)
}

func (c *Client) Status(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, statusOpcodes, rasmonitoring.PrefixStatus, methodIndirect)
}

func (c *Client) PcieInfo(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, pciInfoOpcodes, rasmonitoring.PrefixPcieInfo, methodIndirect)
}

func (c *Client) EthernetStatus(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {

	return c.getMetrics(ctx, log, ethernetStatusOpcodes, rasmonitoring.PrefixEthernetStatus, methodIndirect)

}

func (c *Client) EthernetInfo(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, ethernetInfoOpcodes, rasmonitoring.PrefixEthernetInfo, methodIndirect)
}

func (c *Client) Temperature(ctx context.Context, log *logrus.Entry) []rasmonitoring.Metric {
	return c.getMetrics(ctx, log, temperatureOpcodes, rasmonitoring.PrefixTemperature, methodIndirect)
}

// getMetrics gets the metrics according to the opcodes and decode them according to their expected type.
func (c *Client) getMetrics(ctx context.Context, log *logrus.Entry, opcodes map[string]Opcode, prefix, method string) []rasmonitoring.Metric {

	var metrics []rasmonitoring.Metric
	ch := make(chan result)
	done := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for oam := 0; oam < c.Oams; oam++ {

			for fieldName, op := range opcodes {

				val, err := c.decodeOpcode(ctx, op, oam, fieldName, method)
				if err != nil {
					ch <- result{
						err: err,
					}
					continue
				}

				if op.Decoder == nil {

					switch op.ExpectedType {

					// For string we add it a label to our metric.
					case ReverseString, String, AsciiString:
						customLabels := map[string]string{
							fieldName: val,
						}
						ch <- result{
							metric: rasmonitoring.Metric{
								MetricName:   fmt.Sprintf("%s_%s", prefix, fieldName),
								Hostname:     c.Hostname,
								Oam:          fmt.Sprintf("%d", oam),
								MetricValue:  0,
								CustomLabels: customLabels,
							}}

					case Int, Bin, BinArray:
						value, err := strconv.Atoi(val)
						if err != nil {
							ch <- result{
								err: err,
							}
							return
						}
						ch <- result{
							metric: rasmonitoring.Metric{
								MetricName:  fmt.Sprintf("%s_%s", prefix, fieldName),
								Hostname:    c.Hostname,
								Oam:         fmt.Sprintf("%d", oam),
								MetricValue: value,
							}}
					}

				} else {

					// Call a custom func to decode the data.
					m, err := op.Decoder(val, oam, c.Hostname, fieldName, prefix)
					if err != nil {
						log.WithError(err).Error()
						continue
					}

					// send the metrics
					for _, i := range m {
						ch <- result{
							metric: i,
						}
					}
				}

			}

		}

	}()

	go func() {
		wg.Wait()
		done <- struct{}{}
		close(ch)
	}()

	for {
		select {
		case <-done:
			return metrics
		case data := <-ch:
			if data.err != nil {

				log.WithError(data.err).Error()
				continue
			}
			metrics = append(metrics, data.metric)
		}
	}

}
