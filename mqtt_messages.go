package eadmin

import "time"

type Ping struct {
	ID       string  `json:"id"`
	Disc     string  `json:"disc"`
	IDType   int     `json:"idType"`
	Rssi1M   int     `json:"rssi@1m"`
	Rssi     int     `json:"rssi"`
	Raw      float64 `json:"raw"`
	Distance float64 `json:"distance"`
	Speed    float64 `json:"speed"`
	Mac      string  `json:"mac"`
	Interval int     `json:"interval"`
	Recieved time.Time
}
