package rasmonitoringapi

import (
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"habana_bmc_exporter/pkg/ipmi"
	"testing"
)

func Test_oamFromSensorNum(t *testing.T) {
	type args struct {
		a ipmi.SelInfo
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test oam0",
			args: args{
				a: ipmi.SelInfo{
					SensorNumber: "14",
				},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "test oam7",
			args: args{
				a: ipmi.SelInfo{
					SensorNumber: "1b",
				},
			},
			want:    7,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := oamFromSensorNum(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("oamFromSensorNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("oamFromSensorNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_alertTypeFromEventData(t *testing.T) {
	type args struct {
		a ipmi.SelInfo
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test double ecc error",
			args: args{
				a: ipmi.SelInfo{
					EventData: "050000",
				},
			},
			want:    rasmonitoring.AlertTypeDoubleECCErr,
			wantErr: false,
		},
		{
			name: "test Security Agent Uboot Err",
			args: args{
				a: ipmi.SelInfo{
					EventData: "2b0000",
				},
			},
			want:    rasmonitoring.AlertTypeSecurityAgentUbootErr,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := alertTypeFromEventData(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("alertTypeFromEventData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("alertTypeFromEventData() = %v, want %v", got, tt.want)
			}
		})
	}
}
