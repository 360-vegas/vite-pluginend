# Vite Pluginend 项目打包工具 (Windows PowerShell 版本)
param(
    [string]$Target = "all",
    [string]$Version = "",
    [switch]$Clean,
    [switch]$NoDeps,
    [switch]$Help
)

# 项目信息
$ProjectName = "vite-pluginend"
$BuildDir = "./build"
$DistDir = "./dist"
$ReleaseDir = "./releases"

if (-not $Version) {
    $Version = Get-Date -Format "yyyyMMdd-HHmmss"
}

# 显示帮助信息
function Show-Help {
    Write-Host "🎯 Vite Pluginend 项目打包工具" -ForegroundColor Blue
    Write-Host "=================================" -ForegroundColor Blue
    Write-Host ""
    Write-Host "使用方法: .\package.ps1 [选项]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "目标:" -ForegroundColor Yellow
    Write-Host "  -Target frontend    仅打包前端项目"
    Write-Host "  -Target backend     仅打包后端项目"
    Write-Host "  -Target installer   仅打包安装向导"
    Write-Host "  -Target all         打包完整项目（默认）"
    Write-Host "  -Target release     创建发布版本"
    Write-Host ""
    Write-Host "选项:" -ForegroundColor Yellow
    Write-Host "  -Version <版本号>   指定版本号"
    Write-Host "  -Clean              构建前清理输出目录"
    Write-Host "  -NoDeps             跳过依赖安装"
    Write-Host "  -Help               显示帮助信息"
    Write-Host ""
    Write-Host "示例:" -ForegroundColor Yellow
    Write-Host "  .\package.ps1 -Target frontend"
    Write-Host "  .\package.ps1 -Target all -Clean"
    Write-Host "  .\package.ps1 -Target release -Version 1.2.0"
}

if ($Help) {
    Show-Help
    exit 0
}

# 检查依赖
function Test-Dependencies {
    Write-Host "🔍 检查构建依赖..." -ForegroundColor Blue
    
    $missingDeps = @()
    
    if (-not (Get-Command node -ErrorAction SilentlyContinue)) {
        $missingDeps += "Node.js"
    }
    
    if (-not (Get-Command npm -ErrorAction SilentlyContinue)) {
        $missingDeps += "npm"
    }
    
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        $missingDeps += "Go"
    }
    
    if ($missingDeps.Count -gt 0) {
        Write-Host "❌ 缺少以下依赖: $($missingDeps -join ', ')" -ForegroundColor Red
        Write-Host "请先安装缺失的依赖后再继续" -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host "✅ 所有依赖检查通过" -ForegroundColor Green
}

# 清理构建目录
function Clear-BuildDirs {
    Write-Host "🧹 清理构建目录..." -ForegroundColor Blue
    
    if (Test-Path $BuildDir) { Remove-Item $BuildDir -Recurse -Force }
    if (Test-Path $DistDir) { Remove-Item $DistDir -Recurse -Force }
    if (Test-Path $ReleaseDir) { Remove-Item $ReleaseDir -Recurse -Force }
    
    # 清理各模块的构建产物
    if (Test-Path "./node_modules/.vite") { Remove-Item "./node_modules/.vite" -Recurse -Force }
    if (Test-Path "./backend/dist") { Remove-Item "./backend/dist" -Recurse -Force }
    if (Test-Path "./installer/dist") { Remove-Item "./installer/dist" -Recurse -Force }
    if (Test-Path "./installer/web/dist") { Remove-Item "./installer/web/dist" -Recurse -Force }
    
    Write-Host "✅ 构建目录清理完成" -ForegroundColor Green
}

# 安装依赖
function Install-Dependencies {
    if ($NoDeps) {
        Write-Host "⏭️  跳过依赖安装" -ForegroundColor Yellow
        return
    }
    
    Write-Host "📦 安装项目依赖..." -ForegroundColor Blue
    
    # 安装主项目依赖
    Write-Host "安装前端依赖..." -ForegroundColor Blue
    npm install
    if ($LASTEXITCODE -ne 0) { throw "前端依赖安装失败" }
    
    # 安装后端依赖
    Write-Host "安装后端依赖..." -ForegroundColor Blue
    Push-Location backend
    go mod download
    if ($LASTEXITCODE -ne 0) { throw "后端依赖安装失败" }
    Pop-Location
    
    # 安装安装向导依赖
    Write-Host "安装安装向导依赖..." -ForegroundColor Blue
    Push-Location installer
    go mod download
    if ($LASTEXITCODE -ne 0) { throw "安装向导Go依赖安装失败" }
    
    Push-Location web
    npm install
    if ($LASTEXITCODE -ne 0) { throw "安装向导前端依赖安装失败" }
    Pop-Location
    Pop-Location
    
    Write-Host "✅ 所有依赖安装完成" -ForegroundColor Green
}

# 打包前端项目
function Build-Frontend {
    Write-Host "🎨 构建前端项目..." -ForegroundColor Blue
    
    # 创建构建目录
    New-Item -Path "$BuildDir/frontend" -ItemType Directory -Force | Out-Null
    
    # 生成版本信息
    node scripts/version-manager.js build
    if ($LASTEXITCODE -ne 0) { throw "版本信息生成失败" }
    
    # 构建前端
    npm run build
    if ($LASTEXITCODE -ne 0) { throw "前端构建失败" }
    
    # 复制构建产物
    Copy-Item dist/* "$BuildDir/frontend/" -Recurse
    
    # 创建压缩包
    $frontendPath = "$BuildDir/frontend"
    Compress-Archive -Path $frontendPath -DestinationPath "$BuildDir/frontend-$Version.zip" -Force
    
    Write-Host "✅ 前端项目构建完成" -ForegroundColor Green
    Write-Host "   输出位置: $BuildDir/frontend-$Version.zip" -ForegroundColor Yellow
}

# 打包后端项目
function Build-Backend {
    Write-Host "⚙️  构建后端项目..." -ForegroundColor Blue
    
    # 创建构建目录
    New-Item -Path "$BuildDir/backend" -ItemType Directory -Force | Out-Null
    
    Push-Location backend
    
    # 构建不同平台的二进制文件
    Write-Host "构建 Linux x64..." -ForegroundColor Blue
    $env:GOOS = "linux"; $env:GOARCH = "amd64"
    go build -ldflags="-s -w -X main.version=$Version" -o "../$BuildDir/backend/vite-pluginend-linux-amd64" cmd/server/main.go
    if ($LASTEXITCODE -ne 0) { throw "Linux版本构建失败" }
    
    Write-Host "构建 Windows x64..." -ForegroundColor Blue
    $env:GOOS = "windows"; $env:GOARCH = "amd64"
    go build -ldflags="-s -w -X main.version=$Version" -o "../$BuildDir/backend/vite-pluginend-windows-amd64.exe" cmd/server/main.go
    if ($LASTEXITCODE -ne 0) { throw "Windows版本构建失败" }
    
    Write-Host "构建 macOS x64..." -ForegroundColor Blue
    $env:GOOS = "darwin"; $env:GOARCH = "amd64"
    go build -ldflags="-s -w -X main.version=$Version" -o "../$BuildDir/backend/vite-pluginend-darwin-amd64" cmd/server/main.go
    if ($LASTEXITCODE -ne 0) { throw "macOS x64版本构建失败" }
    
    Write-Host "构建 macOS ARM64..." -ForegroundColor Blue
    $env:GOOS = "darwin"; $env:GOARCH = "arm64"
    go build -ldflags="-s -w -X main.version=$Version" -o "../$BuildDir/backend/vite-pluginend-darwin-arm64" cmd/server/main.go
    if ($LASTEXITCODE -ne 0) { throw "macOS ARM64版本构建失败" }
    
    # 重置环境变量
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    
    Pop-Location
    
    # 复制配置文件和文档
    Copy-Item backend/README.md "$BuildDir/backend/"
    Copy-Item backend/go.mod "$BuildDir/backend/"
    
    # 创建Windows启动脚本
    @'
@echo off
echo 启动 Vite Pluginend 服务器...

REM 检测架构
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set ARCH=amd64
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set ARCH=arm64
) else (
    echo 不支持的架构: %PROCESSOR_ARCHITECTURE%
    pause
    exit /b 1
)

set BINARY=vite-pluginend-windows-%ARCH%.exe

if not exist "%BINARY%" (
    echo 找不到适用于 windows-%ARCH% 的二进制文件
    pause
    exit /b 1
)

"%BINARY%" %*
'@ | Out-File -FilePath "$BuildDir/backend/start.bat" -Encoding ASCII
    
    # 创建压缩包
    $backendPath = "$BuildDir/backend"
    Compress-Archive -Path $backendPath -DestinationPath "$BuildDir/backend-$Version.zip" -Force
    
    Write-Host "✅ 后端项目构建完成" -ForegroundColor Green
    Write-Host "   输出位置: $BuildDir/backend-$Version.zip" -ForegroundColor Yellow
}

# 打包安装向导
function Build-Installer {
    Write-Host "🛠️  构建安装向导..." -ForegroundColor Blue
    
    # 使用安装向导自己的构建脚本 (Windows版本)
    Push-Location installer
    
    # 手动执行构建步骤（因为PowerShell可能无法直接运行.sh脚本）
    Write-Host "构建安装向导前端..." -ForegroundColor Blue
    Push-Location web
    npm install
    npm run build
    Pop-Location
    
    Write-Host "构建安装向导后端..." -ForegroundColor Blue
    New-Item -Path "dist" -ItemType Directory -Force | Out-Null
    
    go mod tidy
    
    # 构建不同平台的安装程序
    $env:GOOS = "windows"; $env:GOARCH = "amd64"
    go build -ldflags="-s -w" -o "dist/installer-windows-amd64.exe" cmd/installer/main.go
    
    $env:GOOS = "darwin"; $env:GOARCH = "amd64"
    go build -ldflags="-s -w" -o "dist/installer-macos-amd64" cmd/installer/main.go
    
    $env:GOOS = "darwin"; $env:GOARCH = "arm64"
    go build -ldflags="-s -w" -o "dist/installer-macos-arm64" cmd/installer/main.go
    
    $env:GOOS = "linux"; $env:GOARCH = "amd64"
    go build -ldflags="-s -w" -o "dist/installer-linux-amd64" cmd/installer/main.go
    
    # 重置环境变量
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    
    Pop-Location
    
    # 复制到统一的构建目录
    New-Item -Path "$BuildDir/installer" -ItemType Directory -Force | Out-Null
    Copy-Item installer/dist/* "$BuildDir/installer/" -Recurse
    
    # 创建压缩包
    $installerPath = "$BuildDir/installer"
    Compress-Archive -Path $installerPath -DestinationPath "$BuildDir/installer-$Version.zip" -Force
    
    Write-Host "✅ 安装向导构建完成" -ForegroundColor Green
    Write-Host "   输出位置: $BuildDir/installer-$Version.zip" -ForegroundColor Yellow
}

# 创建完整发布包
function New-Release {
    Write-Host "📦 创建发布版本..." -ForegroundColor Blue
    
    # 创建发布目录
    New-Item -Path "$ReleaseDir/v$Version" -ItemType Directory -Force | Out-Null
    
    # 复制所有构建产物
    Copy-Item "$BuildDir/*" "$ReleaseDir/v$Version/" -Recurse
    
    # 创建发布说明
    $releaseNotes = @"
# Vite Pluginend v$Version

## 📦 发布内容

### 前端应用
- ``frontend-$Version.zip`` - 前端应用包

### 后端服务
- ``backend-$Version.zip`` - 跨平台后端包

包含的二进制文件：
- ``vite-pluginend-linux-amd64`` - Linux x64
- ``vite-pluginend-windows-amd64.exe`` - Windows x64  
- ``vite-pluginend-darwin-amd64`` - macOS x64 (Intel)
- ``vite-pluginend-darwin-arm64`` - macOS ARM64 (Apple Silicon)

### 安装向导
- ``installer-$Version.zip`` - 跨平台安装向导

包含的安装程序：
- ``installer-linux-amd64`` - Linux x64 安装向导
- ``installer-windows-amd64.exe`` - Windows x64 安装向导
- ``installer-macos-amd64`` - macOS x64 安装向导
- ``installer-macos-arm64`` - macOS ARM64 安装向导

## 🚀 快速开始

### 使用安装向导（推荐）
``````
# 解压安装向导
Expand-Archive installer-$Version.zip

# 运行安装向导
cd installer
.\installer-windows-amd64.exe
``````

### 手动部署
``````
# 1. 部署前端
Expand-Archive frontend-$Version.zip
# 将 frontend/ 目录内容部署到 Web 服务器

# 2. 部署后端
Expand-Archive backend-$Version.zip
cd backend
.\start.bat
``````

## 📋 系统要求

- **内存**: 最少 2GB RAM
- **磁盘**: 最少 5GB 可用空间
- **网络**: 需要网络连接用于下载依赖

### 软件依赖
- Node.js >= 18.0
- Go >= 1.21 (仅开发环境)
- MySQL >= 8.0
- MongoDB >= 6.0

---
构建时间: $(Get-Date)
构建版本: $Version
"@
    
    $releaseNotes | Out-File -FilePath "$ReleaseDir/v$Version/RELEASE_NOTES.md" -Encoding UTF8
    
    # 创建完整发布包
    $releasePath = "$ReleaseDir/v$Version"
    Compress-Archive -Path $releasePath -DestinationPath "$ReleaseDir/vite-pluginend-complete-$Version.zip" -Force
    
    Write-Host "✅ 发布版本创建完成" -ForegroundColor Green
    Write-Host "   发布目录: $ReleaseDir/v$Version/" -ForegroundColor Yellow
    Write-Host "   完整包: $ReleaseDir/vite-pluginend-complete-$Version.zip" -ForegroundColor Yellow
}

# 显示构建结果
function Show-Results {
    Write-Host ""
    Write-Host "🎉 打包完成！" -ForegroundColor Green
    Write-Host "=================================" -ForegroundColor Blue
    
    if (Test-Path $BuildDir) {
        Write-Host "📁 构建产物:" -ForegroundColor Yellow
        Get-ChildItem $BuildDir -Filter "*.zip" | ForEach-Object {
            $size = [math]::Round($_.Length / 1MB, 2)
            Write-Host "   ✓ $($_.Name) ($size MB)" -ForegroundColor Green
        }
    }
    
    if (Test-Path $ReleaseDir) {
        Write-Host ""
        Write-Host "📦 发布文件:" -ForegroundColor Yellow
        Get-ChildItem $ReleaseDir -Recurse -Filter "*.zip" | ForEach-Object {
            $size = [math]::Round($_.Length / 1MB, 2)
            Write-Host "   ✓ $($_.Name) ($size MB)" -ForegroundColor Green
        }
    }
    
    Write-Host ""
    Write-Host "📖 使用说明:" -ForegroundColor Blue
    Write-Host "   前端部署: 解压 frontend-$Version.zip 到 Web 服务器" -ForegroundColor Yellow
    Write-Host "   后端部署: 解压 backend-$Version.zip 并运行 start.bat" -ForegroundColor Yellow
    Write-Host "   一键安装: 使用 installer-$Version.zip 中的安装向导" -ForegroundColor Yellow
}

# 主执行流程
function main {
    try {
        Write-Host "🎯 开始打包 $Target (版本: $Version)..." -ForegroundColor Blue
        
        # 检查是否在项目根目录
        if (-not (Test-Path "package.json") -or -not (Test-Path "backend")) {
            Write-Host "❌ 请在项目根目录中运行此脚本" -ForegroundColor Red
            exit 1
        }
        
        # 检查依赖
        Test-Dependencies
        
        # 清理构建目录
        if ($Clean) {
            Clear-BuildDirs
        }
        
        # 创建构建目录
        New-Item -Path $BuildDir -ItemType Directory -Force | Out-Null
        
        # 安装依赖
        Install-Dependencies
        
        # 根据目标执行相应的构建
        switch ($Target) {
            "frontend" { Build-Frontend }
            "backend" { Build-Backend }
            "installer" { Build-Installer }
            "all" {
                Build-Frontend
                Build-Backend
                Build-Installer
            }
            "release" {
                Build-Frontend
                Build-Backend
                Build-Installer
                New-Release
            }
            default {
                Write-Host "❌ 未知的构建目标: $Target" -ForegroundColor Red
                exit 1
            }
        }
        
        # 显示结果
        Show-Results
        
    } catch {
        Write-Host "❌ 构建失败: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
}

# 运行主函数
main 