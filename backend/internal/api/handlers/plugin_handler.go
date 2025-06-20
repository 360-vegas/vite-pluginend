package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"vite-pluginend/internal/models"
	"vite-pluginend/internal/services"
	customerrors "vite-pluginend/pkg/errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// PluginHandler 处理插件相关的HTTP请求
type PluginHandler struct {
	pluginService     *services.PluginService
	dependencyService *services.DependencyService
}

// NewPluginHandler 创建新的插件处理器
func NewPluginHandler(pluginService *services.PluginService, dependencyService *services.DependencyService) *PluginHandler {
	return &PluginHandler{
		pluginService:     pluginService,
		dependencyService: dependencyService,
	}
}

// normalizePluginName 标准化插件名称，确保有plugin-前缀
func (h *PluginHandler) normalizePluginName(pluginKey string) string {
	if strings.HasPrefix(pluginKey, "plugin-") {
		return pluginKey
	}
	return "plugin-" + pluginKey
}

// CreatePlugin 创建插件
func (h *PluginHandler) CreatePlugin(c *gin.Context) {
	var req services.Plugin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, customerrors.NewError("无效的请求参数", http.StatusBadRequest))
		return
	}

	if err := h.pluginService.CreatePlugin(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("创建插件失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "插件创建成功"})
}

// GetPlugin 获取插件信息
func (h *PluginHandler) GetPlugin(c *gin.Context) {
	id := c.Param("id")
	plugin, err := h.pluginService.GetPlugin(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("获取插件失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, plugin)
}

// ListPlugins 获取插件列表
func (h *PluginHandler) ListPlugins(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	plugins, total, err := h.pluginService.ListPlugins(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("获取插件列表失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":   total,
		"plugins": plugins,
	})
}

// UpdatePlugin 更新插件
func (h *PluginHandler) UpdatePlugin(c *gin.Context) {
	id := c.Param("id")
	var req services.Plugin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, customerrors.NewError("无效的请求参数", http.StatusBadRequest))
		return
	}

	update := bson.M{
		"name":        req.Name,
		"version":     req.Version,
		"description": req.Description,
		"author":      req.Author,
		"file_id":     req.FileID,
	}

	if err := h.pluginService.UpdatePlugin(c.Request.Context(), id, update); err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("更新插件失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "插件更新成功"})
}

// DeletePlugin 删除插件
func (h *PluginHandler) DeletePlugin(c *gin.Context) {
	pluginKey := c.Param("id")

	// 获取当前工作目录
	workDir, _ := os.Getwd()

	// 构建插件目录路径 - 从 backend/cmd/server 回到项目根目录
	// 当前在 backend/cmd/server，需要上3级到项目根目录，然后进入 src/plugins
	fullPluginName := h.normalizePluginName(pluginKey)
	pluginDir := filepath.Join(workDir, "..", "..", "..", "src", "plugins", fullPluginName)

	// 清理路径以移除 .. 引用
	pluginDir = filepath.Clean(pluginDir)

	// 添加调试信息
	fmt.Printf("🗑️ 删除插件请求: %s\n", pluginKey)
	fmt.Printf("标准化插件名: %s\n", fullPluginName)
	fmt.Printf("当前工作目录: %s\n", workDir)
	fmt.Printf("插件目录: %s\n", pluginDir)

	// 检查插件目录是否存在
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		fmt.Printf("❌ 插件目录不存在: %s\n", pluginDir)
		c.JSON(http.StatusNotFound, customerrors.NewError("插件不存在", http.StatusNotFound))
		return
	}

	// 删除插件目录
	if err := os.RemoveAll(pluginDir); err != nil {
		fmt.Printf("❌ 删除插件目录失败: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, customerrors.NewError("删除插件失败: "+err.Error(), http.StatusInternalServerError))
		return
	}

	fmt.Printf("✅ 插件删除成功: %s\n", pluginKey)

	// 从数据库删除（如果存在）
	if err := h.pluginService.DeletePlugin(c.Request.Context(), pluginKey); err != nil {
		// 文件已删除，数据库删除失败只记录警告，不返回错误
		// 因为插件文件系统删除是主要操作
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "插件删除成功",
		"data": gin.H{
			"pluginKey": pluginKey,
		},
	})
}

// PackageConfig 定义打包配置
type PackageConfig struct {
	IncludeDemoData bool `json:"includeDemoData"`
	IncludeDocs     bool `json:"includeDocs"`
	IncludeTests    bool `json:"includeTests"`
	Minify          bool `json:"minify"`
}

// PackagePluginRequest 定义打包请求
type PackagePluginRequest struct {
	PluginName string        `json:"pluginName"`
	Config     PackageConfig `json:"config"`
}

// PackagePlugin 处理插件打包请求
func (h *PluginHandler) PackagePlugin(c *gin.Context) {
	var req PackagePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 使用绝对路径构建插件目录 - 从 backend/cmd/server 回到项目根目录
	workDir, _ := os.Getwd()

	// 验证插件是否存在
	pluginDir := filepath.Join(workDir, "..", "..", "..", "src", "plugins", req.PluginName)
	pluginDir = filepath.Clean(pluginDir)
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "插件不存在"})
		return
	}

	// 创建输出目录
	outputDir := filepath.Join(workDir, "..", "..", "..", "dist", "plugins", req.PluginName)
	outputDir = filepath.Clean(outputDir)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建输出目录失败"})
		return
	}

	// 打包插件
	zipFile := filepath.Join(outputDir, req.PluginName+".zip")
	if err := h.createPluginPackage(pluginDir, zipFile, req.Config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打包失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "打包成功",
		"data": gin.H{
			"outputPath": zipFile,
		},
	})
}

// createPluginPackage 创建插件包
func (h *PluginHandler) createPluginPackage(srcDir, destZip string, config PackageConfig) error {
	// 创建 zip 文件
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历源目录
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过不需要的文件和目录
		if !config.IncludeDemoData && strings.Contains(path, "demo") {
			return nil
		}
		if !config.IncludeDocs && strings.Contains(path, "docs") {
			return nil
		}
		if !config.IncludeTests && strings.Contains(path, "test") {
			return nil
		}

		// 获取相对路径
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// 创建 zip 头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 写入文件内容
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

// ScanPlugins 扫描插件目录
func (h *PluginHandler) ScanPlugins(c *gin.Context) {
	// 使用绝对路径构建插件目录 - 从 backend/cmd/server 回到项目根目录
	workDir, _ := os.Getwd()
	pluginsDir := filepath.Join(workDir, "..", "..", "..", "src", "plugins")
	pluginsDir = filepath.Clean(pluginsDir)

	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "扫描插件目录失败",
			"workDir":     workDir,
			"pluginsDir":  pluginsDir,
			"errorDetail": err.Error(),
		})
		return
	}

	var plugins []gin.H
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "plugin-") {
			// 读取插件的 meta.ts 文件获取版本信息
			metaPath := filepath.Join(pluginsDir, entry.Name(), "meta.ts")
			version := "0.1.0" // 默认版本

			if metaContent, err := os.ReadFile(metaPath); err == nil {
				// 简单解析版本信息
				if matches := regexp.MustCompile(`version:\s*['"]([^'"]+)['"]`).FindSubmatch(metaContent); len(matches) > 1 {
					version = string(matches[1])
				}
			}

			plugins = append(plugins, gin.H{
				"name":      entry.Name(),
				"version":   version,
				"status":    "ready",
				"directory": filepath.Join("src/plugins", entry.Name()),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"plugins": plugins,
	})
}

// DownloadPlugin 处理插件下载请求
func (h *PluginHandler) DownloadPlugin(c *gin.Context) {
	// 使用绝对路径构建下载路径 - 从 backend/cmd/server 回到项目根目录
	workDir, _ := os.Getwd()

	pluginName := c.Param("name")
	zipPath := filepath.Join(workDir, "..", "..", "..", "dist", "plugins", pluginName, pluginName+".zip")
	zipPath = filepath.Clean(zipPath)

	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "插件包不存在"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", pluginName))
	c.Header("Content-Type", "application/zip")
	c.File(zipPath)
}

// GeneratePluginRequest 定义生成插件请求
type GeneratePluginRequest struct {
	Name        string `json:"name"`
	Key         string `json:"key"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Category    string `json:"category"`
	Pages       []struct {
		Title string `json:"title"`
		Key   string `json:"key"`
		Type  string `json:"type"`
	} `json:"pages"`
	WithTests bool     `json:"withTests"`
	WithDocs  bool     `json:"withDocs"`
	Features  []string `json:"features"`
}

// GeneratePlugin 处理插件生成请求
func (h *PluginHandler) GeneratePlugin(c *gin.Context) {
	var req GeneratePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 使用绝对路径构建插件目录 - 从 backend/cmd/server 回到项目根目录
	workDir, _ := os.Getwd()

	// 创建插件目录
	pluginDir := filepath.Join(workDir, "..", "..", "..", "src", "plugins", "plugin-"+req.Key)
	pluginDir = filepath.Clean(pluginDir)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建插件目录失败"})
		return
	}

	// 生成插件文件
	if err := h.generatePluginFiles(pluginDir, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成插件文件失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "插件生成成功",
		"data": gin.H{
			"pluginKey": "plugin-" + req.Key,
			"path":      pluginDir,
		},
	})
}

// generatePluginFiles 生成插件文件
func (h *PluginHandler) generatePluginFiles(pluginDir string, req GeneratePluginRequest) error {
	// 创建子目录
	directories := []string{"pages", "components", "assets"}
	for _, dir := range directories {
		if err := os.MkdirAll(filepath.Join(pluginDir, dir), 0755); err != nil {
			return err
		}
	}

	// 创建用于模板的数据结构
	templateData := struct {
		Name        string
		Version     string
		Author      string
		Description string
		Category    string
		Key         string
		Pages       []struct {
			Key   string
			Title string
		}
	}{
		Name:        req.Name,
		Version:     req.Version,
		Author:      req.Author,
		Description: req.Description,
		Category:    req.Category,
		Key:         req.Key,
	}

	for _, page := range req.Pages {
		templateData.Pages = append(templateData.Pages, struct {
			Key   string
			Title string
		}{
			Key:   page.Key,
			Title: page.Title,
		})
	}

	// 生成 meta.ts 文件
	metaTemplate := `export default {
  name: '{{.Name}}',
  version: '{{.Version}}',
  author: '{{.Author}}',
  description: '{{.Description}}',
  category: '{{.Category}}',
  mainNav: {
    key: '{{.Key}}',
    title: '{{.Name}}',
    icon: 'Box',
    path: '/{{.Key}}'
  },
  subNav: [
    {{range .Pages}}{
      key: '{{.Key}}',
      title: '{{.Title}}',
      icon: 'Document',
      path: '/{{$.Key}}/{{.Key}}'
    },{{end}}
  ],
  pages: [
    {{range .Pages}}{
      key: '{{.Key}}',
      title: '{{.Title}}',
      path: '/{{$.Key}}/{{.Key}}',
      name: '{{$.Key}}-{{.Key}}',
      icon: 'Document',
      description: '{{.Title}}页面',
      component: () => import('./pages/{{.Key}}.vue')
    },{{end}}
  ]
}
`

	if err := h.writeTemplateFile(filepath.Join(pluginDir, "meta.ts"), metaTemplate, templateData); err != nil {
		return err
	}

	// 生成 index.ts 文件
	indexTemplate := `import type { App } from 'vue'
import meta from './meta'

export default {
  install(app: App) {
    // 插件安装逻辑
    console.log('{{.Name}} 插件已安装')
  },
  meta
}

export { meta }
`

	if err := h.writeTemplateFile(filepath.Join(pluginDir, "index.ts"), indexTemplate, req); err != nil {
		return err
	}

	// 为每个页面生成 Vue 文件
	pageTemplate := `<template>
  <div class="plugin-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <h1 class="page-title">{{.PageTitle}}</h1>
      <p class="page-description">{{.PluginName}} 插件</p>
    </div>

    <!-- 页面内容 -->
    <div class="page-content">
      <div class="content-grid">
        <!-- 欢迎卡片 -->
        <el-card class="welcome-card" shadow="never">
          <div class="welcome-content">
            <el-icon class="welcome-icon"><Box /></el-icon>
            <h3>欢迎使用 {{.PluginName}}</h3>
            <p>开始构建您的功能模块</p>
          </div>
        </el-card>

        <!-- 功能区域 -->
        <el-card class="feature-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>功能模块</span>
              <el-tag size="small">{{.PluginName}}</el-tag>
            </div>
          </template>
          
          <div class="feature-list">
            <div class="feature-item">
              <el-icon class="feature-icon"><Setting /></el-icon>
              <div class="feature-info">
                <h4>基础配置</h4>
                <p>配置插件的基本参数和选项</p>
              </div>
            </div>
            
            <div class="feature-item">
              <el-icon class="feature-icon"><DataLine /></el-icon>
              <div class="feature-info">
                <h4>数据管理</h4>
                <p>管理和处理插件相关数据</p>
              </div>
            </div>
            
            <div class="feature-item">
              <el-icon class="feature-icon"><Tools /></el-icon>
              <div class="feature-info">
                <h4>工具集成</h4>
                <p>集成外部工具和服务</p>
              </div>
            </div>
          </div>

          <div class="action-buttons">
            <el-button type="primary" :icon="Plus" @click="handleAction('create')">
              创建新项目
            </el-button>
            <el-button :icon="Setting" @click="handleAction('settings')">
              插件设置
            </el-button>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElMessage } from 'element-plus'
import { 
  Box, 
  Setting, 
  DataLine, 
  Tools,
  Plus 
} from '@element-plus/icons-vue'

// 操作处理
const handleAction = (action: string) => {
  console.log('执行操作:', action)
  ElMessage.success(` + "`${action} 功能开发中...`" + `)
}
</script>

<style scoped>
.plugin-page {
  padding: 20px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 120px);
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.page-description {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.page-content {
  max-width: 1200px;
}

.content-grid {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 20px;
}

.welcome-card {
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.welcome-content {
  text-align: center;
  padding: 40px 20px;
}

.welcome-icon {
  font-size: 48px;
  color: #409eff;
  margin-bottom: 16px;
}

.welcome-content h3 {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
}

.welcome-content p {
  font-size: 14px;
  color: #909399;
  margin: 0;
}

.feature-card {
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
  color: #303133;
}

.feature-list {
  margin-bottom: 24px;
}

.feature-item {
  display: flex;
  align-items: center;
  padding: 16px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.feature-item:hover {
  background-color: #f5f7fa;
}

.feature-icon {
  font-size: 20px;
  color: #409eff;
  margin-right: 12px;
  flex-shrink: 0;
}

.feature-info h4 {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 4px 0;
}

.feature-info p {
  font-size: 12px;
  color: #909399;
  margin: 0;
  line-height: 1.4;
}

.action-buttons {
  display: flex;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid #e4e7ed;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .plugin-page {
    padding: 16px;
  }
  
  .content-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .welcome-content {
    padding: 24px 16px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .action-buttons .el-button {
    width: 100%;
  }
}
</style>
`

	for _, page := range req.Pages {
		pageData := struct {
			PluginName string
			PluginKey  string
			PageTitle  string
			PageKey    string
		}{
			PluginName: req.Name,
			PluginKey:  req.Key,
			PageTitle:  page.Title,
			PageKey:    page.Key,
		}

		if err := h.writeTemplateFile(filepath.Join(pluginDir, "pages", page.Key+".vue"), pageTemplate, pageData); err != nil {
			return err
		}
	}

	return nil
}

// TogglePlugin 启用/禁用插件
func (h *PluginHandler) TogglePlugin(c *gin.Context) {
	pluginKey := c.Param("id")
	var req struct {
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, customerrors.NewError("无效的请求参数", http.StatusBadRequest))
		return
	}

	// 这里可以添加插件状态管理逻辑
	// 目前简单返回成功
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "插件状态切换成功",
		"data": gin.H{
			"pluginKey": pluginKey,
			"enabled":   req.Enabled,
		},
	})
}

// GetPluginDetail 获取插件详情
func (h *PluginHandler) GetPluginDetail(c *gin.Context) {
	pluginKey := c.Param("id")

	// 这里可以添加从文件系统读取插件信息的逻辑
	// 目前返回模拟数据
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"key":         pluginKey,
			"name":        "插件 " + pluginKey,
			"version":     "1.0.0",
			"description": "这是插件 " + pluginKey + " 的描述",
			"author":      "插件作者",
			"enabled":     true,
		},
	})
}

// ExportPlugin 导出插件
func (h *PluginHandler) ExportPlugin(c *gin.Context) {
	pluginKey := c.Param("id")

	// 获取当前工作目录
	workDir, _ := os.Getwd()

	// 构建插件目录路径 - 从 backend/cmd/server 回到项目根目录
	fullPluginName := h.normalizePluginName(pluginKey)
	pluginDir := filepath.Join(workDir, "..", "..", "..", "src", "plugins", fullPluginName)
	pluginDir = filepath.Clean(pluginDir)

	// 检查插件目录是否存在
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, customerrors.NewError("插件不存在", http.StatusNotFound))
		return
	}

	// 创建临时zip文件
	tempZip := filepath.Join(os.TempDir(), pluginKey+".zip")
	defer os.Remove(tempZip) // 清理临时文件

	// 创建插件包
	config := PackageConfig{
		IncludeDemoData: true,
		IncludeDocs:     true,
		IncludeTests:    true,
		Minify:          false,
	}

	if err := h.createPluginPackage(pluginDir, tempZip, config); err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("导出插件失败", http.StatusInternalServerError))
		return
	}

	// 返回文件
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+pluginKey+".zip")
	c.Header("Content-Type", "application/octet-stream")
	c.File(tempZip)
}

// InstallPlugin 安装插件包
func (h *PluginHandler) InstallPlugin(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("❌ 获取上传文件失败: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请选择要安装的插件包",
		})
		return
	}

	fmt.Printf("📦 安装插件包: %s (大小: %d bytes)\n", file.Filename, file.Size)

	// 验证文件类型
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".zip") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "只支持 .zip 格式的插件包",
		})
		return
	}

	// 创建临时文件
	tempFile := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tempFile); err != nil {
		fmt.Printf("❌ 保存临时文件失败: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存上传文件失败",
		})
		return
	}
	defer os.Remove(tempFile) // 清理临时文件

	// 解压并安装插件
	pluginKey, err := h.extractAndInstallPlugin(tempFile)
	if err != nil {
		fmt.Printf("❌ 安装插件失败: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "安装插件失败: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ 插件安装成功: %s\n", pluginKey)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "插件安装成功",
		"data": gin.H{
			"pluginKey": pluginKey,
		},
	})
}

// extractAndInstallPlugin 解压并安装插件包
func (h *PluginHandler) extractAndInstallPlugin(zipPath string) (string, error) {
	// 打开zip文件
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("打开zip文件失败: %w", err)
	}
	defer reader.Close()

	// 获取插件目录路径
	workDir, _ := os.Getwd()
	pluginsDir := filepath.Join(workDir, "..", "..", "..", "src", "plugins")
	pluginsDir = filepath.Clean(pluginsDir)

	var pluginKey string
	var hasMetaFile bool

	// 验证插件包结构并提取插件key
	fmt.Printf("🔍 分析zip包结构:\n")
	for _, file := range reader.File {
		fmt.Printf("  文件: %s\n", file.Name)
		if strings.HasSuffix(file.Name, "meta.ts") {
			hasMetaFile = true

			// 尝试多种方式提取插件key
			parts := strings.Split(file.Name, "/")

			// 方式1: 如果meta.ts在根目录，尝试从文件名推断
			if len(parts) == 1 {
				// 直接在根目录的meta.ts，可能需要读取文件内容来获取插件名
				fmt.Printf("  meta.ts在根目录，需要读取内容获取插件名\n")
			} else if len(parts) > 1 {
				// 方式2: 从目录路径提取
				possibleKey := parts[0]
				if strings.HasPrefix(possibleKey, "plugin-") {
					pluginKey = strings.TrimPrefix(possibleKey, "plugin-")
					fmt.Printf("  从目录名提取插件key: %s\n", pluginKey)
				} else {
					// 方式3: 目录名不是plugin-开头，使用目录名作为key
					pluginKey = possibleKey
					fmt.Printf("  使用目录名作为插件key: %s\n", pluginKey)
				}
			}
			break
		}
	}

	// 如果还没有找到key，尝试从文件名中提取
	if pluginKey == "" && hasMetaFile {
		// 从zip文件名中提取插件key
		zipFileName := filepath.Base(zipPath)
		if strings.HasPrefix(zipFileName, "plugin-") && strings.HasSuffix(zipFileName, ".zip") {
			pluginKey = strings.TrimSuffix(strings.TrimPrefix(zipFileName, "plugin-"), ".zip")
			fmt.Printf("  从zip文件名提取插件key: %s\n", pluginKey)
		}
	}

	if !hasMetaFile {
		return "", fmt.Errorf("无效的插件包：缺少 meta.ts 文件")
	}

	if pluginKey == "" {
		return "", fmt.Errorf("无效的插件包：无法确定插件名称")
	}

	// 检查插件是否已存在
	fullPluginName := h.normalizePluginName(pluginKey)
	targetDir := filepath.Join(pluginsDir, fullPluginName)
	if _, err := os.Stat(targetDir); err == nil {
		return "", fmt.Errorf("插件 %s 已存在，请先卸载后再安装", pluginKey)
	}

	// 创建插件目录
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("创建插件目录失败: %w", err)
	}

	fmt.Printf("📁 创建插件目录: %s\n", targetDir)

	// 解压文件
	for _, file := range reader.File {
		// 跳过目录项
		if file.FileInfo().IsDir() {
			continue
		}

		// 处理文件路径，移除顶级目录
		relativePath := file.Name
		if strings.Contains(relativePath, "/") {
			parts := strings.Split(relativePath, "/")
			if len(parts) > 1 {
				relativePath = strings.Join(parts[1:], "/")
			}
		}

		if relativePath == "" {
			continue
		}

		// 构建目标文件路径
		destPath := filepath.Join(targetDir, relativePath)

		// 创建目录
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return "", fmt.Errorf("创建目录失败: %w", err)
		}

		// 打开zip中的文件
		rc, err := file.Open()
		if err != nil {
			return "", fmt.Errorf("打开zip中的文件失败: %w", err)
		}

		// 创建目标文件
		destFile, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			return "", fmt.Errorf("创建目标文件失败: %w", err)
		}

		// 复制文件内容
		_, err = io.Copy(destFile, rc)
		rc.Close()
		destFile.Close()

		if err != nil {
			return "", fmt.Errorf("复制文件内容失败: %w", err)
		}

		fmt.Printf("📄 解压文件: %s\n", relativePath)
	}

	return pluginKey, nil
}

// writeTemplateFile 写入模板文件
func (h *PluginHandler) writeTemplateFile(filePath, templateStr string, data interface{}) error {
	tmpl, err := template.New("plugin").Parse(templateStr)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// CheckPluginDependencies 检查插件依赖
func (h *PluginHandler) CheckPluginDependencies(c *gin.Context) {
	pluginKey := c.Param("id")
	if pluginKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "插件ID不能为空",
		})
		return
	}

	// 这里应该从插件的依赖配置文件中读取依赖信息
	// 为了演示，我们创建一个示例依赖配置
	dependencies := h.getPluginDependencyConfig(pluginKey)

	// 检查依赖
	ctx := context.Background()
	result, err := h.dependencyService.CheckPluginDependencies(ctx, pluginKey, dependencies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "检查依赖失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// SetupPluginDatabase 设置插件数据库
func (h *PluginHandler) SetupPluginDatabase(c *gin.Context) {
	pluginKey := c.Param("id")
	if pluginKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "插件ID不能为空",
		})
		return
	}

	var config models.DatabaseSetupOptions
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的配置数据",
		})
		return
	}

	// 设置数据库
	ctx := context.Background()
	if err := h.dependencyService.SetupDatabase(ctx, pluginKey, config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "设置数据库失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "数据库设置成功",
	})
}

// getPluginDependencyConfig 获取插件依赖配置
func (h *PluginHandler) getPluginDependencyConfig(pluginKey string) *models.PluginDependency {
	// 这里应该从插件目录中的依赖配置文件读取
	// 现在先返回一个示例配置，特别是对外链插件

	if pluginKey == "wailki" || strings.Contains(pluginKey, "wailki") {
		return &models.PluginDependency{
			PluginKey: pluginKey,
			Database: &models.DatabaseRequirement{
				Type:         "mongodb",
				DatabaseName: "plugin_external_links",
				Required:     true,
				Collections: []models.CollectionInfo{
					{
						Name:        "external_links",
						Description: "外链数据表",
						Indexes: []models.IndexInfo{
							{
								Name:   "url_index",
								Fields: map[string]int{"url": 1},
								Unique: true,
							},
							{
								Name:   "created_at_index",
								Fields: map[string]int{"created_at": -1},
							},
						},
					},
					{
						Name:        "link_stats",
						Description: "外链统计数据表",
						Indexes: []models.IndexInfo{
							{
								Name:   "link_id_index",
								Fields: map[string]int{"link_id": 1},
							},
						},
					},
				},
			},
			Environment: []models.EnvironmentVariable{
				{
					Name:        "MONGODB_URI",
					Required:    true,
					Type:        "url",
					Description: "MongoDB连接字符串",
					Default:     "mongodb://localhost:27017",
				},
			},
			Dependencies: []models.Dependency{
				{
					Name:        "go.mongodb.org/mongo-driver",
					Version:     ">=1.10.0",
					Type:        "go_module",
					Required:    true,
					Description: "MongoDB Go驱动",
				},
			},
		}
	}

	// 默认配置
	return &models.PluginDependency{
		PluginKey:    pluginKey,
		Dependencies: []models.Dependency{},
	}
}
