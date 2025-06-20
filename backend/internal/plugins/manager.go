package plugins

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"vite-pluginend/internal/plugins/external_links"
	"vite-pluginend/pkg/cache"
)

// Plugin 插件接口
type Plugin interface {
	Register(r *gin.RouterGroup)
	GetInfo() map[string]interface{}
}

// Manager 插件管理器
type Manager struct {
	plugins map[string]Plugin
}

// NewManager 创建插件管理器实例
func NewManager(db *mongo.Database, cache cache.Cache) *Manager {
	manager := &Manager{
		plugins: make(map[string]Plugin),
	}

	// 注册内置插件
	manager.RegisterPlugin("external-links", external_links.NewPlugin(db, cache))

	return manager
}

// RegisterPlugin 注册插件
func (m *Manager) RegisterPlugin(name string, plugin Plugin) {
	m.plugins[name] = plugin
}

// GetPlugin 获取插件
func (m *Manager) GetPlugin(name string) (Plugin, bool) {
	plugin, exists := m.plugins[name]
	return plugin, exists
}

// GetPlugins 获取所有插件
func (m *Manager) GetPlugins() map[string]Plugin {
	return m.plugins
}

// RegisterRoutes 注册所有插件的路由
func (m *Manager) RegisterRoutes(r *gin.RouterGroup) {
	for _, plugin := range m.plugins {
		plugin.Register(r)
	}
}

// GetPluginsInfo 获取所有插件信息
func (m *Manager) GetPluginsInfo() map[string]map[string]interface{} {
	info := make(map[string]map[string]interface{})
	for name, plugin := range m.plugins {
		info[name] = plugin.GetInfo()
	}
	return info
} 