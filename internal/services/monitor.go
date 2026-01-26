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

// Iniciar monitoramento contínuo
func (m *MonitorService) Start() {
	if m.isRunning {
		log.Println("⚠️ Monitoramento já está rodando")
		return
	}

	m.isRunning = true
	log.Println("🚀 Iniciando serviço de monitoramento...")

	// Goroutine para monitoramento contínuo
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Check a cada 30 segundos
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.checkAllSites()
			case <-m.stopChan:
				log.Println("⏹️ Parando serviço de monitoramento...")
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

// Verificar todos os sites ativos
func (m *MonitorService) checkAllSites() {
	db := database.GetDB()

	var sites []models.Site
	if err := db.Where("active = ?", true).Find(&sites).Error; err != nil {
		log.Printf("❌ Erro ao buscar sites: %v", err)
		return
	}

	log.Printf("🔍 Verificando %d sites...", len(sites))

	for _, site := range sites {
		go m.checkSite(site) // Verificação paralela
	}
}

// Verificar um site específico
func (m *MonitorService) checkSite(site models.Site) {
	startTime := time.Now()

	// Fazer requisição HTTP
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(site.URL)

	// Calcular tempo de resposta
	responseTime := time.Since(startTime).Milliseconds()

	// Criar log entry
	monitorLog := models.MonitorLog{
		SiteID:       site.ID,
		ResponseTime: responseTime,
		CheckedAt:    time.Now(),
	}

	if err != nil {
		// Erro na requisição
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

	// Salvar no banco
	db := database.GetDB()
	if err := db.Create(&monitorLog).Error; err != nil {
		log.Printf("❌ Erro ao salvar log para %s: %v", site.Name, err)
	}
}

// Verificar site individual (para API)
func (m *MonitorService) CheckSiteNow(siteID uint) *models.MonitorLog {
	db := database.GetDB()

	var site models.Site
	if err := db.First(&site, siteID).Error; err != nil {
		log.Printf("❌ Site não encontrado: %d", siteID)
		return nil
	}

	m.checkSite(site)

	// Retornar último log
	var lastLog models.MonitorLog
	db.Where("site_id = ?", siteID).Order("checked_at desc").First(&lastLog)

	return &lastLog
}

func (m *MonitorService) IsRunning() bool {
	return m.isRunning
}
