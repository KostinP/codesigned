#!/bin/bash

# VPS requirements check script

set -e

source ./scripts/env.sh .env.prod

check_required_vars VPS_IP VPS_USER

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ssh_vps() {
    if [ -n "$VPS_SSH_KEY_PATH" ] && [ -f "$VPS_SSH_KEY_PATH" ]; then
        ssh -o StrictHostKeyChecking=no -i "$VPS_SSH_KEY_PATH" ${VPS_USER}@${VPS_IP} "$@"
    elif [ -n "$VPS_PASSWORD" ]; then
        sshpass -p "$VPS_PASSWORD" ssh -o StrictHostKeyChecking=no ${VPS_USER}@${VPS_IP} "$@"
    else
        echo -e "${RED}❌ No SSH key or password provided${NC}"
        exit 1
    fi
}

check_connection() {
    echo -e "${BLUE}🔗 Checking VPS connection...${NC}"
    if ssh_vps "echo 'Connected successfully'"; then
        echo -e "${GREEN}✅ VPS connection: OK${NC}"
        return 0
    else
        echo -e "${RED}❌ VPS connection: FAILED${NC}"
        return 1
    fi
}

check_os() {
    echo -e "${BLUE}🖥️  Checking OS...${NC}"
    if ssh_vps "test -f /etc/os-release"; then
        ssh_vps "source /etc/os-release && echo \"✅ OS: \$NAME \$VERSION\""
        return 0
    else
        echo -e "${RED}❌ Unsupported OS${NC}"
        return 1
    fi
}

check_docker() {
    echo -e "${BLUE}🐳 Checking Docker...${NC}"
    if ssh_vps "command -v docker &> /dev/null"; then
        ssh_vps "docker --version"
        echo -e "${GREEN}✅ Docker: INSTALLED${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  Docker: NOT INSTALLED${NC}"
        return 1
    fi
}

check_docker_compose() {
    echo -e "${BLUE}🎯 Checking Docker Compose...${NC}"
    if ssh_vps "command -v docker compose &> /dev/null || command -v docker &> /dev/null && docker compose version &> /dev/null"; then
        if ssh_vps "command -v docker compose &> /dev/null"; then
            ssh_vps "docker compose --version"
        else
            ssh_vps "docker compose version"
        fi
        echo -e "${GREEN}✅ Docker Compose: INSTALLED${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  Docker Compose: NOT INSTALLED${NC}"
        return 1
    fi
}

check_nginx() {
    echo -e "${BLUE}🌐 Checking Nginx...${NC}"
    if ssh_vps "command -v nginx &> /dev/null"; then
        ssh_vps "nginx -v"
        echo -e "${GREEN}✅ Nginx: INSTALLED${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠️  Nginx: NOT INSTALLED${NC}"
        return 1
    fi
}

check_resources() {
    echo -e "${BLUE}📊 Checking system resources...${NC}"
    echo -e "Memory:"
    ssh_vps "free -h"
    echo -e "\nDisk:"
    ssh_vps "df -h /"
    echo -e "\nCPU:"
    ssh_vps "nproc"
}

check_ports() {
    echo -e "${BLUE}🔌 Checking open ports...${NC}"
    ssh_vps "ss -tulpn | grep -E ':(80|443|22)' || true"
}

check_ufw() {
    echo -e "${BLUE}🔥 Checking firewall...${NC}"
    if ssh_vps "command -v ufw &> /dev/null"; then
        ssh_vps "ufw status"
    else
        echo -e "${YELLOW}⚠️  UFW: NOT INSTALLED${NC}"
    fi
}

generate_report() {
    local missing_tools=()
    
    echo -e "\n${BLUE}📋 SYSTEM CHECK REPORT${NC}"
    echo "===================="
    
    if ! check_connection; then
        echo -e "${RED}❌ Cannot proceed - VPS connection failed${NC}"
        exit 1
    fi
    
    check_os
    check_resources
    
    echo -e "\n${BLUE}REQUIRED TOOLS:${NC}"
    if ! check_docker; then missing_tools+=("docker"); fi
    if ! check_docker_compose; then missing_tools+=("docker-compose"); fi
    if ! check_nginx; then missing_tools+=("nginx"); fi
    
    check_ports
    check_ufw
    
    if [ ${#missing_tools[@]} -eq 0 ]; then
        echo -e "\n${GREEN}🎉 All requirements satisfied!${NC}"
        echo -e "Your VPS is ready for deployment."
    else
        echo -e "\n${YELLOW}⚠️  Missing tools: ${missing_tools[*]}${NC}"
        echo -e "Run './scripts/setup_vps.sh' to install missing components."
    fi
}

main() {
    echo -e "${BLUE}VPS Requirements Check${NC}"
    echo "=========================="
    generate_report
}

main "$@"