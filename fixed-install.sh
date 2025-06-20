#!/bin/bash

set -e

echo "🔧 Vite Pluginend 修复版安装脚本"
echo "=================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 检测系统
if [ -f /etc/os-release ]; then
    . /etc/os-release
    echo -e "${BLUE}检测到系统: $NAME $VERSION_ID${NC}"
else
    echo -e "${RED}❌ 无法检测操作系统${NC}"
    exit 1
fi

# 更新系统包
echo -e "${BLUE}📦 更新系统包...${NC}"
if command -v apt &> /dev/null; then
    apt update
    apt install -y curl wget git build-essential ca-certificates
elif command -v yum &> /dev/null; then
    yum update -y
    yum install -y curl wget git gcc gcc-c++ make ca-certificates
elif command -v dnf &> /dev/null; then
    dnf update -y
    dnf install -y curl wget git gcc gcc-c++ make ca-certificates
else
    echo -e "${RED}❌ 不支持的包管理器${NC}"
    exit 1
fi

# 强制安装最新版Node.js
echo -e "${BLUE}📦 安装最新版 Node.js...${NC}"
NODE_VERSION="v20.10.0"
ARCH=$(uname -m)
case $ARCH in
    x86_64) NODE_ARCH="x64" ;;
    aarch64) NODE_ARCH="arm64" ;;
    *) echo -e "${RED}❌ 不支持的架构: $ARCH${NC}"; exit 1 ;;
esac

# 移除旧版本Node.js
rm -rf /usr/local/bin/node /usr/local/bin/npm /usr/local/bin/npx
rm -rf /usr/local/lib/node_modules
rm -rf /usr/local/include/node

cd /tmp
echo -e "${BLUE}下载 Node.js ${NODE_VERSION}...${NC}"
wget "https://nodejs.org/dist/${NODE_VERSION}/node-${NODE_VERSION}-linux-${NODE_ARCH}.tar.xz"
tar -xf "node-${NODE_VERSION}-linux-${NODE_ARCH}.tar.xz"

echo -e "${BLUE}安装 Node.js...${NC}"
cp -r "node-${NODE_VERSION}-linux-${NODE_ARCH}"/* /usr/local/
rm -rf "node-${NODE_VERSION}-linux-${NODE_ARCH}"*

# 更新PATH
export PATH="/usr/local/bin:$PATH"

# 验证Node.js安装
echo -e "${BLUE}验证 Node.js 安装...${NC}"
/usr/local/bin/node --version
/usr/local/bin/npm --version

# 安装最新版Go
echo -e "${BLUE}📦 安装最新版 Go...${NC}"
GO_VERSION="1.21.5"
ARCH=$(uname -m)
case $ARCH in
    x86_64) GO_ARCH="amd64" ;;
    aarch64) GO_ARCH="arm64" ;;
    *) echo -e "${RED}❌ 不支持的架构: $ARCH${NC}"; exit 1 ;;
esac

cd /tmp
echo -e "${BLUE}下载 Go ${GO_VERSION}...${NC}"
wget "https://golang.org/dl/go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"

echo -e "${BLUE}安装 Go...${NC}"
rm -rf /usr/local/go
tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
rm "go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"

# 设置Go环境变量
export PATH="/usr/local/go/bin:$PATH"
export GOPATH="/root/go"
export GOPROXY="https://goproxy.io,direct"
export GOSUMDB="sum.golang.org"

# 添加到环境变量
echo 'export PATH="/usr/local/go/bin:$PATH"' >> /etc/profile
echo 'export GOPATH="/root/go"' >> /etc/profile
echo 'export GOPROXY="https://goproxy.io,direct"' >> /etc/profile

# 验证Go安装
echo -e "${BLUE}验证 Go 安装...${NC}"
/usr/local/go/bin/go version

# 设置npm配置
echo -e "${BLUE}⚙️ 配置npm...${NC}"
/usr/local/bin/npm config set registry https://registry.npmjs.org/
/usr/local/bin/npm config set fund false
/usr/local/bin/npm config set audit-level moderate

# 项目设置
PROJECT_DIR="/opt/vite-pluginend"
echo -e "${BLUE}📁 项目目录: $PROJECT_DIR${NC}"

# 尝试克隆项目或创建基本结构
echo -e "${BLUE}📥 获取项目...${NC}"
if [ -d "$PROJECT_DIR" ]; then
    cd "$PROJECT_DIR"
    echo -e "${YELLOW}⚠️ 项目目录已存在，尝试更新...${NC}"
    
    # 如果是git仓库，尝试拉取
    if [ -d ".git" ]; then
        git pull || echo -e "${YELLOW}⚠️ Git拉取失败，继续使用现有代码${NC}"
    fi
else
    # 尝试克隆，如果失败则创建基本结构
    if git clone https://github.com/360-vegas/vite-pluginend.git "$PROJECT_DIR" 2>/dev/null; then
        echo -e "${GREEN}✅ 成功克隆GitHub仓库${NC}"
        cd "$PROJECT_DIR"
    else
        echo -e "${YELLOW}⚠️ GitHub克隆失败，创建基本项目结构...${NC}"
        mkdir -p "$PROJECT_DIR"
        cd "$PROJECT_DIR"
        
        # 创建现代的package.json
        cat > package.json << 'PACKAGE_EOF'
{
  "name": "vite-pluginend",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite --host 0.0.0.0 --port 3000",
    "build": "vite build",
    "preview": "vite preview --host 0.0.0.0 --port 3000"
  },
  "dependencies": {
    "vue": "^3.3.0",
    "vue-router": "^4.2.0",
    "element-plus": "^2.4.0",
    "axios": "^1.6.0"
  },
  "devDependencies": {
    "@vitejs/plugin-vue": "^4.5.0",
    "vite": "^5.0.0",
    "typescript": "^5.2.0"
  }
}
PACKAGE_EOF

        # 创建基本的vite.config.js
        cat > vite.config.js << 'VITE_EOF'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0',
    port: 3000
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets'
  }
})
VITE_EOF

        # 创建基本的index.html
        mkdir -p public
        cat > index.html << 'HTML_EOF'
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Vite Pluginend</title>
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.js"></script>
  </body>
</html>
HTML_EOF

        # 创建基本的Vue应用
        mkdir -p src
        cat > src/main.js << 'MAIN_EOF'
import { createApp } from 'vue'
import App from './App.vue'

createApp(App).mount('#app')
MAIN_EOF

        cat > src/App.vue << 'APP_EOF'
<template>
  <div id="app">
    <h1>🚀 Vite Pluginend</h1>
    <p>项目安装成功！</p>
    <p>当前时间: {{ currentTime }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const currentTime = ref('')

onMounted(() => {
  currentTime.value = new Date().toLocaleString()
  setInterval(() => {
    currentTime.value = new Date().toLocaleString()
  }, 1000)
})
</script>

<style>
#app {
  text-align: center;
  margin-top: 50px;
  font-family: Arial, sans-serif;
}
</style>
APP_EOF

        # 创建后端目录和文件
        mkdir -p backend/cmd/server
        
        cat > backend/go.mod << 'GOMOD_EOF'
module vite-pluginend

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/gin-contrib/cors v1.4.0
)
GOMOD_EOF

        cat > backend/cmd/server/main.go << 'MAIN_GO_EOF'
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    r := gin.Default()
    
    // 配置CORS
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
    config.AllowHeaders = []string{"*"}
    r.Use(cors.New(config))
    
    // 健康检查
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
            "message": "Vite Pluginend Server is running!",
        })
    })
    
    // API路由
    api := r.Group("/api")
    {
        api.GET("/status", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{
                "status": "success",
                "message": "API is working",
            })
        })
    }
    
    r.Run(":8081")
}
MAIN_GO_EOF
        
        echo -e "${GREEN}✅ 创建了基本项目结构${NC}"
    fi
fi

# 构建前端
echo -e "${BLUE}🎨 构建前端...${NC}"
if [ -f "package.json" ]; then
    echo -e "${BLUE}安装前端依赖...${NC}"
    /usr/local/bin/npm install --no-audit --no-fund
    
    echo -e "${BLUE}构建前端应用...${NC}"
    /usr/local/bin/npm run build
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ 前端构建成功${NC}"
    else
        echo -e "${YELLOW}⚠️ 前端构建失败，但继续后端构建${NC}"
    fi
else
    echo -e "${YELLOW}⚠️ 没有找到package.json，跳过前端构建${NC}"
fi

# 构建后端
echo -e "${BLUE}🔧 构建后端...${NC}"
cd backend

# 清理go.sum文件以避免校验和问题
rm -f go.sum

# 初始化Go模块
/usr/local/go/bin/go mod tidy

# 下载依赖
echo -e "${BLUE}下载Go依赖...${NC}"
/usr/local/go/bin/go mod download

# 构建二进制文件
echo -e "${BLUE}编译Go应用...${NC}"
/usr/local/go/bin/go build -o ../vite-pluginend-server cmd/server/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 后端构建成功${NC}"
else
    echo -e "${RED}❌ 后端构建失败${NC}"
    exit 1
fi

cd ..

# 启动服务
echo -e "${BLUE}🚀 启动服务...${NC}"
chmod +x vite-pluginend-server

# 停止之前的服务（如果存在）
if [ -f "vite-pluginend.pid" ]; then
    OLD_PID=$(cat vite-pluginend.pid)
    if kill -0 "$OLD_PID" 2>/dev/null; then
        echo -e "${YELLOW}停止旧服务 (PID: $OLD_PID)...${NC}"
        kill "$OLD_PID"
        sleep 2
    fi
fi

# 启动新服务
nohup ./vite-pluginend-server > vite-pluginend.log 2>&1 &
NEW_PID=$!
echo $NEW_PID > vite-pluginend.pid

sleep 3

# 检查服务状态
if kill -0 $NEW_PID 2>/dev/null; then
    echo -e "${GREEN}✅ 服务启动成功 (PID: $NEW_PID)${NC}"
    
    # 获取服务器IP
    SERVER_IP=$(hostname -I | awk '{print $1}' | head -n1)
    if [ -z "$SERVER_IP" ]; then
        SERVER_IP="localhost"
    fi
    
    echo -e "${BLUE}📍 访问地址:${NC}"
    echo -e "   🌐 前端: http://$SERVER_IP:3000"
    echo -e "   🔧 后端: http://$SERVER_IP:8081"
    echo -e "   💚 健康检查: http://$SERVER_IP:8081/health"
    
    echo -e "\n${BLUE}📋 管理命令:${NC}"
    echo -e "   查看日志: tail -f $PROJECT_DIR/vite-pluginend.log"
    echo -e "   停止服务: kill $NEW_PID"
    echo -e "   重启服务: cd $PROJECT_DIR && ./vite-pluginend-server &"
    
    # 测试服务
    echo -e "\n${BLUE}🧪 测试服务...${NC}"
    sleep 2
    if curl -s "http://localhost:8081/health" > /dev/null; then
        echo -e "${GREEN}✅ 服务响应正常${NC}"
    else
        echo -e "${YELLOW}⚠️ 服务可能需要更多时间启动${NC}"
    fi
    
else
    echo -e "${RED}❌ 服务启动失败${NC}"
    echo -e "${YELLOW}查看错误日志:${NC}"
    tail -20 vite-pluginend.log
    exit 1
fi

echo -e "\n${GREEN}🎉 安装完成！${NC}"
echo -e "${GREEN}项目已成功部署到 $PROJECT_DIR${NC}" 