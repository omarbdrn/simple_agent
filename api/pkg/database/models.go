package database

import (
	"time"

	"gorm.io/gorm"
)

type Share struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"`
	IPRanges []IPRange

	CreatedAt time.Time
}

type IPRange struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey"`
	IPRange string `gorm:"unique"`
	Taken   bool   `gorm:"default:false"`

	CreatedAt     time.Time
	LastScannedAt time.Time
	ShareID       uint
}

type Question struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	IPRange  string `gorm:"unique"`
	Answered bool   `gorm:"default:false"`

	CreatedAt     time.Time
	LastScannedAt time.Time
}

type Host struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	IP       string `gorm:"unique"`
	Services []Service

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Service struct {
	gorm.Model
	ID       uint `gorm:"primaryKey"`
	Port     int
	Protocol string
	Version  string
	Metadata string

	CreatedAt time.Time
	UpdatedAt time.Time
	HostID    uint
}
