package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/kostinp/edu-platform-backend/internal/user/entity"
)

// AnalyticsEvent расширяет базовое событие для аналитики
type AnalyticsEvent struct {
	entity.VisitorEvent
	PageURL          string `json:"page_url,omitempty"`
	Referrer         string `json:"referrer,omitempty"`
	UTMSource        string `json:"utm_source,omitempty"`
	UTMMedium        string `json:"utm_medium,omitempty"`
	UTMCampaign      string `json:"utm_campaign,omitempty"`
	ScreenResolution string `json:"screen_resolution,omitempty"`
	Language         string `json:"language,omitempty"`
}

// VisitorEventRepository - основной интерфейс репозитория событий посетителей
type VisitorEventRepository interface {
	Create(ctx context.Context, event *entity.VisitorEvent) error
	GetEventsByVisitor(ctx context.Context, visitorID string) ([]*entity.VisitorEvent, error)
}

// AnalyticsRepository - интерфейс для аналитических методов
type AnalyticsRepository interface {
	GetPageViews(ctx context.Context, days int) ([]map[string]interface{}, error)
	GetConversionRate(ctx context.Context, eventType string, days int) (float64, error)
	GetUTMStats(ctx context.Context, days int) ([]map[string]interface{}, error)
}

// PrimaryVisitorEventRepo - интерфейс для основного репозитория (ClickHouse)
type PrimaryVisitorEventRepo interface {
	VisitorEventRepository
	AnalyticsRepository
}

// ClickHouseVisitorEventRepo - ClickHouse реализация с аналитическими методами
type ClickHouseVisitorEventRepo struct {
	conn clickhouse.Conn
}

// Убедитесь, что ClickHouseVisitorEventRepo реализует PrimaryVisitorEventRepo
var _ PrimaryVisitorEventRepo = (*ClickHouseVisitorEventRepo)(nil)

// NewClickHouseVisitorEventRepo создает новый экземпляр ClickHouse репозитория
func NewClickHouseVisitorEventRepo(conn clickhouse.Conn) *ClickHouseVisitorEventRepo {
	return &ClickHouseVisitorEventRepo{conn: conn}
}

// Create сохраняет событие в ClickHouse
func (r *ClickHouseVisitorEventRepo) Create(ctx context.Context, event *entity.VisitorEvent) error {
	if r.conn == nil {
		return fmt.Errorf("clickhouse connection is not available")
	}
	eventJSON, err := json.Marshal(event.EventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}
	if r.conn == nil {
		return fmt.Errorf("analytics disabled")
	}
	query := `
        INSERT INTO visitor_events (
            id, visitor_id, event_type, event_data, created_at
        ) VALUES (?, ?, ?, ?, ?)
    `
	err = r.conn.Exec(ctx, query,
		event.ID,
		event.VisitorID,
		event.EventType,
		string(eventJSON),
		event.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert into ClickHouse: %w", err)
	}
	return nil
}

// CreateAnalyticsEvent сохраняет аналитическое событие в ClickHouse
func (r *ClickHouseVisitorEventRepo) CreateAnalyticsEvent(ctx context.Context, event *AnalyticsEvent) error {
	if r.conn == nil {
		return fmt.Errorf("clickhouse connection is not available")
	}
	eventJSON, err := json.Marshal(event.EventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}
	query := `
		INSERT INTO visitor_events (
			id, visitor_id, event_type, event_data, created_at, page_url, referrer, utm_source, utm_medium, utm_campaign, screen_resolution, language) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
	batch, err := r.conn.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}
	err = batch.Append(
		event.ID,
		event.VisitorID,
		event.EventType,
		string(eventJSON),
		event.CreatedAt,
		event.PageURL,
		event.Referrer,
		event.UTMSource,
		event.UTMMedium,
		event.UTMCampaign,
		event.ScreenResolution,
		event.Language,
	)
	if err != nil {
		return fmt.Errorf("failed to append to batch: %w", err)
	}
	return batch.Send()
}

// GetEventsByVisitor возвращает события из ClickHouse (базовая реализация)
func (r *ClickHouseVisitorEventRepo) GetEventsByVisitor(ctx context.Context, visitorID string) ([]*entity.VisitorEvent, error) {
	if r.conn == nil {
		return nil, fmt.Errorf("clickhouse connection is not available")
	}
	if r.conn == nil {
		return nil, fmt.Errorf("analytics disabled")
	}
	query := `
		SELECT id, visitor_id, event_type, event_data, created_at
		FROM visitor_events
		WHERE visitor_id = ?
		ORDER BY created_at DESC
		`
	rows, err := r.conn.Query(ctx, query, visitorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query visitor events: %w", err)
	}
	defer rows.Close()
	var events []*entity.VisitorEvent
	for rows.Next() {
		var event entity.VisitorEvent
		var eventJSON string
		err := rows.Scan(
			&event.ID,
			&event.VisitorID,
			&event.EventType,
			&eventJSON,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan visitor event: %w", err)
		}
		// Декодируем JSON данные
		var eventData map[string]interface{}
		if err := json.Unmarshal([]byte(eventJSON), &eventData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event data: %w", err)
		}
		event.EventData = eventData
		events = append(events, &event)
	}
	return events, nil
}

// Аналитические методы для ClickHouse
// GetPageViews возвращает статистику просмотров страниц
func (r *ClickHouseVisitorEventRepo) GetPageViews(ctx context.Context, days int) ([]map[string]interface{}, error) {
	if r.conn == nil {
		return nil, fmt.Errorf("clickhouse connection is not available")
	}
	if r.conn == nil {
		return nil, fmt.Errorf("analytics disabled")
	}
	query := `
		SELECT 
			page_url,
			count(*) as views,
			count(distinct visitor_id) as unique_visitors
		FROM visitor_events
		WHERE event_type = 'page_view'
		AND created_at >= now() - INTERVAL ? DAY
		AND page_url != ''
		GROUP BY page_url
		ORDER BY views DESC
		`
	rows, err := r.conn.Query(ctx, query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query page views: %w", err)
	}
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		var pageURL string
		var views, uniqueVisitors uint64
		err := rows.Scan(&pageURL, &views, &uniqueVisitors)
		if err != nil {
			return nil, fmt.Errorf("failed to scan page view: %w", err)
		}
		results = append(results, map[string]interface{}{
			"page_url":        pageURL,
			"views":           views,
			"unique_visitors": uniqueVisitors,
		})
	}
	return results, nil
}

// GetConversionRate возвращает коэффициент конверсии
func (r *ClickHouseVisitorEventRepo) GetConversionRate(ctx context.Context, eventType string, days int) (float64, error) {
	if r.conn == nil {
		return 0, fmt.Errorf("clickhouse connection is not available")
	}
	if r.conn == nil {
		return 0, fmt.Errorf("analytics disabled")
	}
	query := `
		SELECT count(distinct case when event_type = ? then visitor_id end) * 100.0 / count(distinct visitor_id) as conversion_rate
		FROM visitor_events
		WHERE created_at >= now() - INTERVAL ? DAY
		`
	var conversionRate float64
	err := r.conn.QueryRow(ctx, query, eventType, days).Scan(&conversionRate)
	if err != nil {
		return 0, fmt.Errorf("failed to query conversion rate: %w", err)
	}
	return conversionRate, nil
}

// GetUTMStats возвращает статистику UTM меток
func (r *ClickHouseVisitorEventRepo) GetUTMStats(ctx context.Context, days int) ([]map[string]interface{}, error) {
	if r.conn == nil {
		return nil, fmt.Errorf("clickhouse connection is not available")
	}
	if r.conn == nil {
		return nil, fmt.Errorf("analytics disabled")
	}
	query := `
		SELECT
			utm_source,
			utm_medium,
			utm_campaign,
			count(*) as events,
			count(distinct visitor_id) as unique_visitors
		FROM visitor_events
		WHERE created_at >= now() - INTERVAL ? DAY
		AND utm_source != ''
		GROUP BY utm_source, utm_medium, utm_campaign
		ORDER BY events DESC
		`
	rows, err := r.conn.Query(ctx, query, days)
	if err != nil {
		return nil, fmt.Errorf("failed to query UTM stats: %w", err)
	}
	defer rows.Close()
	var results []map[string]interface{}
	for rows.Next() {
		var source, medium, campaign string
		var events, uniqueVisitors uint64
		err := rows.Scan(&source, &medium, &campaign, &events, &uniqueVisitors)
		if err != nil {
			return nil, fmt.Errorf("failed to scan UTM stat: %w", err)
		}
		results = append(results, map[string]interface{}{
			"utm_source":      source,
			"utm_medium":      medium,
			"utm_campaign":    campaign,
			"events":          events,
			"unique_visitors": uniqueVisitors,
		})
	}
	return results, nil
}

// VisitorEventRepoWithFallback - репозиторий (упрощен, только primary)
type VisitorEventRepoWithFallback struct {
	primary PrimaryVisitorEventRepo // ClickHouse
}

// NewVisitorEventRepoWithFallback создает новый экземпляр (без fallback)
func NewVisitorEventRepoWithFallback(primary PrimaryVisitorEventRepo, fallback interface{}) *VisitorEventRepoWithFallback {
	return &VisitorEventRepoWithFallback{
		primary: primary,
	}
}

// Create сохраняет событие
func (r *VisitorEventRepoWithFallback) Create(ctx context.Context, event *entity.VisitorEvent) error {
	if r.primary == nil {
		return ErrNoRepositoryAvailable
	}
	if r.primary == nil {
		return fmt.Errorf("analytics disabled")
	}
	return r.primary.Create(ctx, event)
}

// GetEventsByVisitor возвращает события
func (r *VisitorEventRepoWithFallback) GetEventsByVisitor(ctx context.Context, visitorID string) ([]*entity.VisitorEvent, error) {
	if r.primary == nil {
		return nil, ErrNoRepositoryAvailable
	}
	if r.primary == nil {
		return nil, fmt.Errorf("analytics disabled")
	}
	return r.primary.GetEventsByVisitor(ctx, visitorID)
}

// Ошибки
var (
	ErrNoRepositoryAvailable = fmt.Errorf("no repository available for visitor events")
)
