#!/bin/bash

# Vite Pluginend 快速安装脚本
# 使用方法: curl -fsSL https://raw.githubusercontent.com/360-vegas/vite-pluginend/main/quick-install.sh | bash

set -e

GITHUB_REPO="https://github.com/360-vegas/vite-pluginend.git"
INSTALL_SCRIPT_URL="https://raw.githubusercontent.com/360-vegas/vite-pluginend/main/install-from-github.sh"

echo "🚀 Vite Pluginend 快速安装"
echo "=========================="
echo "正在下载完整安装脚本..."

# 下载完整安装脚本
curl -fsSL "$INSTALL_SCRIPT_URL" -o /tmp/install-from-github.sh

# 给脚本执行权限
chmod +x /tmp/install-from-github.sh

# 执行安装脚本
echo "开始执行安装..."
/tmp/install-from-github.sh

# 清理临时文件
rm -f /tmp/install-from-github.sh

echo "✅ 快速安装完成！" 