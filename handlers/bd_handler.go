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
		c.PureJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	batchUuid := c.Query("batchUuid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	resp, err := h.service.GetPaginatedImages(c.Request.Context(), token, batchUuid, page, limit)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(http.StatusOK, resp)
}

func (h *BdHandler) GetUserBatchesWithCovers(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.PureJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	resp, err := h.service.GetUserBatchesWithCovers(c.Request.Context(), token)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(http.StatusOK, resp)
}

func (h *BdHandler) GetImageMetrics(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.PureJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	imageUuid := c.Param("image_uuid")
	if imageUuid == "" {
		c.PureJSON(http.StatusBadRequest, gin.H{"error": "image_uuid is required"})
		return
	}

	resp, err := h.service.GetImageMetrics(c.Request.Context(), token, imageUuid)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(http.StatusOK, resp)
}
