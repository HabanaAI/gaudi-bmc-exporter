package redfish

import (
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSramDoubleEccAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *SRAMDoubleECCAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Red fish returns valid metrics",
			resp: &SRAMDoubleECCAlert{
				{
					BlockIndex: 102, IndexInBlock: 100, Address: 10, Syndrome: 10, Timestamp: 1632294282,
				},
				{
					BlockIndex: 101, IndexInBlock: 10, Address: 250, Syndrome: 11, Timestamp: 1632294282,
				},
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "SRAM Double ECC Block Error. block index: 102, index in block: 100, address: 10, syndrome: 10",
						"error_type":  bmcmonitoring.AlertTypeDoubleECCBlockErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "SRAM Double ECC Block Error. block index: 101, index in block: 10, address: 250, syndrome: 11",
						"error_type":  bmcmonitoring.AlertTypeDoubleECCBlockErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := sramDoubleEccAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestSramCorrectableAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *SRAMCorrectableAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				return metrics
			},
		},
		{
			name: "Red fish returns valid metrics",
			resp: &SRAMCorrectableAlert{
				SingleECC: []SingleEcc{
					{
						BlockIndex: 1,
						Index:      100,
						Address:    10,
						Syndrome:   10,
						Timestamp:  1632294282,
					},
				},
				MultipleECC: []MultipleECC{
					{
						BlockIndex: 2,
						Index:      100,
						Address:    250,
						Cause:      "Multi SERR",
						Timestamp:  1632294282,
					},
				},
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "SRAM Single ECC Block Error. block index: 1,index: 100, address: 10, syndrome: 10",
						"error_type":  bmcmonitoring.AlertTypeSingleECCBlockErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Fatal SRAM ECC. block index: 2,index: 100, address: 250, cause: Multi SERR",
						"error_type":  bmcmonitoring.AlertTypeFatalSramEccErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Fatal",
					},
				})
				return metrics
			},
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := sramCorrectableAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestPcieAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *PCIeAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				return metrics
			},
		},
		{
			name: "No alerts",
			resp: &PCIeAlert{
				Fatal:       "False",
				Correctable: "False",
				NonFatal:    "False",
				Timestamp:   1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				return metrics
			},
		},
		{
			name: "Red fish returns valid metrics",
			resp: &PCIeAlert{
				Fatal:       "True",
				Correctable: "True",
				NonFatal:    "True",
				Timestamp:   1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// Fatal alerts

				// PCIE_UNCORRECTABLE_FATAL_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Uncorrectable Fatal Error",
						"error_type":  bmcmonitoring.AlertTypePcieUncorrectableFatalErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// AlertTypePcieReceiverErr
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Receiver Error",
						"error_type":  bmcmonitoring.AlertTypePcieReceiverErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// PCIE_RECEIVER_OVERFLOW_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Receiver Overflow",
						"error_type":  bmcmonitoring.AlertTypePcieReceiverOverflowErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// PCIE_FLOW_CONTROL_PROTOCOL_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Flow Control Protocol Error",
						"error_type":  bmcmonitoring.AlertTypePcieFlowControlProtocolErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// PCIE_SURPRISE_LINK_DOWN_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Surprise Link Down Error",
						"error_type":  bmcmonitoring.AlertTypePcieSurpriseLinkDownErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// PCIE_MALFUNCTION_TLP_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Malfunction TLP Error",
						"error_type":  bmcmonitoring.AlertTypePcieMalfunctionTlpErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// PCIE_DLLP_PROTOCOL_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe DLLP Protocol Error",
						"error_type":  bmcmonitoring.AlertTypePcieDllpProtocolErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// Non fatal alerts

				// PCIE_ECRC_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe ECRC Error",
						"error_type":  bmcmonitoring.AlertTypePcieEcrcErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// PCIE_COMPLETION_TIMEOUT_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Completion Timeout Error",
						"error_type":  bmcmonitoring.AlertTypePcieCompletionTimeoutErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// Correctable alerts

				// PCIE_CORRECTED_INTERNAL_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Corrected Internal Error",
						"error_type":  bmcmonitoring.AlertTypePcieCorrectedInternalErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// PCIE_REPLAY_BUFFER_ROLLOVER_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Replay Buffer Rollover Error",
						"error_type":  bmcmonitoring.AlertTypePcieReplayBufferRolloverErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// PCIE_REPLAY_TIMER_TIMEOUT_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Replay Timer Timeout Error",
						"error_type":  bmcmonitoring.AlertTypePcieReplayTimerTimeoutErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// PCIE_BAD_DLLP_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Bad DLLP Error",
						"error_type":  bmcmonitoring.AlertTypePcieBadDllpErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// PCIE_BAD_TLP_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Bad TLP Error",
						"error_type":  bmcmonitoring.AlertTypePcieBadTlpErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// PCIE_RECEIVER_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Receiver Error",
						"error_type":  bmcmonitoring.AlertTypePcieReceiverErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// PCIE_LCRD_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe LCRD Error",
						"error_type":  bmcmonitoring.AlertTypePcieLcrdErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// PCIE_RX_NAK_DLLP_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe RX NAK DLLP Error",
						"error_type":  bmcmonitoring.AlertTypePcieRxNakDllpErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// PCIE_TX_NAK_DLLP_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe TX NAK DLLP Error",
						"error_type":  bmcmonitoring.AlertTypePcieTxNakDllpErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// PCIE_RETRY_TLP_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PCIe Retry TLP Error",
						"error_type":  bmcmonitoring.AlertTypePcieRetryTlpErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})
				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := pcieAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestHbmAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *HBMAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				return metrics
			},
		},
		{
			name: "Red fish returns valid alerts",
			resp: &HBMAlert{
				DERR: []Derr{
					{
						HBMIndex:  1,
						Address:   10,
						Syndrome:  11,
						PCIndex:   1,
						Timestamp: 1632294282,
					},
				},
				MultiSERR: []MultiSERR{
					{
						HBMIndex:  2,
						Address:   1000,
						Syndrome:  1,
						PCIndex:   3,
						Timestamp: 1632294282,
					},
				},
				CATTRIP: "True",
				Parity: []Parity{
					{
						HBMIndex:  0,
						PCIndex:   1,
						Parity:    "Command/Address",
						Timestamp: 1632294282,
					},
				},
				SameAddressMultiSERR: []SameAddressMultiSERR{
					{
						HBMIndex:  1,
						PCIndex:   2,
						SIDIndex:  3,
						BankIndex: 4,
						Address:   5,
						Count:     6,
						Timestamp: 1632294282,
					},
				},
				Timestamp: 1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// DOUBLE_ECC_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "HBM Double ECC Error. hbm index: 1, address: 10,syndrome: 11, pci index: 1",
						"error_type":  bmcmonitoring.AlertTypeDoubleECCErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// SINGLE_ECC_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "HBM Single ECC Error. hbm index: 2, address: 1000,syndrome: 1, pci index: 3",
						"error_type":  bmcmonitoring.AlertTypeSingleECCErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// HBM_CATASTROPHIC_TEMP_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "HBM Catastrophic Temperature",
						"error_type":  bmcmonitoring.AlertTypeHbmCatastrophicTempErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Fatal",
					},
				})

				// FATAL_HBM_PARITY_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Fatal HBM Parity. hbm index: 0, pci index: 1, parity: Command/Address",
						"error_type":  bmcmonitoring.AlertTypeFatalHbmParityErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Fatal",
					},
				})

				// FATAL_HBM_ECC_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Fatal HBM ECC. hbm index: 1, pci index: 2, sid index: 3, bank index: 4, address: 5, count: 6",
						"error_type":  bmcmonitoring.AlertTypeFatalHbmEccErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Fatal",
					},
				})
				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := hbmAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)
			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestWdCauseAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *WDAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Red fish returns valid alerts",
			resp: &WDAlert{
				Heartbeat: "True",
				TDR:       "True",
				HW:        "True",
				Timestamp: 1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// WD_HEARTBEAT_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Heartbeat Watchdog Timed Out",
						"error_type":  bmcmonitoring.AlertTypeWdHeartbeatErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// WD_TDR_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "TDR Watchdog Timed Out",
						"error_type":  bmcmonitoring.AlertTypeWdTdrErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// WD_HW_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Firmware Watchdog Timed Out",
						"error_type":  bmcmonitoring.AlertTypeWdHwErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})
				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := wdCauseAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestTemperatureAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *TemperatureAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Red fish returns valid alerts",
			resp: &TemperatureAlert{
				ThermalTemperatureCrossing: &ThermalTemperatureCrossing{
					Crossed:       "True",
					MaxHistorical: 120,
				},
				RiseTimeViolation: "True",
				Timestamp:         1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// THERMAL_THRESHOLD_CROSSING_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Thermal Threshold Crossing. crossed: True, max historical: 120",
						"error_type":  bmcmonitoring.AlertTypeThermalThresholdCrossingErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// THERMAL_TEMP_RISE_TIME_VIOLATION_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Thermal Rise Time Violation",
						"error_type":  bmcmonitoring.AlertTypeThermalTempRiseTimeViolationErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})
				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := temperatureAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestSecurityAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *SecurityAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Red fish returns valid alerts",
			resp: &SecurityAlert{
				SPIFailure: &SPIFailure{
					RevokedPublicKey: "True",
					SVN:              "True",
					Signature:        "False",
				},
				AgentFailure: &AgentFailure{
					RevokedPublicKey: "True",
					SVN:              "True",
					Signature:        "True",
					UBoot:            "True",
					Linux:            "True",
					Zephyr:           "True",
				},
				PermissionAccessError: "True",
				Timestamp:             1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// SPIFailure

				// SECURITY_SPI_REVOKED_PUB_KEY_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Flash - Public Key Revoked",
						"error_type":  bmcmonitoring.AlertTypeSecuritySpiRevokedPubKeyErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// SECURITY_SPI_SVN_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Flash - SVN Error",
						"error_type":  bmcmonitoring.AlertTypeSecuritySpiSvnErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// SECURITY_SPI_SIGNATURE_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Flash - Signature Error",
						"error_type":  bmcmonitoring.AlertTypeSecuritySpiSignatureErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// AgentFailure

				// SECURITY_AGENT_REVOKED_PUB_KEY_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Agent - Public Key Revoked",
						"error_type":  bmcmonitoring.AlertTypeSecurityAgentRevokedPubKeyErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// SECURITY_AGENT_SVN_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Agent - SVN (softwart version number) mismatch",
						"error_type":  bmcmonitoring.AlertTypeSecurityAgentSvnErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// SECURITY_AGENT_SIGNATURE_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Agent - Signature validition failure",
						"error_type":  bmcmonitoring.AlertTypeSecurityAgentSignatureErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// SECURITY_AGENT_UBOOT_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Agent - uBoot failure",
						"error_type":  bmcmonitoring.AlertTypeSecurityAgentUbootErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// SECURITY_AGENT_LINUX_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Agent - Linux failure",
						"error_type":  bmcmonitoring.AlertTypeSecurityAgentLinuxErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// SECURITY_AGENT_ZEPHYR_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Agent - Zephyr failure",
						"error_type":  bmcmonitoring.AlertTypeSecurityAgentZephyrErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// SECURITY_PERMISSION_ACCESS_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Security - Permission Access",
						"error_type":  bmcmonitoring.AlertTypeSecurityPermissionAccessErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := securityAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestRmaAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *RMAAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Red fish returns valid alerts",
			resp: &RMAAlert{
				FatalSRAMDoubleECC: []FatalSRAMDoubleECC{
					{
						Index:     100,
						Address:   250,
						Cause:     "Multi SERR",
						Count:     2,
						Timestamp: 1632294282,
					},
				},
				HBMECC: []HBMECC{
					{
						HBMIndex:  1,
						PCIndex:   2,
						SIDIndex:  3,
						BankIndex: 4,
						Cause:     "Multi SERR",
						Address:   1000,
						Count:     5,
						Timestamp: 1632294282,
					},
				},
				HBMParity: []HBMParity{
					{
						HBMIndex:  1,
						PCIndex:   2,
						Cause:     "CA",
						Count:     1,
						Timestamp: 1632294282,
					},
				},
				HBMRowReplacement: "Row repair: scrubbing failed",
				HBMInitialization: "Pre-training fail",
				SRAMRepairFailure: &SRAMRepairFailure{
					NonRepairable: []NonRepairable{
						{
							RingNumber: 2,
							SubServer:  4,
						},
					},
					Repairable: []Repairable{
						{
							RingNumber: 5,
							SubServer:  6,
						},
					},
				},
				HBMMBISTUnrepairable: &HBMMBISTUnrepairable{
					InFieldFailure: "reset fail",
					HBMIndex:       1,
				},
				Timestamp: 1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// FATAL_SRAM_ECC_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Fatal SRAM ECC. index: 100, address: 250, cause: Multi SERR, count: 2",
						"error_type":  bmcmonitoring.AlertTypeFatalSramEccErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Fatal",
					},
				})

				// FATAL_HBM_ECC_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Fatal HBM ECC. hbm index: 1, pci index: 2, sid index: 3, bank index: 4, cause: Multi SERR, address: 1000, count: 5",
						"error_type":  bmcmonitoring.AlertTypeFatalHbmEccErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Fatal",
					},
				})

				// FATAL_HBM_PARITY_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Fatal HBM Parity. hbm index: 1, pci index: 2, cause: CA, count: 1",
						"error_type":  bmcmonitoring.AlertTypeFatalHbmParityErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Fatal",
					},
				})

				// FATAL_ROW_REPLACEMENT_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Fatal Row Replacement. Row repair: scrubbing failed",
						"error_type":  bmcmonitoring.AlertTypeFatalRowReplacementErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Fatal",
					},
				})

				// FATAL_HBM_INIT_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Fatal HBM Initialization. Pre-training fail",
						"error_type":  bmcmonitoring.AlertTypeFatalHbmInitErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Fatal",
					},
				})

				// SRAM_UNREPAIRABLE_RINGS_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "SRAM Rings Unrepairable. ring number: 2, sub server 4",
						"error_type":  bmcmonitoring.AlertTypeSramUnrepairableRingsErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// SRAM_REPAIRED_RINGS_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "SRAM Rings Repaired. ring number: 5, sub server 6",
						"error_type":  bmcmonitoring.AlertTypeSramRepairedRingsErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// HBM_MBIST_REPAIR_FAILURE_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "HBM MBIST Repair Failure. in field failure: reset fail, hbm index: 1",
						"error_type":  bmcmonitoring.AlertTypeHbmMbistRepairFailureErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Fatal",
					},
				})
				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := rmaAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestVoltageMonitorAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *VoltageMonitorAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Red fish returns valid alerts",
			resp: &VoltageMonitorAlert{
				Threshold1Crossing: "True",
				Threshold2Crossing: "True",
				Timestamp:          1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// VOLTAGE_MONITOR_THRESHOLD_1_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Voltage Monitor - Warning Threshold crossed",
						"error_type":  bmcmonitoring.AlertTypeVoltageMonitorThreshold1Err,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// VOLTAGE_MONITOR_THRESHOLD_2_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Voltage Monitor - Error Threshold crossed",
						"error_type":  bmcmonitoring.AlertTypeVoltageMonitorThreshold2Err,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := voltageMonitorAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestNicAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *NICAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Red fish returns valid alerts",
			resp: &NICAlert{
				HighBER:                       "True",
				MultipleNAC:                   "True",
				MultipleCRCError:              "True",
				MultipleLinkReinitializations: "True",
				MultipleLinkRetransmissions:   "True",
				Timestamp:                     1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// NIC_BER_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "NIC Bit Error Rate",
						"error_type":  bmcmonitoring.AlertTypeNicHighBerErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// NIC_MULTIPLE_NACK_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "NIC Multiple NACK (no acknowledge)",
						"error_type":  bmcmonitoring.AlertTypeNicMultipleNackErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				// NIC_MULTIPLE_CRC_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "NIC Multiple CRC",
						"error_type":  bmcmonitoring.AlertTypeNicMultipleCrcErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// NIC_MULTIPLE_LINK_RE_INITIALIZATIONS_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "NIC Multiple Link Re-initializations",
						"error_type":  bmcmonitoring.AlertTypeNicMultipleLinkReInitializationsErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Critical",
					},
				})

				// NIC_MULTIPLE_LINK_RE_TRANSMISSIONS_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "NIC Multiple Link Re-transmissions",
						"error_type":  bmcmonitoring.AlertTypeNicMultipleLinkReTransmissionsErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Notice",
					},
				})

				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := nicAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestChecksumAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *ChecksumAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Red fish returns valid alerts",
			resp: &ChecksumAlert{
				RMATables: &RMATables{
					HBMParity:            "Fail",
					HBMECC:               "Fail",
					SRAMECC:              "Fail",
					RowReplacement:       "Fail",
					RMACauseFailure:      "Fail",
					SpareRowAvailability: "Fail",
				},
				SPIChecksum: &SPIChecksum{
					PpBootPrimary:    "Fail",
					PpBootSecondary:  "Fail",
					PreBootPrimary:   "Fail",
					PreBootSecondary: "Fail",
				},
				CPLDStatus: "Fail",
				EEPROM:     "I/O fail",
				Timestamp:  1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// CHECKSUM_HBM_RMA_PERITY_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - HBM Parity table section",
						"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaPerityErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_HBM_RMA_ECC_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - HBM ECC table section",
						"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaEccErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_HBM_RMA_SRAM_ECC_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - SRAM ECC table section",
						"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaSramEccErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_HBM_RMA_ROW_REPLACEMENT_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - HBM Row replacement table section",
						"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaRowReplacementErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_HBM_RMA_CAUSE_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - RMA Cause table section",
						"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaCauseErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_HBM_RMA_SPARE_ROW_AVAILABILITY_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "pending embedded description",
						"error_type":  bmcmonitoring.AlertTypeChecksumHbmRmaSpareRowAvailabilityErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_PPBOOT_PRIMARY_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - Pre-PreBoot Primary mismatch",
						"error_type":  bmcmonitoring.AlertTypeChecksumPpbootPrimaryErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_PPBOOT_SECONDARY_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - Pre-PreBoot Secondary mismatch",
						"error_type":  bmcmonitoring.AlertTypeChecksumPpbootSecondaryErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_PREBOOT_PRIMARY_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - PreBoot Primary mismatch",
						"error_type":  bmcmonitoring.AlertTypeChecksumPrebootPrimaryErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_PREBOOT_SECONDARY_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - PreBoot Secondary mismatch",
						"error_type":  bmcmonitoring.AlertTypeChecksumPrebootSecondaryErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_CPLD_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - CPLD mismatch",
						"error_type":  bmcmonitoring.AlertTypeChecksumCpldErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// CHECKSUM_EEPROM_DEVICE_ACCESS_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "Checksum - EEPROM Device Access Error",
						"error_type":  bmcmonitoring.AlertTypeChecksumEepromDeviceAccessErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := checksumAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestPllAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *PLLAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Red fish returns valid alerts",
			resp: &PLLAlert{
				"Dcore0HbmPll": "Locked",
				"Dcore1HbmPll": "Not Locked",
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// PLL_LOCK_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "PLL Lock Error. Dcore1HbmPll is not locked",
						"error_type":  bmcmonitoring.AlertTypePllLockErr,
						"time":        "",
						"severity":    "Critical",
					},
				})

				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := pllAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}

func TestFitAlerts(t *testing.T) {
	log := logger.New().WithField("test", "true")

	alertsData, err := bmcmonitoring.InitAlertsDescription("./testdata/alert_description.json")
	require.NoError(t, err)

	type testCase struct {
		name     string
		resp     *FITAlert
		expected func() []bmcmonitoring.Metric
	}

	tests := []testCase{
		{
			name: "Red fish doesn't provide this field in the response",
			resp: nil,
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric
				return metrics
			},
		},
		{
			name: "Red fish returns valid alerts",
			resp: &FITAlert{
				InvalidImageFormat:              "True",
				FWUpgradeTargetAddressViolation: "True",
				FITNotRunnable:                  "True",
				Timestamp:                       1632294282,
			},
			expected: func() []bmcmonitoring.Metric {
				var metrics []bmcmonitoring.Metric

				// FIT_INVALID_IMAGE_FORMAT_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "FIT Invalid Image Format",
						"error_type":  bmcmonitoring.AlertTypeFitInvalidImageFormatErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// FIT_TARGET_ADDREES_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "FIT Target Address Error",
						"error_type":  bmcmonitoring.AlertTypeFitTargetAddreesErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})

				// FIT_NOT_RUNABLE_ERR
				metrics = append(metrics, bmcmonitoring.Metric{
					Hostname:   "hostname",
					Oam:        "0",
					MetricName: bmcmonitoring.PrefixAlerts,
					CustomLabels: map[string]string{
						"description": "FIT Not Runnable",
						"error_type":  bmcmonitoring.AlertTypeFitNotRunableErr,
						"time":        "2021-09-22 07:04:42 +0000 UTC",
						"severity":    "Error",
					},
				})
				return metrics
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			got := fitAlerts(log, alertsData, test.resp, "hostname", 0)
			require.Equal(t, test.expected(), got)

			for _, alert := range got {
				bmcmonitoring.VerifyAlerts(alert)
			}
		})
	}
}
