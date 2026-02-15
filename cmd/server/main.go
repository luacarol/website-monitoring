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
			return true // Allow all origins in development
		},
	}
)

func main() {
	// Initialize database
	database.InitDatabase()

	// Initialize monitoring service
	monitorService = services.NewMonitorService()
	monitorService.Start()

	// Configure Gin
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // React dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API Routes
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

	// WebSocket for real-time updates
	router.GET("/ws", handleWebSocket)

	// Serve React static files (when built)
	router.Static("/static", "./web/build/static")
	router.StaticFile("/", "./web/build/index.html")
	router.StaticFile("/favicon.ico", "./web/build/favicon.ico")

	// Capture system signals for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("\n🛑 Received stop signal...")
		monitorService.Stop()
		os.Exit(0)
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server started on port %s", port)
	log.Printf("📊 Dashboard: http://localhost:%s", port)
	log.Printf("🔗 API: http://localhost:%s/api", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

// Handler to check site immediately
func checkSiteNow(c *gin.Context) {
	siteID := c.Param("id")
	if siteID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Site ID is required"})
		return
	}

	// Convert string to uint
	var id uint
	if _, err := fmt.Sscanf(siteID, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result := monitorService.CheckSiteNow(id)
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Check performed",
		"result":  result,
	})
}

// Handler for monitor status
func getMonitorStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"running":   monitorService.IsRunning(),
		"timestamp": time.Now(),
	})
}

// WebSocket for real-time updates
func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("❌ Error upgrading WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("🔗 New WebSocket connection established")

	// Loop to keep connection alive and send updates
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Send updated stats
		stats := getStatsData()
		if err := conn.WriteJSON(stats); err != nil {
			log.Printf("❌ Error sending WebSocket data: %v", err)
			return
		}
	}
}

// Helper to get stats data
func getStatsData() map[string]interface{} {
	// Implement logic similar to handlers.GetStats
	// For simplicity, returning mock data here
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
