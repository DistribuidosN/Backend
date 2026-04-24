package handlers

import (
	"Backend/models/interfaces/ports"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BdHandler struct {
	service ports.BdService
}

func NewBdHandler(s ports.BdService) *BdHandler {
	return &BdHandler{service: s}
}

func (h *BdHandler) GetPaginatedImages(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	batchUuid := c.Query("batchUuid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	resp, err := h.service.GetPaginatedImages(c.Request.Context(), token, batchUuid, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *BdHandler) GetUserBatchesWithCovers(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	resp, err := h.service.GetUserBatchesWithCovers(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
