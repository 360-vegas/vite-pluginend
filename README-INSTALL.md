# Linux 一键安装指南

## 🚀 快速安装

### 方法一：直接安装（推荐）

```bash
# 一键安装命令
bash <(curl -fsSL https://gist.githubusercontent.com/raw/vite-pluginend-installer.sh)
```

### 方法二：手动安装

```bash
# 1. 安装基础依赖
sudo apt update && sudo apt install -y git curl wget

# 2. 克隆项目
git clone https://github.com/360-vegas/vite-pluginend.git
cd vite-pluginend

# 3. 安装Node.js (如果未安装)
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt-get install -y nodejs

# 4. 安装Go (如果未安装)
wget https://golang.org/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 5. 安装数据库 (可选)
sudo apt install -y mysql-server mongodb

# 6. 构建项目
npm install
npm run build
cd backend
go mod download
go build -o ../vite-pluginend-server cmd/server/main.go
cd ..

# 7. 启动服务
./vite-pluginend-server
```

### 方法三：Docker 安装

```bash
# 使用Docker (如果您有Docker环境)
docker run -d \
  --name vite-pluginend \
  -p 3000:3000 \
  -p 8081:8081 \
  ubuntu:20.04 \
  bash -c "
    apt update && 
    apt install -y git curl nodejs npm golang-go mysql-server mongodb &&
    git clone https://github.com/360-vegas/vite-pluginend.git /app &&
    cd /app &&
    npm install && npm run build &&
    cd backend && go build -o ../server cmd/server/main.go &&
    cd .. && ./server
  "
```

## ⚡ 超级简化版安装

如果您只想快速体验，可以使用这个最简单的命令：

```bash
# 一行命令完成所有安装
curl -sSL https://git.io/vite-pluginend | bash
```

## 🔧 手动配置

### 配置数据库

```bash
# MySQL
sudo mysql -e "CREATE DATABASE vite_pluginend; CREATE USER 'pluginend'@'localhost' IDENTIFIED BY 'pluginend123'; GRANT ALL ON vite_pluginend.* TO 'pluginend'@'localhost';"

# MongoDB (无需特殊配置，直接启动即可)
sudo systemctl start mongod
```

### 启动服务

```bash
# 前端 (开发模式)
npm run dev

# 后端
cd backend
go run cmd/server/main.go

# 或者使用构建后的二进制文件
./vite-pluginend-server
```

## 📍 访问地址

- 前端：http://localhost:3000
- 后端API：http://localhost:8081
- 插件管理：http://localhost:3000/app-market

## 🆘 故障排除

如果遇到问题，请检查：

1. **端口占用**：确保3000和8081端口未被占用
2. **权限问题**：确保有sudo权限
3. **网络连接**：确保能访问GitHub和包管理器
4. **系统要求**：Ubuntu 18.04+, 2GB RAM, 5GB磁盘空间 