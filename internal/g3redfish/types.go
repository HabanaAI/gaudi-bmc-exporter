package g3redfish

type V1SessionService struct {
	Name           string `json:"Name"`
	ID             string `json:"Id"`
	Description    string `json:"Description"`
	ServiceEnabled bool   `json:"ServiceEnabled"`
	SessionTimeout int    `json:"SessionTimeout"`
}

type V1Collection struct {
	Name    string               `json:"Name"`
	Members []V1CollectionMember `json:"Members"`
}

type V1CollectionMember struct {
	Id string `json:"@odata.id"`
}

type V1SensorReading struct {
	Name       string   `json:"Name"`
	Value      *float64 `json:"Reading"`
	Type       string   `json:"ReadingType"`
	Unit       string   `json:"ReadingUnits"`
	Thresholds struct {
		UpperCritical V1SensorThreshold
		UpperCaution  V1SensorThreshold
		LowerCaution  V1SensorThreshold
		LowerCritical V1SensorThreshold
	} `json:"Thresholds,omitempty"`
}

type V1SensorThreshold struct {
	Value *float32 `json:"Reading"`
}
