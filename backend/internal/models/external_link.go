package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ExternalLink 表示一个外部链接
type ExternalLink struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	URL       string            `bson:"url" json:"url"`
	Category  string            `bson:"category" json:"category"`
	Clicks    int               `bson:"clicks" json:"clicks"`
	Status    bool              `bson:"status" json:"status"`
	IsValid   bool              `bson:"is_valid" json:"is_valid"`
	IsActive  bool              `bson:"is_active" json:"is_active"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}

// ExternalLinkQuery 外链查询参数
type ExternalLinkQuery struct {
	Page        int    `form:"page"`
	PerPage     int    `form:"per_page"`
	Keyword     string `form:"keyword"`
	Category    string `form:"category"`
	Status      string `form:"status"`
	IsValid     *bool  `form:"is_valid"`
	MinPriority int    `form:"min_priority"`
	MaxPriority int    `form:"max_priority"`
	MinClicks   int    `form:"min_clicks"`
	OnlyValid   bool   `form:"only_valid"`
	Popular     bool   `form:"popular"`
	SortField   string `form:"sort_field"`
	SortOrder   string `form:"sort_order"`
}

// ExternalLinkResponse 外链响应结构
type ExternalLinkResponse struct {
	Data []ExternalLink `json:"data"`
	Meta struct {
		Total       int `json:"total"`
		PerPage     int `json:"per_page"`
		CurrentPage int `json:"current_page"`
		LastPage    int `json:"last_page"`
	} `json:"meta"`
}

// ExternalStatistics 外链统计信息
type ExternalStatistics struct {
	TotalLinks    int            `json:"total_links"`
	ActiveLinks   int            `json:"active_links"`
	ExpiredLinks  int            `json:"expired_links"`
	InvalidLinks  int            `json:"invalid_links"`
	TotalClicks   int            `json:"total_clicks"`
	AverageClicks float64        `json:"average_clicks"`
	Categories    map[string]int `json:"categories"`
	Tags          map[string]int `json:"tags"`
}

// ExternalTrend 外链趋势数据
type ExternalTrend struct {
	Date        string `json:"date"`
	NewLinks    int    `json:"new_links"`
	TotalClicks int    `json:"total_clicks"`
	ActiveLinks int    `json:"active_links"`
}

// LinkCheckResult 链接检测结果
type LinkCheckResult struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	IsValid      bool   `json:"is_valid"`
	Message      string `json:"message"`
	ErrorMessage string `json:"error_message,omitempty"`
	CheckedAt    time.Time `json:"checked_at"`
} 