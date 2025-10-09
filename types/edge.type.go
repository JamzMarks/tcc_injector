package types

type EdgeData struct {
	DeviceID string `json:"deviceId"`
	Location struct {
		To   int `json:"to"`
		From int `json:"from"`
	} `json:"location"`
	Data struct {
		Confiability float32 `json:"confiability"`
		Flow         float32 `json:"flow"`
	} `json:"data"`
}
