package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vite-pluginend/internal/services"
	customerrors "vite-pluginend/pkg/errors"
)

// UserHandler 处理用户相关的HTTP请求
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler 创建新的用户处理器
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req services.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, customerrors.NewError("无效的请求参数", http.StatusBadRequest))
		return
	}

	if err := h.userService.Register(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("注册失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "注册成功"})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, customerrors.NewError("无效的请求参数", http.StatusBadRequest))
		return
	}

	token, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, customerrors.NewError("用户名或密码错误", http.StatusUnauthorized))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// GetUser 获取用户信息
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("获取用户信息失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser 更新用户信息
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req services.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, customerrors.NewError("无效的请求参数", http.StatusBadRequest))
		return
	}

	if err := h.userService.UpdateUser(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusInternalServerError, customerrors.NewError("更新用户信息失败", http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
} 