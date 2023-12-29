package database

import (
	"gorm.io/gorm"
)

type DatabaseReporter struct {
	DB *gorm.DB
}

// func (dpr *DatabaseReporter) GetAssetDetails(ssetDomain string) []Subdomain{
// 	dpr.DB = GetDB()
// 	var subdomains []Subdomain
// 	var result *gorm.DB

// 	result = db.Model(&Subdomain{}).Where("TLD = ?", AssetDomain).Preload("Network").Preload("Network.Services").Preload("WebServer").Preload("WebServer.Technologies").Preload("Paths").Find(&subdomains)

// 	if result.Error != nil || result.RowsAffected == 0 {
// 		return []Subdomain{}
// 	}

// 	return subdomains
// }
