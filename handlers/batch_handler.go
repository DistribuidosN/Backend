package handlers

import (
	"Backend/models/interfaces/ports"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BatchHandler struct {
	batchService ports.BatchService
}

func NewBatchHandler(s ports.BatchService) *BatchHandler {
	return &BatchHandler{batchService: s}
}

func (h *BatchHandler) DownloadBatch(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.PureJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	batchId := c.Param("id")
	if batchId == "" {
		c.PureJSON(http.StatusBadRequest, gin.H{"error": "missing batch id"})
		return
	}

	resp, err := h.batchService.DownloadBatch(c.Request.Context(), token, batchId)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(http.StatusOK, resp)
}
