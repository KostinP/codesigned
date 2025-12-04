// cmd/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/kostinp/edu-platform-backend/docs"

	"github.com/kostinp/edu-platform-backend/internal/shared/abac"
	"github.com/kostinp/edu-platform-backend/internal/shared/config"
	"github.com/kostinp/edu-platform-backend/internal/shared/db"
	"github.com/kostinp/edu-platform-backend/internal/shared/logger"
	"github.com/kostinp/edu-platform-backend/internal/user/usecase"
)

// @title Edu Platform API
// @version 1.0
// @description Backend for Edu Platform with gamification
// @termsOfService https://edu-platform.com/terms
// @contact.name Support Team
// @contact.email support@edu-platform.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	fmt.Printf("=== CONFIG DEBUG ===\n")
	fmt.Printf("DB Host: %s\n", cfg.Database.Host)
	fmt.Printf("DB Port: %s\n", cfg.Database.Port)
	fmt.Printf("DB User: %s\n", cfg.Database.User)
	fmt.Printf("DB Password: %s\n", "***") // Не выводите реальный пароль
	fmt.Printf("DB Name: %s\n", cfg.Database.Name)
	fmt.Printf("====================\n")

	log.Printf("=== STARTING SERVER DEBUG ===")
	log.Printf("Postgres: %s:%s", cfg.Database.Host, cfg.Database.Port)
	fmt.Printf("Postgres DB User: %s\n", cfg.Database.User)
	fmt.Printf("Postgres DB Name: %s\n", cfg.Database.Name)
	log.Printf("ClickHouse: %s:%s", cfg.Clickhouse.Host, cfg.Clickhouse.Port)
	log.Printf("JWT Secret length: %d", len(cfg.JWT.Secret))
	log.Printf("Analytics enabled: %v", cfg.Analytics.Enabled)
	log.Printf("=== SERVER DEBUG END ===")

	chConn := db.ConnectClickhouse(cfg)
	if chConn == nil && cfg.Analytics.Enabled {
		logger.Info("Аналитика отключена: ClickHouse недоступен")
		// Здесь можно set flag или inject nil в Wire для repo
	}

	pool := db.ConnectPostgres(cfg)
	// Создаём ABAC движок
	abacEngine := abac.NewABACEngine()
	// Загружаем политики из БД
	if err := abac.LoadPoliciesFromDB(abacEngine, pool); err != nil {
		logger.Fatal("Не удалось загрузить ABAC политики", err)
	}

	if len(abacEngine.Policies) == 0 {
		for _, p := range abac.GetDefaultPolicies() {
			abacEngine.AddPolicy(p)
		}
		logger.Info("Загружены дефолтные ABAC-политики")
	}

	// Инициализируем usecase отдельно
	sessionUsecase, err := InitializeSessionUsecase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Запускаем задачу очистки
	go usecase.StartSessionCleanupTask(context.Background(), sessionUsecase, time.Hour*24)

	// Теперь инициализируем сервер с готовым engine
	server, err := InitializeServer(cfg, abacEngine)
	if err != nil {
		logger.Fatal("Ошибка инициализации сервера", err)
	}
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Info(fmt.Sprintf("🚀 Запуск сервера на %s", addr))
	if err := server.Start(addr); err != nil {
		logger.Fatal("❌ Не удалось запустить сервер", err)
	}
}
