package bmcmonitoring

import (
	"fmt"
	"testing"
)

func Test_basicVerification(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Empty hostname",
			args: args{
				metric: Metric{
					Oam:        "0",
					MetricName: "some name",
				},
			},
			wantErr: true,
		},
		{
			name: "Empty metric name",
			args: args{
				metric: Metric{
					Oam:      "0",
					Hostname: "some hostname",
				},
			},
			wantErr: true,
		},
		{
			name: "With metric name and hostname",
			args: args{
				metric: Metric{
					Oam:        "0",
					Hostname:   "some hostname",
					MetricName: "some name",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := basicVerification(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("basicVerification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyDirect(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without label", DirectMetricPCIeVendorID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricPCIeVendorID),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", DirectMetricPCIeVendorID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricPCIeVendorID),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						DirectMetricPCIeVendorID: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", DirectMetricASICSerialNumber),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricASICSerialNumber),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", DirectMetricASICSerialNumber),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricASICSerialNumber),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						DirectMetricASICSerialNumber: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", DirectMetricFWVersionMajor),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricFWVersionMajor),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", DirectMetricFWVersionMinor),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricFWVersionMinor),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", DirectMetricFWVersionPatch),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricFWVersionPatch),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", DirectMetricCoreVDD),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricCoreVDD),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", DirectMetricHBMVDDq),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricHBMVDDq),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", DirectMetricV12),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricV12),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", DirectMetricApiVersion),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricApiVersion),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},

		{
			name: fmt.Sprintf("%s without label", DirectMetricIBAccessState),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricIBAccessState),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", DirectMetricIBAccessState),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricIBAccessState),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"state": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", DirectMetricOOBAccessState),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricOOBAccessState),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", DirectMetricOOBAccessState),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixDirect, DirectMetricOOBAccessState),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"state": "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyDirect(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyDirect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyInfo(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without label", InfoMetricDeviceID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricDeviceID),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", InfoMetricDeviceID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricDeviceID),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						InfoMetricDeviceID: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", InfoMetricSubsystemDeviceID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricSubsystemDeviceID),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", InfoMetricSubsystemDeviceID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricSubsystemDeviceID),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						InfoMetricSubsystemDeviceID: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", InfoMetricSubsystemVendorID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricSubsystemVendorID),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", InfoMetricSubsystemVendorID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricSubsystemVendorID),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						InfoMetricSubsystemVendorID: "",
					},
				},
			},
			wantErr: false,
		},

		{
			name: fmt.Sprintf("%s without label", InfoMetricASICSerialNumber),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricASICSerialNumber),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", InfoMetricASICSerialNumber),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricASICSerialNumber),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						InfoMetricASICSerialNumber: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", InfoMetricBoardSerialNumber),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricBoardSerialNumber),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", InfoMetricBoardSerialNumber),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricBoardSerialNumber),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						InfoMetricBoardSerialNumber: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", InfoMetricSRAMSize),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricSRAMSize),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", InfoMetricHBMSize),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricHBMSize),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", InfoMetricUUID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricUUID),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", InfoMetricUUID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixInfo, InfoMetricUUID),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						InfoMetricUUID: "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyInfo(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyStatus(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without label", StatusMetricBootStage),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricBootStage),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", StatusMetricBootStage),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricBootStage),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						StatusMetricBootStage: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", StatusMetricEmergencyPowerReduction),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricEmergencyPowerReduction),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", StatusMetricEmergencyPowerReduction),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricEmergencyPowerReduction),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						StatusMetricEmergencyPowerReduction: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", StatusMetricClockThrottling),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricClockThrottling),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", StatusMetricClockThrottling),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricClockThrottling),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						StatusMetricClockThrottling: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", StatusMetricLastClockThrottlingDuration),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricLastClockThrottlingDuration),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", StatusMetricPowerState),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricPowerState),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", StatusMetricPowerState),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricPowerState),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						StatusMetricPowerState: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", StatusMetricTotalClockThrottlingDuration),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricTotalClockThrottlingDuration),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", StatusMetricGlobalTimeFromReset),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricGlobalTimeFromReset),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", StatusMetricChipStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricChipStatus),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", StatusMetricChipStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricChipStatus),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						StatusMetricChipStatus: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", StatusMetricDeviceActivity),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricDeviceActivity),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", StatusMetricDeviceActivity),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricDeviceActivity),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						StatusMetricDeviceActivity: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", StatusMetricDeviceActivityCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixStatus, StatusMetricDeviceActivityCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyStatus(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyEthernetInfo(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without label", EthernetInfoMetricSerDesAvailability),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetInfo, EthernetInfoMetricSerDesAvailability),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetInfoMetricSerDesAvailability),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetInfo, EthernetInfoMetricSerDesAvailability),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"state": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", EthernetInfoMetricPortMaxSpeed),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetInfo, EthernetInfoMetricPortMaxSpeed),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetInfoMetricANLTStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetInfo, EthernetInfoMetricANLTStatus),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetInfoMetricANLTStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetInfo, EthernetInfoMetricANLTStatus),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"state": "",
						"port":  "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", EthernetInfoMetricNumberOfLanes),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetInfo, EthernetInfoMetricNumberOfLanes),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", EthernetInfoMetricNumberOfLinks),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetInfo, EthernetInfoMetricNumberOfLinks),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", EthernetInfoMetricLinkSpeed),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetInfo, EthernetInfoMetricLinkSpeed),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyEthernetInfo(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyEthernetInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyEthernetStatus(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricExternalLinkStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricExternalLinkStatus),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricExternalLinkStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricExternalLinkStatus),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"state": "",
						"link":  "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricLinkStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricLinkStatus),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricLinkStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricLinkStatus),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"state": "",
						"link":  "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricPHYStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricPHYStatus),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricPHYStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricPHYStatus),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"state": "",
						"phy":   "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricPortMapping),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricPortMapping),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricPortMapping),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricPortMapping),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
						"type": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricBERCorrectable),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricBERCorrectable),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricBERCorrectable),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricBERCorrectable),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricBERUncorrectable),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricBERUncorrectable),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricBERUncorrectable),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricBERUncorrectable),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricNack),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricNack),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricNack),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricNack),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricCRC),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricCRC),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricCRC),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricCRC),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricRetransmissionTimeout),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricRetransmissionTimeout),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricRetransmissionTimeout),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricRetransmissionTimeout),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricLinkRetrainingDueToBER),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricLinkRetrainingDueToBER),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricLinkRetrainingDueToBER),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricLinkRetrainingDueToBER),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricStateTogglingCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricStateTogglingCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricStateTogglingCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricStateTogglingCounter),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricMACRemote),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricMACRemote),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricMACRemote),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricMACRemote),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricRetransmission),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricRetransmission),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricRetransmission),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricRetransmission),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},

		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricRetraining),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricRetraining),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricRetraining),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricRetraining),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricSERPreFEC),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricSERPreFEC),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricSERPreFEC),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricSERPreFEC),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricSERPostFEC),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricSERPostFEC),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricSERPostFEC),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricSERPostFEC),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricLatency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricLatency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricLatency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricLatency),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", EthernetStatusMetricThroughput),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricThroughput),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", EthernetStatusMetricThroughput),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixEthernetStatus, EthernetStatusMetricThroughput),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"port": "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyEthernetStatus(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyEthernetStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyPcieInfo(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without label", PcieInfoMetricMaxPCIeLinkSpeed),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricMaxPCIeLinkSpeed),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", PcieInfoMetricMaxPCIeLinkSpeed),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricMaxPCIeLinkSpeed),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						PcieInfoMetricMaxPCIeLinkSpeed: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", PcieInfoMetricCurrentPCIeLinkSpeed),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricCurrentPCIeLinkSpeed),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", PcieInfoMetricCurrentPCIeLinkSpeed),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricCurrentPCIeLinkSpeed),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						PcieInfoMetricCurrentPCIeLinkSpeed: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricMaxPCIeLinkWidth),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricMaxPCIeLinkWidth),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricCurrentPCIeLinkWidth),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricCurrentPCIeLinkWidth),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricPCIeDeviceID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeDeviceID),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricPCIeSubsystemID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeSubsystemID),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", PcieInfoMetricPCIeSubsystemVendorID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeSubsystemVendorID),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", PcieInfoMetricPCIeSubsystemVendorID),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeSubsystemVendorID),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						PcieInfoMetricPCIeSubsystemVendorID: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", PcieInfoMetricPCIeBusAndDevice),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeBusAndDevice),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", PcieInfoMetricPCIeBusAndDevice),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeBusAndDevice),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						PcieInfoMetricPCIeBusAndDevice: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricCorrectedInternalErrorStatus),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricCorrectedInternalErrorStatus),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricReplayBufferNumRolloverError),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricReplayBufferNumRolloverError),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricReplayTimerTimeoutError),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricReplayTimerTimeoutError),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricBadTLPCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricBadTLPCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricBadDLLPCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricBadDLLPCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricReceiverErrorCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricReceiverErrorCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricLCRCErrorCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricLCRCErrorCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricECRCErrorCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricECRCErrorCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricCompletionTimeoutIndication),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricCompletionTimeoutIndication),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricUncorrectableInternalErrorIndication),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricUncorrectableInternalErrorIndication),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricReceiverOverflowIndication),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricReceiverOverflowIndication),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricFlowControlProtocolErrorIndication),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricFlowControlProtocolErrorIndication),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricSurpriseLinkDownIndication),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricSurpriseLinkDownIndication),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricMalfunctionTLPErrorIndication),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricMalfunctionTLPErrorIndication),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricDLLPProtocolErrorIndication),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricDLLPProtocolErrorIndication),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricRXNakDLLPCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricRXNakDLLPCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricTxNakDLLPCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricTxNakDLLPCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricRetryTLPcounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricRetryTLPcounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricPWRBRKindication),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPWRBRKindication),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricPCIeRXMemoryWriteCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeRXMemoryWriteCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricPCIeRXMemoryReadCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeRXMemoryReadCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricPCIeTXMemoryWriteCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeTXMemoryWriteCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricPCIeTXMemoryReadCounter),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeTXMemoryReadCounter),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricAERCapabilityControlOffset),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricAERCapabilityControlOffset),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricAERerrorlog),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricAERerrorlog),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PcieInfoMetricPCIeFWversion),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPcieInfo, PcieInfoMetricPCIeFWversion),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyPcieInfo(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyPcieInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifySecurity(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without labels", SecurityMetricCurrentPublicKeyHashIndex),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricCurrentPublicKeyHashIndex),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SecurityMetricCurrentSVNversion),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricCurrentSVNversion),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", SecurityMetricKey0Revocation),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricKey0Revocation),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", SecurityMetricKey0Revocation),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricKey0Revocation),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						SecurityMetricKey0Revocation: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", SecurityMetricKey1Revocation),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricKey1Revocation),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", SecurityMetricKey1Revocation),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricKey1Revocation),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						SecurityMetricKey1Revocation: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", SecurityMetricKey2Revocation),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricKey2Revocation),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", SecurityMetricKey2Revocation),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricKey2Revocation),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						SecurityMetricKey2Revocation: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", SecurityMetricKey3Revocation),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricKey3Revocation),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", SecurityMetricKey3Revocation),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricKey3Revocation),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						SecurityMetricKey3Revocation: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", SecurityMetricKey4Revocation),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricKey4Revocation),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", SecurityMetricKey4Revocation),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricKey4Revocation),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						SecurityMetricKey4Revocation: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", SecurityMetricMinimalSVNindex),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricMinimalSVNindex),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", SecurityMetricMinimalSVNindex),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricMinimalSVNindex),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						SecurityMetricMinimalSVNindex: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", SecurityMetricFWImageSource),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricFWImageSource),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", SecurityMetricFWImageSource),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricFWImageSource),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						SecurityMetricFWImageSource: "",
					},
				},
			},
			wantErr: false,
		},

		{
			name: fmt.Sprintf("%s without label", SecurityMetricCPLDVersion),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricCPLDVersion),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", SecurityMetricCPLDVersion),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricCPLDVersion),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						SecurityMetricCPLDVersion: "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", SecurityMetricCPLDVersionTimestamp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricCPLDVersionTimestamp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", SecurityMetricCPLDVersionTimestamp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSecurity, SecurityMetricCPLDVersionTimestamp),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						SecurityMetricCPLDVersionTimestamp: "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifySecurity(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifySecurity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyTemperature(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricCurrentBoardTemp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricCurrentBoardTemp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricCurrentVRMTemp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricCurrentVRMTemp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricCurrentDRAMTemp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricCurrentDRAMTemp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricCurrentOnDieTemp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricCurrentOnDieTemp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricHistoricalBoardTemp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricHistoricalBoardTemp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricHistoricalVRMTemp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricHistoricalVRMTemp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricHistoricalDRAMTemp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricHistoricalDRAMTemp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricHistoricalOnDieTemp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricHistoricalOnDieTemp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricMaxTempRiseTime),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricMaxTempRiseTime),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricMaxSocTempErrorThreshold),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricMaxSocTempErrorThreshold),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricMaxSocTempWarmingThreshold),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricMaxSocTempWarmingThreshold),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricMaxHbmTempThreshold),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricMaxHbmTempThreshold),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricCurrentSocTempErrorThreshold),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricCurrentSocTempErrorThreshold),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricCurrentSocTempWarningThreshold),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricCurrentSocTempWarningThreshold),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", TemperatureMetricCurrentHbmTempThreshold),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixTemperature, TemperatureMetricCurrentHbmTempThreshold),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyTemperature(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyTemperature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyFrequency(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricHBMFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricHBMFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricMaxTPCFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricMaxTPCFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricMaxMMEFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricMaxMMEFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricMaxDMAFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricMaxDMAFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricMaxMediaFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricMaxMediaFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricMaxPCIeFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricMaxPCIeFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricMaxARMFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricMaxARMFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricMaxNICFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricMaxNICFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricMaxNoCFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricMaxNoCFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricMaxNoCFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricMaxNoCFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricCurrentTPCFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricCurrentTPCFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricCurrentMMEFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricCurrentMMEFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricCurrentDMAFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricCurrentDMAFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricCurrentMediaFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricCurrentMediaFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricCurrentPCIeFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricCurrentPCIeFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricCurrentARMFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricCurrentARMFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricCurrentNICFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricCurrentNICFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricCurrentNoCFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricCurrentNoCFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricCurrentSRAMFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricCurrentSRAMFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", FrequencyMetricCurrentMSSFrequency),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixFrequency, FrequencyMetricCurrentMSSFrequency),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyFrequency(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyFrequency() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyPower(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without labels", PowerMetricCurrentPowerConsumption),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, PowerMetricCurrentPowerConsumption),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", PowerMetricPeakPowerConsumption),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, PowerMetricPeakPowerConsumption),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyPower(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyPower() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifySensorTemperature(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricOnDie0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricOnDie0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricOnDie1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricOnDie1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricOnDie2),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricOnDie2),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricOnDie3),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricOnDie3),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricHBM0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricHBM0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricHBM1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricHBM1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricHBM2),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricHBM2),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricHBM3),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricHBM3),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricHBM4),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricHBM4),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricHBM5),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricHBM5),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricCPLDLocal),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricCPLDLocal),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricCPLD0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricCPLD0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricCPLD1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricCPLD1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricCPLD2),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricCPLD2),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricCPLD3),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricCPLD3),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricOnboard0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricOnboard0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricOnboard1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricOnboard1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricOnboard2),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricOnboard2),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricOnboard3),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricOnboard3),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricCPLDTemp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricCPLDTemp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricPSUStage1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricPSUStage1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorTemperatureMetricPSUStage2),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixPower, SensorTemperatureMetricPSUStage2),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifySensorTemperature(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifySensorTemperature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifySensorVoltage(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVADC54),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVADC54),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVRM1in),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVRM1in),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVRM1out),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVRM1out),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVRM2in),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVRM2in),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVRM2VDDout),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVRM2VDDout),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVRM2HBMout),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVRM2HBMout),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVMONPCIEVPH1P8V),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVMONPCIEVPH1P8V),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVMON1P8HBMVAA),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVMON1P8HBMVAA),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVMON2P5),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVMON2P5),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVMON48VHIMON),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVMON48VHIMON),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVMONP5V),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVMONP5V),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVMON12V1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVMON12V1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVMONHBM),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVMONHBM),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageVMONCore),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageVMONCore),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageCPLDHIMON1P8NIC),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltage, SensorVoltageCPLDHIMON1P8NIC),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifySensorVoltage(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifySensorVoltage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifySensorCurrent(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without labels", SensorCurrentVin54),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorCurrent, SensorCurrentVin54),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorCurrentP1Vin12),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorCurrent, SensorCurrentP1Vin12),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorCurrentStage154Vin),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorCurrent, SensorCurrentStage154Vin),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorCurrentStage113P5VOut),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorCurrent, SensorCurrentStage113P5VOut),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorCurrentStage213P5Vin),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorCurrent, SensorCurrentStage213P5Vin),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorCurrentStage2CoreOut),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorCurrent, SensorCurrentStage2CoreOut),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorCurrentStage2HBMout),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorCurrent, SensorCurrentStage2HBMout),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifySensorCurrent(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifySensorCurrent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyCTemperature(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without labels", PrefixCTemperature),
			args: args{
				metric: Metric{
					MetricName: PrefixSensorCurrent,
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyCTemperature(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyCTemperature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifySensorVoltageMonitor(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwVm),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwVm),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwCpeEu_0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwCpeEu_0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwCpeEu_1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwCpeEu_1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwCpeEu_2),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwCpeEu_2),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwcpeEu_3),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwcpeEu_3),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwCpeHbm),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwCpeHbm),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwTpcSb),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwTpcSb),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwTft),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwTft),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwCnt),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwCnt),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwMcHbm0_0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwMcHbm0_0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwMcHbm0_1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwMcHbm0_1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwHconHbm0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwHconHbm0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwPpw_0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwPpw_0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwPpw_1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwPpw_1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwL2cMacro),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwL2cMacro),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSwVcdMacro),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSwVcdMacro),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeVm),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeVm),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSePe0_0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSePe0_0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSePe0_1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSePe0_1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeEuCore),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeEuCore),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeBpyramid),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeBpyramid),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeRx),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeRx),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeMx),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeMx),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeHconHbm1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeHconHbm1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeMcHbm1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeMcHbm1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeGasket),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeGasket),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeMmeCtrl),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeMmeCtrl),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeMmeQman),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeMmeQman),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeSbte),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeSbte),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeRtrDn),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeRtrDn),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeRtrUp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeRtrUp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorSeSram),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorSeSram),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwVm),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwVm),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwPe0_0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwPe0_0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwPe0_1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwPe0_1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwPe0_2),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwPe0_2),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwPe0_3),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwPe0_3),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwEuCore),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwEuCore),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwBpyramid),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwBpyramid),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwMcHbm4_0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwMcHbm4_0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwMcHbm4_1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwMcHbm4_1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwHconHbm4),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwHconHbm4),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwAcc),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwAcc),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwWap),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwWap),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwTif),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwTif),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwRtrDn),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwRtrDn),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwRtrUp),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwRtrUp),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},

		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNwSram),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNwSram),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeVM),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeVM),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeTx),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeTx),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeMx_0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeMx_0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeMx_1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeMx_1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeHconHbm5),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeHconHbm5),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeMcHbm5),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeMcHbm5),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeSob),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeSob),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeQnt),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeQnt),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeCnt_0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeCnt_0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeCnt_1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeCnt_1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeCpeEu),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeCpeEu),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeCntHbm),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeCntHbm),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeSbte_0),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeSbte_0),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeSbte_1),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeSbte_1),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeMif),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeMif),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", SensorVoltageMonitorNeRtrDn),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixSensorVoltageMonitor, SensorVoltageMonitorNeRtrDn),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifySensorVoltageMonitor(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifySensorVoltageMonitor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyHbm(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without labels", HBMMetricEccErrors),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixHbm, HBMMetricEccErrors),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", HBMMetricNumOfRepairedLanes),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixHbm, HBMMetricNumOfRepairedLanes),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without labels", HBMMetricNumOfReplacedRows),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixHbm, HBMMetricNumOfReplacedRows),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", HBMMetricMbistRepair),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixHbm, HBMMetricMbistRepair),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", HBMMetricMbistRepair),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixHbm, HBMMetricMbistRepair),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"state": "",
						"index": "",
					},
				},
			},
			wantErr: false,
		},
		{
			name: fmt.Sprintf("%s without label", HBMMetricGlobalECC),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixHbm, HBMMetricGlobalECC),
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: true,
		},
		{
			name: fmt.Sprintf("%s with label", HBMMetricGlobalECC),
			args: args{
				metric: Metric{
					MetricName: fmt.Sprintf("%s_%s", PrefixHbm, HBMMetricGlobalECC),
					Hostname:   "some hostname",
					Oam:        "0",
					CustomLabels: map[string]string{
						"index": "",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyHbm(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyHbm() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyBmcState(t *testing.T) {
	type args struct {
		metric Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: fmt.Sprintf("%s without label", PrefixBMCState),
			args: args{
				metric: Metric{
					MetricName: PrefixBMCState,
					Hostname:   "some hostname",
					Oam:        "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyBmcState(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("verifyBmcState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
