package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 表示用户信息
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string            `bson:"username" json:"username"`
	Password  string            `bson:"password" json:"-"` // 不在JSON中显示密码
	Email     string            `bson:"email" json:"email"`
	Role      string            `bson:"role" json:"role"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// UserResponse 用户响应
type UserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    *User  `json:"data,omitempty"`
	Token   string `json:"token,omitempty"`
} 