package handlers

import (
	"Backend/models/interfaces/ports"
	"Backend/models/node"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NodeHandler struct {
	nodeService ports.NodeService
}

func NewNodeHandler(s ports.NodeService) *NodeHandler {
	return &NodeHandler{nodeService: s}
}

func (h *NodeHandler) UploadImages(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	var req node.ImageUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.nodeService.UploadImages(c.Request.Context(), token, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
