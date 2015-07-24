package main
import (
	"time"
)

/*
"imei": "353188020901357",
"name": "Samsung S2",
"connected_hostname": "host100",
"vnc_port": 5901,
"ssh_port": 6000,
"adb_name": "G1002de5a082",
"start_time": "2015-07-20T11:12:59.625592627-07:00",
"stop_time": "0001-01-01T00:00:00Z"
*/

type Device struct {
	IMEI string `json: "imei"`
	Name string	`json: "name"`
	ConnectedHostname string `json: "connected_hostname"`
	VNCPort int `json: "vnc_port"`
	SSHPort int `json: "ssh_port"`
	ADBName string `json: "adb_name"`
	StartTime time.Time `json:"start_time"`
	StopTime time.Time `json:"stop_time"`
}

type Devices []Device


type Inventroies []Device
