package models

import (
	"time"

	"gorm.io/gorm"
)

type MonitorLog struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	SiteID       uint           `json:"site_id" gorm:"not null"`
	StatusCode   int            `json:"status_code"`
	ResponseTime int64          `json:"response_time"` // em millisegundos
	IsOnline     bool           `json:"is_online"`
	ErrorMessage string         `json:"error_message"`
	CheckedAt    time.Time      `json:"checked_at"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relacionamento
	Site Site `json:"site,omitempty" gorm:"foreignKey:SiteID"`
}

type LogsQuery struct {
	SiteID    uint   `form:"site_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Status    string `form:"status"` // "online", "offline", "all"
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
}

type StatsResponse struct {
	TotalSites    int       `json:"total_sites"`
	OnlineSites   int       `json:"online_sites"`
	OfflineSites  int       `json:"offline_sites"`
	OverallUptime float64   `json:"overall_uptime"`
	LastUpdate    time.Time `json:"last_update"`
}
