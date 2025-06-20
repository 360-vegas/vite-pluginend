package services

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"vite-pluginend/pkg/cache"
	customerrors "vite-pluginend/pkg/errors"
	"vite-pluginend/pkg/logger"
	"vite-pluginend/pkg/utils"
)

// User 用户模型
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string            `bson:"username" json:"username"`
	Password  string            `bson:"password" json:"-"`
	Role      string            `bson:"role" json:"role"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}

// UserService 用户服务
type UserService struct {
	db    *mongo.Database
	cache *cache.RedisCache
}

// NewUserService 创建用户服务
func NewUserService(db *mongo.Database, cache *cache.RedisCache) *UserService {
	return &UserService{
		db:    db,
		cache: cache,
	}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, user *User) error {
	var existing User
	err := s.db.Collection("users").FindOne(ctx, bson.M{"username": user.Username}).Decode(&existing)
	if err == nil {
		return customerrors.NewError("用户名已存在", http.StatusBadRequest)
	} else if err != mongo.ErrNoDocuments {
		logger.Error("检查用户名失败", zap.Error(err))
		return customerrors.NewError("数据库错误", http.StatusInternalServerError)
	}

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		logger.Error("密码加密失败", zap.Error(err))
		return customerrors.NewError("密码加密失败", http.StatusInternalServerError)
	}
	user.Password = hashedPassword

	_, err = s.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		logger.Error("创建用户失败", zap.Error(err))
		return customerrors.NewError("数据库错误", http.StatusInternalServerError)
	}

	return nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, username, password string) (string, error) {
	var user User
	err := s.db.Collection("users").FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", customerrors.NewError("用户名或密码错误", http.StatusUnauthorized)
		}
		logger.Error("查找用户失败", zap.Error(err))
		return "", customerrors.NewError("数据库错误", http.StatusInternalServerError)
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", customerrors.NewError("用户名或密码错误", http.StatusUnauthorized)
	}

	token, err := utils.GenerateToken(user.ID.Hex(), user.Username, user.Role)
	if err != nil {
		logger.Error("生成token失败", zap.Error(err))
		return "", customerrors.NewError("生成token失败", http.StatusInternalServerError)
	}

	return token, nil
}

// GetUser 获取用户信息
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	cacheKey := "users:" + id
	var user User
	if err := s.cache.Get(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, customerrors.NewError("无效的用户ID", http.StatusBadRequest)
	}

	err = s.db.Collection("users").FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, customerrors.NewError("用户不存在", http.StatusNotFound)
		}
		logger.Error("获取用户失败", zap.Error(err))
		return nil, customerrors.NewError("数据库错误", http.StatusInternalServerError)
	}

	s.cache.Set(ctx, cacheKey, user, 30*time.Minute)

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, id string, update User) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return customerrors.NewError("无效的用户ID", http.StatusBadRequest)
	}

	updateMap := bson.M{
		"username":   update.Username,
		"updated_at": time.Now(),
	}

	result, err := s.db.Collection("users").UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": updateMap},
	)
	if err != nil {
		logger.Error("更新用户失败", zap.Error(err))
		return customerrors.NewError("数据库错误", http.StatusInternalServerError)
	}

	if result.MatchedCount == 0 {
		return customerrors.NewError("用户不存在", http.StatusNotFound)
	}

	s.cache.Delete(ctx, "users:"+id)

	return nil
} 