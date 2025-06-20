package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// CreateZipFile 创建zip文件
func CreateZipFile(sourceDir, targetPath string) error {
	// 创建zip文件
	zipFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	// 创建zip写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历源目录
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 创建相对路径
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}

		// 创建zip文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create zip header: %v", err)
		}
		header.Name = relPath

		// 创建zip文件写入器
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip writer: %v", err)
		}

		// 打开源文件
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open source file: %v", err)
		}
		defer file.Close()

		// 复制文件内容
		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("failed to copy file content: %v", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk source directory: %v", err)
	}

	return nil
}

// EnsureDir 确保目录存在
func EnsureDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	return nil
}

// CleanPath 清理文件路径
func CleanPath(path string) string {
	return filepath.Clean(strings.TrimSpace(path))
}

// GetFileSize 获取文件大小
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %v", err)
	}
	return info.Size(), nil
}

// IsValidFileType 检查文件类型是否有效
func IsValidFileType(ext string) bool {
	// 允许的文件类型
	allowedTypes := map[string]bool{
		".js":    true,
		".ts":    true,
		".jsx":   true,
		".tsx":   true,
		".json":  true,
		".css":   true,
		".scss":  true,
		".less":  true,
		".html":  true,
		".md":    true,
		".png":   true,
		".jpg":   true,
		".jpeg":  true,
		".gif":   true,
		".svg":   true,
		".ico":   true,
		".woff":  true,
		".woff2": true,
		".ttf":   true,
		".eot":   true,
		".zip":   true,
		".tar":   true,
		".gz":    true,
		".rar":   true,
		".7z":    true,
	}

	return allowedTypes[strings.ToLower(ext)]
} 