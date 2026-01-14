package bmcmonitoring

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type ExporterStore struct {
	log       *logrus.Entry
	Exporters []Exporter
}

func NewExporterStore(exporters []Exporter, log *logrus.Entry) *ExporterStore {
	return &ExporterStore{
		Exporters: exporters,
		log:       log,
	}
}

// Append will add more exporters to the store.
func (s *ExporterStore) Append(exs []Exporter) {
	s.Exporters = append(s.Exporters, exs...)
}

// Close will logout each client.
func (s *ExporterStore) Close() error {
	for _, e := range s.Exporters {
		bmc := e.Name()
		err := e.Logout()
		if err != nil {
			s.log.WithError(fmt.Errorf("%s failed to logout: %w", bmc, err)).Error()
			continue
		}
		s.log.Infof("%s logged out", bmc)
	}

	return nil
}

func (s *ExporterStore) DeleteAllTokens(ctx context.Context) error {
	for _, e := range s.Exporters {
		bmc := e.Name()
		err := e.DeleteAllTokens(ctx)
		if err != nil {
			s.log.WithError(fmt.Errorf("%s failed to delete tokens: %w", bmc, err)).Error()
			continue
		}
		s.log.Infof("%s deleted tokens", bmc)
	}

	return nil
}

func (s *ExporterStore) RefreshToken(ctx context.Context) {
	var wg sync.WaitGroup
	for _, e := range s.Exporters {
		wg.Add(1)
		go func(e Exporter) {
			defer wg.Done()
			bmc := e.Name()
			err := e.RefreshToken(ctx)
			if err != nil {
				s.log.WithError(fmt.Errorf("%s failed refresh token: %w", bmc, err)).Error()
				return
			}
			s.log.Infof("%s refreshed token", bmc)
		}(e)

	}

	wg.Wait()

}

func (s *ExporterStore) Temperature(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixTemperature,
				"uuid":     uuid,
			})

			// e.RefreshToken()
			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			tempMetrics := e.Temperature(ctx, ll)
			var mets []Metric
			for _, m := range tempMetrics {
				if err := VerifyTemperature(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) EthernetInfo(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixEthernetInfo,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			ethernetInfoMetrics := e.EthernetInfo(ctx, ll)
			var mets []Metric
			for _, m := range ethernetInfoMetrics {

				if err := VerifyEthernetInfo(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) EthernetStatus(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixEthernetStatus,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			ethernetStatusMetrics := e.EthernetStatus(ctx, ll)

			var mets []Metric
			for _, m := range ethernetStatusMetrics {

				if err := VerifyEthernetStatus(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) PcieInfo(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixPcieInfo,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			pcieInfoMetrics := e.PcieInfo(ctx, ll)
			var mets []Metric
			for _, m := range pcieInfoMetrics {

				if err := VerifyPcieInfo(m); err != nil {
					ll.WithError(err).Error()
					continue
				}

				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) Status(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixStatus,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			statusMetrics := e.Status(ctx, ll)
			var mets []Metric
			for _, m := range statusMetrics {

				if err := VerifyStatus(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) Frequency(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixFrequency,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			frequencyMetrics := e.Frequency(ctx, ll)
			var mets []Metric
			for _, m := range frequencyMetrics {

				if err := VerifyFrequency(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) Power(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixPower,
				"uuid":     uuid,
			})
			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			powerMetrics := e.Power(ctx, ll)
			var mets []Metric
			for _, m := range powerMetrics {

				if err := VerifyPower(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) SensorTemperature(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixSensorTemperature,
				"uuid":     uuid,
			})
			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			sensorTemperatureMetrics := e.SensorTemperature(ctx, ll)
			var mets []Metric
			for _, m := range sensorTemperatureMetrics {

				if err := VerifySensorTemperature(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) SensorVoltage(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixSensorVoltage,
				"uuid":     uuid,
			})
			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			sensorVoltageMetrics := e.SensorVoltage(ctx, ll)
			var mets []Metric
			for _, m := range sensorVoltageMetrics {

				if err := VerifySensorVoltage(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) SensorCurrent(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixSensorCurrent,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			sensorCurrentMetrics := e.SensorCurrent(ctx, ll)
			var mets []Metric
			for _, m := range sensorCurrentMetrics {

				if err := VerifySensorCurrent(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) Info(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixInfo,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			infoMetrics := e.Info(ctx, ll)
			var mets []Metric

			for _, m := range infoMetrics {
				if err := VerifyInfo(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) Security(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixSecurity,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			securityMetrics := e.Security(ctx, ll)
			var mets []Metric
			for _, m := range securityMetrics {

				if err := VerifySecurity(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) Direct(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixDirect,
				"uuid":     uuid,
			})
			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			directMetrics := e.Direct(ctx, ll)
			var mets []Metric

			for _, m := range directMetrics {
				if err := VerifyDirect(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) CTemperature(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixCTemperature,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			cTemperatureMetrics := e.CTemperature(ctx, ll)
			var mets []Metric

			for _, m := range cTemperatureMetrics {
				if err := VerifyCTemperature(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) SensorVoltageMonitor(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixSensorVoltageMonitor,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			sensorVoltageMonitorMetrics := e.SensorVoltageMonitor(ctx, ll)
			var mets []Metric
			for _, m := range sensorVoltageMonitorMetrics {

				if err := VerifySensorVoltageMonitor(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) Hbm(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixHbm,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			hbmMetrics := e.Hbm(ctx, ll)

			var mets []Metric
			for _, m := range hbmMetrics {

				if err := VerifyHbm(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) EthernetStatusCounters(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {

			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixEthernetStatus,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			hbmMetrics := e.EthernetStatusCounters(ctx, ll)

			var mets []Metric
			for _, m := range hbmMetrics {

				if err := verifyEthernetStatusCounters(m); err != nil {
					ll.WithField("hostname", e.Name()).WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) LaneInfo(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {

			// e.RefreshToken()

			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixLaneInfo,
				"uuid":     uuid,
			})

			if reason, valid := e.ValidToken(ctx, ll); !valid {
				ll.Infof("invalid token: %s, %s", e.Name(), reason)
				ch <- []Metric{}
				return
			}

			metrics := e.LaneInfo(ctx, ll)

			var mets []Metric
			for _, m := range metrics {

				if err := VerifyLaneInfo(m); err != nil {
					ll.WithField("hostname", e.Name()).WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) BmcState(ctx context.Context, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixBMCState,
				"uuid":     uuid,
			})
			bmcStateMetrics := e.BmcState(ctx, ll)

			var mets []Metric
			for _, m := range bmcStateMetrics {

				if err := verifyBmcState(m); err != nil {
					s.log.WithField("hostname", e.Name()).WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}

func (s *ExporterStore) Alerts(ctx context.Context, alertsData AlertsData, uuid string) []Metric {

	ch := make(chan []Metric)
	var metrics []Metric
	for _, e := range s.Exporters {
		go func(e Exporter) {
			ll := s.log.WithFields(logrus.Fields{
				"hostname": e.Name(),
				"metric":   PrefixAlerts,
				"uuid":     uuid,
			})
			alertsMetrics := e.Alerts(ctx, alertsData, ll)

			var mets []Metric
			for _, m := range alertsMetrics {

				if err := VerifyAlerts(m); err != nil {
					ll.WithError(err).Error()
					continue
				}
				mets = append(mets, m)

			}

			ch <- mets

		}(e)
	}

	for i := 0; i < len(s.Exporters); i++ {
		m := <-ch
		metrics = append(metrics, m...)
	}

	return metrics
}
