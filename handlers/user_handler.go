package handlers

import (
	"Backend/models/interfaces/ports"
	"Backend/models/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(s ports.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	profile, err := h.userService.GetProfile(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"username": profile.Username,
		"status":   profile.Status,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	token := extractToken(c)
	var data user.UserProfile
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.userService.UpdateProfile(c.Request.Context(), token, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.UserUpdateResponse{
		Message:  "profile updated successfully",
		Username: data.Username,
		Valid:    true,
	})
}

func (h *UserHandler) GetActivity(c *gin.Context) {
	token := extractToken(c)
	activities, err := h.userService.GetActivity(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

func (h *UserHandler) SearchUser(c *gin.Context) {
	token := extractToken(c)
	uid := c.Query("uid")
	if uid == "" {
		uid = c.Query("username")
	}
	if uid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uid or username parameter required"})
		return
	}

	user, err := h.userService.SearchUser(c.Request.Context(), token, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	token := extractToken(c)
	if err := h.userService.DeleteAccount(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}

func (h *UserHandler) GetStatistics(c *gin.Context) {
	token := extractToken(c)
	stats, err := h.userService.GetStatistics(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
