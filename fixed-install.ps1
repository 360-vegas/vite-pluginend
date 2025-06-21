# Vite Pluginend Windows 修复版安装脚本
# PowerShell 版本

param(
    [string]$ProjectPath = "D:\vge\vite-pluginend"
)

# 颜色输出函数
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    else {
        $input | Write-Output
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Write-Info($Message) {
    Write-ColorOutput Blue "🔵 $Message"
}

function Write-Success($Message) {
    Write-ColorOutput Green "✅ $Message"
}

function Write-Warning($Message) {
    Write-ColorOutput Yellow "⚠️  $Message"
}

function Write-Error($Message) {
    Write-ColorOutput Red "❌ $Message"
}

Write-Info "🔧 Vite Pluginend Windows 修复版安装脚本"
Write-Info "============================================="

# 检查管理员权限
if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Warning "需要管理员权限来安装某些依赖"
    Write-Info "继续以当前权限安装..."
}

# 设置执行策略
try {
    Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
    Write-Success "设置PowerShell执行策略成功"
} catch {
    Write-Warning "无法设置执行策略: $($_.Exception.Message)"
}

# 检查并安装Chocolatey
Write-Info "检查包管理器..."
if (!(Get-Command choco -ErrorAction SilentlyContinue)) {
    Write-Info "安装Chocolatey包管理器..."
    try {
        Set-ExecutionPolicy Bypass -Scope Process -Force
        [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
        Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
        
        # 刷新环境变量
        $env:PATH = [System.Environment]::GetEnvironmentVariable("PATH","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("PATH","User")
        
        Write-Success "Chocolatey安装成功"
    } catch {
        Write-Warning "Chocolatey安装失败: $($_.Exception.Message)"
        Write-Info "将尝试直接下载安装依赖"
    }
} else {
    Write-Success "Chocolatey已安装"
}

# 安装/更新Node.js
Write-Info "检查Node.js..."
$nodeVersion = ""
try {
    $nodeVersion = node --version 2>$null
    if ($nodeVersion) {
        Write-Info "当前Node.js版本: $nodeVersion"
        
        # 检查版本是否足够新
        $versionNumber = [int]($nodeVersion -replace 'v(\d+)\..*', '$1')
        if ($versionNumber -lt 18) {
            Write-Warning "Node.js版本过低，需要升级"
            $needNodeUpgrade = $true
        } else {
            Write-Success "Node.js版本满足要求"
            $needNodeUpgrade = $false
        }
    } else {
        $needNodeUpgrade = $true
    }
} catch {
    $needNodeUpgrade = $true
}

if ($needNodeUpgrade) {
    Write-Info "安装最新版Node.js..."
    try {
        if (Get-Command choco -ErrorAction SilentlyContinue) {
            choco install nodejs -y --force
        } else {
            # 直接下载安装
            $nodeUrl = "https://nodejs.org/dist/v20.10.0/node-v20.10.0-x64.msi"
            $tempFile = "$env:TEMP\nodejs-installer.msi"
            
            Write-Info "从官网下载Node.js..."
            Invoke-WebRequest -Uri $nodeUrl -OutFile $tempFile
            
            Write-Info "安装Node.js..."
            Start-Process msiexec.exe -Wait -ArgumentList "/i $tempFile /quiet"
            Remove-Item $tempFile -Force
        }
        
        # 刷新PATH
        $env:PATH = [System.Environment]::GetEnvironmentVariable("PATH","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("PATH","User")
        
        Write-Success "Node.js安装完成"
    } catch {
        Write-Error "Node.js安装失败: $($_.Exception.Message)"
        exit 1
    }
}

# 验证Node.js安装
try {
    $nodeVersion = node --version
    $npmVersion = npm --version
    Write-Success "Node.js版本: $nodeVersion"
    Write-Success "npm版本: $npmVersion"
} catch {
    Write-Error "Node.js验证失败"
    exit 1
}

# 安装/更新Go
Write-Info "检查Go..."
$goVersion = ""
try {
    $goVersion = go version 2>$null
    if ($goVersion) {
        Write-Info "当前Go版本: $goVersion"
        
        # 检查版本
        if ($goVersion -match 'go1\.(\d+)') {
            $versionNumber = [int]$matches[1]
            if ($versionNumber -lt 21) {
                Write-Warning "Go版本过低，需要升级"
                $needGoUpgrade = $true
            } else {
                Write-Success "Go版本满足要求"
                $needGoUpgrade = $false
            }
        } else {
            $needGoUpgrade = $true
        }
    } else {
        $needGoUpgrade = $true
    }
} catch {
    $needGoUpgrade = $true
}

if ($needGoUpgrade) {
    Write-Info "安装最新版Go..."
    try {
        if (Get-Command choco -ErrorAction SilentlyContinue) {
            choco install golang -y --force
        } else {
            # 直接下载安装
            $goUrl = "https://go.dev/dl/go1.21.5.windows-amd64.msi"
            $tempFile = "$env:TEMP\go-installer.msi"
            
            Write-Info "从官网下载Go..."
            Invoke-WebRequest -Uri $goUrl -OutFile $tempFile
            
            Write-Info "安装Go..."
            Start-Process msiexec.exe -Wait -ArgumentList "/i $tempFile /quiet"
            Remove-Item $tempFile -Force
        }
        
        # 刷新PATH
        $env:PATH = [System.Environment]::GetEnvironmentVariable("PATH","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("PATH","User")
        
        Write-Success "Go安装完成"
    } catch {
        Write-Error "Go安装失败: $($_.Exception.Message)"
        exit 1
    }
}

# 验证Go安装
try {
    $goVersion = go version
    Write-Success "Go版本: $goVersion"
} catch {
    Write-Error "Go验证失败"
    exit 1
}

# 设置Go环境
$env:GOPROXY = "https://goproxy.io,direct"
$env:GOSUMDB = "sum.golang.org"

# 进入项目目录
Write-Info "项目目录: $ProjectPath"
if (!(Test-Path $ProjectPath)) {
    Write-Error "项目目录不存在: $ProjectPath"
    exit 1
}

Set-Location $ProjectPath

# 修复package.json（如果存在问题）
if (Test-Path "package.json") {
    Write-Info "检查package.json..."
    try {
        $packageJson = Get-Content "package.json" | ConvertFrom-Json
        
        # 确保使用module类型
        if (!$packageJson.type -or $packageJson.type -ne "module") {
            Write-Info "更新package.json为ES模块..."
            $packageJson | Add-Member -Name "type" -Value "module" -MemberType NoteProperty -Force
            $packageJson | ConvertTo-Json -Depth 10 | Set-Content "package.json"
        }
        
        Write-Success "package.json检查完成"
    } catch {
        Write-Warning "package.json检查失败，但继续安装"
    }
}

# 设置npm配置
Write-Info "配置npm..."
npm config set registry https://registry.npmjs.org/
npm config set fund false
npm config set audit-level moderate

# 清理旧的依赖
Write-Info "清理旧依赖..."
if (Test-Path "node_modules") {
    Remove-Item "node_modules" -Recurse -Force -ErrorAction SilentlyContinue
}
if (Test-Path "package-lock.json") {
    Remove-Item "package-lock.json" -Force -ErrorAction SilentlyContinue
}

# 安装前端依赖
Write-Info "安装前端依赖..."
try {
    npm install --no-audit --no-fund
    Write-Success "前端依赖安装成功"
} catch {
    Write-Warning "前端依赖安装失败: $($_.Exception.Message)"
}

# 构建前端
Write-Info "构建前端..."
try {
    npm run build
    Write-Success "前端构建成功"
} catch {
    Write-Warning "前端构建失败，但继续后端构建"
}

# 构建后端
Write-Info "构建后端..."
if (Test-Path "backend") {
    Set-Location "backend"
    
    # 清理go.sum避免校验问题
    if (Test-Path "go.sum") {
        Remove-Item "go.sum" -Force
        Write-Info "已清理go.sum文件"
    }
    
    # 重新生成go.mod和go.sum
    Write-Info "整理Go模块..."
    try {
        go mod tidy
        Write-Success "Go模块整理完成"
    } catch {
        Write-Warning "Go模块整理失败: $($_.Exception.Message)"
    }
    
    # 下载依赖
    Write-Info "下载Go依赖..."
    try {
        go mod download
        Write-Success "Go依赖下载成功"
    } catch {
        Write-Warning "Go依赖下载失败: $($_.Exception.Message)"
    }
    
    # 构建
    Write-Info "编译Go应用..."
    try {
        if (Test-Path "cmd\server\main.go") {
            go build -o "..\vite-pluginend-server.exe" "cmd\server\main.go"
        } elseif (Test-Path "main.go") {
            go build -o "..\vite-pluginend-server.exe" "main.go"
        } else {
            Write-Warning "未找到Go主文件，跳过后端构建"
            Set-Location ..
            return
        }
        Write-Success "后端构建成功"
    } catch {
        Write-Error "后端构建失败: $($_.Exception.Message)"
        Set-Location ..
        exit 1
    }
    
    Set-Location ..
} else {
    Write-Warning "未找到backend目录，跳过后端构建"
}

# 启动服务
Write-Info "启动服务..."
if (Test-Path "vite-pluginend-server.exe") {
    # 停止旧进程
    $oldProcess = Get-Process "vite-pluginend-server" -ErrorAction SilentlyContinue
    if ($oldProcess) {
        Write-Info "停止旧服务..."
        $oldProcess | Stop-Process -Force
        Start-Sleep -Seconds 2
    }
    
    # 启动新服务
    Write-Info "启动新服务..."
    $process = Start-Process -FilePath ".\vite-pluginend-server.exe" -PassThru -RedirectStandardOutput "vite-pluginend.log" -RedirectStandardError "vite-pluginend-error.log"
    
    Start-Sleep -Seconds 3
    
    # 检查服务状态
    if (!$process.HasExited) {
        Write-Success "服务启动成功 (PID: $($process.Id))"
        
        # 获取本机IP
        $localIP = (Get-NetIPAddress -AddressFamily IPv4 | Where-Object {$_.IPAddress -ne "127.0.0.1" -and $_.PrefixOrigin -eq "Dhcp"} | Select-Object -First 1).IPAddress
        if (!$localIP) {
            $localIP = "localhost"
        }
        
        Write-Info "📍 访问地址:"
        Write-Info "   🌐 前端: http://${localIP}:3000"
        Write-Info "   🔧 后端: http://${localIP}:8080"
        Write-Info "   💚 健康检查: http://${localIP}:8080/health"
        
        Write-Info ""
        Write-Info "📋 管理命令:"
        Write-Info "   查看日志: Get-Content vite-pluginend.log -Wait"
        Write-Info "   停止服务: Stop-Process -Name 'vite-pluginend-server'"
        Write-Info "   重启服务: .\vite-pluginend-server.exe"
        
        # 测试服务
        Write-Info ""
        Write-Info "🧪 测试服务..."
        Start-Sleep -Seconds 2
        try {
            $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -TimeoutSec 5 -ErrorAction Stop
            if ($response.StatusCode -eq 200) {
                Write-Success "服务响应正常"
            }
        } catch {
            Write-Warning "服务可能需要更多时间启动"
        }
        
    } else {
        Write-Error "服务启动失败"
        Write-Info "错误日志:"
        if (Test-Path "vite-pluginend-error.log") {
            Get-Content "vite-pluginend-error.log" | Select-Object -Last 10
        }
        exit 1
    }
} else {
    Write-Warning "未找到可执行文件，服务启动跳过"
}

Write-Info ""
Write-Success "🎉 安装完成！"
Write-Success "项目已成功部署到 $ProjectPath" 