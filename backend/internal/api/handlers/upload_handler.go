package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"vite-pluginend/internal/services"
	customerrors "vite-pluginend/pkg/errors"
)

// UploadHandler 文件上传处理器
type UploadHandler struct {
	uploadService *services.UploadService
}

// NewUploadHandler 创建文件上传处理器实例
func NewUploadHandler(uploadService *services.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

// UploadFile 上传文件
func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, customerrors.NewError("无效的请求参数", http.StatusBadRequest))
		return
	}

	if file == nil {
		c.JSON(http.StatusBadRequest, customerrors.NewError("未找到文件", http.StatusBadRequest))
		return
	}

	filename, err := h.uploadService.SaveFile(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("保存文件失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "文件上传成功",
		"filename": filename,
	})
}

// GetFile 获取文件
func (h *UploadHandler) GetFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, customerrors.NewError("文件名不能为空", http.StatusBadRequest))
		return
	}

	_, data, err := h.uploadService.GetFile(c.Request.Context(), filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("获取文件失败", http.StatusInternalServerError))
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", data)
}

// DeleteFile 删除文件
func (h *UploadHandler) DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, customerrors.NewError("文件名不能为空", http.StatusBadRequest))
		return
	}

	if err := h.uploadService.DeleteFile(c.Request.Context(), filename); err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("删除文件失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件删除成功"})
} 