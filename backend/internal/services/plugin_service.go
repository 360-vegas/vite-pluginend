package services

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"vite-pluginend/pkg/cache"
	customerrors "vite-pluginend/pkg/errors"
	"vite-pluginend/pkg/logger"
)

// Plugin 插件模型
type Plugin struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Version     string             `bson:"version" json:"version"`
	Description string             `bson:"description" json:"description"`
	Author      string             `bson:"author" json:"author"`
	FileID      string             `bson:"file_id" json:"file_id"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// PluginService 插件服务
type PluginService struct {
	db    *mongo.Database
	cache *cache.RedisCache
}

// NewPluginService 创建插件服务
func NewPluginService(db *mongo.Database, cache *cache.RedisCache) *PluginService {
	return &PluginService{
		db:    db,
		cache: cache,
	}
}

// CreatePlugin 创建插件
func (s *PluginService) CreatePlugin(ctx context.Context, plugin *Plugin) error {
	plugin.ID = primitive.NewObjectID()
	plugin.CreatedAt = time.Now()
	plugin.UpdatedAt = time.Now()

	_, err := s.db.Collection("plugins").InsertOne(ctx, plugin)
	if err != nil {
		logger.Error("创建插件失败", zap.Error(err))
		return customerrors.NewError("创建插件失败", http.StatusInternalServerError)
	}

	if err := s.cache.Delete(ctx, "plugins:list"); err != nil {
		logger.Error("清除缓存失败", zap.Error(err))
	}

	return nil
}

// GetPlugin 获取插件
func (s *PluginService) GetPlugin(ctx context.Context, id string) (*Plugin, error) {
	cacheKey := "plugins:" + id
	var plugin Plugin
	if err := s.cache.Get(ctx, cacheKey, &plugin); err == nil {
		return &plugin, nil
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, customerrors.NewError("无效的插件ID", http.StatusBadRequest)
	}

	err = s.db.Collection("plugins").FindOne(ctx, bson.M{"_id": objectID}).Decode(&plugin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, customerrors.NewError("插件不存在", http.StatusNotFound)
		}
		logger.Error("获取插件失败", zap.Error(err))
		return nil, customerrors.NewError("获取插件失败", http.StatusInternalServerError)
	}

	if err := s.cache.Set(ctx, cacheKey, plugin, 30*time.Minute); err != nil {
		logger.Error("更新缓存失败", zap.Error(err))
	}

	return &plugin, nil
}

// ListPlugins 获取插件列表
func (s *PluginService) ListPlugins(ctx context.Context, page, pageSize int) ([]*Plugin, int64, error) {
	total, err := s.db.Collection("plugins").CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Error("获取插件总数失败", zap.Error(err))
		return nil, 0, customerrors.NewError("获取插件总数失败", http.StatusInternalServerError)
	}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	findOptions := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := s.db.Collection("plugins").Find(ctx, bson.M{}, findOptions)
	if err != nil {
		logger.Error("获取插件列表失败", zap.Error(err))
		return nil, 0, customerrors.NewError("获取插件列表失败", http.StatusInternalServerError)
	}
	defer cursor.Close(ctx)

	var plugins []*Plugin
	if err := cursor.All(ctx, &plugins); err != nil {
		logger.Error("解析插件列表失败", zap.Error(err))
		return nil, 0, customerrors.NewError("解析插件列表失败", http.StatusInternalServerError)
	}

	return plugins, total, nil
}

// UpdatePlugin 更新插件
func (s *PluginService) UpdatePlugin(ctx context.Context, id string, update bson.M) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return customerrors.NewError("无效的插件ID", http.StatusBadRequest)
	}

	update["updated_at"] = time.Now()
	result, err := s.db.Collection("plugins").UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		logger.Error("更新插件失败", zap.Error(err))
		return customerrors.NewError("更新插件失败", http.StatusInternalServerError)
	}

	if result.MatchedCount == 0 {
		return customerrors.NewError("插件不存在", http.StatusNotFound)
	}

	if err := s.cache.Delete(ctx, "plugins:"+id); err != nil {
		logger.Error("清除缓存失败", zap.Error(err))
	}
	if err := s.cache.Delete(ctx, "plugins:list"); err != nil {
		logger.Error("清除缓存失败", zap.Error(err))
	}

	return nil
}

// DeletePlugin 删除插件
func (s *PluginService) DeletePlugin(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return customerrors.NewError("无效的插件ID", http.StatusBadRequest)
	}

	result, err := s.db.Collection("plugins").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		logger.Error("删除插件失败", zap.Error(err))
		return customerrors.NewError("删除插件失败", http.StatusInternalServerError)
	}

	if result.DeletedCount == 0 {
		return customerrors.NewError("插件不存在", http.StatusNotFound)
	}

	if err := s.cache.Delete(ctx, "plugins:"+id); err != nil {
		logger.Error("清除缓存失败", zap.Error(err))
	}
	if err := s.cache.Delete(ctx, "plugins:list"); err != nil {
		logger.Error("清除缓存失败", zap.Error(err))
	}

	return nil
}