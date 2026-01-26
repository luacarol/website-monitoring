package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luacarol/website-monitoring/internal/database"
	"github.com/luacarol/website-monitoring/internal/models"
)

// GET /api/sites
func GetSites(c *gin.Context) {
	var sites []models.Site
	db := database.GetDB()

	// Buscar sites com informações do último check
	if err := db.Find(&sites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar sites"})
		return
	}

	// Para cada site, buscar informações do último log
	var sitesResponse []models.SiteResponse
	for _, site := range sites {
		var lastLog models.MonitorLog
		db.Where("site_id = ?", site.ID).Order("checked_at desc").First(&lastLog)

		siteResponse := models.SiteResponse{
			Site:       site,
			LastStatus: lastLog.StatusCode,
			LastCheck:  lastLog.CheckedAt,
			Uptime:     calculateUptime(site.ID), // Implementar função
		}
		sitesResponse = append(sitesResponse, siteResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"sites": sitesResponse,
		"total": len(sitesResponse),
	})
}

// POST /api/sites
func CreateSite(c *gin.Context) {
	var request models.SiteRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	site := models.Site{
		Name:   request.Name,
		URL:    request.URL,
		Active: true,
	}

	db := database.GetDB()
	if err := db.Create(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar site"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Site criado com sucesso",
		"site":    site,
	})
}

// DELETE /api/sites/:id
func DeleteSite(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	db := database.GetDB()

	// Verificar se site existe
	var site models.Site
	if err := db.First(&site, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site não encontrado"})
		return
	}

	// Soft delete
	if err := db.Delete(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar site"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Site removido com sucesso"})
}

// PUT /api/sites/:id/toggle
func ToggleSite(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	db := database.GetDB()
	var site models.Site

	if err := db.First(&site, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site não encontrado"})
		return
	}

	site.Active = !site.Active
	if err := db.Save(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar site"})
		return
	}

	status := "ativado"
	if !site.Active {
		status = "desativado"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Site " + status + " com sucesso",
		"site":    site,
	})
}

// Função helper para calcular uptime
func calculateUptime(siteID uint) float64 {
	db := database.GetDB()

	var totalLogs int64
	var onlineLogs int64

	// Contar logs dos últimos 24h
	since := time.Now().Add(-24 * time.Hour)
	db.Model(&models.MonitorLog{}).Where("site_id = ? AND checked_at >= ?", siteID, since).Count(&totalLogs)
	db.Model(&models.MonitorLog{}).Where("site_id = ? AND checked_at >= ? AND is_online = ?", siteID, since, true).Count(&onlineLogs)

	if totalLogs == 0 {
		return 0.0
	}

	return (float64(onlineLogs) / float64(totalLogs)) * 100
}
