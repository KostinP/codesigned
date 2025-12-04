package user

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/kostinp/edu-platform-backend/internal/shared/config"
	"github.com/kostinp/edu-platform-backend/internal/shared/db"
)

func ProvideBotToken(cfg *config.Config) config.BotToken {
	return config.BotToken(cfg.Telegram.Token)
}

func ProvideJwtSecret(cfg *config.Config) config.JwtSecret {
	return config.JwtSecret(cfg.JWT.Secret)
}

// ProvideClickHouseConn предоставляет подключение к ClickHouse
func ProvideClickHouseConn(cfg *config.Config) clickhouse.Conn {
	return db.ConnectClickhouse(cfg)
}
