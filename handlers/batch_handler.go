package handlers

import (
	"Backend/models"
	"Backend/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BatchHandler is the handler for batch requests
type BatchHandler struct {
	batchService *services.BatchService
}

// NewBatchHandler creates a new BatchHandler
func NewBatchHandler(batchService *services.BatchService) *BatchHandler {
	return &BatchHandler{batchService: batchService}
}

// UploadBatch is the handler for the POST /api/v1/batch endpoint
func (h *BatchHandler) UploadBatch(c *gin.Context) {
	var req models.BatchRestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.batchService.ProcessBatch(req)
	if err != nil {
		if apiErr, ok := err.(*models.ApiError); ok {
			c.JSON(apiErr.Status, gin.H{"error": apiErr.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
