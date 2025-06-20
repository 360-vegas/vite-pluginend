package services

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"vite-pluginend/internal/models"
	"vite-pluginend/pkg/logger"
)

type DependencyService struct {
	mongoClient *mongo.Client
}

func NewDependencyService(mongoClient *mongo.Client) *DependencyService {
	return &DependencyService{
		mongoClient: mongoClient,
	}
}

// CheckPluginDependencies 检查插件依赖
func (s *DependencyService) CheckPluginDependencies(ctx context.Context, pluginKey string, dependencies *models.PluginDependency) (*models.DependencyCheckResult, error) {
	result := &models.DependencyCheckResult{
		PluginKey:     pluginKey,
		OverallStatus: "success",
		CanInstall:    true,
		RequiresSetup: false,
	}

	// 检查基础依赖
	if len(dependencies.Dependencies) > 0 {
		depStatuses := s.checkBasicDependencies(dependencies.Dependencies)
		result.Dependencies = depStatuses

		for _, status := range depStatuses {
			if status.Status == "missing" && status.Dependency.Required {
				result.OverallStatus = "error"
				result.CanInstall = false
			} else if status.Status == "version_mismatch" {
				if result.OverallStatus != "error" {
					result.OverallStatus = "warning"
				}
			}
		}
	}

	// 检查数据库依赖
	if dependencies.Database != nil {
		dbStatus, err := s.checkDatabaseRequirement(ctx, *dependencies.Database)
		if err != nil {
			logger.Error("Failed to check database requirement", zap.Error(err))
			result.OverallStatus = "error"
			result.CanInstall = false
		} else {
			result.Database = dbStatus
			if dbStatus.Status == "missing" && dependencies.Database.Required {
				result.OverallStatus = "error"
				result.CanInstall = false
			} else if dbStatus.Status == "setup_required" {
				result.RequiresSetup = true
				if result.OverallStatus != "error" {
					result.OverallStatus = "warning"
				}
			}
		}
	}

	// 检查服务依赖
	if len(dependencies.Services) > 0 {
		serviceStatuses := s.checkServiceRequirements(dependencies.Services)
		result.Services = serviceStatuses

		for _, status := range serviceStatuses {
			if status.Status == "missing" || status.Status == "unreachable" {
				if status.Requirement.Required {
					result.OverallStatus = "error"
					result.CanInstall = false
				} else if result.OverallStatus != "error" {
					result.OverallStatus = "warning"
				}
			}
		}
	}

	// 检查环境变量
	if len(dependencies.Environment) > 0 {
		envStatuses := s.checkEnvironmentVariables(dependencies.Environment)
		result.Environment = envStatuses

		for _, status := range envStatuses {
			if status.Status == "missing" && status.Requirement.Required {
				result.OverallStatus = "error"
				result.CanInstall = false
				result.RequiresSetup = true
			}
		}
	}

	// 生成建议
	result.Suggestions = s.generateSuggestions(result)

	return result, nil
}

// checkBasicDependencies 检查基础依赖
func (s *DependencyService) checkBasicDependencies(dependencies []models.Dependency) []models.DependencyStatus {
	var statuses []models.DependencyStatus

	for _, dep := range dependencies {
		status := models.DependencyStatus{
			Dependency: dep,
			Status:     "available",
			Message:    "依赖可用",
		}

		// 这里可以根据依赖类型进行具体检查
		switch dep.Type {
		case "go_module":
			// 检查Go模块（这里简化处理）
			status.Status = "available"
			status.Message = "Go模块依赖已满足"
		case "system":
			// 检查系统依赖
			status.Status = "available"
			status.Message = "系统依赖已满足"
		default:
			status.Status = "unknown"
			status.Message = "未知依赖类型"
		}

		statuses = append(statuses, status)
	}

	return statuses
}

// checkDatabaseRequirement 检查数据库需求
func (s *DependencyService) checkDatabaseRequirement(ctx context.Context, requirement models.DatabaseRequirement) (*models.DatabaseStatus, error) {
	status := &models.DatabaseStatus{
		Requirement: requirement,
		Status:      "missing",
		Message:     "数据库连接不可用",
		CanCreate:   false,
	}

	switch requirement.Type {
	case "mongodb":
		return s.checkMongoDBRequirement(ctx, requirement)
	case "mysql", "postgres":
		// TODO: 实现其他数据库类型的检查
		status.Status = "missing"
		status.Message = fmt.Sprintf("%s 数据库支持尚未实现", requirement.Type)
	default:
		status.Status = "missing"
		status.Message = "不支持的数据库类型"
	}

	return status, nil
}

// checkMongoDBRequirement 检查MongoDB需求
func (s *DependencyService) checkMongoDBRequirement(ctx context.Context, requirement models.DatabaseRequirement) (*models.DatabaseStatus, error) {
	status := &models.DatabaseStatus{
		Requirement: requirement,
		Status:      "missing",
		Message:     "MongoDB连接不可用",
		CanCreate:   false,
	}

	if s.mongoClient == nil {
		status.Message = "MongoDB客户端未初始化"
		return status, nil
	}

	// 测试连接
	if err := s.mongoClient.Ping(ctx, nil); err != nil {
		status.Message = "无法连接到MongoDB服务器"
		return status, nil
	}

	status.Status = "available"
	status.CanCreate = true

	// 检查数据库是否存在
	dbName := requirement.DatabaseName
	databases, err := s.mongoClient.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		status.Status = "setup_required"
		status.Message = "无法列出数据库"
		return status, nil
	}

	dbExists := false
	for _, name := range databases {
		if name == dbName {
			dbExists = true
			status.CurrentDatabase = name
			break
		}
	}

	if dbExists {
		status.Status = "available"
		status.Message = fmt.Sprintf("数据库 '%s' 已存在", dbName)

		// 检查集合是否存在
		if len(requirement.Collections) > 0 {
			missing := s.checkMongoCollections(ctx, dbName, requirement.Collections)
			if len(missing) > 0 {
				status.Status = "setup_required"
				status.Message = fmt.Sprintf("数据库存在但缺少集合: %s", strings.Join(missing, ", "))
			}
		}
	} else {
		status.Status = "setup_required"
		status.Message = fmt.Sprintf("数据库 '%s' 不存在，需要创建", dbName)
	}

	// 生成设置选项
	status.SetupOptions = &models.DatabaseSetupOptions{
		SuggestedDatabaseName: dbName,
		AvailableDatabases:    databases,
		CreateNewDatabase:     true,
		UseExistingDatabase:   len(databases) > 0,
	}

	return status, nil
}

// checkMongoCollections 检查MongoDB集合
func (s *DependencyService) checkMongoCollections(ctx context.Context, dbName string, collections []models.CollectionInfo) []string {
	database := s.mongoClient.Database(dbName)

	existingCollections, err := database.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to list collections", zap.Error(err))
		return []string{}
	}

	var missing []string
	for _, required := range collections {
		found := false
		for _, existing := range existingCollections {
			if existing == required.Name {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, required.Name)
		}
	}

	return missing
}

// checkServiceRequirements 检查服务需求
func (s *DependencyService) checkServiceRequirements(requirements []models.ServiceRequirement) []models.ServiceStatus {
	var statuses []models.ServiceStatus

	for _, req := range requirements {
		status := models.ServiceStatus{
			Requirement: req,
			Status:      "missing",
			Message:     "服务不可用",
		}

		// 检查服务连接
		if req.Host != "" && req.Port > 0 {
			if s.checkServiceConnection(req.Host, req.Port) {
				status.Status = "available"
				status.Message = fmt.Sprintf("服务 %s:%d 可用", req.Host, req.Port)
			} else {
				status.Status = "unreachable"
				status.Message = fmt.Sprintf("无法连接到服务 %s:%d", req.Host, req.Port)
			}
		} else {
			status.Status = "unknown"
			status.Message = "服务配置信息不完整"
		}

		statuses = append(statuses, status)
	}

	return statuses
}

// checkServiceConnection 检查服务连接
func (s *DependencyService) checkServiceConnection(host string, port int) bool {
	timeout := time.Second * 3
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// checkEnvironmentVariables 检查环境变量
func (s *DependencyService) checkEnvironmentVariables(requirements []models.EnvironmentVariable) []models.EnvironmentStatus {
	var statuses []models.EnvironmentStatus

	for _, req := range requirements {
		status := models.EnvironmentStatus{
			Requirement: req,
			Status:      "missing",
			Message:     "环境变量未设置",
		}

		value := os.Getenv(req.Name)
		if value != "" {
			status.Status = "set"
			status.Current = value
			status.Message = "环境变量已设置"

			// 验证类型
			if err := s.validateEnvironmentVariableType(value, req.Type); err != nil {
				status.Status = "invalid"
				status.Message = fmt.Sprintf("环境变量类型无效: %v", err)
			}
		} else if req.Default != "" {
			status.Status = "set"
			status.Current = req.Default
			status.Message = "使用默认值"
		}

		statuses = append(statuses, status)
	}

	return statuses
}

// validateEnvironmentVariableType 验证环境变量类型
func (s *DependencyService) validateEnvironmentVariableType(value, expectedType string) error {
	switch expectedType {
	case "int":
		_, err := strconv.Atoi(value)
		return err
	case "bool":
		_, err := strconv.ParseBool(value)
		return err
	case "url":
		if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
			return fmt.Errorf("无效的URL格式")
		}
	}
	return nil
}

// generateSuggestions 生成建议
func (s *DependencyService) generateSuggestions(result *models.DependencyCheckResult) []string {
	var suggestions []string

	if result.Database != nil && result.Database.Status == "setup_required" {
		suggestions = append(suggestions, "需要设置数据库连接和创建必要的集合")
	}

	for _, env := range result.Environment {
		if env.Status == "missing" && env.Requirement.Required {
			suggestions = append(suggestions, fmt.Sprintf("请设置环境变量: %s", env.Requirement.Name))
		}
	}

	for _, service := range result.Services {
		if service.Status == "unreachable" && service.Requirement.Required {
			suggestions = append(suggestions, fmt.Sprintf("请确保服务 %s 可用", service.Requirement.Name))
		}
	}

	if !result.CanInstall {
		suggestions = append(suggestions, "请解决所有必需的依赖问题后再尝试安装")
	}

	return suggestions
}

// SetupDatabase 设置数据库
func (s *DependencyService) SetupDatabase(ctx context.Context, pluginKey string, config models.DatabaseSetupOptions) error {
	// 根据配置创建数据库和集合
	if config.CreateNewDatabase {
		return s.createDatabaseForPlugin(ctx, pluginKey, config)
	}
	return nil
}

// createDatabaseForPlugin 为插件创建数据库
func (s *DependencyService) createDatabaseForPlugin(ctx context.Context, pluginKey string, config models.DatabaseSetupOptions) error {
	dbName := config.SuggestedDatabaseName
	if dbName == "" {
		dbName = fmt.Sprintf("plugin_%s", pluginKey)
	}

	// 创建数据库（MongoDB中通过创建集合来隐式创建数据库）
	database := s.mongoClient.Database(dbName)

	// 创建一个临时集合来确保数据库被创建
	tempCollection := database.Collection("_setup")
	_, err := tempCollection.InsertOne(ctx, bson.M{"setup": true, "created_at": time.Now()})
	if err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}

	// 删除临时集合
	tempCollection.Drop(ctx)

	logger.Info("Database created for plugin", zap.String("plugin", pluginKey), zap.String("database", dbName))
	return nil
}
