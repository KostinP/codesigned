package usecase

import (
	"context"
	"log"
	"time"
)

func StartSessionCleanupTask(ctx context.Context, sessionUC SessionUsecase, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := sessionUC.DeleteExpiredSessions(ctx); err != nil {
					log.Printf("session cleanup error: %v", err)
				} else {
					log.Printf("session cleanup completed")
				}
			case <-ctx.Done():
				log.Printf("session cleanup stopped")
				return
			}
		}
	}()
}
