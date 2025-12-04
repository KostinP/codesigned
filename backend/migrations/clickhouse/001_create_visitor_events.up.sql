CREATE TABLE IF NOT EXISTS analytics.visitor_events (
    id UUID,
    visitor_id UUID,
    event_type String,
    event_data String,
    created_at DateTime64(3)
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(created_at)
ORDER BY (created_at, id);

-- Индексы для быстрых запросов
ALTER TABLE analytics.visitor_events ADD INDEX idx_visitor_id visitor_id TYPE bloom_filter GRANULARITY 1;
ALTER TABLE analytics.visitor_events ADD INDEX idx_event_type event_type TYPE bloom_filter GRANULARITY 1;
ALTER TABLE analytics.visitor_events ADD INDEX idx_created_at created_at TYPE minmax GRANULARITY 1;