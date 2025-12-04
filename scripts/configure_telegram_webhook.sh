#!/bin/bash

# Configure Telegram webhook with .env support

set -e

# Load environment
source ./scripts/env.sh .env.prod

check_required_vars TELEGRAM_BOT_TOKEN DOMAIN

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

configure_webhook() {
    local token="$TELEGRAM_BOT_TOKEN"
    local webhook_url="https://${DOMAIN}/api/telegram/webhook"
    
    log "Setting Telegram webhook..."
    echo "  Token: ${token:0:10}..."
    echo "  URL: $webhook_url"
    echo "  Domain: $DOMAIN"
    
    # Test if the webhook endpoint is accessible
    log "Testing webhook endpoint..."
    if curl -s -f "https://${DOMAIN}/health" > /dev/null; then
        echo "  ✅ Application is reachable"
    else
        echo "  ⚠️  Application may not be ready yet"
    fi
    
    # Set webhook
    log "Configuring webhook with Telegram API..."
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "{\"url\": \"$webhook_url\"}" \
        "https://api.telegram.org/bot${token}/setWebhook")
    
    echo "  Response: $response"
    
    if echo "$response" | grep -q "\"ok\":true"; then
        echo -e "${GREEN}✅ Webhook configured successfully!${NC}"
        
        # Get webhook info
        log "Getting webhook info..."
        webhook_info=$(curl -s "https://api.telegram.org/bot${token}/getWebhookInfo")
        
        if command -v jq &> /dev/null; then
            echo "$webhook_info" | jq .
        else
            echo "$webhook_info"
        fi
    else
        echo -e "${RED}❌ Failed to configure webhook${NC}"
        echo "Please check:"
        echo "  1. Telegram bot token is correct"
        echo "  2. Domain is accessible via HTTPS"
        echo "  3. SSL certificate is valid"
        exit 1
    fi
}

main() {
    log "Starting Telegram webhook configuration..."
    configure_webhook
}

main "$@"