#!/bin/bash

set -e

# Load environment
source ./scripts/env.sh .env.prod

check_required_vars VPS_IP VPS_USER

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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

check_vps_connection() {
    log "Checking VPS connection..."
    if ssh_vps "echo 'Connection successful'"; then
        success "VPS connection verified"
        return 0
    else
        error "Cannot connect to VPS"
        return 1
    fi
}

update_system() {
    log "Updating system packages..."
    ssh_vps "
        export DEBIAN_FRONTEND=noninteractive
        apt-get update
        apt-get upgrade -y
        apt-get autoremove -y
        apt-get clean
    "
    success "System updated"
}

install_basic_tools() {
    log "Installing basic tools..."
    ssh_vps "
        export DEBIAN_FRONTEND=noninteractive
        apt-get install -y \
            curl \
            wget \
            git \
            htop \
            nload \
            tree \
            ncdu \
            zip \
            unzip \
            fail2ban \
            ufw \
            cron \
            logrotate
    "
    success "Basic tools installed"
}

install_docker() {
    log "Installing Docker..."
    ssh_vps "
        export DEBIAN_FRONTEND=noninteractive
        
        # Remove any conflicting packages
        apt-get remove --purge -y docker docker-engine docker.io containerd runc 2>/dev/null || true
        
        # Install prerequisites
        apt-get install -y \
            apt-transport-https \
            ca-certificates \
            gnupg \
            lsb-release
        
        # Add Docker repository
        mkdir -p /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
        echo \
            \"deb [arch=\$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
            \$(lsb_release -cs) stable\" | tee /etc/apt/sources.list.d/docker.list > /dev/null
        
        # Install Docker
        apt-get update
        apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin
        
        # Start and enable Docker
        systemctl enable docker
        systemctl start docker
    "
    success "Docker installed"
}

install_docker_compose() {
    log "Installing Docker Compose..."
    ssh_vps "
        # Install Docker Compose (standalone for compatibility)
        curl -L \"https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-\$(uname -s)-\$(uname -m)\" -o /usr/local/bin/docker-compose
        chmod +x /usr/local/bin/docker-compose
        ln -sf /usr/local/bin/docker-compose /usr/bin/docker-compose
    "
    success "Docker Compose installed"
}

install_nginx() {
    log "Installing Nginx..."
    ssh_vps "
        export DEBIAN_FRONTEND=noninteractive
        apt-get install -y nginx
        systemctl enable nginx
        systemctl start nginx
    "
    success "Nginx installed"
}

setup_firewall() {
    log "Setting up firewall..."
    ssh_vps "
        # Reset UFW to defaults
        ufw --force reset
        
        # Set default policies
        ufw default deny incoming
        ufw default allow outgoing
        
        # Allow SSH, HTTP and HTTPS
        ufw allow 22/tcp
        ufw allow 80/tcp
        ufw allow 443/tcp
        
        # Enable UFW
        ufw --force enable
    "
    success "Firewall configured"
}

setup_swap() {
    log "Setting up swap file..."
    ssh_vps "
        # Check if swap already exists
        if [ ! -f /swapfile ]; then
            # Create 2GB swap file
            fallocate -l 2G /swapfile
            chmod 600 /swapfile
            mkswap /swapfile
            swapon /swapfile
            
            # Make permanent
            echo '/swapfile none swap sw 0 0' >> /etc/fstab
            
            # Optimize system settings
            echo 'vm.swappiness=10' >> /etc/sysctl.conf
            echo 'net.ipv4.tcp_congestion_control=bbr' >> /etc/sysctl.conf
            sysctl -p
        fi
    "
    success "Swap configured"
}

setup_docker_user() {
    log "Setting up Docker user permissions..."
    ssh_vps "
        # Add user to docker group
        usermod -aG docker $VPS_USER 2>/dev/null || true
        
        # Create necessary directories
        mkdir -p /opt/edu-platform
        mkdir -p /opt/edu-platform/logs
        mkdir -p /opt/edu-platform/data
    "
    success "Docker user setup completed"
}

setup_monitoring() {
    log "Setting up basic monitoring..."
    ssh_vps "
        # Enable and start fail2ban
        systemctl enable fail2ban
        systemctl start fail2ban
    "
    success "Monitoring setup completed"
}

setup_ssl_directory() {
    log "Setting up SSL directory..."
    ssh_vps "
        mkdir -p /etc/nginx/ssl
        chmod 700 /etc/nginx/ssl
    "
    success "SSL directory created"
}

show_system_info() {
    log "System information:"
    ssh_vps "
        echo '=== OS Information ==='
        lsb_release -a 2>/dev/null || cat /etc/os-release
        
        echo ''
        echo '=== Docker Information ==='
        docker --version
        docker-compose --version
        
        echo ''
        echo '=== Nginx Information ==='
        nginx -v 2>&1
        
        echo ''
        echo '=== Disk Usage ==='
        df -h
        
        echo ''
        echo '=== Memory Usage ==='
        free -h
        
        echo ''
        echo '=== UFW Status ==='
        ufw status
    "
}

main() {
    log "Starting VPS setup process..."
    echo "Using environment: prod"
    
    # Check connection first
    check_vps_connection
    
    # Perform setup steps
    update_system
    install_basic_tools
    install_docker
    install_docker_compose
    install_nginx
    setup_swap
    setup_firewall
    setup_docker_user
    setup_monitoring
    setup_ssl_directory
    
    # Show system information
    show_system_info
    
    success "🎉 VPS setup completed successfully!"
    echo ""
    echo "Next steps:"
    echo "  1. Run: make prod-upload-certs"
    echo "  2. Run: make prod-deploy" 
    echo "  3. Run: make prod-setup-webhook"
    echo ""
    echo "Or run everything at once: make prod-full-deploy"
    echo ""
    echo "Your VPS is now ready for application deployment!"
}

main "$@"