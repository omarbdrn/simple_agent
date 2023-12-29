package database

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type DatabaseReporter struct {
	DB *gorm.DB
}

func NewDatabaseReporter() *DatabaseReporter {
	return &DatabaseReporter{}
}

func (dpr *DatabaseReporter) GetHosts() []Host {
	dpr.DB = GetDB()
	var hosts []Host

	result := db.Model(&Host{}).Preload("Services").Find(&hosts)
	if result.Error != nil || result.RowsAffected == 0 {
		return []Host{}
	}

	return hosts

}

func (dpr *DatabaseReporter) AddOrUpdateService(host Host, service ServiceDTO) {
	var ServiceModel Service
	var result *gorm.DB

	metadata, err := json.Marshal(service.Raw)
	if err != nil {
		metadata = []byte{}
	}
	service.Metadata = string(metadata)

	result = dpr.DB.First(&ServiceModel, "port = ? AND host_id = ?", service.Port, host.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ServiceModel := &Service{Port: service.Port, Protocol: service.Protocol, Version: service.Version, Metadata: service.Metadata, HostID: host.ID}
			db.Create(ServiceModel)
		}
	} else {
		ServiceModel.Protocol = service.Protocol
		ServiceModel.Version = service.Version
		ServiceModel.Metadata = service.Metadata
	}
}

func (dpr *DatabaseReporter) AddHost(hostService *HostServicesDTO) {
	dpr.DB = GetDB()
	var host Host
	var result *gorm.DB

	result = db.First(&host, "ip = ?", hostService.IP)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			host := &Host{IP: hostService.IP}
			for _, service := range hostService.Ports {
				metadata, err := json.Marshal(service.Raw)
				if err != nil {
					metadata = []byte{}
				}
				service.Metadata = string(metadata)
				serviceModel := Service{Port: service.Port, Protocol: service.Protocol, Version: service.Version, Metadata: service.Metadata}
				host.Services = append(host.Services, serviceModel)
			}
			db.Create(host)
		}
	} else {
		for _, service := range hostService.Ports {
			dpr.AddOrUpdateService(host, service)
		}
	}

}
