package types

type EdgeData struct {
	DeviceId   string `json:"deviceId"`
	DeviceType string `json:"deviceType"`
	Data       struct {
		Confiability float64 `json:"confiability"`
		Flow         float64 `json:"flow"`
	} `json:"data"`
	TS string `json:"ts"`
}
