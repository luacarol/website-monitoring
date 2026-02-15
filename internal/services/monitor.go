package services

import (
	"log"
	"net/http"
	"time"

	"github.com/luacarol/website-monitoring/internal/database"
	"github.com/luacarol/website-monitoring/internal/models"
)

type MonitorService struct {
	isRunning bool
	stopChan  chan bool
}

func NewMonitorService() *MonitorService {
	return &MonitorService{
		isRunning: false,
		stopChan:  make(chan bool),
	}
}

// Start continuous monitoring
func (m *MonitorService) Start() {
	if m.isRunning {
		log.Println("⚠️ Monitoring is already running")
		return
	}

	m.isRunning = true
	log.Println("🚀 Starting monitoring service...")

	// Goroutine for continuous monitoring
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Check a cada 30 segundos
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.checkAllSites()
			case <-m.stopChan:
				log.Println("⏹️ Stopping monitoring service...")
				m.isRunning = false
				return
			}
		}
	}()
}

// Parar monitoramento
func (m *MonitorService) Stop() {
	if m.isRunning {
		m.stopChan <- true
	}
}

// Check all active sites
func (m *MonitorService) checkAllSites() {
	db := database.GetDB()

	var sites []models.Site
	if err := db.Where("active = ?", true).Find(&sites).Error; err != nil {
		log.Printf("❌ Error fetching sites: %v", err)
		return
	}

	log.Printf("🔍 Verificando %d sites...", len(sites))

	for _, site := range sites {
		go m.checkSite(site) // Parallel checking
	}
}

// Check a specific site
func (m *MonitorService) checkSite(site models.Site) {
	startTime := time.Now()

	// Make HTTP request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(site.URL)

	// Calculate response time
	responseTime := time.Since(startTime).Milliseconds()

	// Create log entry
	monitorLog := models.MonitorLog{
		SiteID:       site.ID,
		ResponseTime: responseTime,
		CheckedAt:    time.Now(),
	}

	if err != nil {
		// Request error
		monitorLog.IsOnline = false
		monitorLog.StatusCode = 0
		monitorLog.ErrorMessage = err.Error()

		log.Printf("❌ %s - OFFLINE: %v", site.Name, err)
	} else {
		defer response.Body.Close()

		monitorLog.StatusCode = response.StatusCode

		if response.StatusCode >= 200 && response.StatusCode < 400 {
			monitorLog.IsOnline = true
			log.Printf("✅ %s - ONLINE (%d) - %dms", site.Name, response.StatusCode, responseTime)
		} else {
			monitorLog.IsOnline = false
			log.Printf("⚠️ %s - PROBLEMA (%d) - %dms", site.Name, response.StatusCode, responseTime)
		}
	}

	// Save to database
	db := database.GetDB()
	if err := db.Create(&monitorLog).Error; err != nil {
		log.Printf("❌ Error saving log for %s: %v", site.Name, err)
	}
}

// Check individual site (for API)
func (m *MonitorService) CheckSiteNow(siteID uint) *models.MonitorLog {
	db := database.GetDB()

	var site models.Site
	if err := db.First(&site, siteID).Error; err != nil {
		log.Printf("❌ Site not found: %d", siteID)
		return nil
	}

	m.checkSite(site)

	// Return last log
	var lastLog models.MonitorLog
	db.Where("site_id = ?", siteID).Order("checked_at desc").First(&lastLog)

	return &lastLog
}

func (m *MonitorService) IsRunning() bool {
	return m.isRunning
}
