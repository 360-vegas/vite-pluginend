#!/bin/bash

set -e

echo "🏗️  构建跨平台安装向导程序..."

# 检查依赖
check_dependency() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ $1 未安装，请先安装 $1"
        exit 1
    fi
}

echo "🔍 检查构建依赖..."
check_dependency "go"
check_dependency "node"
check_dependency "npm"

# 清理之前的构建
echo "🧹 清理构建目录..."
rm -rf dist/
mkdir -p dist/

# 构建前端
echo "📦 构建前端界面..."
cd web
npm install
npm run build
cd ..

# 创建资源目录
echo "📁 准备资源文件..."
mkdir -p assets

# 设置Go模块
echo "📦 下载Go依赖..."
go mod tidy

# 构建不同平台的安装程序
echo "🔨 构建Go程序..."

# Windows 64位
echo "📦 构建 Windows x64..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/installer-windows-amd64.exe cmd/installer/main.go

# macOS 64位 (Intel)
echo "📦 构建 macOS x64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/installer-macos-amd64 cmd/installer/main.go

# macOS ARM64 (Apple Silicon)
echo "📦 构建 macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/installer-macos-arm64 cmd/installer/main.go

# Linux 64位
echo "📦 构建 Linux x64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/installer-linux-amd64 cmd/installer/main.go

# 生成README
echo "📝 生成使用说明..."
cat > dist/README.md << 'EOF'
# 插件系统安装程序

## 使用方法

### Windows
```bash
# 下载并运行
installer-windows-amd64.exe
```

### macOS
```bash
# 赋予执行权限
chmod +x installer-macos-amd64

# 运行安装程序
./installer-macos-amd64
```

### Linux
```bash
# 赋予执行权限
chmod +x installer-linux-amd64

# 运行安装程序
./installer-linux-amd64
```

## 选项

- `-port 8888` : 指定Web界面端口
- `-silent` : 静默安装模式
- `-help` : 显示帮助信息

## 系统要求

### 最低要求
- Windows 10 / macOS 10.15 / Ubuntu 18.04
- 2GB RAM
- 5GB 可用磁盘空间
- 网络连接

### 自动安装的软件
- Node.js (≥18.0)
- Go (≥1.21)
- MySQL (≥8.0)
- MongoDB (≥6.0)
- Git (≥2.0)

## 故障排除

### 端口冲突
如果默认端口被占用，使用 `-port` 参数指定其他端口：
```bash
./installer -port 9999
```

### 权限问题
某些操作需要管理员权限，请以管理员身份运行。

### 网络问题
确保能够访问以下服务：
- GitHub (下载依赖)
- NPM Registry (npm包)
- Go Module Proxy (Go模块)

## 支持

如遇问题，请访问项目GitHub页面提交Issue。
EOF

echo "✅ 构建完成！"
echo ""
echo "📁 构建文件位置: dist/"
echo "🚀 安装程序已准备就绪！"
echo ""

# 显示构建结果
ls -la dist/ 