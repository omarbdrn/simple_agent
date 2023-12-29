package models

import (
	"github.com/omarbdrn/fingerprintx/pkg/plugins"
)

type HostServices struct {
	IP    string    `json:"ip"`
	Ports []Service `json:"ports"`
}

type Service struct {
	Port     int
	Protocol string
	Version  string
	Metadata plugins.Metadata
}
