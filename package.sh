#!/bin/bash

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# 项目信息
PROJECT_NAME="vite-pluginend"
VERSION=${VERSION:-$(date +%Y%m%d-%H%M%S)}
BUILD_DIR="./build"
DIST_DIR="./dist"
RELEASE_DIR="./releases"

echo -e "${BLUE}🎯 Vite Pluginend 项目打包工具${NC}"
echo -e "${BLUE}=================================${NC}"

# 显示帮助信息
show_help() {
    echo -e "${YELLOW}使用方法: $0 [选项] [目标]${NC}"
    echo ""
    echo -e "${YELLOW}目标:${NC}"
    echo "  frontend         - 仅打包前端项目"
    echo "  backend          - 仅打包后端项目"
    echo "  installer        - 仅打包安装向导"
    echo "  all             - 打包完整项目（默认）"
    echo "  release         - 创建发布版本"
    echo ""
    echo -e "${YELLOW}选项:${NC}"
    echo "  -v, --version    - 指定版本号"
    echo "  -h, --help      - 显示帮助信息"
    echo "  --clean         - 构建前清理输出目录"
    echo "  --no-deps       - 跳过依赖安装"
    echo ""
    echo -e "${YELLOW}示例:${NC}"
    echo "  $0 frontend                    # 仅打包前端"
    echo "  $0 all --clean                # 完整打包并清理"
    echo "  $0 release -v 1.2.0           # 创建1.2.0发布版本"
}

# 检查依赖
check_dependencies() {
    echo -e "${BLUE}🔍 检查构建依赖...${NC}"
    
    local missing_deps=()
    
    if ! command -v node &> /dev/null; then
        missing_deps+=("Node.js")
    fi
    
    if ! command -v npm &> /dev/null; then
        missing_deps+=("npm")
    fi
    
    if ! command -v go &> /dev/null; then
        missing_deps+=("Go")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        echo -e "${RED}❌ 缺少以下依赖: ${missing_deps[*]}${NC}"
        echo -e "${YELLOW}请先安装缺失的依赖后再继续${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ 所有依赖检查通过${NC}"
}

# 清理构建目录
clean_build_dirs() {
    echo -e "${BLUE}🧹 清理构建目录...${NC}"
    
    rm -rf "$BUILD_DIR"
    rm -rf "$DIST_DIR"
    rm -rf "$RELEASE_DIR"
    
    # 清理各模块的构建产物
    rm -rf "./node_modules/.vite"
    rm -rf "./backend/dist"
    rm -rf "./installer/dist"
    rm -rf "./installer/web/dist"
    
    echo -e "${GREEN}✅ 构建目录清理完成${NC}"
}

# 安装依赖
install_dependencies() {
    if [ "$SKIP_DEPS" = "true" ]; then
        echo -e "${YELLOW}⏭️  跳过依赖安装${NC}"
        return
    fi
    
    echo -e "${BLUE}📦 安装项目依赖...${NC}"
    
    # 安装主项目依赖
    echo -e "${BLUE}安装前端依赖...${NC}"
    npm install
    
    # 安装后端依赖
    echo -e "${BLUE}安装后端依赖...${NC}"
    cd backend
    go mod download
    cd ..
    
    # 安装安装向导依赖
    echo -e "${BLUE}安装安装向导依赖...${NC}"
    cd installer
    go mod download
    cd web
    npm install
    cd ../..
    
    echo -e "${GREEN}✅ 所有依赖安装完成${NC}"
}

# 打包前端项目
build_frontend() {
    echo -e "${BLUE}🎨 构建前端项目...${NC}"
    
    # 创建构建目录
    mkdir -p "$BUILD_DIR/frontend"
    
    # 生成版本信息
    node scripts/version-manager.js build
    
    # 构建前端
    npm run build
    
    # 复制构建产物
    cp -r dist/* "$BUILD_DIR/frontend/"
    
    # 创建前端打包文件
    cd "$BUILD_DIR"
    tar -czf "frontend-${VERSION}.tar.gz" frontend/
    zip -r "frontend-${VERSION}.zip" frontend/
    cd ..
    
    echo -e "${GREEN}✅ 前端项目构建完成${NC}"
    echo -e "${YELLOW}   输出位置: ${BUILD_DIR}/frontend-${VERSION}.*${NC}"
}

# 打包后端项目
build_backend() {
    echo -e "${BLUE}⚙️  构建后端项目...${NC}"
    
    # 创建构建目录
    mkdir -p "$BUILD_DIR/backend"
    
    cd backend
    
    # 构建不同平台的二进制文件
    echo -e "${BLUE}构建 Linux x64...${NC}"
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o "../$BUILD_DIR/backend/vite-pluginend-linux-amd64" cmd/server/main.go
    
    echo -e "${BLUE}构建 Windows x64...${NC}"
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o "../$BUILD_DIR/backend/vite-pluginend-windows-amd64.exe" cmd/server/main.go
    
    echo -e "${BLUE}构建 macOS x64...${NC}"
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o "../$BUILD_DIR/backend/vite-pluginend-darwin-amd64" cmd/server/main.go
    
    echo -e "${BLUE}构建 macOS ARM64...${NC}"
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=${VERSION}" -o "../$BUILD_DIR/backend/vite-pluginend-darwin-arm64" cmd/server/main.go
    
    cd ..
    
    # 复制配置文件和文档
    cp backend/README.md "$BUILD_DIR/backend/"
    cp backend/go.mod "$BUILD_DIR/backend/"
    
    # 创建启动脚本
    cat > "$BUILD_DIR/backend/start.sh" << 'EOF'
#!/bin/bash
# 启动脚本

PLATFORM=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "不支持的架构: $ARCH"; exit 1 ;;
esac

BINARY="vite-pluginend-${PLATFORM}-${ARCH}"
if [ "$PLATFORM" = "windows" ]; then
    BINARY="${BINARY}.exe"
fi

if [ ! -f "$BINARY" ]; then
    echo "找不到适用于 ${PLATFORM}-${ARCH} 的二进制文件"
    exit 1
fi

echo "启动 Vite Pluginend 服务器..."
chmod +x "$BINARY"
./"$BINARY" "$@"
EOF
    
    chmod +x "$BUILD_DIR/backend/start.sh"
    
    # 创建后端打包文件
    cd "$BUILD_DIR"
    tar -czf "backend-${VERSION}.tar.gz" backend/
    zip -r "backend-${VERSION}.zip" backend/
    cd ..
    
    echo -e "${GREEN}✅ 后端项目构建完成${NC}"
    echo -e "${YELLOW}   输出位置: ${BUILD_DIR}/backend-${VERSION}.*${NC}"
}

# 打包安装向导
build_installer() {
    echo -e "${BLUE}🛠️  构建安装向导...${NC}"
    
    # 使用安装向导自己的构建脚本
    cd installer
    chmod +x scripts/build.sh
    ./scripts/build.sh
    cd ..
    
    # 复制到统一的构建目录
    mkdir -p "$BUILD_DIR/installer"
    cp -r installer/dist/* "$BUILD_DIR/installer/"
    
    # 创建安装向导打包文件
    cd "$BUILD_DIR"
    tar -czf "installer-${VERSION}.tar.gz" installer/
    zip -r "installer-${VERSION}.zip" installer/
    cd ..
    
    echo -e "${GREEN}✅ 安装向导构建完成${NC}"
    echo -e "${YELLOW}   输出位置: ${BUILD_DIR}/installer-${VERSION}.*${NC}"
}

# 创建完整发布包
create_release() {
    echo -e "${BLUE}📦 创建发布版本...${NC}"
    
    # 创建发布目录
    mkdir -p "$RELEASE_DIR/v${VERSION}"
    
    # 复制所有构建产物
    cp -r "$BUILD_DIR"/* "$RELEASE_DIR/v${VERSION}/"
    
    # 创建发布说明
    cat > "$RELEASE_DIR/v${VERSION}/RELEASE_NOTES.md" << EOF
# Vite Pluginend v${VERSION}

## 📦 发布内容

### 前端应用
- \`frontend-${VERSION}.tar.gz\` - Linux/macOS 前端包
- \`frontend-${VERSION}.zip\` - Windows 前端包

### 后端服务
- \`backend-${VERSION}.tar.gz\` - 跨平台后端包 (Linux/macOS)
- \`backend-${VERSION}.zip\` - 跨平台后端包 (Windows)

包含的二进制文件：
- \`vite-pluginend-linux-amd64\` - Linux x64
- \`vite-pluginend-windows-amd64.exe\` - Windows x64  
- \`vite-pluginend-darwin-amd64\` - macOS x64 (Intel)
- \`vite-pluginend-darwin-arm64\` - macOS ARM64 (Apple Silicon)

### 安装向导
- \`installer-${VERSION}.tar.gz\` - 跨平台安装向导 (Linux/macOS)
- \`installer-${VERSION}.zip\` - 跨平台安装向导 (Windows)

包含的安装程序：
- \`installer-linux-amd64\` - Linux x64 安装向导
- \`installer-windows-amd64.exe\` - Windows x64 安装向导
- \`installer-macos-amd64\` - macOS x64 安装向导
- \`installer-macos-arm64\` - macOS ARM64 安装向导

## 🚀 快速开始

### 使用安装向导（推荐）
\`\`\`bash
# 下载并解压安装向导
tar -xzf installer-${VERSION}.tar.gz
cd installer

# 运行安装向导
chmod +x installer-linux-amd64
./installer-linux-amd64
\`\`\`

### 手动部署
\`\`\`bash
# 1. 部署前端
tar -xzf frontend-${VERSION}.tar.gz
# 将 frontend/ 目录内容部署到 Web 服务器

# 2. 部署后端
tar -xzf backend-${VERSION}.tar.gz
cd backend
./start.sh
\`\`\`

## 📋 系统要求

- **内存**: 最少 2GB RAM
- **磁盘**: 最少 5GB 可用空间
- **网络**: 需要网络连接用于下载依赖

### 软件依赖
- Node.js >= 18.0
- Go >= 1.21 (仅开发环境)
- MySQL >= 8.0
- MongoDB >= 6.0

## 🔧 配置

详细配置说明请参考各组件的 README 文件。

## 📝 更新日志

$(cd .. && node scripts/version-manager.js changelog || echo "无更新记录")

---
构建时间: $(date)
构建版本: ${VERSION}
EOF
    
    # 创建完整发布包
    cd "$RELEASE_DIR"
    tar -czf "vite-pluginend-complete-${VERSION}.tar.gz" "v${VERSION}/"
    zip -r "vite-pluginend-complete-${VERSION}.zip" "v${VERSION}/"
    cd ..
    
    echo -e "${GREEN}✅ 发布版本创建完成${NC}"
    echo -e "${YELLOW}   发布目录: ${RELEASE_DIR}/v${VERSION}/${NC}"
    echo -e "${YELLOW}   完整包: ${RELEASE_DIR}/vite-pluginend-complete-${VERSION}.*${NC}"
}

# 显示构建结果
show_results() {
    echo -e "\n${GREEN}🎉 打包完成！${NC}"
    echo -e "${BLUE}=================================${NC}"
    
    if [ -d "$BUILD_DIR" ]; then
        echo -e "${YELLOW}📁 构建产物:${NC}"
        ls -la "$BUILD_DIR"/ | grep -E '\.(tar\.gz|zip)$' | while read -r line; do
            echo -e "   ${GREEN}✓${NC} $line"
        done
    fi
    
    if [ -d "$RELEASE_DIR" ]; then
        echo -e "\n${YELLOW}📦 发布文件:${NC}"
        find "$RELEASE_DIR" -name "*.tar.gz" -o -name "*.zip" | while read -r file; do
            size=$(du -h "$file" | cut -f1)
            echo -e "   ${GREEN}✓${NC} $(basename "$file") (${size})"
        done
    fi
    
    echo -e "\n${BLUE}📖 使用说明:${NC}"
    echo -e "   ${YELLOW}前端部署${NC}: 解压 frontend-${VERSION}.* 到 Web 服务器"
    echo -e "   ${YELLOW}后端部署${NC}: 解压 backend-${VERSION}.* 并运行 start.sh"
    echo -e "   ${YELLOW}一键安装${NC}: 使用 installer-${VERSION}.* 中的安装向导"
}

# 解析命令行参数
CLEAN=false
SKIP_DEPS=false
TARGET="all"

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        --clean)
            CLEAN=true
            shift
            ;;
        --no-deps)
            SKIP_DEPS=true
            shift
            ;;
        frontend|backend|installer|all|release)
            TARGET="$1"
            shift
            ;;
        *)
            echo -e "${RED}未知选项: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# 主执行流程
main() {
    echo -e "${BLUE}开始打包 ${TARGET} (版本: ${VERSION})...${NC}"
    
    # 检查依赖
    check_dependencies
    
    # 清理构建目录
    if [ "$CLEAN" = "true" ]; then
        clean_build_dirs
    fi
    
    # 创建构建目录
    mkdir -p "$BUILD_DIR"
    
    # 安装依赖
    install_dependencies
    
    # 根据目标执行相应的构建
    case $TARGET in
        frontend)
            build_frontend
            ;;
        backend)
            build_backend
            ;;
        installer)
            build_installer
            ;;
        all)
            build_frontend
            build_backend
            build_installer
            ;;
        release)
            build_frontend
            build_backend
            build_installer
            create_release
            ;;
        *)
            echo -e "${RED}未知的构建目标: $TARGET${NC}"
            exit 1
            ;;
    esac
    
    # 显示结果
    show_results
}

# 检查是否在项目根目录
if [ ! -f "package.json" ] || [ ! -d "backend" ]; then
    echo -e "${RED}❌ 请在项目根目录中运行此脚本${NC}"
    exit 1
fi

# 运行主函数
main "$@" 