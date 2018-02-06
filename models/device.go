package models

type Device struct {
	Id   string `json:"deviceId"`
	Imei string `json:"imeiNumber"`
	Name string `json:"description"`
}
