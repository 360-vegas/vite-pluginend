package models

import "time"

// PluginDependency 插件依赖配置
type PluginDependency struct {
	ID           string                  `json:"id" bson:"_id,omitempty"`
	PluginKey    string                  `json:"plugin_key" bson:"plugin_key"`
	Dependencies []Dependency            `json:"dependencies" bson:"dependencies"`
	Database     *DatabaseRequirement    `json:"database,omitempty" bson:"database,omitempty"`
	Services     []ServiceRequirement    `json:"services,omitempty" bson:"services,omitempty"`
	Environment  []EnvironmentVariable   `json:"environment,omitempty" bson:"environment,omitempty"`
	Permissions  []PermissionRequirement `json:"permissions,omitempty" bson:"permissions,omitempty"`
	CreatedAt    time.Time               `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at" bson:"updated_at"`
}

// Dependency 基础依赖
type Dependency struct {
	Name        string `json:"name" bson:"name"`
	Version     string `json:"version" bson:"version"`
	Type        string `json:"type" bson:"type"` // package, service, database, etc.
	Required    bool   `json:"required" bson:"required"`
	Description string `json:"description" bson:"description"`
}

// DatabaseRequirement 数据库需求
type DatabaseRequirement struct {
	Type         string            `json:"type" bson:"type"` // mongodb, mysql, postgres, sqlite
	Version      string            `json:"version,omitempty" bson:"version,omitempty"`
	DatabaseName string            `json:"database_name" bson:"database_name"`
	Collections  []CollectionInfo  `json:"collections,omitempty" bson:"collections,omitempty"`
	Tables       []TableInfo       `json:"tables,omitempty" bson:"tables,omitempty"`
	Config       map[string]string `json:"config,omitempty" bson:"config,omitempty"`
	Required     bool              `json:"required" bson:"required"`
}

// CollectionInfo MongoDB集合信息
type CollectionInfo struct {
	Name        string            `json:"name" bson:"name"`
	Indexes     []IndexInfo       `json:"indexes,omitempty" bson:"indexes,omitempty"`
	Schema      map[string]string `json:"schema,omitempty" bson:"schema,omitempty"`
	Description string            `json:"description" bson:"description"`
}

// TableInfo SQL表信息
type TableInfo struct {
	Name        string       `json:"name" bson:"name"`
	Columns     []ColumnInfo `json:"columns" bson:"columns"`
	Indexes     []string     `json:"indexes,omitempty" bson:"indexes,omitempty"`
	Description string       `json:"description" bson:"description"`
}

// IndexInfo 索引信息
type IndexInfo struct {
	Name    string                 `json:"name" bson:"name"`
	Fields  map[string]int         `json:"fields" bson:"fields"` // field: 1 (asc) or -1 (desc)
	Unique  bool                   `json:"unique" bson:"unique"`
	Options map[string]interface{} `json:"options,omitempty" bson:"options,omitempty"`
}

// ColumnInfo 列信息
type ColumnInfo struct {
	Name        string `json:"name" bson:"name"`
	Type        string `json:"type" bson:"type"`
	Nullable    bool   `json:"nullable" bson:"nullable"`
	Default     string `json:"default,omitempty" bson:"default,omitempty"`
	Description string `json:"description" bson:"description"`
}

// ServiceRequirement 服务需求
type ServiceRequirement struct {
	Name        string            `json:"name" bson:"name"`
	Type        string            `json:"type" bson:"type"` // api, queue, cache, etc.
	Version     string            `json:"version,omitempty" bson:"version,omitempty"`
	Host        string            `json:"host,omitempty" bson:"host,omitempty"`
	Port        int               `json:"port,omitempty" bson:"port,omitempty"`
	Config      map[string]string `json:"config,omitempty" bson:"config,omitempty"`
	Required    bool              `json:"required" bson:"required"`
	Description string            `json:"description" bson:"description"`
}

// EnvironmentVariable 环境变量需求
type EnvironmentVariable struct {
	Name        string `json:"name" bson:"name"`
	Required    bool   `json:"required" bson:"required"`
	Default     string `json:"default,omitempty" bson:"default,omitempty"`
	Description string `json:"description" bson:"description"`
	Type        string `json:"type" bson:"type"` // string, int, bool, url, etc.
}

// PermissionRequirement 权限需求
type PermissionRequirement struct {
	Name        string `json:"name" bson:"name"`
	Type        string `json:"type" bson:"type"` // read, write, admin, etc.
	Resource    string `json:"resource" bson:"resource"`
	Description string `json:"description" bson:"description"`
}

// DependencyCheckResult 依赖检查结果
type DependencyCheckResult struct {
	PluginKey     string              `json:"plugin_key"`
	OverallStatus string              `json:"overall_status"` // success, warning, error
	Dependencies  []DependencyStatus  `json:"dependencies"`
	Database      *DatabaseStatus     `json:"database,omitempty"`
	Services      []ServiceStatus     `json:"services,omitempty"`
	Environment   []EnvironmentStatus `json:"environment,omitempty"`
	Permissions   []PermissionStatus  `json:"permissions,omitempty"`
	Suggestions   []string            `json:"suggestions,omitempty"`
	CanInstall    bool                `json:"can_install"`
	RequiresSetup bool                `json:"requires_setup"`
}

// DependencyStatus 依赖状态
type DependencyStatus struct {
	Dependency Dependency `json:"dependency"`
	Status     string     `json:"status"` // available, missing, version_mismatch
	Message    string     `json:"message"`
	Current    string     `json:"current,omitempty"`
}

// DatabaseStatus 数据库状态
type DatabaseStatus struct {
	Requirement     DatabaseRequirement   `json:"requirement"`
	Status          string                `json:"status"` // available, missing, accessible, setup_required
	Message         string                `json:"message"`
	CurrentDatabase string                `json:"current_database,omitempty"`
	CanCreate       bool                  `json:"can_create"`
	SetupOptions    *DatabaseSetupOptions `json:"setup_options,omitempty"`
}

// DatabaseSetupOptions 数据库设置选项
type DatabaseSetupOptions struct {
	SuggestedDatabaseName string            `json:"suggested_database_name"`
	AvailableDatabases    []string          `json:"available_databases"`
	CreateNewDatabase     bool              `json:"create_new_database"`
	UseExistingDatabase   bool              `json:"use_existing_database"`
	Config                map[string]string `json:"config,omitempty"`
}

// ServiceStatus 服务状态
type ServiceStatus struct {
	Requirement ServiceRequirement `json:"requirement"`
	Status      string             `json:"status"` // available, missing, unreachable
	Message     string             `json:"message"`
}

// EnvironmentStatus 环境变量状态
type EnvironmentStatus struct {
	Requirement EnvironmentVariable `json:"requirement"`
	Status      string              `json:"status"` // set, missing, invalid
	Message     string              `json:"message"`
	Current     string              `json:"current,omitempty"`
}

// PermissionStatus 权限状态
type PermissionStatus struct {
	Requirement PermissionRequirement `json:"requirement"`
	Status      string                `json:"status"` // granted, denied, unknown
	Message     string                `json:"message"`
}
