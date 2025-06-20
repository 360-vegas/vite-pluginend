#!/bin/bash

set -e

# GitHub仓库信息
GITHUB_REPO="https://github.com/360-vegas/vite-pluginend.git"
PROJECT_NAME="vite-pluginend"
INSTALL_DIR="/opt/vite-pluginend"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo -e "${CYAN}"
cat << 'EOF'
 ____   ____.__  __             __________.__                 .__                        .___
\   \ /   /|__|/  |_  ____     \______   \  |  __ __  ____  |__| ____   ____   ____   __| _/
 \   Y   / |  \   __\/ __ \     |     ___/  | |  |  \/ ___\ |  |/    \_/ __ \ /    \ / __ | 
  \     /  |  ||  | \  ___/     |    |   |  |_|  |  / /_/  >|  |   |  \  ___/|   |  / /_/ | 
   \___/   |__||__|  \___  >    |____|   |____/____/\___  / |__|___|  /\___  >___|  \____ | 
                         \/                        /_____/          \/     \/     \/     \/ 
                                                                                              
         🚀 一键安装脚本 - 自动从GitHub下载并安装
EOF
echo -e "${NC}"

echo -e "${BLUE}==================================================================${NC}"
echo -e "${YELLOW}GitHub仓库: $GITHUB_REPO${NC}"
echo -e "${YELLOW}安装目录: $INSTALL_DIR${NC}"
echo -e "${BLUE}==================================================================${NC}"

# 检测Linux发行版
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$NAME
        DISTRO=$ID
        VERSION=$VERSION_ID
    else
        echo -e "${RED}无法检测操作系统版本${NC}"
        exit 1
    fi
    echo -e "${BLUE}检测到系统: $OS $VERSION${NC}"
}

# 检查依赖是否已安装
check_dependency() {
    if command -v $1 &> /dev/null; then
        echo -e "${GREEN}✓ $1 已安装${NC}"
        return 0
    else
        echo -e "${YELLOW}⚠ $1 未安装${NC}"
        return 1
    fi
}

# 检查并安装基础依赖
install_basic_deps() {
    echo -e "${BLUE}🔧 检查并安装基础依赖...${NC}"
    
    local need_update=false
    
    # 检查Git
    if ! check_dependency "git"; then
        need_update=true
        echo -e "${BLUE}安装 Git...${NC}"
        case $DISTRO in
            ubuntu|debian)
                if [ "$need_update" = true ]; then
                    sudo apt update
                fi
                sudo apt install -y git
                ;;
            centos|rhel)
                sudo yum install -y git
                ;;
            fedora)
                sudo dnf install -y git
                ;;
            arch)
                sudo pacman -S git --noconfirm
                ;;
            *)
                echo -e "${RED}不支持的发行版，请手动安装Git${NC}"
                exit 1
                ;;
        esac
    fi
    
    # 检查curl和wget
    if ! check_dependency "curl"; then
        echo -e "${BLUE}安装 curl...${NC}"
        case $DISTRO in
            ubuntu|debian)
                if [ "$need_update" = true ]; then
                    sudo apt update
                fi
                sudo apt install -y curl
                ;;
            centos|rhel)
                sudo yum install -y curl
                ;;
            fedora)
                sudo dnf install -y curl
                ;;
            arch)
                sudo pacman -S curl --noconfirm
                ;;
        esac
    fi
    
    if ! check_dependency "wget"; then
        echo -e "${BLUE}安装 wget...${NC}"
        case $DISTRO in
            ubuntu|debian)
                if [ "$need_update" = true ]; then
                    sudo apt update
                fi
                sudo apt install -y wget
                ;;
            centos|rhel)
                sudo yum install -y wget
                ;;
            fedora)
                sudo dnf install -y wget
                ;;
            arch)
                sudo pacman -S wget --noconfirm
                ;;
        esac
    fi
    
    echo -e "${GREEN}✅ 基础依赖检查完成${NC}"
}

# 克隆项目
clone_project() {
    echo -e "${BLUE}📥 从GitHub克隆项目...${NC}"
    
    # 如果目标目录已存在，询问是否删除
    if [ -d "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}目录 $INSTALL_DIR 已存在${NC}"
        read -p "是否删除现有目录并重新安装? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${BLUE}删除现有目录...${NC}"
            sudo rm -rf "$INSTALL_DIR"
        else
            echo -e "${YELLOW}取消安装${NC}"
            exit 0
        fi
    fi
    
    # 创建安装目录的父目录
    sudo mkdir -p $(dirname "$INSTALL_DIR")
    
    # 克隆项目
    echo -e "${BLUE}正在克隆仓库...${NC}"
    sudo git clone "$GITHUB_REPO" "$INSTALL_DIR"
    
    # 设置目录权限
    sudo chown -R $USER:$USER "$INSTALL_DIR"
    
    echo -e "${GREEN}✅ 项目克隆完成${NC}"
}

# 安装Node.js
install_nodejs() {
    echo -e "${BLUE}📦 安装 Node.js...${NC}"
    
    # 使用NodeSource仓库安装最新LTS版本
    case $DISTRO in
        ubuntu|debian)
            curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
            sudo apt-get install -y nodejs
            ;;
        centos|rhel)
            curl -fsSL https://rpm.nodesource.com/setup_lts.x | sudo bash -
            sudo yum install -y nodejs npm
            ;;
        fedora)
            sudo dnf install -y nodejs npm
            ;;
        arch)
            sudo pacman -S nodejs npm --noconfirm
            ;;
        *)
            echo -e "${RED}不支持的发行版，请手动安装Node.js${NC}"
            exit 1
            ;;
    esac
    
    # 验证安装
    node_version=$(node --version)
    npm_version=$(npm --version)
    echo -e "${GREEN}✓ Node.js $node_version 安装成功${NC}"
    echo -e "${GREEN}✓ npm $npm_version 安装成功${NC}"
}

# 安装Go
install_go() {
    echo -e "${BLUE}🔧 安装 Go...${NC}"
    
    # 检查是否已经安装了合适版本的Go
    if command -v go &> /dev/null; then
        go_version=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
        if awk "BEGIN {exit !($go_version >= 1.21)}"; then
            echo -e "${GREEN}✓ Go $go_version 已安装且版本符合要求${NC}"
            return
        fi
    fi
    
    # 下载最新的Go版本
    GO_VERSION="1.21.5"
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64)
            GO_ARCH="amd64"
            ;;
        aarch64|arm64)
            GO_ARCH="arm64"
            ;;
        *)
            echo -e "${RED}不支持的架构: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    # 下载并安装Go
    echo -e "${BLUE}下载 Go ${GO_VERSION}...${NC}"
    cd /tmp
    wget -q "https://golang.org/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    
    echo -e "${BLUE}安装 Go...${NC}"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    rm "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    
    # 设置环境变量
    if ! grep -q '/usr/local/go/bin' ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export GOBIN=$GOPATH/bin' >> ~/.bashrc
    fi
    
    if ! grep -q '/usr/local/go/bin' /etc/profile; then
        echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
    fi
    
    # 当前会话中设置环境变量
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    export GOBIN=$GOPATH/bin
    
    # 验证安装
    go_version=$(/usr/local/go/bin/go version)
    echo -e "${GREEN}✓ $go_version 安装成功${NC}"
}

# 安装数据库
install_databases() {
    echo -e "${BLUE}🗄️  安装数据库...${NC}"
    
    # 询问是否安装MySQL
    read -p "是否安装MySQL? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}安装 MySQL...${NC}"
        case $DISTRO in
            ubuntu|debian)
                sudo apt update
                sudo apt install -y mysql-server mysql-client
                ;;
            centos|rhel)
                sudo yum install -y mysql-server mysql
                ;;
            fedora)
                sudo dnf install -y mysql-server mysql
                ;;
            arch)
                sudo pacman -S mysql --noconfirm
                ;;
        esac
        
        # 启动MySQL服务
        sudo systemctl start mysql || sudo systemctl start mysqld
        sudo systemctl enable mysql || sudo systemctl enable mysqld
        
        echo -e "${GREEN}✓ MySQL 安装完成${NC}"
        echo -e "${YELLOW}请稍后运行 sudo mysql_secure_installation 来配置MySQL安全设置${NC}"
    fi
    
    # 询问是否安装MongoDB
    read -p "是否安装MongoDB? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}安装 MongoDB...${NC}"
        case $DISTRO in
            ubuntu)
                wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
                echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu $(lsb_release -cs)/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
                sudo apt update
                sudo apt install -y mongodb-org
                ;;
            debian)
                wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
                echo "deb http://repo.mongodb.org/apt/debian $(lsb_release -cs)/mongodb-org/6.0 main" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
                sudo apt update
                sudo apt install -y mongodb-org
                ;;
            centos|rhel)
                cat > /tmp/mongodb-org-6.0.repo << 'EOF'
[mongodb-org-6.0]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/$releasever/mongodb-org/6.0/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-6.0.asc
EOF
                sudo mv /tmp/mongodb-org-6.0.repo /etc/yum.repos.d/
                sudo yum install -y mongodb-org
                ;;
            fedora)
                sudo dnf install -y mongodb-server
                ;;
            arch)
                sudo pacman -S mongodb --noconfirm
                ;;
        esac
        
        # 启动MongoDB服务
        sudo systemctl start mongod
        sudo systemctl enable mongod
        
        echo -e "${GREEN}✓ MongoDB 安装完成${NC}"
    fi
}

# 配置数据库
setup_databases() {
    echo -e "${BLUE}⚙️  配置数据库...${NC}"
    
    # 检查MySQL是否运行
    if systemctl is-active --quiet mysql || systemctl is-active --quiet mysqld; then
        echo -e "${BLUE}配置MySQL数据库...${NC}"
        read -p "是否创建项目数据库和用户? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "请输入MySQL root密码:"
            mysql -u root -p << 'EOF'
CREATE DATABASE IF NOT EXISTS vite_pluginend;
CREATE USER IF NOT EXISTS 'pluginend'@'localhost' IDENTIFIED BY 'pluginend123';
GRANT ALL PRIVILEGES ON vite_pluginend.* TO 'pluginend'@'localhost';
FLUSH PRIVILEGES;
SELECT 'MySQL数据库配置完成' AS status;
EOF
            echo -e "${GREEN}✓ MySQL数据库配置完成${NC}"
        fi
    fi
    
    # 检查MongoDB是否运行
    if systemctl is-active --quiet mongod; then
        echo -e "${GREEN}✓ MongoDB 运行正常，无需额外配置${NC}"
    fi
}

# 安装项目依赖并构建
build_project() {
    echo -e "${BLUE}🔨 构建项目...${NC}"
    
    cd "$INSTALL_DIR"
    
    # 安装前端依赖
    echo -e "${BLUE}安装前端依赖...${NC}"
    npm install
    
    # 构建前端
    echo -e "${BLUE}构建前端...${NC}"
    npm run build
    
    # 安装后端依赖
    echo -e "${BLUE}安装后端依赖...${NC}"
    cd backend
    go mod download
    
    # 构建后端
    echo -e "${BLUE}构建后端...${NC}"
    go build -o ../vite-pluginend-server cmd/server/main.go
    cd ..
    
    echo -e "${GREEN}✅ 项目构建完成${NC}"
}

# 创建系统服务
create_service() {
    echo -e "${BLUE}⚙️  创建系统服务...${NC}"
    
    read -p "是否创建systemd服务以自动启动? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        return
    fi
    
    # 创建环境配置文件
    sudo tee /etc/vite-pluginend.env > /dev/null << EOF
PORT=8081
GIN_MODE=release
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=vite_pluginend
UPLOAD_DIR=$INSTALL_DIR/uploads
EOF
    
    # 创建systemd服务文件
    sudo tee /etc/systemd/system/vite-pluginend.service > /dev/null << EOF
[Unit]
Description=Vite Pluginend Server
After=network.target mongod.service mysql.service

[Service]
Type=simple
User=$USER
Group=$USER
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/vite-pluginend-server
EnvironmentFile=/etc/vite-pluginend.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
    
    # 重新加载systemd并启用服务
    sudo systemctl daemon-reload
    sudo systemctl enable vite-pluginend
    
    echo -e "${GREEN}✅ 系统服务创建完成${NC}"
}

# 启动服务
start_services() {
    echo -e "${BLUE}🚀 启动服务...${NC}"
    
    # 启动数据库服务
    echo -e "${BLUE}启动数据库服务...${NC}"
    if systemctl list-unit-files | grep -q mysql; then
        sudo systemctl start mysql
    elif systemctl list-unit-files | grep -q mysqld; then
        sudo systemctl start mysqld
    fi
    
    if systemctl list-unit-files | grep -q mongod; then
        sudo systemctl start mongod
    fi
    
    # 启动应用服务
    if systemctl list-unit-files | grep -q vite-pluginend; then
        echo -e "${BLUE}启动应用服务...${NC}"
        sudo systemctl start vite-pluginend
        
        # 检查服务状态
        sleep 3
        if sudo systemctl is-active --quiet vite-pluginend; then
            echo -e "${GREEN}✅ 应用服务启动成功${NC}"
        else
            echo -e "${RED}⚠ 应用服务启动失败，查看日志:${NC}"
            sudo systemctl status vite-pluginend
        fi
    else
        echo -e "${BLUE}手动启动应用...${NC}"
        cd "$INSTALL_DIR"
        nohup ./vite-pluginend-server > vite-pluginend.log 2>&1 &
        echo $! > vite-pluginend.pid
        echo -e "${GREEN}✅ 应用已在后台启动${NC}"
    fi
}

# 显示安装结果
show_results() {
    echo -e "\n${GREEN}🎉 安装完成！${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo -e "${YELLOW}📁 项目目录: $INSTALL_DIR${NC}"
    echo -e "${YELLOW}🌐 前端访问: http://localhost:3000${NC}"
    echo -e "${YELLOW}🔧 后端API: http://localhost:8081${NC}"
    echo -e "${YELLOW}📦 插件管理: http://localhost:3000/app-market${NC}"
    echo -e "${CYAN}================================================${NC}"
    
    echo -e "\n${BLUE}📋 常用命令:${NC}"
    if systemctl list-unit-files | grep -q vite-pluginend; then
        echo -e "  ${YELLOW}查看服务状态:${NC} sudo systemctl status vite-pluginend"
        echo -e "  ${YELLOW}重启服务:${NC} sudo systemctl restart vite-pluginend"
        echo -e "  ${YELLOW}查看日志:${NC} sudo journalctl -u vite-pluginend -f"
        echo -e "  ${YELLOW}停止服务:${NC} sudo systemctl stop vite-pluginend"
    else
        echo -e "  ${YELLOW}查看日志:${NC} tail -f $INSTALL_DIR/vite-pluginend.log"
        echo -e "  ${YELLOW}停止服务:${NC} kill \$(cat $INSTALL_DIR/vite-pluginend.pid)"
    fi
    
    echo -e "\n${BLUE}🔧 配置文件:${NC}"
    echo -e "  ${YELLOW}环境配置:${NC} /etc/vite-pluginend.env"
    echo -e "  ${YELLOW}项目配置:${NC} $INSTALL_DIR/backend/config/"
    
    echo -e "\n${GREEN}安装完成！请在浏览器中访问上述地址开始使用。${NC}"
}

# 主安装流程
main() {
    echo -e "${BLUE}🚀 开始自动安装流程...${NC}"
    
    # 检查是否以root权限运行
    if [[ $EUID -eq 0 ]]; then
        echo -e "${RED}❌ 请不要以root用户运行此脚本${NC}"
        echo -e "${YELLOW}脚本会在需要时自动请求sudo权限${NC}"
        exit 1
    fi
    
    # 检测操作系统
    detect_os
    
    # 安装基础依赖
    install_basic_deps
    
    # 克隆项目
    clone_project
    
    # 检查并安装依赖
    echo -e "\n${BLUE}📦 检查和安装系统依赖...${NC}"
    
    if ! check_dependency "node"; then
        install_nodejs
    fi
    
    if ! check_dependency "go"; then
        install_go
    fi
    
    # 安装数据库
    install_databases
    
    # 配置数据库
    setup_databases
    
    # 构建项目
    build_project
    
    # 创建服务
    create_service
    
    # 启动服务
    start_services
    
    # 显示结果
    show_results
}

# 捕获中断信号
trap 'echo -e "\n${RED}安装被中断${NC}"; exit 1' INT

# 运行主函数
main "$@" 