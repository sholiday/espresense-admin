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

var idTypes = map[int]string{
	-10: "ECHO_LOST",
	-5:  "MISC_APPLE",
	0:   "PE",
	1:   "RAND_MAC",
	10:  "AD",
	15:  "SD",
	20:  "MD",
	30:  "MISC",
	32:  "FINDMY",
	35:  "NAME",
	40:  "MSFT",
	50:  "UNIQUE",
	55:  "PUBLIC_MAC",
	105: "SONOS",
	107: "GARMIN",
	110: "MITHERM",
	115: "MIFIT",
	120: "EXPOSURE",
	121: "SMARTTAG",
	125: "ITAG",
	127: "ITRACK",
	128: "NUT",
	130: "TRACKR",
	135: "TILE",
	140: "MEATER",
	142: "TRACTIVE",
	145: "VANMOOF",
	150: "APPLE_NEARBY",
	155: "QUERY_MODEL",
	160: "QUERY_NAME",
	165: "RM_ASST",
	170: "EBEACON",
	175: "ABEACON",
	180: "IBEACON",
	200: "KNOWN_IRK",
	210: "KNOWN_MAC",
	250: "ALIAS",
}
