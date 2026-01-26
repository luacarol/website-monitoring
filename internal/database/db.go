package database

import (
	"log"

	"github.com/luacarol/website-monitoring/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	var err error

	// Configurar GORM com SQLite
	DB, err = gorm.Open(sqlite.Open("monitor.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Falha ao conectar com o banco de dados:", err)
	}

	// Auto-migrate dos models
	err = DB.AutoMigrate(
		&models.Site{},
		&models.MonitorLog{},
	)

	if err != nil {
		log.Fatal("Falha ao executar migrations:", err)
	}

	log.Println("✅ Banco de dados configurado com sucesso!")
}

func GetDB() *gorm.DB {
	return DB
}
