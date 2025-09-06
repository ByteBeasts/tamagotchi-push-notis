package common

type Payload struct {
	Message string   `json:"message"`
	AppID   string   `json:"app_id"`
	Addresses []string `json:"addresses"`
}