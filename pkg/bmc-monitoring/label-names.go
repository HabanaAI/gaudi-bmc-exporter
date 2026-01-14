package bmcmonitoring

const (
	RMALabelHBMDeviceIndex = "HBM device Index"
	RMALabelHBMPCIndex     = "HBM PC Index"
	RMALabelHBMSIDIndex    = "HBM SID Index"
	RMALabelHBMBankIndex   = "HBM Bank Index"
	RMALabelErrorCause     = "error cause"
	RMALabelHBMRowAddress  = "HBM Row Address"
	RMALabelCount          = "count"
)

// Label Names
const (
	// HBM
	// replace Rows
	HBMMetricReplaceRowsLabelHBMIndex   = "hbm_index"
	HBMMetricReplaceRowsLabelPCIndex    = "pc_index"
	HBMMetricReplaceRowsLabelStackID    = "stack_id"
	HBMMetricReplaceRowsLabelBankIndex  = "bank_index"
	HBMMetricReplaceRowsLabelCause      = "cause"
	HBMMetricReplaceRowsLabelRowAddress = "row_address"

	// repaired Lanes
	HBMMetricRepairedLanesLabelHBMIndex  = "hbm_index"
	HBMMetricRepairedLanesLabelMCChannel = "mc_channel"

	// mbistRepair
	HBMMetricMbistRepairLabelState = "state"
	HBMMetricMbistRepairLabelIndex = "index"

	// global ECC
	HBMMetricGlobalECCLabelIndex = "index"
)
