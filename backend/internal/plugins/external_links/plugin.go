package external_links

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"vite-pluginend/internal/api/handlers"
	"vite-pluginend/internal/services"
	"vite-pluginend/pkg/cache"
)

// Plugin 外链插件
type Plugin struct {
	db    *mongo.Database
	cache cache.Cache
}

// NewPlugin 创建外链插件实例
func NewPlugin(db *mongo.Database, cache cache.Cache) *Plugin {
	return &Plugin{
		db:    db,
		cache: cache,
	}
}

// Register 注册插件路由
func (p *Plugin) Register(r *gin.RouterGroup) {
	// 创建服务和处理器
	externalLinkService := services.NewExternalLinkService(p.db, p.cache)
	externalLinkHandler := handlers.NewExternalLinkHandler(externalLinkService)

	// 注册路由
	externalLinks := r.Group("/external-links")
	{
		externalLinks.POST("", externalLinkHandler.CreateExternalLink)
		externalLinks.GET("/:id", externalLinkHandler.GetExternalLink)
		externalLinks.PUT("/:id", externalLinkHandler.UpdateExternalLink)
		externalLinks.DELETE("/:id", externalLinkHandler.DeleteExternalLink)
		externalLinks.GET("", externalLinkHandler.ListExternalLinks)
		externalLinks.GET("/all", externalLinkHandler.GetAllExternalLinks)
		externalLinks.GET("/invalid", externalLinkHandler.GetInvalidExternalLinks)
		externalLinks.DELETE("/batch", externalLinkHandler.BatchDeleteExternalLinks)
		externalLinks.DELETE("/invalid/batch", externalLinkHandler.BatchDeleteInvalidExternalLinks)
		externalLinks.POST("/batch-check", externalLinkHandler.BatchCheckExternalLinks)
		externalLinks.GET("/statistics", externalLinkHandler.GetExternalStatistics)
		externalLinks.GET("/trends", externalLinkHandler.GetExternalTrends)
		externalLinks.POST("/:id/clicks", externalLinkHandler.IncrementClicks)
	}
}

// GetInfo 获取插件信息
func (p *Plugin) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":        "external-links",
		"version":     "1.0.0",
		"description": "外链管理插件",
		"author":      "System",
		"routes": []string{
			"/api/external-links",
			"/api/external-links/:id",
			"/api/external-links/statistics",
			"/api/external-links/trends",
		},
	}
} 