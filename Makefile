# Makefile para Website Monitor

.PHONY: help build run-cli run-server dev-frontend build-frontend clean

help: ## Mostrar ajuda
	@echo "Comandos disponíveis:"
	@echo "  run-cli       - Executar versão CLI"
	@echo "  run-server    - Executar servidor web"
	@echo "  dev-frontend  - Iniciar frontend em modo desenvolvimento"
	@echo "  build         - Compilar binários"
	@echo "  clean         - Limpar arquivos temporários"

# Backend
run-cli: ## Executar CLI atual
	go run cmd/cli/main.go

run-server: ## Executar servidor web
	go run cmd/server/main.go

build: ## Compilar binários
	go build -o bin/monitor-cli cmd/cli/main.go
	go build -o bin/monitor-server cmd/server/main.go

# Frontend
setup-frontend: ## Configurar frontend React
	cd web && npm install

dev-frontend: ## Iniciar desenvolvimento frontend
	cd web && npm start

build-frontend: ## Build frontend para produção
	cd web && npm run build

# Database
migrate: ## Executar migrations (GORM faz automático)
	@echo "GORM AutoMigrate executado automaticamente"

# Utils
clean: ## Limpar arquivos temporários
	rm -rf bin/
	rm -f monitor.db
	cd web && rm -rf build/

deps: ## Instalar dependências Go
	go mod tidy
	go mod download

dev: ## Ambiente completo de desenvolvimento
	@echo "🚀 Iniciando ambiente de desenvolvimento..."
	@echo "📊 Backend: http://localhost:8080"
	@echo "🎨 Frontend: http://localhost:3000"
	make run-server &
	sleep 3
	make dev-frontend

.DEFAULT_GOAL := help
