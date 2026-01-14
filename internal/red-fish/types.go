package redfish

import "fmt"

type DirectResponse struct {
	Response DirectResp `json:"DIRECT"`
}

type DirectResp struct {
	PCIeVendorID     string `json:"PCIeVendorID"`
	AsicSerialNumber string `json:"ASICSerialNumber"`
	OSStage          string `json:"OSStage"`
	IBAccessState    string `json:"IBAccessState"`
	OOBAccessState   string `json:"OOBAccessState"`
	FWVersionMajor   int    `json:"FWVersionMajor"`
	FWVersionMinor   int    `json:"FWVersionMinor"`
	FWVersionPatch   int    `json:"FWVersionPatch"`
	CoreVDD          int    `json:"CoreVDD"`
	HBMVddq          int    `json:"HBMVDDq"`
	V12              int    `json:"12V"`
	APIVersion       int    `json:"APIVersion"`
}

type InfoResponse struct {
	Response InfoResp `json:"INFO"`
}

type InfoResp struct {
	DeviceID          string `json:"DeviceID"`
	SubSystemDeviceID string `json:"SubsystemDeviceID"`
	SubSystemVendorID string `json:"SubsystemVendorID"`
	ASICSerialNumber  string `json:"ASICSerialNumber"`
	BoardSerialNumber string `json:"BoardSerialNumber"`
	UUID              string `json:"UUID"`
	SRAMSize          int    `json:"SRAMSize"`
	HBMSize           int    `json:"HBMSize"`
}

type StatusResponse struct {
	Response StatusResp `json:"STATUS"`
}
type StatusResp struct {
	BootStage                    string `json:"BootStage"`
	EmergencyPowerReduction      string `json:"EmergencyPowerReduction"`
	ClockThrottling              string `json:"ClockThrottling"`
	PowerState                   string `json:"PowerState"`
	ChipStatus                   string `json:"ChipStatus"`
	DeviceActivity               string `json:"DeviceActivity"`
	DevicePowerReduction         string `json:"DevicePowerReduction"`
	LastClockThrottlingDuration  int    `json:"LastClockThrottlingDuration"`
	TotalClockThrottlingDuration int    `json:"TotalClockThrottlingDuration"`
	GlobalTimeFromReset          int    `json:"GlobalTimeFromReset"`
	DeviceActivityCounter        int    `json:"DeviceActivityCounter"`
	LastPowerReductionDuration   int    `json:"LastPowerReductionDuration"`
}

type CTemperatureResponse struct {
	Response CTemperatureResp `json:"COMPOSITE_TEMPERATURE"`
}

type CTemperatureResp struct {
	Temperature int `json:"Temperature"`
}

type TemperatureResponse struct {
	Response TemperatureResp `json:"TEMPERATURE"`
}

type TemperatureResp struct {
	CurrentBoardTemperature               int `json:"CurrentBoardTemperature"`
	CurrentVRMTemperature                 int `json:"CurrentVRMTemperature"`
	CurrentDRAMTemperature                int `json:"CurrentDRAMTemperature"`
	CurrentOndieTemperature               int `json:"CurrentOndieTemperature"`
	HistoricalBoardTemperature            int `json:"HistoricalBoardTemperature"`
	HistoricalVRMTemperature              int `json:"HistoricalVRMTemperature"`
	HistoricalDRAMTemperature             int `json:"HistoricalDRAMTemperature"`
	HistoricalOndieTemperature            int `json:"HistoricalOndieTemperature"`
	MaxTemperatureRiseTime                int `json:"MaxTemperatureRiseTime"`
	MaxSOCTemperatureErrorThreshold       int `json:"MaxSOCTemperatureErrorThreshold"`
	MaxSOCTemperatureWarningThreshold     int `json:"MaxSOCTemperatureWarningThreshold"`
	MaxHBMTemperatureThreshold            int `json:"MaxHBMTemperatureThreshold"`
	CurrentSOCTemperatureErrorThreshold   int `json:"CurrentSOCTemperatureErrorThreshold"`
	CurrentSOCTemperatureWarningThreshold int `json:"CurrentSOCTemperatureWarningThreshold"`
	CurrentHBMTemperatureThreshold        int `json:"CurrentHBMTemperatureThreshold"`
}

type FrequencyResponse struct {
	Response FrequencyResp `json:"FREQUENCY"`
}

type FrequencyResp struct {
	HBMFrequency          int `json:"HBMFrequency"`
	MaxTPCFrequency       int `json:"MaxTPCFrequency"`
	MaxMMEFrequency       int `json:"MaxMMEFrequency"`
	MaxDMAFrequency       int `json:"MaxDMAFrequency"`
	MaxMediaFrequency     int `json:"MaxMediaFrequency"`
	MaxPCIeFrequency      int `json:"MaxPCIeFrequency"`
	MaxARMFrequency       int `json:"MaxARMFrequency"`
	MaxNICFrequency       int `json:"MaxNICFrequency"`
	MaxNoCFrequency       int `json:"MaxNoCFrequency"`
	CurrentTPCFrequency   int `json:"CurrentTPCFrequency"`
	CurrentMMEFrequency   int `json:"CurrentMMEFrequency"`
	CurrentDMAFrequency   int `json:"CurrentDMAFrequency"`
	CurrentMediaFrequency int `json:"CurrentMediaFrequency"`
	CurrentPCIeFrequency  int `json:"CurrentPCIeFrequency"`
	CurrentARMFrequency   int `json:"CurrentARMFrequency"`
	CurrentNICFrequency   int `json:"CurrentNICFrequency"`
	CurrentNoCFrequency   int `json:"CurrentNoCFrequency"`
	CurrentSRAMFrequency  int `json:"CurrentSRAMFrequency"`
	CurrentMSSFrequency   int `json:"CurrentMSSFrequency"`
}

type PowerResponse struct {
	Response PowerResp `json:"POWER"`
}

type PowerResp struct {
	CurrentPowerConsumption int `json:"CurrentPowerConsumption"`
	PeakPowerConsumption    int `json:"PeakPowerConsumption"`
}

type HbmResponse struct {
	Response HbmResp `json:"HBM"`
}

type HbmResp struct {
	RepairedLanes   []RepairedLane `json:"RepairedLane"`
	ReplacedRows    []ReplacedRow  `json:"ReplacedRow"`
	HBMRepairStatus []RepairStatus `json:"HBMRepairStatus"`
	EccErrors       int            `json:"ECCErrors"`
}

type RepairStatus struct {
	MBISTRepair string `json:"MBISTRepair"`
	HBMIndex    int    `json:"HBMIndex"`
	GlobalECC   int    `json:"GlobalECC"`
}
type ReplacedRow struct {
	Cause      string `json:"Cause"`
	HBMIndex   int    `json:"HBMIndex"`
	PCIndex    int    `json:"PCIndex"`
	StackID    int    `json:"StackID"`
	BankIndex  int    `json:"BankIndex"`
	RowAddress int    `json:"RowAddress"`
}
type RepairedLane struct {
	HBMIndex  int `json:"HBMIndex"`
	MCChannel int `json:"MCChannel"`
}

type SecurityResponse struct {
	Response SecurityResp `json:"SECURITY"`
}

type SecurityResp struct {
	Key0Revocation            string `json:"Key0Revocation"`
	Key1Revocation            string `json:"Key1Revocation"`
	Key2Revocation            string `json:"Key2Revocation"`
	Key3Revocation            string `json:"Key3Revocation"`
	Key4Revocation            string `json:"Key4Revocation"`
	MinimalSVNIndex           string `json:"MinimalSVNIndex"`
	FWImageSource             string `json:"FWImageSource"`
	TpmPcrPpboot              string `json:"TpmPcrPpboot"`
	TpmPcrPreboot             string `json:"TpmPcrPreboot"`
	TpmPcrUboot               string `json:"TpmPcrUboot"`
	TpmPcrLinux               string `json:"TpmPcrLinux"`
	CPLDVersion               string `json:"CPLDVersion"`
	CurrentPublicKeyHashIndex int    `json:"CurrentPublicKeyHashIndex"`
	CurrentSVNVersion         int    `json:"CurrentSVNVersion"`
	CPLDVersionTimestamp      int    `json:"CPLDVersionTimestamp"`
}

type SensorCurrentResponse struct {
	Response SensorCurrentResp `json:"SENSORS_CURRENT"`
}

type SensorCurrentResp struct {
	VIn54          int `json:"54VIn"`
	P1VIn12        int `json:"12P1VIn"`
	Stage154VIn    int `json:"Stage154VIn"`
	Stage113P5VOut int `json:"Stage113P5VOut"`
	Stage213P5VIn  int `json:"Stage213P5VIn"`
	Stage2CoreOut  int `json:"Stage2CoreOut"`
	Stage2HBMOut   int `json:"Stage2HBMOut"`
}

type SensorVoltageMonitorResponse struct {
	Response SensorVoltageMonitorResp `json:"SENSORS_VOLTAGE_MONITOR"`
}

type SensorVoltageMonitorResp struct {
	SwVm       int `json:"SwVm"`
	SwCpeEu0   int `json:"SwCpeEu_0"`
	SwCpeEu1   int `json:"SwCpeEu_1"`
	SwCpeEu2   int `json:"SwCpeEu_2"`
	SwcpeEu3   int `json:"SwcpeEu_3"`
	SwCpeHbm   int `json:"SwCpeHbm"`
	SwTpcSb    int `json:"SwTpcSb"`
	SwTft      int `json:"SwTft"`
	SwCnt      int `json:"SwCnt"`
	SwMcHbm00  int `json:"SwMcHbm0_0"`
	SwMcHbm01  int `json:"SwMcHbm0_1"`
	SwHconHbm0 int `json:"SwHconHbm0"`
	SwPpw0     int `json:"SwPpw_0"`
	SwPpw1     int `json:"SwPpw_1"`
	SwL2cMacro int `json:"SwL2cMacro"`
	SwVcdMacro int `json:"SwVcdMacro"`
	SeVm       int `json:"SeVm"`
	SePe00     int `json:"SePe0_0"`
	SePe01     int `json:"SePe0_1"`
	SeEuCore   int `json:"SeEuCore"`
	SeBpyramid int `json:"SeBpyramid"`
	SeRx       int `json:"SeRx"`
	SeMx       int `json:"SeMx"`
	SeHconHbm1 int `json:"SeHconHbm1"`
	SeMcHbm1   int `json:"SeMcHbm1"`
	SeGasket   int `json:"SeGasket"`
	SeMmeCtrl  int `json:"SeMmeCtrl"`
	SeMmeQman  int `json:"SeMmeQman"`
	SeSbte     int `json:"SeSbte"`
	SeRtrDn    int `json:"SeRtrDn"`
	SeRtrUp    int `json:"SeRtrUp"`
	SeSram     int `json:"SeSram"`
	NwVm       int `json:"NwVm"`
	NwPe00     int `json:"NwPe0_0"`
	NwPe01     int `json:"NwPe0_1"`
	NwPe02     int `json:"NwPe0_2"`
	NwPe03     int `json:"NwPe0_3"`
	NwEuCore   int `json:"NwEuCore"`
	NwBpyramid int `json:"NwBpyramid"`
	NwMcHbm40  int `json:"NwMcHbm4_0"`
	NwMcHbm41  int `json:"NwMcHbm4_1"`
	NwHconHbm4 int `json:"NwHconHbm4"`
	NwAcc      int `json:"NwAcc"`
	NwWap      int `json:"NwWap"`
	NwTif      int `json:"NwTif"`
	NwRtrDn    int `json:"NwRtrDn"`
	NwRtrUp    int `json:"NwRtrUp"`
	NwSram     int `json:"NwSram"`
	NeVM       int `json:"NeVM"`
	NeTx       int `json:"NeTx"`
	NeMx0      int `json:"NeMx_0"`
	NeMx1      int `json:"NeMx_1"`
	NeHconHbm5 int `json:"NeHconHbm5"`
	NeMcHbm5   int `json:"NeMcHbm5"`
	NeSob      int `json:"NeSob"`
	NeQnt      int `json:"NeQnt"`
	NeCnt0     int `json:"NeCnt_0"`
	NeCnt1     int `json:"NeCnt_1"`
	NeCpeEu    int `json:"NeCpeEu"`
	NeCntHbm   int `json:"NeCntHbm"`
	NeSbte0    int `json:"NeSbte_0"`
	NeSbte1    int `json:"NeSbte_1"`
	NeMif      int `json:"NeMif"`
	NeRtrDn    int `json:"NeRtrDn"`
}

type SensorVoltageResponse struct {
	Response SensorVoltageResp `json:"SENSORS_VOLTAGE"`
}

type SensorVoltageResp struct {
	VADC54          int `json:"54VADC"`
	Vrm1In          int `json:"Vrm1In"`
	Vrm1Out         int `json:"Vrm1Out"`
	Vrm2In          int `json:"Vrm2In"`
	Vrm2VddOut      int `json:"Vrm2VddOut"`
	Vrm2HbmOut      int `json:"Vrm2HbmOut"`
	VmonPcieVph1P8V int `json:"VmonPcieVph1P8V"`
	Vmon1P8HbmVaa   int `json:"Vmon1P8HbmVaa"`
	Vmon2P5         int `json:"Vmon2P5"`
	Vmon48VHimon    int `json:"Vmon48VHimon"`
	VmonP5V         int `json:"VmonP5V"`
	Vmon12V1        int `json:"Vmon12V1"`
	VmonHbm         int `json:"VmonHbm"`
	VmonCore        int `json:"VmonCore"`
	CpldHimon1P8NIC int `json:"CpldHimon1P8NIC"`
}

type SensorTemperatureResponse struct {
	Response SensorTemperatureResp `json:"RESPONSE"`
}

type SensorTemperatureResp struct {
	OnDie0    int `json:"OnDie0"`
	OnDie1    int `json:"OnDie1"`
	OnDie2    int `json:"OnDie2"`
	OnDie3    int `json:"OnDie3"`
	HBM0      int `json:"HBM0"`
	HBM1      int `json:"HBM1"`
	HBM2      int `json:"HBM2"`
	HBM3      int `json:"HBM3"`
	HBM4      int `json:"HBM4"`
	HBM5      int `json:"HBM5"`
	CPLDLocal int `json:"CPLDLocal"`
	CPLD0     int `json:"CPLD0"`
	CPLD1     int `json:"CPLD1"`
	CPLD2     int `json:"CPLD2"`
	CPLD3     int `json:"CPLD3"`
	OnBoard0  int `json:"OnBoard0"`
	OnBoard1  int `json:"OnBoard1"`
	OnBoard2  int `json:"OnBoard2"`
	OnBoard3  int `json:"OnBoard3"`
	CPLDTemp  int `json:"CPLDTemp"`
	PSUStage1 int `json:"PSUStage1"`
	PSUStage2 int `json:"PSUStage2"`
}

type PcieInfoResponse struct {
	Response PcieInfoResp `json:"PCIE_INFO"`
}

type PcieInfoResp struct {
	MaxPCIeLinkSpeed                     string `json:"MaxPCIeLinkSpeed"`
	CurrentPCIeLinkSpeed                 string `json:"CurrentPCIeLinkSpeed"`
	PCIeSubsystemVendorID                string `json:"PCIeSubsystemVendorID"`
	PCIeBusAndDevice                     string `json:"PCIeBusAndDevice"`
	MaxPCIeLinkWidth                     int    `json:"MaxPCIeLinkWidth"`
	CurrentPCIeLinkWidth                 int    `json:"CurrentPCIeLinkWidth"`
	PCIeDeviceID                         int    `json:"PCIeDeviceID"`
	PCIeSubsystemID                      int    `json:"PCIeSubsystemID"`
	CorrectedInternalErrorStatus         int    `json:"CorrectedInternalErrorStatus"`
	ReplayBufferNumRolloverError         int    `json:"ReplayBufferNumRolloverError "`
	ReplayTimerTimeoutError              int    `json:"ReplayTimerTimeoutError "`
	BadTLPCounter                        int    `json:"BadTLPCounter"`
	BadDLLPCounter                       int    `json:"BadDLLPCounter"`
	ReceiverErrorCounter                 int    `json:"ReceiverErrorCounter"`
	LCRCErrorCounter                     int    `json:"LCRCErrorCounter"`
	ECRCErrorCounter                     int    `json:"ECRCErrorCounter"`
	CompletionTimeoutIndication          int    `json:"CompletionTimeoutIndication"`
	UncorrectableInternalErrorIndication int    `json:"UncorrectableInternalErrorIndication"`
	ReceiverOverflowIndication           int    `json:"ReceiverOverflowIndication"`
	FlowControlProtocolErrorIndication   int    `json:"FlowControlProtocolErrorIndication"`
	SurpriseLinkDownIndication           int    `json:"SurpriseLinkDown Indication"`
	MalfunctionTLPErrorIndication        int    `json:"MalfunctionTLPErrorIndication"`
	DLLPProtocolErrorIndication          int    `json:"DLLPProtocolErrorIndication"`
	RxNakDLLPCounter                     int    `json:"RxNakDLLPCounter"`
	TxNakDLLPCounter                     int    `json:"TxNakDLLPCounter"`
	RetryTLPCounter                      int    `json:"RetryTLPCounter"`
	PWRBRKIndication                     int    `json:"PWRBRKIndication"`
	PCIeRxMemoryWriteCounter             int    `json:"PCIeRxMemoryWriteCounter"`
	PCIeRxMemoryReadCounter              int    `json:"PCIeRxMemoryReadCounter"`
	PCIeTxMemoryWriteCounter             int    `json:"PCIeTxMemoryWriteCounter"`
	PCIeTxMemoryReadCounter              int    `json:"PCIeTxMemoryReadCounter"`
	AERCapabilityControlOffset           int    `json:"AERCapabilityControlOffset"`
	AERErrorLog                          int    `json:"AERErrorLog"`
	PCIeFWVersion                        int    `json:"PCIeFWVersion"`
}

type EthernetStatusResponse struct {
	Response EthernetStatusResp `json:"ETHERNET_STATUS"`
}

type EthernetStatusResp struct {
	PortMap                string `json:"PortMap"`
	ExternalLinkStatus     string `json:"ExternalLinkStatus"`
	LinkStatus             string `json:"LinkStatus"`
	PHYStatus              string `json:"PHYStatus"`
	PortNum                int
	ToggleCount            int `json:"ToggleCount"`
	BERCorrectable         int `json:"BERCorrectable"`
	BERUncorrectable       int `json:"BERUncorrectable"`
	Nack                   int `json:"Nack"`
	CRC                    int `json:"CRC"`
	RetransmissionTimeout  int `json:"RetransmissionTimeout"`
	LinkRetrainingDueToBER int `json:"LinkRetrainingDueToBER"`
	MACRemote              int `json:"MACRemote"`
	Retransmission         int `json:"Retransmission"`
	Retraining             int `json:"Retraining"`
	SERPreFEC              int `json:"SERPreFEC"`
	SERPostFEC             int `json:"SERPostFEC"`
	Latency                int `json:"Latency"`
	Throughput             int `json:"Throughput"`
}

type EthernetInfoResponse struct {
	Response EthernetInfoResp `json:"ETHERNET_INFO"`
}

type EthernetInfoResp struct {
	SerdesAvailability string   `json:"SerdesAvailability"`
	ANLTStatus         []string `json:"ANLTStatus"`
	PortMap            []string `json:"PortMap"`
	ToggleCount        []int    `json:"ToggleCount"`
	ExternalLinkStatus []string `json:"ExternalLinkStatus"`
	LinkStatus         []string `json:"LinkStatus"`
	PHYStatus          []string `json:"PHYStatus"`
	PortMaxSpeed       int      `json:"PortMaxSpeed"`
	NumberOfLanes      int      `json:"NumberOfLanes"`
	NumberOfLinks      int      `json:"NumberOfLinks"`
	LinkSpeed          int      `json:"LinkSpeed"`
}

type LaneInfoResponse struct {
	// TODO: verify this json name
	Response LaneInfoResp `json:"LANE_INFO"`
}

type LaneInfoResp struct {
	LaneNum             int
	EBUFOverflow        int `json:"EBUFOverflow"`
	EBUFUnderRun        int `json:"EBUFUnderRun"`
	DecodeError         int `json:"DecodeError"`
	RunningDisplayError int `json:"RunningDisplayError"`
	SKPOSParityError    int `json:"SKPOSParityError"`
	SYNCHeaderError     int `json:"SYNCHeaderError"`
	DeskewError         int `json:"DeskewError"`
}

type BmcStateResponse struct {
	PowerState string `json:"PowerState"`
}

type SessionInformation struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Status      struct {
		State  string `json:"State"`
		Health string `json:"Health"`
	} `json:"Status"`
	ServiceEnabled bool `json:"ServiceEnabled"`
	SessionTimeout int  `json:"SessionTimeout"`
}

type AlertsResponse struct {
	Response AlertsResp `json:"ALERT"`
}

type AlertsResp struct {
	SRAMDoubleECC   *SRAMDoubleECCAlert   `json:"SRAMDoubleECC,omitempty"`
	SRAMCorrectable *SRAMCorrectableAlert `json:"SRAMCorrectable,omitempty"`
	PCIe            *PCIeAlert            `json:"PCIe,omitempty"`
	HBM             *HBMAlert             `json:"HBM,omitempty"`
	WD              *WDAlert              `json:"WD,omitempty"`
	PowerThreshold  *PowerThresholdAlert  `json:"PowerThreshold,omitempty"`
	Temperature     *TemperatureAlert     `json:"Temperature,omitempty"`
	Security        *SecurityAlert        `json:"Security,omitempty"`
	RMA             *RMAAlert             `json:"RMA,omitempty"`
	VoltageMonitor  *VoltageMonitorAlert  `json:"VoltageMonitor,omitempty"`
	NIC             *NICAlert             `json:"NIC,omitempty"`
	Checksum        *ChecksumAlert        `json:"Checksum,omitempty"`
	PLL             *PLLAlert             `json:"PLL,omitempty"`
	FIT             *FITAlert             `json:"FIT,omitempty"`
}

type SRAMDoubleECCAlert []SRAMDoubleECC

type SRAMDoubleECC struct {
	BlockIndex   int   `json:"BlockIndex"`
	IndexInBlock int   `json:"IndexInBlock"`
	Address      int   `json:"Address"`
	Syndrome     int   `json:"Syndrome"`
	Timestamp    int64 `json:"Timestamp"`
}

func (s SRAMDoubleECC) Description() string {
	return fmt.Sprintf("block index: %d, index in block: %d, address: %d, syndrome: %d",
		s.BlockIndex, s.IndexInBlock, s.Address, s.Syndrome)
}

type SingleEcc struct {
	BlockIndex int   `json:"BlockIndex"`
	Index      int   `json:"Index"`
	Address    int   `json:"Address"`
	Syndrome   int   `json:"Syndrome"`
	Timestamp  int64 `json:"Timestamp"`
}

func (s SingleEcc) Description() string {
	return fmt.Sprintf("block index: %d,index: %d, address: %d, syndrome: %d", s.BlockIndex, s.Index, s.Address, s.Syndrome)
}

type MultipleECC struct {
	Cause      string `json:"Cause"`
	BlockIndex int    `json:"BlockIndex"`
	Index      int    `json:"Index"`
	Address    int    `json:"Address"`
	Timestamp  int64  `json:"Timestamp"`
}

func (m MultipleECC) Description() string {
	return fmt.Sprintf("block index: %d,index: %d, address: %d, cause: %s", m.BlockIndex, m.Index, m.Address, m.Cause)
}

type PCIeAlert struct {
	Fatal       string `json:"Fatal"`
	Correctable string `json:"Correctable"`
	NonFatal    string `json:"NonFatal"`
	Timestamp   int64  `json:"Timestamp"`
}

type Derr struct {
	HBMIndex  int   `json:"HBMIndex"`
	Address   int   `json:"Address"`
	Syndrome  int   `json:"Syndrome"`
	PCIndex   int   `json:"PCIndex"`
	Timestamp int64 `json:"Timestamp"`
}

func (d Derr) Description() string {
	return fmt.Sprintf("hbm index: %d, address: %d,syndrome: %d, pci index: %d",
		d.HBMIndex, d.Address, d.Syndrome, d.PCIndex)
}

type MultiSERR struct {
	HBMIndex  int   `json:"HBMIndex"`
	Address   int   `json:"Address"`
	Syndrome  int   `json:"Syndrome"`
	PCIndex   int   `json:"PCIndex"`
	Timestamp int64 `json:"Timestamp"`
}

func (m MultiSERR) Description() string {
	return fmt.Sprintf("hbm index: %d, address: %d,syndrome: %d, pci index: %d",
		m.HBMIndex, m.Address, m.Syndrome, m.PCIndex)
}

type Parity struct {
	Parity    string `json:"Parity"`
	HBMIndex  int    `json:"HBMIndex"`
	PCIndex   int    `json:"PCIndex"`
	Timestamp int64  `json:"Timestamp"`
}

func (p Parity) Description() string {
	return fmt.Sprintf("hbm index: %d, pci index: %d, parity: %s",
		p.HBMIndex, p.PCIndex, p.Parity)
}

type SameAddressMultiSERR struct {
	HBMIndex  int   `json:"HBMIndex"`
	PCIndex   int   `json:"PCIndex"`
	SIDIndex  int   `json:"SIDIndex"`
	BankIndex int   `json:"BankIndex"`
	Address   int   `json:"Address"`
	Count     int   `json:"Count"`
	Timestamp int64 `json:"Timestamp"`
}

func (s SameAddressMultiSERR) Description() string {
	return fmt.Sprintf("hbm index: %d, pci index: %d, sid index: %d, bank index: %d, address: %d, count: %d",
		s.HBMIndex, s.PCIndex, s.SIDIndex, s.BankIndex, s.Address, s.Count)
}

type WDAlert struct {
	Heartbeat string `json:"Heartbeat"`
	TDR       string `json:"TDR"`
	HW        string `json:"HW"`
	Timestamp int64  `json:"Timestamp"`
}

type PowerThresholdAlert struct {
	Crossed   string `json:"Crossed"`
	Timestamp int64  `json:"Timestamp"`
}

type ThermalTemperatureCrossing struct {
	Crossed       string `json:"Crossed"`
	MaxHistorical int    `json:"MaxHistorical"`
}

func (t ThermalTemperatureCrossing) Description() string {
	return fmt.Sprintf("crossed: %s, max historical: %d", t.Crossed, t.MaxHistorical)
}

type SPIFailure struct {
	RevokedPublicKey string `json:"RevokedPublicKey"`
	SVN              string `json:"SVN"`
	Signature        string `json:"Signature"`
}

type AgentFailure struct {
	RevokedPublicKey string `json:"RevokedPublicKey"`
	SVN              string `json:"SVN"`
	Signature        string `json:"Signature"`
	UBoot            string `json:"uBoot"`
	Linux            string `json:"Linux"`
	Zephyr           string `json:"Zephyr"`
}
type FatalSRAMDoubleECC struct {
	Cause     string `json:"Cause"`
	Index     int    `json:"Index"`
	Address   int    `json:"Address"`
	Count     int    `json:"Count"`
	Timestamp int64  `json:"Timestamp"`
}

func (f FatalSRAMDoubleECC) Description() string {
	return fmt.Sprintf("index: %d, address: %d, cause: %s, count: %d",
		f.Index, f.Address, f.Cause, f.Count)
}

type HBMECC struct {
	Cause     string `json:"Cause"`
	HBMIndex  int    `json:"HBMIndex"`
	PCIndex   int    `json:"PCIndex"`
	SIDIndex  int    `json:"SIDIndex"`
	BankIndex int    `json:"BankIndex"`
	Address   int    `json:"Address"`
	Count     int    `json:"Count"`
	Timestamp int64  `json:"Timestamp"`
}

func (h HBMECC) Description() string {
	return fmt.Sprintf("hbm index: %d, pci index: %d, sid index: %d, bank index: %d, cause: %s, address: %d, count: %d",
		h.HBMIndex, h.PCIndex, h.SIDIndex, h.BankIndex, h.Cause, h.Address, h.Count)
}

type HBMParity struct {
	Cause     string `json:"Cause"`
	HBMIndex  int    `json:"HBMIndex"`
	PCIndex   int    `json:"PCIndex"`
	Count     int    `json:"Count"`
	Timestamp int64  `json:"Timestamp"`
}

func (h HBMParity) Description() string {
	return fmt.Sprintf("hbm index: %d, pci index: %d, cause: %s, count: %d",
		h.HBMIndex, h.PCIndex, h.Cause, h.Count)
}

type NonRepairable struct {
	RingNumber int `json:"RingNumber"`
	SubServer  int `json:"SubServer"`
}

func (n NonRepairable) Description() string {
	return fmt.Sprintf("ring number: %d, sub server %d",
		n.RingNumber, n.SubServer)
}

type Repairable struct {
	RingNumber int `json:"RingNumber"`
	SubServer  int `json:"SubServer"`
}

func (r Repairable) Description() string {
	return fmt.Sprintf("ring number: %d, sub server %d",
		r.RingNumber, r.SubServer)
}

type HBMMBISTUnrepairable struct {
	InFieldFailure string `json:"InFieldFailure"`
	HBMIndex       int    `json:"HBMIndex"`
}

func (h HBMMBISTUnrepairable) Description() string {
	return fmt.Sprintf("in field failure: %s, hbm index: %d",
		h.InFieldFailure, h.HBMIndex)
}

type VoltageMonitorAlert struct {
	Threshold1Crossing string `json:"Threshold1Crossing"`
	Threshold2Crossing string `json:"Threshold2Crossing"`
	Timestamp          int64  `json:"Timestamp"`
}
type NICAlert struct {
	HighBER                       string `json:"HighBER"`
	MultipleNAC                   string `json:"MultipleNAC"`
	MultipleCRCError              string `json:"MultipleCRCError"`
	MultipleLinkReinitializations string `json:"MultipleLinkReinitializations"`
	MultipleLinkRetransmissions   string `json:"MultipleLinkRetransmissions"`
	Timestamp                     int64  `json:"Timestamp"`
}

type RMATables struct {
	HBMParity            string `json:"HBMParity"`
	HBMECC               string `json:"HBMECC"`
	SRAMECC              string `json:"SRAMECC"`
	RowReplacement       string `json:"RowReplacement"`
	RMACauseFailure      string `json:"RMACauseFailure"`
	SpareRowAvailability string `json:"SpareRowAvailability"`
}

type SPIChecksum struct {
	PpBootPrimary    string `json:"ppBootPrimary"`
	PpBootSecondary  string `json:"ppBootSecondary"`
	PreBootPrimary   string `json:"preBootPrimary"`
	PreBootSecondary string `json:"preBootSecondary"`
}

type ChecksumAlert struct {
	RMATables   *RMATables   `json:"RMATables,omitempty"`
	SPIChecksum *SPIChecksum `json:"SPIChecksum,omitempty"`
	CPLDStatus  string       `json:"CPLDStatus"`
	EEPROM      string       `json:"EEPROM"`
	Timestamp   int64        `json:"Timestamp"`
}

type PLLAlert map[string]string

type FITAlert struct {
	InvalidImageFormat              string `json:"InvalidImageFormat"`
	FWUpgradeTargetAddressViolation string `json:"FWUpgradeTargetAddressViolation"`
	FITNotRunnable                  string `json:"FITNotRunnable"`
	Timestamp                       int64  `json:"Timestamp"`
}

type SRAMCorrectableAlert struct {
	SingleECC   []SingleEcc   `json:"SingleECC"`
	MultipleECC []MultipleECC `json:"MultipleECC"`
}

type HBMAlert struct {
	DERR                 []Derr                 `json:"DERR"`
	MultiSERR            []MultiSERR            `json:"MultiSERR"`
	CATTRIP              string                 `json:"CATTRIP"`
	Parity               []Parity               `json:"Parity"`
	SameAddressMultiSERR []SameAddressMultiSERR `json:"SameAddressMultiSERR"`
	Timestamp            int64                  `json:"Timestamp"`
}

type TemperatureAlert struct {
	ThermalTemperatureCrossing *ThermalTemperatureCrossing `json:"ThermalTemperatureCrossing,omitempty"`
	RiseTimeViolation          string                      `json:"RiseTimeViolation"`
	Timestamp                  int64                       `json:"Timestamp"`
}

type SecurityAlert struct {
	SPIFailure            *SPIFailure   `json:"SPIFailure,omitempty"`
	AgentFailure          *AgentFailure `json:"AgentFailure,omitempty"`
	PermissionAccessError string        `json:"PermissionAccessError"`
	Timestamp             int64         `json:"Timestamp"`
}

type RMAAlert struct {
	SRAMRepairFailure    *SRAMRepairFailure    `json:"SRAMRepairFailure,omitempty"`
	HBMMBISTUnrepairable *HBMMBISTUnrepairable `json:"HBMMBISTUnrepairable,omitempty"`
	HBMRowReplacement    string                `json:"HBMRowReplacement"`
	HBMInitialization    string                `json:"HBMInitialization"`
	FatalSRAMDoubleECC   []FatalSRAMDoubleECC  `json:"FatalSRAMDoubleECC"`
	HBMECC               []HBMECC              `json:"HBMECC"`
	HBMParity            []HBMParity           `json:"HBMParity"`
	Timestamp            int64                 `json:"Timestamp"`
}

type SRAMRepairFailure struct {
	NonRepairable []NonRepairable `json:"NonRepairable"`
	Repairable    []Repairable    `json:"Repairable"`
}
