package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vite-pluginend/internal/api/handlers"
	"vite-pluginend/internal/api/middleware"
	"vite-pluginend/internal/plugins"
	"vite-pluginend/internal/services"
	"vite-pluginend/pkg/cache"
	"vite-pluginend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// @title Vite Plugin End API
// @version 1.0
// @description Vite插件管理系统的API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请在值前加上"Bearer "前缀，例如："Bearer abcde12345"

func main() {
	// 设置 Gin 模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化日志
	log := logger.NewLogger()
	defer log.Sync()

	// 连接 MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("连接 MongoDB 失败", zap.Error(err))
	}
	defer client.Disconnect(context.Background())

	// 检查连接
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("MongoDB 连接测试失败", zap.Error(err))
	}

	db := client.Database("vga")

	// 初始化 Redis
	var redisCache *cache.RedisCache
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisCache, err = cache.NewRedisCache(redisAddr, "", 0)
	if err != nil {
		log.Warn("Redis 连接失败，将使用内存缓存", zap.Error(err))
		// 使用内存缓存作为后备
		redisCache = cache.NewMemoryCache()
	}
	defer redisCache.Close()

	// 初始化服务
	userService := services.NewUserService(db, redisCache)
	pluginService := services.NewPluginService(db, redisCache)
	uploadService := services.NewUploadService("uploads")
	dependencyService := services.NewDependencyService(client)

	// 初始化插件管理器
	pluginManager := plugins.NewManager(db, redisCache)

	// 初始化处理器
	userHandler := handlers.NewUserHandler(userService)
	pluginHandler := handlers.NewPluginHandler(pluginService, dependencyService)
	uploadHandler := handlers.NewUploadHandler(uploadService)

	// 创建 Gin 引擎
	r := gin.New() // 使用 New() 而不是 Default()，避免重复的中间件

	// 中间件
	r.Use(middleware.Cors())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 临时测试路由，用于验证服务器是否正常响应
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// API 路由
	api := r.Group("/api")
	{
		// 用户相关路由
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.GET("/users/:id", userHandler.GetUser)
		api.PUT("/users/:id", userHandler.UpdateUser)

		// 插件相关路由
		api.POST("/plugins", pluginHandler.CreatePlugin)
		api.GET("/plugins/:id", pluginHandler.GetPlugin)
		api.GET("/plugins", pluginHandler.ListPlugins)
		api.PUT("/plugins/:id", pluginHandler.UpdatePlugin)
		api.DELETE("/plugins/:id", pluginHandler.DeletePlugin)
		api.POST("/plugins/scan", pluginHandler.ScanPlugins)
		api.POST("/plugins/package", pluginHandler.PackagePlugin)
		api.GET("/plugins/download/:name", pluginHandler.DownloadPlugin)
		api.POST("/create-plugin", pluginHandler.GeneratePlugin)

		// 新增插件管理路由
		api.POST("/plugins/:id/toggle", pluginHandler.TogglePlugin)
		api.GET("/plugins/:id/export", pluginHandler.ExportPlugin)
		api.POST("/plugins/install", pluginHandler.InstallPlugin)

		// 插件依赖检查路由 - 使用 :id 而不是 :key 来避免冲突
		api.GET("/plugins/:id/dependencies/check", pluginHandler.CheckPluginDependencies)
		api.POST("/plugins/:id/dependencies/setup", pluginHandler.SetupPluginDatabase)

		// 文件上传相关路由
		api.POST("/upload", uploadHandler.UploadFile)
		api.GET("/files/:filename", uploadHandler.GetFile)
		api.DELETE("/files/:filename", uploadHandler.DeleteFile)

		// 注册插件路由
		pluginManager.RegisterRoutes(api)
	}

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("启动服务器失败", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器关闭失败", zap.Error(err))
	}
}
