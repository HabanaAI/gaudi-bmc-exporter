package bmcmonitoring

import (
	"encoding/json"
	"os"
)

func InitAlertsDescription(path string) (AlertsData, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return AlertsData{}, err
	}
	var alertsDescription AlertsData
	err = json.Unmarshal(data, &alertsDescription)
	if err != nil {
		return AlertsData{}, err
	}

	return alertsDescription, nil
}
