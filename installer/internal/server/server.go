package server

import (
	"embed"
	"net/http"

	"installer/internal/database"
	"installer/internal/deployer"
	"installer/internal/detector"
	"installer/internal/installer"

	"github.com/gin-gonic/gin"
)

type Server struct {
	webFiles   embed.FS
	assetFiles embed.FS
	detector   detector.Detector
	installer  installer.Manager
	dbManager  database.Manager
	deployer   deployer.Deployer
}

func NewServer(webFiles, assetFiles embed.FS) (*Server, error) {
	return &Server{
		webFiles:   webFiles,
		assetFiles: assetFiles,
		detector:   detector.New(),
		installer:  installer.New(),
		dbManager:  database.New(),
		deployer:   deployer.New(assetFiles),
	}, nil
}

func (s *Server) Router() *gin.Engine {
	r := gin.Default()

	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 静态文件服务
	r.StaticFS("/assets", http.FS(s.webFiles))

	// API路由
	api := r.Group("/api")
	{
		api.GET("/system", s.getSystemInfo)
		api.GET("/dependencies", s.checkDependencies)
		api.POST("/install-dependency", s.installDependency)
		api.POST("/test-database", s.testDatabase)
		api.POST("/setup-database", s.setupDatabase)
		api.POST("/install", s.startInstallation)
		api.GET("/install-progress", s.getInstallProgress)
		api.GET("/install-logs", s.getInstallLogs)
		api.POST("/finish-install", s.finishInstall)
	}

	// 前端路由 (SPA)
	r.NoRoute(func(c *gin.Context) {
		c.FileFromFS("web/dist/index.html", http.FS(s.webFiles))
	})

	return r
}

// API处理函数
func (s *Server) getSystemInfo(c *gin.Context) {
	info, err := s.detector.DetectSystem()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, info)
}

func (s *Server) checkDependencies(c *gin.Context) {
	deps, err := s.detector.CheckDependencies()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, deps)
}

func (s *Server) installDependency(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := s.installer.Install(req.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "安装成功"})
}

func (s *Server) testDatabase(c *gin.Context) {
	var config database.Config
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := s.dbManager.TestConnection(config)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "数据库连接成功"})
}

func (s *Server) setupDatabase(c *gin.Context) {
	var config database.Config
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := s.dbManager.Setup(config)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "数据库配置成功"})
}

func (s *Server) startInstallation(c *gin.Context) {
	var config struct {
		ProjectPath string            `json:"project_path"`
		Database    database.Config   `json:"database"`
		Settings    map[string]string `json:"settings"`
	}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 启动异步安装过程
	go s.deployer.Deploy(config.ProjectPath, config.Database, config.Settings)

	c.JSON(200, gin.H{"message": "安装已开始"})
}

func (s *Server) getInstallProgress(c *gin.Context) {
	progress := s.deployer.GetProgress()
	c.JSON(200, progress)
}

func (s *Server) getInstallLogs(c *gin.Context) {
	logs := s.deployer.GetLogs()
	c.JSON(200, logs)
}

func (s *Server) finishInstall(c *gin.Context) {
	err := s.deployer.Verify()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "安装完成",
		"urls": map[string]string{
			"frontend": "http://localhost:3000",
			"backend":  "http://localhost:8080",
		},
	})
}
