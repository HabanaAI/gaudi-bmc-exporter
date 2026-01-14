package g3redfish

import (
	"fmt"
)

type RedFishErr struct {
	Err GeneralErr `json:"error"`
}

type GeneralErr struct {
	Code         string            `json:"code"`
	Message      string            `json:"message"`
	ExtendedInfo []ExtendedErrInfo `json:"@Message.ExtendedInfo"`
}

type ExtendedErrInfo struct {
	MessageID         string   `json:"MessageId"`
	Severity          string   `json:"Severity"`
	Resolution        string   `json:"Resolution"`
	Message           string   `json:"Message"`
	MessageArgs       []string `json:"MessageArgs"`
	RelatedProperties []string `json:"RelatedProperties"`
}

func (r *RedFishErr) Error() string {
	errMsg := fmt.Sprintf("code: %s, message: %s.", r.Err.Code, r.Err.Message)

	for _, extended := range r.Err.ExtendedInfo {
		errMsg = fmt.Sprintf("%s MessageId: %s, Severity: %s, Resolution: %s,Message: %s",
			errMsg, extended.MessageID, extended.Severity, extended.Resolution, extended.Message)

		for _, msgArg := range extended.MessageArgs {
			errMsg = fmt.Sprintf("%s message arg: %s", errMsg, msgArg)
		}

		for _, relatedTopic := range extended.RelatedProperties {
			errMsg = fmt.Sprintf("%s related topic: %s", errMsg, relatedTopic)
		}
	}
	return errMsg
}
