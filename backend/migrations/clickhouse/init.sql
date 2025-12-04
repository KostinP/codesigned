-- migrations/clickhouse/init.sql
-- Initialize ClickHouse database for analytics
-- Создаем базу данных, если не существует
CREATE DATABASE IF NOT EXISTS analytics;

-- Переключаемся на созданную базу
USE analytics;

-- Таблицы будут созданы через отдельные миграции в /docker-entrypoint-initdb.d/