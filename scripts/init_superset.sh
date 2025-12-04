#!/bin/bash

set -e

# =====================================================================
#                          CONFIGURATION
# =====================================================================

# Default values (will be overridden by .env)
DEFAULT_USERNAME="admin"
DEFAULT_PASSWORD="admin123"
DEFAULT_EMAIL="admin@edu-platform.com"
DEFAULT_LOAD_EXAMPLES="false"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# =====================================================================
#                          UTILITY FUNCTIONS
# =====================================================================

log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

load_environment() {
    local env_file=".env"
    
    if [ -n "$1" ]; then
        env_file=".env.$1"
    fi
    
    if [ ! -f "$env_file" ]; then
        error "Environment file $env_file not found"
    fi
    
    log "Loading environment from $env_file"
    
    # Load environment variables
    set -a
    source "$env_file"
    set +a
    
    # Set defaults if not provided in env
    SUPERSET_USERNAME="${SUPERSET_USERNAME:-$DEFAULT_USERNAME}"
    SUPERSET_PASSWORD="${SUPERSET_PASSWORD:-$DEFAULT_PASSWORD}"
    SUPERSET_EMAIL="${SUPERSET_EMAIL:-$DEFAULT_EMAIL}"
    SUPERSET_LOAD_EXAMPLES="${SUPERSET_LOAD_EXAMPLES:-$DEFAULT_LOAD_EXAMPLES}"
    
    # Use APP_ENV from environment or parameter
    if [ -n "$1" ]; then
        APP_ENV="$1"
    fi
    
    log "Environment: ${APP_ENV:-dev}"
    log "Superset username: $SUPERSET_USERNAME"
    log "Superset email: $SUPERSET_EMAIL"
    log "Load examples: $SUPERSET_LOAD_EXAMPLES"
}

check_docker_compose() {
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        error "Docker Compose is not available"
    fi
}

get_docker_compose_cmd() {
    if command -v docker-compose &> /dev/null; then
        echo "docker-compose"
    else
        echo "docker compose"
    fi
}

# =====================================================================
#                         SUPERSET FUNCTIONS
# =====================================================================

wait_for_superset() {
    local docker_cmd=$(get_docker_compose_cmd)
    local max_attempts=30
    local attempt=1
    
    log "Waiting for Superset to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if $docker_cmd exec superset curl -f http://localhost:8088/health > /dev/null 2>&1; then
            success "Superset is ready!"
            return 0
        fi
        
        log "Attempt $attempt/$max_attempts - Superset not ready yet..."
        sleep 10
        ((attempt++))
    done
    
    error "Superset failed to start within 5 minutes"
}

wait_for_clickhouse() {
    local docker_cmd=$(get_docker_compose_cmd)
    local max_attempts=20
    local attempt=1
    
    log "Waiting for ClickHouse to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if $docker_cmd exec clickhouse clickhouse-client --query "SELECT 1" > /dev/null 2>&1; then
            success "ClickHouse is ready!"
            return 0
        fi
        
        log "Attempt $attempt/$max_attempts - ClickHouse not ready yet..."
        sleep 5
        ((attempt++))
    done
    
    warning "ClickHouse might not be fully ready, but continuing..."
}

init_superset_db() {
    local docker_cmd=$(get_docker_compose_cmd)
    
    log "Initializing Superset database..."
    
    if $docker_cmd exec superset superset db upgrade; then
        success "Database upgraded successfully"
    else
        error "Failed to upgrade database"
    fi
}

create_admin_user() {
    local docker_cmd=$(get_docker_compose_cmd)
    local username="$1"
    local password="$2"
    local email="$3"
    
    log "Creating admin user: $username"
    
    # Check if admin user already exists
    if $docker_cmd exec superset superset fab list-users | grep -q "$username"; then
        warning "Admin user '$username' already exists"
        return 0
    fi
    
    if $docker_cmd exec superset superset fab create-admin \
        --username "$username" \
        --firstname Admin \
        --lastname User \
        --email "$email" \
        --password "$password"; then
        success "Admin user created successfully"
    else
        error "Failed to create admin user"
    fi
}

init_superset() {
    local docker_cmd=$(get_docker_compose_cmd)
    
    log "Initializing Superset with roles and permissions..."
    
    if $docker_cmd exec superset superset init; then
        success "Superset initialized successfully"
    else
        error "Failed to initialize Superset"
    fi
}

create_clickhouse_database() {
    local docker_cmd=$(get_docker_compose_cmd)
    
    log "Creating ClickHouse database connection in Superset..."
    
    # Wait a bit for Superset to be fully ready
    sleep 10
    
    # Build connection string from environment variables - use HTTP port 8123, no extra params
    local clickhouse_uri="clickhousedb://${CLICKHOUSE_USER:-default}:${CLICKHOUSE_PASSWORD}@${CLICKHOUSE_HOST:-clickhouse}:8123/${CLICKHOUSE_DB:-analytics}"
    
    log "ClickHouse connection string: $clickhouse_uri"
    
    # Try to create database connection via CLI (this might not work in all versions)
    # Fallback to manual instructions
    if $docker_cmd exec superset superset set-database-uri --database_name "ClickHouse Analytics" --uri "$clickhouse_uri" 2>/dev/null; then
        success "ClickHouse database connection created"
    else
        warning "Could not create ClickHouse connection automatically"
        log "Please create database connection manually in Superset UI:"
        log "URL: http://localhost:8088/databaseview/list/"
        log "Connection string: $clickhouse_uri"
    fi
}

load_examples() {
    local docker_cmd=$(get_docker_compose_cmd)
    local load_examples="$1"
    
    if [ "$load_examples" = "true" ]; then
        log "Loading example datasets and dashboards..."
        if $docker_cmd exec superset superset load_examples; then
            success "Examples loaded successfully"
        else
            warning "Failed to load examples (this might be normal in some setups)"
        fi
    fi
}

setup_clickhouse_tables() {
    local docker_cmd=$(get_docker_compose_cmd)
    
    log "Setting up ClickHouse tables for analytics..."
    
    # Create useful views in ClickHouse for Superset
    $docker_cmd exec clickhouse clickhouse-client --query "
        CREATE DATABASE IF NOT EXISTS ${CLICKHOUSE_DB:-analytics};
        
        -- Daily page views summary
        CREATE VIEW IF NOT EXISTS ${CLICKHOUSE_DB:-analytics}.daily_page_views AS
        SELECT 
            toDate(created_at) as date,
            page_url,
            count(*) as views,
            count(distinct visitor_id) as unique_visitors
        FROM ${CLICKHOUSE_DB:-analytics}.visitor_events 
        WHERE event_type = 'page_view'
        GROUP BY date, page_url;
        
        -- UTM campaign performance
        CREATE VIEW IF NOT EXISTS ${CLICKHOUSE_DB:-analytics}.utm_performance AS
        SELECT 
            utm_source,
            utm_medium,
            utm_campaign,
            count(*) as total_events,
            count(distinct visitor_id) as unique_visitors,
            countIf(event_type = 'form_submit') as conversions
        FROM ${CLICKHOUSE_DB:-analytics}.visitor_events 
        WHERE utm_source != ''
        GROUP BY utm_source, utm_medium, utm_campaign;
        
        -- User activity summary
        CREATE VIEW IF NOT EXISTS ${CLICKHOUSE_DB:-analytics}.user_activity_summary AS
        SELECT 
            visitor_id,
            min(created_at) as first_seen,
            max(created_at) as last_seen,
            count(*) as total_events,
            count(distinct event_type) as unique_event_types,
            groupUniqArray(event_type) as event_types
        FROM ${CLICKHOUSE_DB:-analytics}.visitor_events 
        GROUP BY visitor_id;
    " && success "ClickHouse views created" || warning "Some ClickHouse views might not have been created"
}

# =====================================================================
#                             MAIN FLOW
# =====================================================================

main() {
    local env_param="$1"
    local username_param="$2"
    local password_param="$3"
    local email_param="$4"
    local load_examples_param="$5"
    
    log "Starting Superset initialization..."
    
    # Load environment first
    load_environment "$env_param"
    
    # Override with command line parameters if provided
    SUPERSET_USERNAME="${username_param:-$SUPERSET_USERNAME}"
    SUPERSET_PASSWORD="${password_param:-$SUPERSET_PASSWORD}"
    SUPERSET_EMAIL="${email_param:-$SUPERSET_EMAIL}"
    SUPERSET_LOAD_EXAMPLES="${load_examples_param:-$SUPERSET_LOAD_EXAMPLES}"
    
    check_docker_compose
    
    # Wait for services
    wait_for_clickhouse
    wait_for_superset
    
    # Initialize Superset
    init_superset_db
    create_admin_user "$SUPERSET_USERNAME" "$SUPERSET_PASSWORD" "$SUPERSET_EMAIL"
    init_superset
    load_examples "$SUPERSET_LOAD_EXAMPLES"
    
    # Setup ClickHouse
    setup_clickhouse_tables
    create_clickhouse_database
    
    success "🎉 Superset initialization completed!"
    echo ""
    echo "Access Superset at: http://localhost:8088"
    echo "Username: $SUPERSET_USERNAME"
    echo "Username: $SUPERSET_PASSWORD"
    echo "Environment: ${APP_ENV:-dev}"
    echo ""
    echo "ClickHouse connection:"
    echo "  Host: ${CLICKHOUSE_HOST:-clickhouse}"
    echo "  Database: ${CLICKHOUSE_DB:-analytics}"
    echo ""
    echo "Next steps:"
    echo "  1. Verify ClickHouse database connection in Superset UI"
    echo "  2. Create datasets from the available views"
    echo "  3. Build dashboards for your analytics"
}

# =====================================================================
#                         USAGE AND HELP
# =====================================================================

show_usage() {
    echo "Usage: $0 [environment] [username] [password] [email] [load_examples]"
    echo ""
    echo "Parameters:"
    echo "  environment    : dev, stage, prod (default: current .env)"
    echo "  username       : Admin username (default: from .env or 'admin')"
    echo "  password       : Admin password (default: from .env or 'admin123')"
    echo "  email          : Admin email (default: from .env or 'admin@edu-platform.com')"
    echo "  load_examples  : true/false to load example data (default: false)"
    echo ""
    echo "Examples:"
    echo "  $0                          # Use current .env"
    echo "  $0 dev                      # Use .env.dev"
    echo "  $0 prod admin pass123 admin@company.com true"
    echo ""
    echo "Environment variables (in .env files):"
    echo "  SUPERSET_USERNAME          # Admin username"
    echo "  SUPERSET_PASSWORD          # Admin password" 
    echo "  SUPERSET_EMAIL             # Admin email"
    echo "  SUPERSET_LOAD_EXAMPLES     # Load examples (true/false)"
    echo "  CLICKHOUSE_*               # ClickHouse connection settings"
}

# Parse command line arguments
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    show_usage
    exit 0
fi

# Run main function with all arguments
main "$@"