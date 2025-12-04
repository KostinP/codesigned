# Edu Platform Backend

Бэкенд для образовательной платформы с интеграцией Telegram бота. Полностью автоматизированное развертывание и управление через Makefile.

## 📋 Оглавление
1. [Быстрый старт](#быстрый-старт)
2. [Автоматизация через Makefile](#автоматизация-через-makefile)
3. [Локальная разработка](#локальная-разработка)
4. [🐳 Docker & Контейнеризация](#-docker--контейнеризация)
5. [Деплой на VPS](#деплой-на-vps)
6. [Управление окружениями](#управление-окружениями)
7. [База данных и миграции](#база-данных-и-миграции)
8. [Мониторинг и логи](#мониторинг-и-логи)
9. [Устранение неисправностей](#устранение-неисправностей)

## 🚀 Быстрый старт

### Вариант 1: Локальная разработка с Docker (рекомендуется)

```bash
# 1. Клонировать проект
git clone https://github.com/KostinP/edu-platform-backend.git
cd edu-platform-backend

# 2. Сделать скрипты исполняемыми
chmod +x scripts/*.sh

# 3. Инициализировать окружения
make init-env

echo "ENABLE_ANALYTICS=false" >> .env.dev
# Включить аналитику позже
echo "ENABLE_ANALYTICS=true" >> .env.dev

# 4. Запустить всё в Docker
make docker-fresh env=dev

# 5. Приложение доступно на http://localhost:8080
```

### Вариант 2: Полное развертывание production сервера

```bash
# 1-3. Те же шаги что выше
# 4. Настроить файлы окружения
nano .env.prod

# 5. Запустить полный деплой
make prod-full-deploy
```

## 🤖 Автоматизация через Makefile

### Управление окружениями

```bash
# Переключение между окружениями
make env-switch env=dev    # разработка
make env-switch env=stage  # staging
make env-switch env=prod   # продакшен

# Показать текущее окружение
make env

# Проверить требуемые переменные окружения
make check-env env=prod

# Инициализировать все файлы окружения
make init-env
```

### Разработка

```bash
# Запуск с hot reload
make run env=dev

# Генерация Swagger документации
make swagger

# Генерация зависимостей Wire
make wire

# Тестирование
make test env=dev

# Сборка приложения
make build env=dev

# Очистка зависимостей
make tidy
```

### База данных

```bash
# Миграции
make migrate-up env=dev         # Применить миграции
make migrate-down env=dev       # Откатить миграции
make migrate-create name=users  # Создать миграцию
make migrate-status env=dev     # Статус миграций
```

### Деплой и управление VPS

```bash
# Полный пайплайн для продакшена
make prod-full-deploy

# Отдельные этапы
make setup-vps env=prod         # Установить ПО на VPS
make deploy env=prod            # Деплой приложения
make upload-certs env=prod      # Загрузить SSL сертификаты
make setup-webhook env=prod     # Настроить Telegram вебхук

# Мониторинг
make vps-logs env=prod          # Логи приложения
make vps-info env=prod          # Информация о системе
```

## 💻 Локальная разработка

### 1. Настройка окружения

```bash
# Создание файлов окружения
make init-env

# Переключение в режим разработки
make env-switch env=dev

# Редактирование настроек
nano .env.dev
```

### 2. Запуск зависимостей

```bash
# Запуск PostgreSQL и ClickHouse
docker compose up -d postgres clickhouse
```

### 3. Запуск приложения

```bash
# Режим разработки с hot reload (air)
make run env=dev

# Или обычный запуск
go run ./cmd
```

### 4. Работа с базой данных

```bash
# Применение миграций
make migrate-up env=dev

# Создание новой миграции
make migrate-create name=add_feature

# Подключение к БД
docker compose exec postgres psql -U postgres -d edu_platform
```

### 5. Генерация кода

```bash
# Swagger документация
make swagger

# Dependency injection
make wire

# Проверка зависимостей
make tidy
```

## 🐳 Docker & Контейнеризация

### Запуск с Docker Compose

```bash
# Полная очистка и перезапуск (рекомендуется для разработки)
make docker-fresh env=dev

# Или пошагово
make docker-build env=dev    # Сборка образа
make docker-up env=dev       # Запуск всех сервисов

# Запуск только базовых сервисов (без аналитики)
make docker-up-core env=dev

# Запуск с аналитикой (ClickHouse + Superset)
make docker-up-with-analytics env=dev

# Остановка контейнеров
make docker-down

# Остановка с очисткой volumes
make docker-clean

# Просмотр логов
make docker-logs

# Статус контейнеров
make docker-ps

# Перезапуск сервисов
make docker-restart env=dev
```

### Структура контейнеров

Сервисы, запускаемые через Docker Compose:

| Сервис | Порт | Назначение | Профиль |
|--------|------|------------||--------|
| **app** | 8080 | Основное приложение (Go) | всегда |
| **postgres** | 5432 | Основная база данных | всегда |
| **clickhouse** | 8123/9000 | Аналитическая база данных | analytics |
| **superset** | 8088 | BI-панель и дашборды | analytics |
| **superset-db** | 5433 | База данных Superset | analytics |

### Управление аналитическими сервисами

Чтобы отключить аналитику и сэкономить ресурсы, установите в .env.dev:

```env
ENABLE_ANALYTICS=false
```

### Экономия ресурсов при отключенной аналитике
При ENABLE_ANALYTICS=false:
- Память: экономия ~2-4GB (ClickHouse + Superset)
- Диск: экономия ~1-2GB (данные аналитики)
- Время запуска: минуты быстрее
- API аналитики: возвращает 503 Service Unavailable

### Переменные окружения для Docker

Убедитесь, что в `.env.dev` установлены:

```env
# Database Configuration
DB_HOST=postgres           # Имя сервиса в Docker сети
DB_PORT=5432
DB_NAME=edu_platform
DB_USER=postgres
DB_PASSWORD=dev_secure_db_password_123

# ClickHouse Configuration  
CLICKHOUSE_HOST=clickhouse # Имя сервиса в Docker сети
CLICKHOUSE_PORT=9000
CLICKHOUSE_DB=analytics
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=clickhouse_password
```

### Решение частых проблем с Docker

**Проблема: Неправильный пароль PostgreSQL**
```bash
# Полная очистка и перезапуск
make docker-clean
APP_ENV=dev docker-compose up --build
```

**Проблема: Конфликт портов**
```bash
# Проверка занятых портов
lsof -i :8080
lsof -i :5432
lsof -i :8088

# Остановка конфликтующих процессов
kill -9 <PID>
```

**Проблема: Ошибки сети Docker**
```bash
# Сброс Docker сетей
docker network prune

# Перезапуск Docker Desktop
```

### Health Checks

Все сервисы имеют health checks:
- ✅ **PostgreSQL**: `pg_isready`
- ✅ **ClickHouse**: HTTP ping на порт 8123  
- ✅ **App**: HTTP запрос на `/health`

### Полезные команды для отладки

```bash
# Подключиться к контейнеру
docker-compose exec app sh
docker-compose exec postgres psql -U postgres -d edu_platform

# Просмотр логов конкретного сервиса
docker-compose logs app -f
docker-compose logs postgres -f

# Проверка переменных окружения в контейнере
docker-compose exec app env | grep DB_

# Мониторинг ресурсов
docker stats
```

### Production vs Development

| Аспект | Development | Production |
|--------|-------------|------------|
| **Пароли БД** | Простые (dev_*) | Сложные случайные |
| **SSL** | Отключен | Включен с LetsEncrypt |
| **Логи** | Подробные (debug) | Только ошибки (warn) |
| **Порты** | Локальные (localhost) | Доменные имена |
| **Миграции** | Автоматические | Ручное управление |

## 🎯 Деплой на VPS

### Подготовка VPS

```bash
# Автоматическая установка (Docker, Nginx, SSL, firewall)
make setup-vps env=prod
```

**Что установит `make setup-vps`:**
- 🐳 Docker & Docker Compose
- 🌐 Nginx + SSL конфигурация
- 🔥 UFW firewall
- 🛠️ Системные утилиты
- 💾 Swap file (2GB)
- 📊 Мониторинг
- 🔐 Безопасность (fail2ban)

### Настройка окружения продакшена

Отредактируйте `.env.prod`:

```env
# VPS Configuration
VPS_IP=your_vps_ip
VPS_USER=root
VPS_PASSWORD=your_vps_password
VPS_SSH_KEY_PATH=~/.ssh/id_rsa

# Application Configuration
APP_ENV=prod
CONFIG_PATH=configs/prod.yaml
LOG_LEVEL=debug

# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_NAME=edu_platform
DB_USER=postgres
DB_PASSWORD=your_very_secure_database_password

# ClickHouse Configuration
CLICKHOUSE_HOST=clickhouse
CLICKHOUSE_PORT=9000
CLICKHOUSE_DB=analytics
CLICKHOUSE_USER=default
CLICKHOUSE_PASSWORD=

# Telegram Configuration
TELEGRAM_BOT_TOKEN=your_telegram_bot_token
TELEGRAM_WEBHOOK_URL=https://your-domain.com/api/telegram/webhook

# JWT Configuration
JWT_SECRET=your_very_secure_jwt_secret

# Superset Configuration
SUPERSET_SECRET_KEY=yoyr_superset_secret_key

# SSL Configuration
SSL_CERT_FILE=example.crt
SSL_KEY_FILE=example.key
SSL_CERT_PATH=/etc/nginx/ssl/server.crt
SSL_KEY_PATH=/etc/nginx/ssl/server.key

# Domain
DOMAIN=your-domain.com

# REPO
REPOLINK=https://github.com/KostinP/edu-platform-backend
```

### Процесс деплоя

```bash
# Полный пайплайн (рекомендуется)
make prod-full-deploy

# Или пошагово
make setup-vps env=prod         # Установка ПО
make upload-certs env=prod      # SSL сертификаты
make deploy env=prod            # Деплой приложения
make setup-webhook env=prod     # Telegram вебхук
```

## 🔄 Управление окружениями

### Структура файлов

```
.env          # Копия активного окружения
.env.dev      # Разработка
.env.stage    # Staging  
.env.prod     # Продакшен
```

### Переключение окружений

```bash
# Показать текущее окружение
make env

# Переключиться
make env-switch env=dev    # -> .env.dev
make env-switch env=stage  # -> .env.stage  
make env-switch env=prod   # -> .env.prod
```

### Переменные окружения

| Переменная | Обязательная | Описание |
|------------|--------------|----------|
| `VPS_IP` | ✅ | IP адрес VPS |
| `VPS_USER` | ✅ | Пользователь VPS |
| `TELEGRAM_BOT_TOKEN` | ✅ | Токен Telegram бота |
| `DB_HOST` | ✅ | Хост PostgreSQL |
| `DB_PORT` | ✅ | Порт PostgreSQL |
| `DB_NAME` | ✅ | Имя базы данных |
| `DB_USER` | ✅ | Пользователь PostgreSQL |
| `DB_PASSWORD` | ✅ | Пароль PostgreSQL |
| `JWT_SECRET` | ✅ | Секрет для JWT |
| `DOMAIN` | ✅ | Домен приложения |
| `REPOLINK` | ✅ | Ссылка на репозиторий |
| `SUPERSET_SECRET_KEY` | ✅ | Секрет для Superset |
| `VPS_SSH_KEY_PATH` | ❌ | Путь к SSH ключу |
| `VPS_PASSWORD` | ❌ | Пароль VPS (если нет ключа) |
| `CLICKHOUSE_HOST` | ❌ | Хост ClickHouse |
| `CLICKHOUSE_PORT` | ❌ | Порт ClickHouse |
| `CLICKHOUSE_DB` | ❌ | Имя базы ClickHouse |
| `CLICKHOUSE_USER` | ❌ | Пользователь ClickHouse |
| `CLICKHOUSE_PASSWORD` | ❌ | Пароль ClickHouse |

## 🗄️ База данных и миграции

### Управление миграциями

```bash
# Создание миграции
make migrate-create name=create_users_table

# Применение миграций
make migrate-up env=prod

# Откат последней миграции  
make migrate-down env=prod

# Просмотр статуса
make migrate-status env=prod
```

### Бэкапы

```bash
# Создание бэкапа
ssh $(VPS_USER)@$(VPS_IP) "docker compose exec postgres pg_dump -U $(DB_USER) $(DB_NAME) > backup.sql"

# Восстановление
ssh $(VPS_USER)@$(VPS_IP) "docker compose exec -T postgres psql -U $(DB_USER) $(DB_NAME) < backup.sql"
```

## 📊 Мониторинг и логи

### Просмотр логов

```bash
# Логи приложения в реальном времени
make vps-logs env=prod

# Логи с фильтрацией
make vps-logs env=prod | grep ERROR
```

### Мониторинг системы

```bash
# Общая информация о системе
make vps-info env=prod

# Статус Docker контейнеров
ssh $(VPS_USER)@$(VPS_IP) "docker ps"

# Использование ресурсов
ssh $(VPS_USER)@$(VPS_IP) "docker stats"
```

### Health checks

Все сервисы имеют health checks:
- PostgreSQL: pg_isready
- ClickHouse: HTTP ping на порт 8123
- App: HTTP запрос на /health
- Superset: HTTP запрос на /health

```bash
# Проверка здоровья приложения
curl https://$(DOMAIN)/health

# Проверка Telegram вебхука
curl -X POST https://$(DOMAIN)/api/telegram/webhook
```

## 🛠️ Устранение неисправностей

### Быстрая диагностика

```bash
# Проверка состояния системы
make vps-info env=prod

# Просмотр логов
make vps-logs env=prod
```

### Частые проблемы

**Приложение не запускается:**
```bash
make vps-logs env=prod
make deploy env=prod  # перезапуск
```

**Проблемы с базой данных:**
```bash
# Проверка подключения
ssh $(VPS_USER)@$(VPS_IP) "docker compose exec postgres psql -U $(DB_USER) -d $(DB_NAME)"

# Запуск миграций
make migrate-up env=prod
```

**Проблемы с SSL:**
```bash
# Перезагрузка Nginx
ssh $(VPS_USER)@$(VPS_IP) "nginx -t && systemctl reload nginx"
```

**Проблемы с вебхуком:**
```bash
# Перенастройка вебхука
make setup-webhook env=prod

# Тестирование endpoint
curl -X POST https://$(DOMAIN)/api/telegram/webhook
```

### Проблемы с Docker

**Контейнеры не запускаются:**
```bash
# Проверка сборки
make docker-build env=dev

# Проверка конфигурации
APP_ENV=dev docker-compose config

# Принудительная пересборка
docker-compose build --no-cache
```

**Ошибки подключения к БД в Docker:**
```bash
# Проверка что контейнеры в одной сети
docker network ls
docker network inspect edu-platform-backend_edu-platform

# Проверка доступности сервисов
docker-compose exec app nslookup postgres
docker-compose exec app nslookup clickhouse
```

**Volume проблемы:**
```bash
# Очистка volumes
docker-compose down -v
docker volume prune

# Проверка volumes
docker volume ls
```

### Полезные команды для отладки

```bash
# Проверка сетевых соединений
ssh $(VPS_USER)@$(VPS_IP) "netstat -tulpn"

# Проверка DNS
nslookup $(DOMAIN)

# Проверка SSL сертификата
openssl s_client -connect $(DOMAIN):443
```

## 🔐 Безопасность

### Рекомендации

1. **Используйте SSH ключи** вместо паролей
2. **Сложные пароли** для баз данных
3. **Регулярно обновляйте** зависимости
4. **Мониторьте логи** на предмет подозрительной активности
5. **Настройте бэкапы**

### Обновление секретов

```bash
# Генерация нового JWT секрета
openssl rand -base64 32

# Обновление в .env.prod
nano .env.prod

# Перезапуск приложения
make deploy env=prod
```

## 📞 Поддержка

### Получение логов для отладки

```bash
# Полные логи приложения
make vps-logs env=prod

# Логи за определенный период
ssh $(VPS_USER)@$(VPS_IP) "docker compose logs --since=1h app"

# Логи Nginx
ssh $(VPS_USER)@$(VPS_IP) "tail -f /var/log/nginx/access.log"
```

### Проверка работоспособности

```bash
# Health check
curl https://$(DOMAIN)/health

# Проверка использования диска
ssh $(VPS_USER)@$(VPS_IP) "df -h"
```

## 📞 Superset

### Использование Superset

```bash
# Инициализация для текущего окружения
make superset-init env=dev

# Полная пересборка с инициализацией Superset
make docker-fresh env=dev
```


---

## 🎯 Чеклист развертывания

- [ ] Клонировать репозиторий
- [ ] `chmod +x scripts/*.sh`
- [ ] `make init-env`
- [ ] Настроить `.env.prod`
- [ ] Положить SSL сертификаты
- [ ] `make prod-full-deploy`
- [ ] Проверить `https://$(DOMAIN)/health`
- [ ] `make setup-webhook env=prod`
- [ ] Протестировать функциональность

## 💡 Все команды Makefile

```bash
# Показать все доступные команды
make help

# Или просто
make
```

## Все файлы через скрипт на python
python3 combined.py --config combined.ini

**Примечание:** Полный список всех команд с описанием можно посмотреть выполнив `make help`.

---

**Примечание по Docker:** Для успешного запуска в Docker обязательно используйте команды с явным указанием `APP_ENV` или `make docker-fresh env=dev`, так как Docker Compose не автоматически подставляет переменные из `.env` файла в YAML конфигурацию.