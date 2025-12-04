#!/bin/bash

set -e

# Load environment
source ./scripts/env.sh .env.prod

check_required_vars() {
    for var in "$@"; do
        if [ -z "${!var}" ]; then
            error "Missing required variable: $var"
        fi
    done
}

check_required_vars VPS_IP VPS_USER TELEGRAM_BOT_TOKEN DOMAIN REPOLINK SUPERSET_SECRET_KEY

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_DIR="/opt/edu-platform"

log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

ssh_vps() {
    if [ -n "$VPS_SSH_KEY_PATH" ] && [ -f "$VPS_SSH_KEY_PATH" ]; then
        ssh -o StrictHostKeyChecking=no -i "$VPS_SSH_KEY_PATH" ${VPS_USER}@${VPS_IP} "$@"
    elif [ -n "$VPS_PASSWORD" ]; then
        sshpass -p "$VPS_PASSWORD" ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_IP} "$@"
    else
        ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_IP} "$@"
    fi
}

scp_to_vps() {
    local local_file=$1
    local remote_path=$2
    
    if [ -n "$VPS_SSH_KEY_PATH" ] && [ -f "$VPS_SSH_KEY_PATH" ]; then
        scp -o StrictHostKeyChecking=no -i "$VPS_SSH_KEY_PATH" "$local_file" "${VPS_USER}@${VPS_IP}:${remote_path}"
    elif [ -n "$VPS_PASSWORD" ]; then
        sshpass -p "$VPS_PASSWORD" scp -o StrictHostKeyChecking=no "$local_file" "${VPS_USER}@${VPS_IP}:${remote_path}"
    else
        scp -o StrictHostKeyChecking=no "$local_file" "${VPS_USER}@${VPS_IP}:${remote_path}"
    fi
}

check_docker() {
    log "Checking Docker..."
    ssh_vps "docker ps" || error "Docker not working"
}

setup_project_structure() {
    log "Setting up project structure..."
    ssh_vps "mkdir -p $PROJECT_DIR $PROJECT_DIR/logs $PROJECT_DIR/data"
}

clone_repository() {
    log "Cloning repository..."
    ssh_vps "
        cd /opt
        rm -rf $PROJECT_DIR
        git clone $REPOLINK edu-platform
    "
}

update_db_password() {
    log "Updating database password if changed..."
    
    ssh_vps "
        cd $PROJECT_DIR
        set -a
        source .env
        set +a
        
        echo '=== Starting PostgreSQL for password update ==='
        docker compose up -d postgres
        
        echo '=== Waiting for PostgreSQL to be ready ==='
        sleep 10
        
        # Проверяем готовность PostgreSQL
        until docker compose exec -T postgres pg_isready -U postgres; do
            echo 'Waiting for PostgreSQL...'
            sleep 5
        done
        
        # Проверяем, изменился ли пароль
        OLD_PASSWORD=\$(docker compose exec -T postgres psql -U postgres -t -c \"SELECT passwd FROM pg_shadow WHERE usename='postgres';\" 2>/dev/null | tr -d ' ' | head -1 || echo '')
        NEW_PASSWORD_HASH=\$(echo -n \"\${DB_PASSWORD}\${DB_USER}\" | md5sum | cut -d' ' -f1)
        
        echo \"Old password hash: \$OLD_PASSWORD\"
        echo \"New password hash: md5\$NEW_PASSWORD_HASH\"
        
        if [ \"\$OLD_PASSWORD\" != \"md5\$NEW_PASSWORD_HASH\" ]; then
            echo 'Password changed, updating...'
            docker compose exec -T postgres psql -U postgres -c \"ALTER USER postgres WITH PASSWORD '\$DB_PASSWORD';\"
            echo '✅ Database password updated'
        else
            echo '✅ Password unchanged'
        fi
    " || warning "Password update may have failed, but continuing deployment..."
}

upload_env_files() {
    log "Uploading environment files..."
    
    # Создаем директорию configs если её нет
    ssh_vps "mkdir -p $PROJECT_DIR/configs"
    
    # Копируем .env.prod для docker-compose
    scp_to_vps ".env.prod" "$PROJECT_DIR/.env.prod"
    
    # Также копируем как .env для backup
    scp_to_vps ".env.prod" "$PROJECT_DIR/.env"
    
    # Копируем конфигурационные файлы
    scp_to_vps "configs/prod.yaml" "$PROJECT_DIR/configs/prod.yaml"
    
    # Проверяем что файлы загружены
    ssh_vps "cd $PROJECT_DIR && ls -la .env* && ls -la configs/"
}

deploy_services() {
    log "Deploying services..."

    # Сначала загружаем env файлы
    upload_env_files
    
    ssh_vps "
        cd $PROJECT_DIR
        set -a
        source .env
        set +a
        
        echo '=== Stopping old services ==='
        docker compose down --remove-orphans || true
        
        echo '=== Building services ==='
        docker compose build
        
        echo '=== Starting services ==='
        docker compose up -d
        
        echo '=== Waiting for services to start ==='
        sleep 30
        
        echo '=== Checking status ==='
        docker compose ps
    "
}

check_environment_files() {
    log "Checking environment files..."
    
    if [ ! -f ".env.prod" ]; then
        error "File .env.prod not found in current directory"
    fi
    
    if [ ! -f "configs/prod.yaml" ]; then
        error "File configs/prod.yaml not found"
    fi
    
    success "Environment files verified"
}

check_status() {
    log "Checking deployment status..."
    
    ssh_vps "
        cd $PROJECT_DIR
        . .env  # Source again for checks
        
        echo '=== Container status ==='
        docker compose ps
        
        echo ''
        echo '=== App logs ==='
        docker compose logs app --tail=20
        
        echo ''
        echo '=== Testing app connectivity ==='
        echo 'Waiting 10 seconds for app to start...'
        sleep 10
        curl -s -m 5 http://localhost:8080/health && echo '✅ App is healthy' || echo '❌ App not responding'
        
        echo ''
        echo '=== Testing database connectivity ==='
        docker compose exec -T postgres pg_isready -U \"\${DB_USER}\" -d \"\${DB_NAME}\" && echo '✅ PostgreSQL is ready' || echo '❌ PostgreSQL not ready'
    "
}

main() {
    log "Starting deployment..."
    
    check_environment_files
    check_docker
    setup_project_structure
    clone_repository
    upload_env_files
    update_db_password
    deploy_services
    check_status
    
    log "Deployment completed!"
    echo ""
    echo "Application should be available at: https://$DOMAIN"
    echo "Check logs with: make vps-logs env=prod"
}

main "$@"