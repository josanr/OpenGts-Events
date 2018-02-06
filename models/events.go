package models

type Events struct {
	Account      string
	Account_desc string
	TimeZone     string
	DeviceList   DeviceList
}

type DeviceList struct {
	Device      string
	Device_desc string
	EventData   []EventData
}

type EventData struct {
	//Device          string
	//Timestamp       string
	//Timestamp_date  string
	//Timestamp_time  string
	//StatusCode      string
	//StatusCode_hex  string
	//StatusCode_desc string
	//GPSPoint        string
	GPSPoint_lat string
	GPSPoint_lon string
	//Speed_kph       string
	//Speed           string
	//Speed_units     string
	//Index           string
}
