-- Расширяем таблицу для лендинг-аналитики
ALTER TABLE analytics.visitor_events 
ADD COLUMN IF NOT EXISTS page_url String,
ADD COLUMN IF NOT EXISTS referrer String,
ADD COLUMN IF NOT EXISTS utm_source String,
ADD COLUMN IF NOT EXISTS utm_medium String,
ADD COLUMN IF NOT EXISTS utm_campaign String,
ADD COLUMN IF NOT EXISTS screen_resolution String,
ADD COLUMN IF NOT EXISTS language String;

-- Индексы для аналитических запросов
ALTER TABLE analytics.visitor_events ADD INDEX idx_page_url page_url TYPE bloom_filter GRANULARITY 1;
ALTER TABLE analytics.visitor_events ADD INDEX idx_utm_source utm_source TYPE bloom_filter GRANULARITY 1;
ALTER TABLE analytics.visitor_events ADD INDEX idx_utm_campaign utm_campaign TYPE bloom_filter GRANULARITY 1;