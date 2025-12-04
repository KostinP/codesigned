#!/bin/bash

# Quick health check for VPS

set -e

source ./scripts/env.sh .env.prod

check_required_vars VPS_IP VPS_USER

ssh_vps() {
    if [ -n "$VPS_SSH_KEY_PATH" ] && [ -f "$VPS_SSH_KEY_PATH" ]; then
        ssh -o StrictHostKeyChecking=no -i "$VPS_SSH_KEY_PATH" ${VPS_USER}@${VPS_IP} "$@"
    elif [ -n "$VPS_PASSWORD" ]; then
        sshpass -p "$VPS_PASSWORD" ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_IP} "$@"
    else
        echo "No SSH credentials configured"
        exit 1
    fi
}

echo "🔍 Quick VPS Health Check"
echo "========================"

# Check services
echo "Services status:"
ssh_vps "
    echo '- Docker:'; systemctl is-active docker
    echo '- Nginx:'; systemctl is-active nginx
    echo '- Containers:'; docker ps --format 'table {{.Names}}\\t{{.Status}}'
"

# Check resources
echo -e "\nResource usage:"
ssh_vps "
    echo '- Memory:'; free -h | grep Mem
    echo '- Disk:'; df -h / | grep -v Filesystem
    echo '- Load:'; uptime | awk '{print \$10 \$11 \$12}'
"

# Check application
echo -e "\nApplication health:"
if ssh_vps "curl -s -f https://localhost/health > /dev/null"; then
    echo "✅ Application is healthy"
else
    echo "❌ Application health check failed"
fi

echo -e "\n🌐 External check:"
if curl -s -f "https://$VPS_IP/health" > /dev/null; then
    echo "✅ Application is accessible externally"
else
    echo "❌ Application is not accessible externally"
fi
