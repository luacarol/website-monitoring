package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luacarol/website-monitoring/internal/database"
	"github.com/luacarol/website-monitoring/internal/models"
)

// GET /api/logs
func GetLogs(c *gin.Context) {
	var query models.LogsQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default values
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 50
	}
	if query.Status == "" {
		query.Status = "all"
	}

	db := database.GetDB()

	// Build query
	dbQuery := db.Model(&models.MonitorLog{}).Preload("Site")

	// Filters
	if query.SiteID != 0 {
		dbQuery = dbQuery.Where("site_id = ?", query.SiteID)
	}

	if query.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", query.StartDate)
		if err == nil {
			dbQuery = dbQuery.Where("checked_at >= ?", startDate)
		}
	}

	if query.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", query.EndDate)
		if err == nil {
			dbQuery = dbQuery.Where("checked_at <= ?", endDate.Add(24*time.Hour))
		}
	}

	if query.Status == "online" {
		dbQuery = dbQuery.Where("is_online = ?", true)
	} else if query.Status == "offline" {
		dbQuery = dbQuery.Where("is_online = ?", false)
	}

	// Count total
	var total int64
	dbQuery.Count(&total)

	// Fetch with pagination
	var logs []models.MonitorLog
	offset := (query.Page - 1) * query.Limit

	if err := dbQuery.Order("checked_at desc").Limit(query.Limit).Offset(offset).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  query.Page,
		"limit": query.Limit,
		"pages": (total + int64(query.Limit) - 1) / int64(query.Limit),
	})
}

// GET /api/stats
func GetStats(c *gin.Context) {
	db := database.GetDB()

	var stats models.StatsResponse

	// Total active sites
	var totalSites int64
	db.Model(&models.Site{}).Where("active = ?", true).Count(&totalSites)
	stats.TotalSites = int(totalSites)

	// Fetch all active sites
	var activeSites []models.Site
	db.Where("active = ?", true).Find(&activeSites)

	// Count online/offline sites based on last check of each site
	onlineCount := 0
	offlineCount := 0

	for _, site := range activeSites {
		var lastLog models.MonitorLog
		err := db.Where("site_id = ?", site.ID).
			Order("checked_at DESC").
			First(&lastLog).Error

		if err == nil {
			// Site has logs
			if lastLog.IsOnline {
				onlineCount++
			} else {
				offlineCount++
			}
		}
		// If no logs yet, don't count in either
	}

	stats.OnlineSites = onlineCount
	stats.OfflineSites = offlineCount

	// Overall uptime (last 24h)
	since := time.Now().Add(-24 * time.Hour)
	var totalChecks int64
	var onlineChecks int64

	db.Model(&models.MonitorLog{}).Where("checked_at >= ?", since).Count(&totalChecks)
	db.Model(&models.MonitorLog{}).Where("checked_at >= ? AND is_online = ?", since, true).Count(&onlineChecks)

	if totalChecks > 0 {
		stats.OverallUptime = (float64(onlineChecks) / float64(totalChecks)) * 100
	}

	stats.LastUpdate = time.Now()

	c.JSON(http.StatusOK, stats)
}
