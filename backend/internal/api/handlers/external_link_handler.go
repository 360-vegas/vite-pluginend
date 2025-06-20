package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"vite-pluginend/internal/models"
	"vite-pluginend/internal/services"
	"vite-pluginend/pkg/errors"
	"vite-pluginend/pkg/logger"
)

// ExternalLinkHandler 外链处理器
type ExternalLinkHandler struct {
	externalLinkService *services.ExternalLinkService
}

// NewExternalLinkHandler 创建外链处理器实例
func NewExternalLinkHandler(externalLinkService *services.ExternalLinkService) *ExternalLinkHandler {
	return &ExternalLinkHandler{
		externalLinkService: externalLinkService,
	}
}

// CreateExternalLink 创建外链
func (h *ExternalLinkHandler) CreateExternalLink(c *gin.Context) {
	log := logger.NewLogger()
	log.Info("CreateExternalLink handler reached!")
	var link models.ExternalLink
	if err := c.ShouldBindJSON(&link); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewError("无效的请求参数", http.StatusBadRequest))
		return
	}

	if err := h.externalLinkService.CreateExternalLink(c.Request.Context(), &link); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("创建外链失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, link)
}

// GetExternalLink 获取外链
func (h *ExternalLinkHandler) GetExternalLink(c *gin.Context) {
	id := c.Param("id")
	link, err := h.externalLinkService.GetExternalLink(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("获取外链失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, link)
}

// UpdateExternalLink 更新外链
func (h *ExternalLinkHandler) UpdateExternalLink(c *gin.Context) {
	id := c.Param("id")
	var update bson.M
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewError("无效的请求参数", http.StatusBadRequest))
		return
	}

	if err := h.externalLinkService.UpdateExternalLink(c.Request.Context(), id, update); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("更新外链失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteExternalLink 删除外链
func (h *ExternalLinkHandler) DeleteExternalLink(c *gin.Context) {
	id := c.Param("id")
	if err := h.externalLinkService.DeleteExternalLink(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("删除外链失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// BatchDeleteExternalLinks 批量删除外链
func (h *ExternalLinkHandler) BatchDeleteExternalLinks(c *gin.Context) {
	var request struct {
		IDs []string `json:"ids"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewError("无效的请求参数", http.StatusBadRequest))
		return
	}

	if len(request.IDs) == 0 {
		c.JSON(http.StatusBadRequest, errors.NewError("删除列表不能为空", http.StatusBadRequest))
		return
	}

	deletedCount, err := h.externalLinkService.BatchDeleteExternalLinks(c.Request.Context(), request.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("批量删除外链失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量删除成功",
		"deleted_count": deletedCount,
	})
}

// ListExternalLinks 获取外链列表
func (h *ExternalLinkHandler) ListExternalLinks(c *gin.Context) {
	var query models.ExternalLinkQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewError("无效的查询参数", http.StatusBadRequest))
		return
	}

	response, err := h.externalLinkService.ListExternalLinks(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("获取外链列表失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAllExternalLinks 获取所有外链（不分页）
func (h *ExternalLinkHandler) GetAllExternalLinks(c *gin.Context) {
	links, err := h.externalLinkService.GetAllExternalLinks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("获取所有外链失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": links,
		"total": len(links),
	})
}

// BatchCheckExternalLinks 批量检测外链
func (h *ExternalLinkHandler) BatchCheckExternalLinks(c *gin.Context) {
	log := logger.NewLogger()
	
	var request struct {
		IDs []string `json:"ids"`
		All bool     `json:"all,omitempty"` // 是否检测所有链接
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("请求参数绑定失败", zap.Error(err))
		c.JSON(http.StatusBadRequest, errors.NewError(fmt.Sprintf("无效的请求参数: %v", err), http.StatusBadRequest))
		return
	}
	
	log.Info("批量检测请求", zap.Any("request", request))

	results, err := h.externalLinkService.BatchCheckExternalLinks(c.Request.Context(), request.IDs, request.All)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("批量检测外链失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量检测完成",
		"results": results,
	})
}

// GetExternalStatistics 获取外链统计信息
func (h *ExternalLinkHandler) GetExternalStatistics(c *gin.Context) {
	stats, err := h.externalLinkService.GetExternalStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("获取统计信息失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetExternalTrends 获取外链趋势数据
func (h *ExternalLinkHandler) GetExternalTrends(c *gin.Context) {
	period := c.DefaultQuery("period", "day")
	limit := 7
	if period == "week" {
		limit = 4
	} else if period == "month" {
		limit = 6
	}

	trends, err := h.externalLinkService.GetExternalTrends(c.Request.Context(), period, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("获取趋势数据失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, trends)
}

// IncrementClicks 增加点击量
func (h *ExternalLinkHandler) IncrementClicks(c *gin.Context) {
	id := c.Param("id")
	if err := h.externalLinkService.IncrementClicks(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("增加点击量失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "点击量增加成功"})
}

// GetInvalidExternalLinks 获取所有不可用的外链
func (h *ExternalLinkHandler) GetInvalidExternalLinks(c *gin.Context) {
	links, err := h.externalLinkService.GetInvalidExternalLinks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("获取不可用外链失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": links,
		"total": len(links),
	})
}

// BatchDeleteInvalidExternalLinks 批量删除所有不可用的外链
func (h *ExternalLinkHandler) BatchDeleteInvalidExternalLinks(c *gin.Context) {
	deletedCount, err := h.externalLinkService.BatchDeleteInvalidExternalLinks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewError("批量删除不可用外链失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "批量删除不可用外链成功",
		"deleted_count": deletedCount,
	})
} 