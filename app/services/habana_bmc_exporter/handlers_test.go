package main

import (
	"bytes"
	"context"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"habana_bmc_exporter/pkg/mock"
	"habana_bmc_exporter/pkg/web"
	"habana_bmc_exporter/pkg/web/checkgrp"
	"habana_bmc_exporter/util"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// fakeTestMetric will return a basic metric, will add custom labels if supplied
func fakeTestMetric(customLabels ...string) bmcmonitoring.Metric {
	metric := bmcmonitoring.Metric{
		Hostname:    "hostname",
		KVMName:     "kvm",
		Oam:         "0",
		MetricValue: 3,
		MetricName:  "some metric name",
		Exporter:    "none",
	}

	if len(customLabels) > 0 {
		metric.CustomLabels = make(map[string]string)

		for _, label := range customLabels {
			metric.CustomLabels[label] = "some label value"
		}
	}
	return metric
}

func TestCtemperature(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter)
		expectedMetrics func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					CTemperature(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric()})
			},
			expectedMetrics: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				metrics = append(metrics, fakeTestMetric())
				return metrics
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

			},
			expectedMetrics: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					CTemperature(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})
			},
			expectedMetrics: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.ctemperature,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics())
		})
	}
}

func TestTemperature(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter)
		expectedMetrics func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					Temperature(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric()})
			},
			expectedMetrics: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				metrics = append(metrics, fakeTestMetric())
				return metrics
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

			},
			expectedMetrics: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					Temperature(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})
			},
			expectedMetrics: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.temperature,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics())
		})
	}
}

func TestEthernetInfo(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				// add EthernetInfoMetricSerDesAvailability metric with it's correct labels
				m := fakeTestMetric("state")
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixEthernetInfo, bmcmonitoring.EthernetInfoMetricSerDesAvailability)
				metrics = append(metrics, m)

				// add EthernetInfoMetricANLTStatus metric with it's correct labels
				m = fakeTestMetric("state", "port")
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixEthernetInfo, bmcmonitoring.EthernetInfoMetricANLTStatus)
				metrics = append(metrics, m)

				exporter.EXPECT().
					EthernetInfo(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					EthernetInfo(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.ethernetInfo,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestEthernetStatus(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				// add EthernetStatusMetricMACRemote metric with it's correct labels
				m := fakeTestMetric("port")
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixEthernetStatus, bmcmonitoring.EthernetStatusMetricMACRemote)
				metrics = append(metrics, m)

				// add EthernetStatusMetricExternalLinkStatus metric with it's correct labels
				m = fakeTestMetric("state", "link")
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixEthernetStatus, bmcmonitoring.EthernetStatusMetricExternalLinkStatus)
				metrics = append(metrics, m)

				// add EthernetStatusMetricPHYStatus metric with it's correct labels
				m = fakeTestMetric("state", "phy")
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixEthernetStatus, bmcmonitoring.EthernetStatusMetricPHYStatus)
				metrics = append(metrics, m)

				// add EthernetStatusMetricPortMapping metric with it's correct labels
				m = fakeTestMetric("port", "type")
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixEthernetStatus, bmcmonitoring.EthernetStatusMetricPortMapping)
				metrics = append(metrics, m)

				exporter.EXPECT().
					EthernetStatus(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					EthernetStatus(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.ethernetStatus,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestEthernetStatusCounters(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric("port"))

				exporter.EXPECT().
					EthernetStatusCounters(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					EthernetStatusCounters(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.ethernetStatusCounters,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestPcieInfo(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				// Add PcieInfoMetricMaxPCIeLinkSpeed metric
				m := fakeTestMetric(bmcmonitoring.PcieInfoMetricMaxPCIeLinkSpeed)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixPcieInfo, bmcmonitoring.PcieInfoMetricMaxPCIeLinkSpeed)
				metrics = append(metrics, m)

				// Add PcieInfoMetricCurrentPCIeLinkSpeed metric
				m = fakeTestMetric(bmcmonitoring.PcieInfoMetricCurrentPCIeLinkSpeed)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixPcieInfo, bmcmonitoring.PcieInfoMetricCurrentPCIeLinkSpeed)
				metrics = append(metrics, m)

				// Add PcieInfoMetricPCIeSubsystemVendorID metric
				m = fakeTestMetric(bmcmonitoring.PcieInfoMetricPCIeSubsystemVendorID)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixPcieInfo, bmcmonitoring.PcieInfoMetricPCIeSubsystemVendorID)
				metrics = append(metrics, m)

				// Add PcieInfoMetricPCIeBusAndDevice metric
				m = fakeTestMetric(bmcmonitoring.PcieInfoMetricPCIeBusAndDevice)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixPcieInfo, bmcmonitoring.PcieInfoMetricPCIeBusAndDevice)
				metrics = append(metrics, m)

				exporter.EXPECT().
					PcieInfo(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					PcieInfo(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.pcieInfo,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestStatus(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				// Add StatusMetricChipStatus metric
				m := fakeTestMetric(bmcmonitoring.StatusMetricChipStatus)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixStatus, bmcmonitoring.StatusMetricChipStatus)
				metrics = append(metrics, m)

				// Add StatusMetricBootStage metric
				m = fakeTestMetric(bmcmonitoring.StatusMetricBootStage)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixStatus, bmcmonitoring.StatusMetricBootStage)
				metrics = append(metrics, m)

				// Add StatusMetricClockThrottling metric
				m = fakeTestMetric(bmcmonitoring.StatusMetricClockThrottling)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixStatus, bmcmonitoring.StatusMetricClockThrottling)
				metrics = append(metrics, m)

				// Add StatusMetricPowerState metric
				m = fakeTestMetric(bmcmonitoring.StatusMetricPowerState)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixStatus, bmcmonitoring.StatusMetricPowerState)
				metrics = append(metrics, m)

				// Add StatusMetricDeviceActivity metric
				m = fakeTestMetric(bmcmonitoring.StatusMetricDeviceActivity)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixStatus, bmcmonitoring.StatusMetricDeviceActivity)
				metrics = append(metrics, m)

				// Add StatusMetricEmergencyPowerReduction metric
				m = fakeTestMetric(bmcmonitoring.StatusMetricEmergencyPowerReduction)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixStatus, bmcmonitoring.StatusMetricEmergencyPowerReduction)
				metrics = append(metrics, m)

				exporter.EXPECT().
					Status(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					Status(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.status,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestFrequency(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				exporter.EXPECT().
					Frequency(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					Frequency(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.frequency,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestPower(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				exporter.EXPECT().
					Power(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					Power(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.power,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestSensorTemperature(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				exporter.EXPECT().
					SensorTemperature(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					SensorTemperature(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.sensorTemperature,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestSensorVoltage(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				exporter.EXPECT().
					SensorVoltage(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					SensorVoltage(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.sensorVoltage,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestSensorCurrent(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				exporter.EXPECT().
					SensorCurrent(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					SensorCurrent(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.sensorCurrent,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestInfo(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				// Add InfoMetricUUID metric
				m := fakeTestMetric(bmcmonitoring.InfoMetricUUID)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixInfo, bmcmonitoring.InfoMetricUUID)
				metrics = append(metrics, m)

				// Add InfoMetricDeviceID metric
				m = fakeTestMetric(bmcmonitoring.InfoMetricDeviceID)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixInfo, bmcmonitoring.InfoMetricDeviceID)
				metrics = append(metrics, m)

				// Add InfoMetricSubsystemDeviceID metric
				m = fakeTestMetric(bmcmonitoring.InfoMetricSubsystemDeviceID)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixInfo, bmcmonitoring.InfoMetricSubsystemDeviceID)
				metrics = append(metrics, m)

				// Add InfoMetricSubsystemVendorID metric
				m = fakeTestMetric(bmcmonitoring.InfoMetricSubsystemVendorID)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixInfo, bmcmonitoring.InfoMetricSubsystemVendorID)
				metrics = append(metrics, m)

				// Add InfoMetricASICSerialNumber metric
				m = fakeTestMetric(bmcmonitoring.InfoMetricASICSerialNumber)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixInfo, bmcmonitoring.InfoMetricASICSerialNumber)
				metrics = append(metrics, m)

				// Add InfoMetricBoardSerialNumber metric
				m = fakeTestMetric(bmcmonitoring.InfoMetricBoardSerialNumber)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixInfo, bmcmonitoring.InfoMetricBoardSerialNumber)
				metrics = append(metrics, m)

				exporter.EXPECT().
					Info(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					Info(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.info,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestSecurity(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				// Add SecurityMetricCPLDVersion metric
				m := fakeTestMetric(bmcmonitoring.SecurityMetricCPLDVersion)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixSecurity, bmcmonitoring.SecurityMetricCPLDVersion)
				metrics = append(metrics, m)

				// Add SecurityMetricMinimalSVNindex metric
				m = fakeTestMetric(bmcmonitoring.SecurityMetricMinimalSVNindex)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixSecurity, bmcmonitoring.SecurityMetricMinimalSVNindex)
				metrics = append(metrics, m)

				// Add SecurityMetricCPLDVersionTimestamp metric
				m = fakeTestMetric(bmcmonitoring.SecurityMetricCPLDVersionTimestamp)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixSecurity, bmcmonitoring.SecurityMetricCPLDVersionTimestamp)
				metrics = append(metrics, m)

				// Add SecurityMetricFWImageSource metric
				m = fakeTestMetric(bmcmonitoring.SecurityMetricFWImageSource)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixSecurity, bmcmonitoring.SecurityMetricFWImageSource)
				metrics = append(metrics, m)

				// Add SecurityMetricKey0Revocation metric
				m = fakeTestMetric(bmcmonitoring.SecurityMetricKey0Revocation)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixSecurity, bmcmonitoring.SecurityMetricKey0Revocation)
				metrics = append(metrics, m)

				// Add SecurityMetricKey1Revocation metric
				m = fakeTestMetric(bmcmonitoring.SecurityMetricKey1Revocation)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixSecurity, bmcmonitoring.SecurityMetricKey1Revocation)
				metrics = append(metrics, m)

				// Add SecurityMetricKey2Revocation metric
				m = fakeTestMetric(bmcmonitoring.SecurityMetricKey2Revocation)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixSecurity, bmcmonitoring.SecurityMetricKey2Revocation)
				metrics = append(metrics, m)

				// Add SecurityMetricKey3Revocation metric
				m = fakeTestMetric(bmcmonitoring.SecurityMetricKey3Revocation)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixSecurity, bmcmonitoring.SecurityMetricKey3Revocation)
				metrics = append(metrics, m)

				// Add SecurityMetricKey4Revocation metric
				m = fakeTestMetric(bmcmonitoring.SecurityMetricKey4Revocation)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixSecurity, bmcmonitoring.SecurityMetricKey4Revocation)
				metrics = append(metrics, m)

				exporter.EXPECT().
					Security(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					Security(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.security,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestDirect(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				// Add DirectMetricPCIeVendorID metric
				m := fakeTestMetric(bmcmonitoring.DirectMetricPCIeVendorID)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixDirect, bmcmonitoring.DirectMetricPCIeVendorID)
				metrics = append(metrics, m)

				// Add DirectMetricASICSerialNumber metric
				m = fakeTestMetric(bmcmonitoring.DirectMetricASICSerialNumber)
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixDirect, bmcmonitoring.DirectMetricASICSerialNumber)
				metrics = append(metrics, m)

				// Add DirectMetricOSStage metric
				m = fakeTestMetric("stage")
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixDirect, bmcmonitoring.DirectMetricOSStage)
				metrics = append(metrics, m)

				exporter.EXPECT().
					Direct(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					Direct(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.direct,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestSensorVoltageMonitor(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				exporter.EXPECT().
					SensorVoltageMonitor(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					SensorVoltageMonitor(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.sensorVoltageMonitor,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestHbm(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

					// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				// Add HBMMetricMbistRepair metric
				m := fakeTestMetric("state", "index")
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixHbm, bmcmonitoring.HBMMetricMbistRepair)
				metrics = append(metrics, m)

				// Add HBMMetricGlobalECC metric
				m = fakeTestMetric("index")
				m.MetricName = fmt.Sprintf("%s_%s", bmcmonitoring.PrefixHbm, bmcmonitoring.HBMMetricGlobalECC)
				metrics = append(metrics, m)

				exporter.EXPECT().
					Hbm(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Invalid token",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("not valid", false)

				return nil

			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {
				exporter.EXPECT().
					ValidToken(gomock.Any(), gomock.Any()).
					Times(1).
					Return("", true)

				exporter.EXPECT().
					Hbm(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.hbm,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestBmcState(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric())

				exporter.EXPECT().
					BmcState(gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.bmcState,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func TestAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	type testCase struct {
		name            string
		buildStubs      func(exporter *mock.MockExporter) []bmcmonitoring.Metric
		expectedMetrics func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Returns valid metrics",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				// add metrics for testing
				var metrics []bmcmonitoring.Metric
				metrics = append(metrics, fakeTestMetric("description", "error_type", "time", "severity"))

				exporter.EXPECT().
					Alerts(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(metrics)

				return metrics
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				return ex
			},
		},
		{
			name: "Unexpected labels in metric",
			buildStubs: func(exporter *mock.MockExporter) []bmcmonitoring.Metric {

				exporter.EXPECT().
					Alerts(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return([]bmcmonitoring.Metric{fakeTestMetric("some_label")})

				return nil
			},
			expectedMetrics: func(ex []bmcmonitoring.Metric) []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			exporter := mock.NewMockExporter(ctrl)

			exporter.EXPECT().Name().AnyTimes()
			app := application{
				log:      log,
				store:    bmcmonitoring.NewExporterStore([]bmcmonitoring.Exporter{exporter}, log),
				exporter: "none",
			}

			ex := test.buildStubs(exporter)

			server := web.NewServer(context.Background(), web.ServerOpts{}, checkgrp.Handlers{},
				[]web.WebHandler{
					{
						Route:   "/metrics",
						Handler: app.alerts,
					},
				})
			recorder := httptest.NewRecorder()

			request, err := http.NewRequest(http.MethodGet, "/metrics", nil)
			require.NoError(t, err)

			server.ServeHTTP(recorder, request)

			require.Equal(t, http.StatusOK, recorder.Code)

			checkBmcMetrics(t, recorder, test.expectedMetrics(ex))
		})
	}
}

func checkBmcMetrics(t *testing.T, recorder *httptest.ResponseRecorder, expected []bmcmonitoring.Metric) {
	if expected == nil {
		require.Empty(t, recorder.Body)
		return
	}

	require.NotEmpty(t, recorder.Body)

	body, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)

	// write the expected metrics to the buffer,
	//  and compare with the result got from the body of the recorder
	buf := bytes.Buffer{}

	for _, metric := range expected {
		err = util.Print(&metric, &buf)
		require.NoError(t, err)
	}

	require.Equal(t, strings.ReplaceAll(buf.String(), `\`, ""), string(body))
}
