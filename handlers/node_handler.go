package handlers

import (
	"Backend/models/interfaces/ports"
	"Backend/models/node"
	"net/http"
	"strings"

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

func (h *NodeHandler) UploadBatch(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "images files are required"})
		return
	}

	filters := normalizeFilters(form.Value["filters"])
	if len(filters) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one filter is required"})
		return
	}

	resp, err := h.nodeService.UploadBatch(c.Request.Context(), token, node.BatchUploadRequest{
		Files:   files,
		Filters: filters,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func normalizeFilters(raw []string) []string {
	var filters []string
	for _, value := range raw {
		for _, item := range strings.Split(value, ",") {
			trimmed := strings.TrimSpace(item)
			if trimmed != "" {
				filters = append(filters, trimmed)
			}
		}
	}
	return filters
}
