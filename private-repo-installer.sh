#!/bin/bash

# Vite Pluginend 私有仓库安装脚本
# 适用于私有GitHub仓库的安装方案

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔒 Vite Pluginend 私有仓库安装向导${NC}"
echo -e "${BLUE}=======================================${NC}"

# 检查是否提供了访问令牌
if [ -z "$GITHUB_TOKEN" ]; then
    echo -e "${YELLOW}⚠️  检测到私有仓库安装${NC}"
    echo -e "${YELLOW}请提供以下信息之一进行认证：${NC}"
    echo
    echo -e "${BLUE}方法一：使用Personal Access Token${NC}"
    echo "1. 访问：https://github.com/settings/tokens"
    echo "2. 创建新的token，权限选择 'repo'"
    echo "3. 设置环境变量：export GITHUB_TOKEN=your_token"
    echo "4. 重新运行此脚本"
    echo
    echo -e "${BLUE}方法二：使用用户名密码${NC}"
    read -p "GitHub用户名: " GITHUB_USER
    read -s -p "GitHub密码或Token: " GITHUB_PASS
    echo
    
    if [ -n "$GITHUB_USER" ] && [ -n "$GITHUB_PASS" ]; then
        GITHUB_AUTH="$GITHUB_USER:$GITHUB_PASS"
    else
        echo -e "${RED}❌ 认证信息不完整${NC}"
        exit 1
    fi
fi

# 设置仓库信息
GITHUB_REPO="https://github.com/360-vegas/vite-pluginend.git"
INSTALL_DIR="/opt/vite-pluginend"

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

# 克隆私有仓库
clone_private_repo() {
    echo -e "${BLUE}📥 克隆私有仓库...${NC}"
    
    # 如果目标目录已存在，询问是否删除
    if [ -d "$INSTALL_DIR" ]; then
        echo -e "${YELLOW}目录 $INSTALL_DIR 已存在${NC}"
        read -p "是否删除现有目录并重新安装? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            sudo rm -rf "$INSTALL_DIR"
        else
            echo -e "${YELLOW}取消安装${NC}"
            exit 0
        fi
    fi
    
    # 创建安装目录
    sudo mkdir -p $(dirname "$INSTALL_DIR")
    
    # 克隆私有仓库
    if [ -n "$GITHUB_TOKEN" ]; then
        # 使用Token
        echo -e "${BLUE}使用Access Token克隆...${NC}"
        sudo git clone https://${GITHUB_TOKEN}@github.com/360-vegas/vite-pluginend.git "$INSTALL_DIR"
    elif [ -n "$GITHUB_AUTH" ]; then
        # 使用用户名密码
        echo -e "${BLUE}使用用户名密码克隆...${NC}"
        sudo git clone https://${GITHUB_AUTH}@github.com/360-vegas/vite-pluginend.git "$INSTALL_DIR"
    else
        echo -e "${RED}❌ 缺少认证信息${NC}"
        exit 1
    fi
    
    # 设置目录权限
    sudo chown -R $USER:$USER "$INSTALL_DIR"
    
    echo -e "${GREEN}✅ 私有仓库克隆完成${NC}"
}

# 安装依赖
install_dependencies() {
    echo -e "${BLUE}📦 安装系统依赖...${NC}"
    
    case $DISTRO in
        ubuntu|debian)
            sudo apt update
            sudo apt install -y curl wget git
            
            # 安装Node.js
            if ! command -v node &> /dev/null; then
                curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
                sudo apt-get install -y nodejs
            fi
            ;;
        centos|rhel)
            sudo yum install -y curl wget git
            
            # 安装Node.js
            if ! command -v node &> /dev/null; then
                curl -fsSL https://rpm.nodesource.com/setup_lts.x | sudo bash -
                sudo yum install -y nodejs npm
            fi
            ;;
        *)
            echo -e "${RED}❌ 不支持的发行版${NC}"
            exit 1
            ;;
    esac
    
    # 安装Go
    if ! command -v go &> /dev/null; then
        echo -e "${BLUE}安装 Go...${NC}"
        wget https://golang.org/dl/go1.21.5.linux-amd64.tar.gz
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
        rm go1.21.5.linux-amd64.tar.gz
        
        # 设置环境变量
        export PATH=$PATH:/usr/local/go/bin
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    fi
    
    echo -e "${GREEN}✅ 依赖安装完成${NC}"
}

# 构建项目
build_project() {
    echo -e "${BLUE}🔨 构建项目...${NC}"
    
    cd "$INSTALL_DIR"
    
    # 构建前端
    echo -e "${BLUE}构建前端...${NC}"
    npm install
    npm run build
    
    # 构建后端
    echo -e "${BLUE}构建后端...${NC}"
    cd backend
    go mod download
    go build -o ../vite-pluginend-server cmd/server/main.go
    cd ..
    
    echo -e "${GREEN}✅ 项目构建完成${NC}"
}

# 启动服务
start_service() {
    echo -e "${BLUE}🚀 启动服务...${NC}"
    
    cd "$INSTALL_DIR"
    
    # 后台启动服务
    nohup ./vite-pluginend-server > vite-pluginend.log 2>&1 &
    echo $! > vite-pluginend.pid
    
    echo -e "${GREEN}✅ 服务启动完成${NC}"
}

# 显示结果
show_results() {
    echo -e "\n${GREEN}🎉 私有仓库安装完成！${NC}"
    echo -e "${BLUE}=======================================${NC}"
    echo -e "${YELLOW}📁 安装目录: $INSTALL_DIR${NC}"
    echo -e "${YELLOW}🌐 访问地址: http://$(hostname -I | awk '{print $1}'):3000${NC}"
    echo -e "${YELLOW}🔧 API地址: http://$(hostname -I | awk '{print $1}'):8081${NC}"
    echo -e "${BLUE}=======================================${NC}"
    
    echo -e "\n${BLUE}📋 管理命令:${NC}"
    echo -e "  查看日志: tail -f $INSTALL_DIR/vite-pluginend.log"
    echo -e "  停止服务: kill \$(cat $INSTALL_DIR/vite-pluginend.pid)"
    echo -e "  重启服务: cd $INSTALL_DIR && ./vite-pluginend-server &"
}

# 主函数
main() {
    detect_os
    clone_private_repo
    install_dependencies
    build_project
    start_service
    show_results
}

# 运行主函数
main "$@" 