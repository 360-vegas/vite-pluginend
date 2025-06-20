package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// GenerateKey 生成随机密钥
func GenerateKey() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为备选方案
		return "key-" + time.Now().Format("20060102150405")
	}
	return "key-" + hex.EncodeToString(bytes)
} 