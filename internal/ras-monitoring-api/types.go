package rasmonitoringapi

import (
	rasmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"net/http"
)

type ExpectedType string

// Expected type the opcode should return.
// This way we know how to format the data.
const (
	Int ExpectedType = "int"

	// we need to reverse the string before reading it.
	ReverseString ExpectedType = "reverse string"

	String ExpectedType = "string"

	AsciiString ExpectedType = "ascii string"

	// bytes array that we need to decode to bin.
	BinArray ExpectedType = "bin array"

	Bin ExpectedType = "bin"
)

// OpcodeData represents the way we get correct data from the api.
type OpcodeData struct {
	Data []byte `json:"data_array"`
}

type Opcode struct {

	// For opcodes with out any clear rule we will create a custom decoder.
	Decoder      func(string, int, string, string, string) ([]rasmonitoring.Metric, error)
	ExpectedType ExpectedType

	OpcodeNumber int
	Offset       int
	Length       int
}

type authDetails struct {
	Cookie         *http.Cookie
	RemoteAddr     string `json:"remote_addr"`
	ServerName     string `json:"server_name"`
	ServerAddr     string `json:"server_addr"`
	CSRFToken      string `json:"CSRFToken"`
	Ok             int    `json:"ok"`
	Privilege      int    `json:"privilege"`
	UserID         int    `json:"user_id"`
	Extendedpriv   int    `json:"extendedpriv"`
	RacsessionID   int    `json:"racsession_id"`
	HTTPSEnabled   int    `json:"HTTPSEnabled"`
	Channel        int    `json:"channel"`
	PasswordStatus int    `json:"passwordStatus"`
}

type result struct {
	metric rasmonitoring.Metric
	err    error
}

type SessionData []struct {
	ClientIP      string `json:"client_ip"`
	UserName      string `json:"user_name"`
	ID            int    `json:"id"`
	SessionID     int    `json:"session_id"`
	SessionType   int    `json:"session_type"`
	UserID        int    `json:"user_id"`
	UserPrivilege int    `json:"user_privilege"`
}
