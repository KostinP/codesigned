#!/bin/bash

set -e

# Load environment
source ./scripts/env.sh .env.prod

check_required_vars VPS_IP VPS_USER SSL_CERT_FILE SSL_KEY_FILE SSL_CERT_PATH SSL_KEY_PATH DOMAIN

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
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
    if [ -n "$VPS_SSH_KEY_PATH" ] && [ -f "$VPS_SSH_KEY_PATH" ]; then
        scp -o StrictHostKeyChecking=no -i "$VPS_SSH_KEY_PATH" "$1" "${VPS_USER}@${VPS_IP}:$2"
    elif [ -n "$VPS_PASSWORD" ]; then
        sshpass -p "$VPS_PASSWORD" scp -o StrictHostKeyChecking=no "$1" "${VPS_USER}@${VPS_IP}:$2"
    else
        scp -o StrictHostKeyChecking=no "$1" "${VPS_USER}@${VPS_IP}:$2"
    fi
}

main() {
    log "Uploading SSL certificates to VPS..."
    
    # Ensure SSL directory exists
    ssh_vps "mkdir -p /etc/nginx/ssl"

    # Check if we have certificate files
    if [ ! -f "$SSL_CERT_FILE" ] || [ ! -f "$SSL_KEY_FILE" ]; then
        error "Certificate files not found: $SSL_CERT_FILE or $SSL_KEY_FILE"
        echo "Please ensure SSL certificate files exist in the project root"
        exit 1
    fi

    log "Found certificate files:"
    echo "  Certificate: $SSL_CERT_FILE"
    echo "  Key: $SSL_KEY_FILE"
    echo "  Target VPS: $VPS_IP"
    echo "  Target paths:"
    echo "    $SSL_CERT_PATH"
    echo "    $SSL_KEY_PATH"

    # Create certificate chain if intermediate is available
    if [ -f "$SSL_CERT_FILE" ]; then
        log "Checking if certificate chain is needed..."
        
        # Test certificate
        if ! openssl verify -CAfile "$SSL_CERT_FILE" "$SSL_CERT_FILE" > /dev/null 2>&1; then
            warning "Certificate chain incomplete, downloading intermediate..."
            
            # Download intermediate certificate
            INTERMEDIATE_CERT="intermediate.crt"
            if [ ! -f "$INTERMEDIATE_CERT" ]; then
                curl -s -o "$INTERMEDIATE_CERT" http://secure.globalsign.com/cacert/gsgccr3dvtlsca2020.crt
            fi
            
            # Create chain
            CHAIN_CERT="prod-chain.crt"
            cat "$SSL_CERT_FILE" "$INTERMEDIATE_CERT" > "$CHAIN_CERT"
            SSL_CERT_FILE="$CHAIN_CERT"
            log "Using certificate chain: $CHAIN_CERT"
        fi
    fi

    # Upload certificate files
    log "Uploading certificate files..."
    scp_to_vps "$SSL_CERT_FILE" "$SSL_CERT_PATH"
    scp_to_vps "$SSL_KEY_FILE" "$SSL_KEY_PATH"

    # Set proper permissions
    log "Setting up certificates on VPS..."
    ssh_vps "
        chmod 644 $SSL_CERT_PATH
        chmod 600 $SSL_KEY_PATH
        echo 'Certificate files installed:'
        ls -la $SSL_CERT_PATH
        ls -la $SSL_KEY_PATH
    "

    # Create Nginx configuration
    log "Configuring Nginx..."
    ssh_vps "
        cat > /etc/nginx/sites-available/default << 'NGINX_CONFIG'
server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;
    return 301 https://\$server_name\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name $DOMAIN www.$DOMAIN;

    ssl_certificate $SSL_CERT_PATH;
    ssl_certificate_key $SSL_KEY_PATH;
    
    # SSL security settings
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    location /api/telegram/webhook {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # Telegram webhook specific settings
        proxy_read_timeout 90;
        proxy_connect_timeout 90;
    }
}
NGINX_CONFIG
    "

    # Test Nginx configuration
    log "Testing Nginx configuration..."
    ssh_vps "nginx -t"

    # Reload Nginx
    log "Reloading Nginx..."
    ssh_vps "systemctl reload nginx"

    log "✅ SSL certificates uploaded and configured successfully!"
}

main "$@"