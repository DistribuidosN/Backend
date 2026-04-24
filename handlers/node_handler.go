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
		c.PureJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	var req node.ImageUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.PureJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resp, err := h.nodeService.UploadImages(c.Request.Context(), token, req)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.PureJSON(http.StatusOK, resp)
}

func (h *NodeHandler) ProcessBatch(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.PureJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.PureJSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.PureJSON(http.StatusBadRequest, gin.H{"error": "no images provided"})
		return
	}

	var transformations []node.Transformation
	for _, f := range form.Value["filters"] {
		transformations = append(transformations, node.Transformation{Name: f})
	}

	resp, err := h.nodeService.ProcessBatch(c.Request.Context(), token, files, transformations)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(http.StatusOK, resp)
}

func (h *NodeHandler) GetBatchStatus(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.PureJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	jobID := c.Param("id")
	if jobID == "" {
		c.PureJSON(http.StatusBadRequest, gin.H{"error": "missing job id"})
		return
	}

	resp, err := h.nodeService.GetBatchStatus(c.Request.Context(), token, jobID)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(http.StatusOK, resp)
}

func (h *NodeHandler) GetBatchResults(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.PureJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	jobID := c.Param("id")
	if jobID == "" {
		c.PureJSON(http.StatusBadRequest, gin.H{"error": "missing job id"})
		return
	}

	resp, err := h.nodeService.GetBatchResults(c.Request.Context(), token, jobID)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(http.StatusOK, resp)
}

func (h *NodeHandler) GetLogsByImage(c *gin.Context) {
	token := extractToken(c)
	imageUuid := c.Param("image_uuid")
	if imageUuid == "" {
		c.PureJSON(http.StatusBadRequest, gin.H{"error": "image_uuid required"})
		return
	}

	logs, err := h.nodeService.GetLogsByImage(c.Request.Context(), token, imageUuid)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(http.StatusOK, logs)
}

func (h *NodeHandler) GetMetricsByNode(c *gin.Context) {
	token := extractToken(c)
	nodeId := c.Param("node_id")
	if nodeId == "" {
		c.PureJSON(http.StatusBadRequest, gin.H{"error": "node_id required"})
		return
	}

	metrics, err := h.nodeService.GetMetricsByNode(c.Request.Context(), token, nodeId)
	if err != nil {
		c.PureJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.PureJSON(http.StatusOK, metrics)
}
