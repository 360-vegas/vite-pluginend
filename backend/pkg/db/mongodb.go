package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"

	"vite-pluginend/pkg/logger"
)

// NewMongoClient 创建MongoDB客户端
func NewMongoClient() (*mongo.Client, error) {
	// 获取连接URI
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	// 配置连接选项
	opts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(5 * time.Minute).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(5 * time.Second)

	// 创建客户端
	client, err := mongo.NewClient(opts)
	if err != nil {
		logger.Error("Failed to create MongoDB client", zap.Error(err))
		return nil, fmt.Errorf("failed to create MongoDB client: %v", err)
	}

	// 连接数据库
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		logger.Error("Failed to connect to MongoDB", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// 测试连接
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Error("Failed to ping MongoDB", zap.Error(err))
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	logger.Info("Connected to MongoDB successfully")
	return client, nil
}

// GetCollection 获取集合
func GetCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}

// CreateIndexes 创建索引
func CreateIndexes(ctx context.Context, collection *mongo.Collection, indexes []mongo.IndexModel) error {
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		logger.Error("Failed to create indexes", zap.Error(err))
		return fmt.Errorf("failed to create indexes: %v", err)
	}
	return nil
}

// CreateUniqueIndex 创建唯一索引
func CreateUniqueIndex(ctx context.Context, collection *mongo.Collection, field string) error {
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		logger.Error("Failed to create unique index", zap.Error(err))
		return fmt.Errorf("failed to create unique index: %v", err)
	}
	return nil
}

// CreateTextIndex 创建文本索引
func CreateTextIndex(ctx context.Context, collection *mongo.Collection, fields ...string) error {
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: "text"})
	}
	index := mongo.IndexModel{
		Keys: keys,
	}
	_, err := collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		logger.Error("Failed to create text index", zap.Error(err))
		return fmt.Errorf("failed to create text index: %v", err)
	}
	return nil
}

// CreateTTLIndex 创建TTL索引
func CreateTTLIndex(ctx context.Context, collection *mongo.Collection, field string, expireAfterSeconds int32) error {
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(expireAfterSeconds),
	}
	_, err := collection.Indexes().CreateOne(ctx, index)
	if err != nil {
		logger.Error("Failed to create TTL index", zap.Error(err))
		return fmt.Errorf("failed to create TTL index: %v", err)
	}
	return nil
} 