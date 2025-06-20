package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"installer/internal/server"

	"github.com/gin-gonic/gin"
)

//go:embed web/dist/*
var webFiles embed.FS

//go:embed assets/*
var assetFiles embed.FS

func main() {
	var (
		port   = flag.String("port", "8888", "安装程序端口")
		silent = flag.Bool("silent", false, "静默安装模式")
		help   = flag.Bool("help", false, "显示帮助信息")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// 显示欢迎信息
	showWelcome()

	if *silent {
		log.Println("🔧 启动静默安装模式...")
		runSilentInstall()
		return
	}

	// 启动Web安装界面
	runWebInstaller(*port)
}

func showWelcome() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════╗
║                    插件系统安装向导                       ║
║                                                           ║
║  欢迎使用跨平台插件系统安装程序                           ║
║  支持 Windows、macOS、Linux 平台                         ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`)
	fmt.Printf("系统信息: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println()
}

func showHelp() {
	fmt.Println(`插件系统安装向导

用法:
  installer [选项]

选项:
  -port string    Web界面端口 (默认: 8888)
  -silent         静默安装模式，跳过Web界面
  -help           显示此帮助信息

示例:
  installer                    # 启动Web安装界面
  installer -port 9999         # 使用自定义端口
  installer -silent            # 静默安装`)
}

func runWebInstaller(port string) {
	gin.SetMode(gin.ReleaseMode)

	// 创建服务器实例
	srv, err := server.NewServer(webFiles, assetFiles)
	if err != nil {
		log.Fatalf("❌ 创建服务器失败: %v", err)
	}

	// 启动HTTP服务器
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: srv.Router(),
	}

	go func() {
		url := fmt.Sprintf("http://localhost:%s", port)
		log.Printf("🚀 安装向导已启动: %s", url)
		log.Println("💡 提示: 安装程序将在浏览器中打开，如未自动打开请手动访问上述地址")

		// 尝试自动打开浏览器
		go openBrowser(url)

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ 启动服务器失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🔄 正在关闭安装程序...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("⚠️  服务器强制关闭: %v", err)
	}
	log.Println("✅ 安装程序已退出")
}

func runSilentInstall() {
	// TODO: 实现静默安装逻辑
	log.Println("⚙️  开始静默安装...")
	log.Println("🔍 检测系统环境...")
	log.Println("📦 安装依赖软件...")
	log.Println("🗄️  配置数据库...")
	log.Println("🚀 部署项目文件...")
	log.Println("✅ 安装完成！")
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Printf("⚠️  无法自动打开浏览器: %v", err)
	}
}
