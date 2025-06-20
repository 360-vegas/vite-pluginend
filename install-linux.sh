#!/bin/bash

set -e

echo "🐧 Linux插件系统手动安装脚本"
echo "=================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# 安装Node.js
install_nodejs() {
    echo -e "${BLUE}安装 Node.js...${NC}"
    
    # 使用NodeSource仓库安装最新LTS版本
    curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
    
    case $DISTRO in
        ubuntu|debian)
            sudo apt-get install -y nodejs
            ;;
        centos|rhel|fedora)
            sudo yum install -y nodejs npm
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
    echo -e "${BLUE}安装 Go...${NC}"
    
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
    wget -q "https://golang.org/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    rm "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    
    # 设置环境变量
    echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
    
    # 验证安装
    go_version=$(/usr/local/go/bin/go version)
    echo -e "${GREEN}✓ $go_version 安装成功${NC}"
}

# 安装MySQL
install_mysql() {
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
        *)
            echo -e "${YELLOW}请手动安装MySQL${NC}"
            return
            ;;
    esac
    
    # 启动MySQL服务
    sudo systemctl start mysql
    sudo systemctl enable mysql
    
    echo -e "${GREEN}✓ MySQL 安装完成${NC}"
    echo -e "${YELLOW}请运行 sudo mysql_secure_installation 来配置MySQL${NC}"
}

# 安装MongoDB
install_mongodb() {
    echo -e "${BLUE}安装 MongoDB...${NC}"
    
    case $DISTRO in
        ubuntu)
            # 添加MongoDB官方仓库
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
        *)
            echo -e "${YELLOW}请手动安装MongoDB${NC}"
            return
            ;;
    esac
    
    # 启动MongoDB服务
    sudo systemctl start mongod
    sudo systemctl enable mongod
    
    echo -e "${GREEN}✓ MongoDB 安装完成${NC}"
}

# 安装Git
install_git() {
    echo -e "${BLUE}安装 Git...${NC}"
    
    case $DISTRO in
        ubuntu|debian)
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
            echo -e "${YELLOW}请手动安装Git${NC}"
            return
            ;;
    esac
    
    git_version=$(git --version)
    echo -e "${GREEN}✓ $git_version 安装成功${NC}"
}

# 配置数据库
setup_databases() {
    echo -e "${BLUE}配置数据库...${NC}"
    
    # 创建MySQL数据库和用户
    echo -e "${YELLOW}配置MySQL数据库...${NC}"
    echo "请输入MySQL root密码:"
    mysql -u root -p << 'EOF'
CREATE DATABASE IF NOT EXISTS vite_pluginend;
CREATE USER IF NOT EXISTS 'pluginend'@'localhost' IDENTIFIED BY 'pluginend123';
GRANT ALL PRIVILEGES ON vite_pluginend.* TO 'pluginend'@'localhost';
FLUSH PRIVILEGES;
EOF
    
    echo -e "${GREEN}✓ MySQL数据库配置完成${NC}"
    echo -e "${GREEN}✓ MongoDB无需额外配置${NC}"
}

# 安装项目
install_project() {
    echo -e "${BLUE}安装项目...${NC}"
    
    # 创建项目目录
    PROJECT_DIR="/opt/vite-pluginend"
    sudo mkdir -p $PROJECT_DIR
    sudo chown $USER:$USER $PROJECT_DIR
    
    # 如果当前目录是项目目录，复制文件
    if [ -f "package.json" ] && [ -d "backend" ]; then
        echo -e "${BLUE}从当前目录安装项目...${NC}"
        cp -r . $PROJECT_DIR/
    else
        echo -e "${YELLOW}请将项目文件放置到 $PROJECT_DIR 目录${NC}"
        echo -e "${YELLOW}或者在项目根目录中运行此脚本${NC}"
        return
    fi
    
    cd $PROJECT_DIR
    
    # 安装前端依赖
    echo -e "${BLUE}安装前端依赖...${NC}"
    npm install
    
    # 安装后端依赖
    echo -e "${BLUE}安装后端依赖...${NC}"
    cd backend
    go mod download
    cd ..
    
    echo -e "${GREEN}✓ 项目依赖安装完成${NC}"
}

# 构建项目
build_project() {
    echo -e "${BLUE}构建项目...${NC}"
    
    cd $PROJECT_DIR
    
    # 构建前端
    echo -e "${BLUE}构建前端...${NC}"
    npm run build
    
    # 构建后端
    echo -e "${BLUE}构建后端...${NC}"
    cd backend
    go build -o ../vite-pluginend-server cmd/server/main.go
    cd ..
    
    echo -e "${GREEN}✓ 项目构建完成${NC}"
}

# 创建systemd服务
create_service() {
    echo -e "${BLUE}创建系统服务...${NC}"
    
    # 创建环境配置文件
    cat > /tmp/vite-pluginend.env << EOF
PORT=8081
GIN_MODE=release
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB=vite_pluginend
UPLOAD_DIR=/opt/vite-pluginend/uploads
EOF
    
    sudo mv /tmp/vite-pluginend.env /etc/vite-pluginend.env
    
    # 创建systemd服务文件
    cat > /tmp/vite-pluginend.service << 'EOF'
[Unit]
Description=Vite Pluginend Server
After=network.target mongod.service mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/vite-pluginend
ExecStart=/opt/vite-pluginend/vite-pluginend-server
EnvironmentFile=/etc/vite-pluginend.env
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF
    
    sudo mv /tmp/vite-pluginend.service /etc/systemd/system/
    
    # 重新加载systemd并启用服务
    sudo systemctl daemon-reload
    sudo systemctl enable vite-pluginend
    
    echo -e "${GREEN}✓ 系统服务创建完成${NC}"
}

# 启动服务
start_services() {
    echo -e "${BLUE}启动服务...${NC}"
    
    # 启动数据库服务
    sudo systemctl start mysql
    sudo systemctl start mongod
    
    # 启动应用服务
    sudo systemctl start vite-pluginend
    
    # 检查服务状态
    if sudo systemctl is-active --quiet vite-pluginend; then
        echo -e "${GREEN}✓ 应用服务启动成功${NC}"
    else
        echo -e "${RED}✗ 应用服务启动失败${NC}"
        sudo systemctl status vite-pluginend
    fi
}

# 主安装流程
main() {
    echo -e "${BLUE}开始安装插件系统...${NC}"
    
    # 检测操作系统
    detect_os
    
    # 检查现有依赖
    echo -e "\n${BLUE}检查系统依赖...${NC}"
    
    # 检查并安装各个组件
    if ! check_dependency "node"; then
        install_nodejs
    fi
    
    if ! check_dependency "go"; then
        install_go
    fi
    
    if ! check_dependency "git"; then
        install_git
    fi
    
    if ! check_dependency "mysql"; then
        read -p "是否安装MySQL? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            install_mysql
        fi
    fi
    
    if ! check_dependency "mongod"; then
        read -p "是否安装MongoDB? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            install_mongodb
        fi
    fi
    
    # 配置数据库
    read -p "是否配置数据库? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        setup_databases
    fi
    
    # 安装项目
    install_project
    
    # 构建项目
    build_project
    
    # 创建服务
    read -p "是否创建系统服务? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        create_service
        start_services
    fi
    
    echo -e "\n${GREEN}🎉 安装完成！${NC}"
    echo -e "${YELLOW}服务信息:${NC}"
    echo -e "  - Web界面: http://localhost:3000"
    echo -e "  - API服务: http://localhost:8081"
    echo -e "  - 项目目录: /opt/vite-pluginend"
    echo -e "\n${YELLOW}常用命令:${NC}"
    echo -e "  - 查看服务状态: sudo systemctl status vite-pluginend"
    echo -e "  - 重启服务: sudo systemctl restart vite-pluginend"
    echo -e "  - 查看日志: sudo journalctl -u vite-pluginend -f"
}

# 检查是否以root权限运行
if [[ $EUID -eq 0 ]]; then
   echo -e "${RED}请不要以root用户运行此脚本${NC}"
   echo -e "${YELLOW}脚本会在需要时自动请求sudo权限${NC}"
   exit 1
fi

# 运行主函数
main "$@" 