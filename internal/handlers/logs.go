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

	// Valores padrão
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

	// Construir query
	dbQuery := db.Model(&models.MonitorLog{}).Preload("Site")

	// Filtros
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

	// Contar total
	var total int64
	dbQuery.Count(&total)

	// Buscar com paginação
	var logs []models.MonitorLog
	offset := (query.Page - 1) * query.Limit

	if err := dbQuery.Order("checked_at desc").Limit(query.Limit).Offset(offset).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar logs"})
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

	// Total de sites ativos
	var totalSites int64
	db.Model(&models.Site{}).Where("active = ?", true).Count(&totalSites)
	stats.TotalSites = int(totalSites)

	// Sites online/offline (baseado no último check)
	// Query mais complexa para obter último status de cada site
	var onlineCount int64
	var offlineCount int64

	// Subquery para pegar o último log de cada site
	subQuery := db.Table("monitor_logs").
		Select("site_id, MAX(checked_at) as last_check").
		Group("site_id")

	// Join para pegar o status do último check
	db.Table("monitor_logs ml1").
		Joins("JOIN (?) ml2 ON ml1.site_id = ml2.site_id AND ml1.checked_at = ml2.last_check", subQuery).
		Joins("JOIN sites s ON ml1.site_id = s.id").
		Where("s.active = ? AND ml1.is_online = ?", true, true).
		Count(&onlineCount)

	db.Table("monitor_logs ml1").
		Joins("JOIN (?) ml2 ON ml1.site_id = ml2.site_id AND ml1.checked_at = ml2.last_check", subQuery).
		Joins("JOIN sites s ON ml1.site_id = s.id").
		Where("s.active = ? AND ml1.is_online = ?", true, false).
		Count(&offlineCount)

	stats.OnlineSites = int(onlineCount)
	stats.OfflineSites = int(offlineCount)

	// Uptime geral (últimas 24h)
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
