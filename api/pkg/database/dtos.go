package database

import "encoding/json"

type HostServicesDTO struct {
	IP    string       `json:"ip"`
	Ports []ServiceDTO `json:"ports"`
}

type ServiceDTO struct {
	Port     int    `json:"Port"`
	Protocol string `json:"Protocol"`
	Version  string `json:"Version"`
	Metadata string
	Raw      json.RawMessage `json:"Metadata"`
}
