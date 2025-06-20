package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Plugin 表示一个插件的基本信息
type Plugin struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string            `bson:"name" json:"name"`
	Key         string            `bson:"key" json:"key"`
	Path        string            `bson:"path" json:"path"`
	Description string            `bson:"description" json:"description"`
	Author      string            `bson:"author" json:"author"`
	Version     string            `bson:"version" json:"version"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

// PluginLog 表示插件的访问日志
type PluginLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PluginKey string            `bson:"plugin_key" json:"plugin_key"`
	UserID    string            `bson:"user_id" json:"user_id"`
	Action    string            `bson:"action" json:"action"`
	Details   map[string]interface{} `bson:"details" json:"details"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
}

// CreatePluginRequest 创建插件的请求结构
type CreatePluginRequest struct {
	Name        string `json:"name" binding:"required"`
	Key         string `json:"key" binding:"required"`
	Path        string `json:"path" binding:"required"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Version     string `json:"version"`
}

// PluginResponse 插件响应结构
type PluginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *Plugin `json:"data,omitempty"`
} 