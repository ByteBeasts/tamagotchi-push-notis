package common

type Payload struct {
	AppID   string   `json:"app_id"`
	Addresses []string `json:"wallet_addresses"`
	Title   string   `json:"title"`
	Message string   `json:"message"`
	MiniAppPath string   `json:"mini_app_path"`
}