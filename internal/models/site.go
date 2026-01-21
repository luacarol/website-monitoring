package models

import (
	"time"

	"gorm.io/gorm"
)

type Site struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	URL       string         `json:"url" gorm:"not null;unique"`
	Active    bool           `json:"active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relacionamento com logs
	Logs []MonitorLog `json:"logs,omitempty" gorm:"foreignKey:SiteID"`

	// Campos computados (não salvos no banco)
	LastStatus *int       `json:"last_status,omitempty" gorm:"-"`
	LastCheck  *time.Time `json:"last_check,omitempty" gorm:"-"`
	Uptime     *float64   `json:"uptime,omitempty" gorm:"-"`
}

type SiteRequest struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required,url"`
}

type SiteResponse struct {
	Site
	LastStatus int       `json:"last_status"`
	LastCheck  time.Time `json:"last_check"`
	Uptime     float64   `json:"uptime"`
}
