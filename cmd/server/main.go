package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/luacarol/website-monitoring/internal/database"
	"github.com/luacarol/website-monitoring/internal/handlers"
	"github.com/luacarol/website-monitoring/internal/services"
)

var (
	monitorService *services.MonitorService
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Permitir todas as origens em desenvolvimento
		},
	}
)

func main() {
	// Inicializar banco de dados
	database.InitDatabase()

	// Inicializar serviço de monitoramento
	monitorService = services.NewMonitorService()
	monitorService.Start()

	// Configurar Gin
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // React dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes da API
	api := router.Group("/api")
	{
		// Sites
		api.GET("/sites", handlers.GetSites)
		api.POST("/sites", handlers.CreateSite)
		api.DELETE("/sites/:id", handlers.DeleteSite)
		api.PUT("/sites/:id/toggle", handlers.ToggleSite)

		// Logs
		api.GET("/logs", handlers.GetLogs)

		// Stats
		api.GET("/stats", handlers.GetStats)

		// Monitor
		api.POST("/monitor/check/:id", checkSiteNow)
		api.GET("/monitor/status", getMonitorStatus)
	}

	// WebSocket para updates em tempo real
	router.GET("/ws", handleWebSocket)

	// Servir arquivos estáticos do React (quando buildado)
	router.Static("/static", "./web/build/static")
	router.StaticFile("/", "./web/build/index.html")
	router.StaticFile("/favicon.ico", "./web/build/favicon.ico")

	// Capturar sinais do sistema para graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("\n🛑 Recebido sinal de parada...")
		monitorService.Stop()
		os.Exit(0)
	}()

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor iniciado na porta %s", port)
	log.Printf("📊 Dashboard: http://localhost:%s", port)
	log.Printf("🔗 API: http://localhost:%s/api", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

// Handler para verificar site imediatamente
func checkSiteNow(c *gin.Context) {
	siteID := c.Param("id")
	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID do site é obrigatório"})
		return
	}

	// Converter string para uint
	var id uint
	if _, err := fmt.Sscanf(siteID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	result := monitorService.CheckSiteNow(id)
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verificação realizada",
		"result":  result,
	})
}

// Handler para status do monitor
func getMonitorStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"running":   monitorService.IsRunning(),
		"timestamp": time.Now(),
	})
}

// WebSocket para updates em tempo real
func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ Erro ao fazer upgrade WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("🔗 Nova conexão WebSocket estabelecida")

	// Loop para manter conexão ativa e enviar updates
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Enviar stats atualizadas
		stats := getStatsData()
		if err := conn.WriteJSON(stats); err != nil {
			log.Printf("❌ Erro ao enviar dados WebSocket: %v", err)
			return
		}
	}
}

// Helper para obter dados de stats
func getStatsData() map[string]interface{} {
	// Implementar lógica similar ao handlers.GetStats
	// Por simplicidade, retornando dados mock aqui
	return map[string]interface{}{
		"type":      "stats_update",
		"timestamp": time.Now(),
		"data": map[string]interface{}{
			"total_sites":    5,
			"online_sites":   4,
			"offline_sites":  1,
			"overall_uptime": 95.5,
		},
	}
}
