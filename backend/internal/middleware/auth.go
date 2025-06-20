package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"vite-pluginend/pkg/errors"
	"vite-pluginend/pkg/utils"
)

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			_, errResp := errors.NewErrorResponse(errors.NewError("未授权访问", http.StatusUnauthorized))
			errResp.Write(c.Writer)
			c.Abort()
			return
		}

		// 检查token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			_, errResp := errors.NewErrorResponse(errors.NewError("无效的token格式", http.StatusUnauthorized))
			errResp.Write(c.Writer)
			c.Abort()
			return
		}

		// 验证token
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			_, errResp := errors.NewErrorResponse(errors.NewError("无效的token", http.StatusUnauthorized))
			errResp.Write(c.Writer)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware 角色中间件
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户角色
		role, exists := c.Get("role")
		if !exists {
			_, errResp := errors.NewErrorResponse(errors.NewError("未授权访问", http.StatusUnauthorized))
			errResp.Write(c.Writer)
			c.Abort()
			return
		}

		// 检查角色权限
		hasRole := false
		for _, r := range roles {
			if r == role.(string) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			_, errResp := errors.NewErrorResponse(errors.NewError("禁止访问", http.StatusForbidden))
			errResp.Write(c.Writer)
			c.Abort()
			return
		}

		c.Next()
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用下一个处理器
		next.ServeHTTP(w, r)
	})
} 