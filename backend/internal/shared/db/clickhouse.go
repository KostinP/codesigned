package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/kostinp/edu-platform-backend/internal/shared/config"
	"github.com/kostinp/edu-platform-backend/internal/shared/logger"
)

// ConnectClickhouse возвращает conn или nil, если подключение не удалось
func ConnectClickhouse(cfg *config.Config) clickhouse.Conn {
	options := &clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", cfg.Clickhouse.Host, cfg.Clickhouse.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Clickhouse.Database,
			Username: cfg.Clickhouse.Username,
			Password: cfg.Clickhouse.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		// Protocol: clickhouse.HTTP,
	}

	conn, err := clickhouse.Open(options)
	if err != nil {
		logger.Error("Не удалось подключиться к ClickHouse", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		logger.Error("ClickHouse недоступен", err)
		return nil
	}

	logger.Info("Подключение к ClickHouse установлено")
	return conn
}
