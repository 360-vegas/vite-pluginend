package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"

	"vite-pluginend/pkg/logger"
	customerrors "vite-pluginend/pkg/errors"
	"vite-pluginend/pkg/utils"
)

// UploadService 文件上传服务
type UploadService struct {
	UploadDir string
}

// NewUploadService 创建文件上传服务实例
func NewUploadService(uploadDir string) *UploadService {
	if err := utils.EnsureDir(uploadDir); err != nil {
		logger.Fatal("Failed to create upload directory", zap.Error(err))
	}

	return &UploadService{UploadDir: uploadDir}
}

// SaveFile 保存上传的文件
func (s *UploadService) SaveFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		logger.Error("打开上传文件失败", zap.Error(err))
		return "", customerrors.NewError("打开上传文件失败", http.StatusInternalServerError)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(s.UploadDir, filename)

	out, err := os.Create(filePath)
	if err != nil {
		logger.Error("创建文件失败", zap.Error(err))
		return "", customerrors.NewError("创建文件失败", http.StatusInternalServerError)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		logger.Error("保存文件失败", zap.Error(err))
		return "", customerrors.NewError("保存文件失败", http.StatusInternalServerError)
	}

	logger.Info("File saved successfully", 
		zap.String("filename", filename),
	)
	return filename, nil
}

// GetFile 获取文件
func (s *UploadService) GetFile(ctx context.Context, filename string) (string, []byte, error) {
	filePath := filepath.Join(s.UploadDir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warn("File not found", 
				zap.String("filename", filename),
			)
			return "", nil, customerrors.NewError("文件未找到", http.StatusNotFound)
		}
		logger.Error("读取文件失败", zap.Error(err))
		return "", nil, customerrors.NewError("读取文件失败", http.StatusInternalServerError)
	}

	logger.Debug("File found", 
		zap.String("filename", filename),
	)
	return filePath, data, nil
}

// DeleteFile 删除文件
func (s *UploadService) DeleteFile(ctx context.Context, filename string) error {
	filePath := filepath.Join(s.UploadDir, filename)
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warn("File not found for deletion", 
				zap.String("filename", filename),
			)
			return customerrors.NewError("文件未找到", http.StatusNotFound)
		}
		logger.Error("Failed to delete file", zap.Error(err))
		return customerrors.NewError("删除文件失败", http.StatusInternalServerError)
	}

	logger.Info("File deleted successfully", 
		zap.String("filename", filename),
	)
	return nil
}

// ListFiles 列出插件目录下的所有文件
func (s *UploadService) ListFiles(pluginKey string) ([]string, error) {
	pluginDir := filepath.Join(s.UploadDir, pluginKey)
	if _, err := os.Stat(pluginDir); os.IsNotExist(err) {
		logger.Warn("Plugin directory not found", zap.String("plugin_key", pluginKey))
		return nil, customerrors.NewError("插件目录未找到", http.StatusNotFound)
	}

	var files []string
	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(pluginDir, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})

	if err != nil {
		logger.Error("Failed to list files", zap.Error(err))
		return nil, customerrors.NewError("列出文件失败", http.StatusInternalServerError)
	}

	logger.Info("Files listed successfully", 
		zap.String("plugin_key", pluginKey),
		zap.Int("count", len(files)),
	)
	return files, nil
}

// GetFileInfo 获取文件信息
func (s *UploadService) GetFileInfo(pluginKey, filename string) (os.FileInfo, error) {
	filePath := filepath.Join(s.UploadDir, pluginKey, filename)
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warn("File not found", 
				zap.String("filename", filename),
				zap.String("plugin_key", pluginKey),
			)
			return nil, customerrors.NewError("文件未找到", http.StatusNotFound)
		}
		logger.Error("Failed to get file info", zap.Error(err))
		return nil, customerrors.NewError("获取文件信息失败", http.StatusInternalServerError)
	}

	logger.Debug("File info retrieved successfully", 
		zap.String("filename", filename),
		zap.String("plugin_key", pluginKey),
	)
	return info, nil
} 