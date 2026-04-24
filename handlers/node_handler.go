package handlers

import (
	"Backend/models/interfaces/ports"
	"Backend/models/node"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func (h *NodeHandler) ProcessBatch(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse multipart form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no images provided"})
		return
	}

	filters := normalizeMultipartFilters(form.Value["filters"])

	resp, err := h.nodeService.ProcessBatch(c.Request.Context(), token, files, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *NodeHandler) GetBatchStatus(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing job id"})
		return
	}

	resp, err := h.nodeService.GetBatchStatus(c.Request.Context(), token, jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *NodeHandler) GetBatchResults(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing job id"})
		return
	}

	resp, err := h.nodeService.GetBatchResults(c.Request.Context(), token, jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *NodeHandler) DownloadBatchResult(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	jobID := c.Param("id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing job id"})
		return
	}

	resp, err := h.nodeService.DownloadBatchResult(c.Request.Context(), token, jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s-results.zip\"", resp.JobID))
	c.Data(http.StatusOK, "application/zip", resp.Content)
}

func normalizeMultipartFilters(raw []string) []string {
	if len(raw) == 0 {
		return nil
	}

	filters := make([]string, 0, len(raw))
	for _, value := range raw {
		for _, part := range strings.Split(value, "\n") {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}
			filters = append(filters, normalizeSingleFilter(trimmed)...)
		}
	}
	return filters
}

type filterPayload struct {
	Name   string `json:"name"`
	Params any    `json:"params"`
}

func normalizeSingleFilter(raw string) []string {
	var payload filterPayload
	if err := json.Unmarshal([]byte(raw), &payload); err != nil || strings.TrimSpace(payload.Name) == "" {
		return []string{raw}
	}

	params := parseFilterParams(payload.Params)
	name := strings.ToLower(strings.TrimSpace(payload.Name))

	switch name {
	case "grayscale", "ocr", "inference":
		return []string{name}
	case "flip":
		return []string{fmt.Sprintf("flip:%s", firstNonEmpty(
			paramString(params, "direction"),
			"horizontal",
		))}
	case "blur":
		return []string{fmt.Sprintf("blur:%s", firstNonEmpty(
			paramString(params, "radius"),
			paramString(params, "sigma"),
			"1.5",
		))}
	case "sharpen":
		return []string{fmt.Sprintf("sharpen:%s", firstNonEmpty(
			paramString(params, "factor"),
			"2.0",
		))}
	case "rotate":
		return []string{fmt.Sprintf(
			"rotate:%s,%s",
			firstNonEmpty(paramString(params, "angle"), "0"),
			firstNonEmpty(paramString(params, "expand"), "true"),
		)}
	case "brightness":
		return []string{fmt.Sprintf("brightness:%s", firstNonEmpty(
			paramString(params, "brightness"),
			"1.0",
		))}
	case "contrast":
		return []string{fmt.Sprintf("contrast:%s", firstNonEmpty(
			paramString(params, "contrast"),
			"1.0",
		))}
	case "brightness_contrast":
		return []string{fmt.Sprintf(
			"brightness_contrast:%s,%s",
			firstNonEmpty(paramString(params, "brightness"), "1.0"),
			firstNonEmpty(paramString(params, "contrast"), "1.0"),
		)}
	case "resize":
		return []string{fmt.Sprintf(
			"resize:%sx%s",
			firstNonEmpty(paramString(params, "width"), "0"),
			firstNonEmpty(paramString(params, "height"), "0"),
		)}
	case "crop":
		return []string{fmt.Sprintf(
			"crop:%s,%s,%s,%s",
			firstNonEmpty(paramString(params, "left"), "0"),
			firstNonEmpty(paramString(params, "upper"), "0"),
			firstNonEmpty(paramString(params, "right"), "0"),
			firstNonEmpty(paramString(params, "lower"), "0"),
		)}
	case "watermark", "watermark_text":
		return []string{fmt.Sprintf(
			"watermark_text:%s|%s|%s|%s|%s|%s|%s|%s|%s|%s",
			firstNonEmpty(paramString(params, "text"), ""),
			firstNonEmpty(paramString(params, "x"), "16"),
			firstNonEmpty(paramString(params, "y"), "16"),
			firstNonEmpty(paramString(params, "fill"), paramString(params, "color"), "white"),
			firstNonEmpty(paramString(params, "size"), "36"),
			firstNonEmpty(paramString(params, "stroke_width"), "2"),
			firstNonEmpty(paramString(params, "opacity"), "96"),
			firstNonEmpty(paramString(params, "angle"), "-28"),
			firstNonEmpty(paramString(params, "spacing_x"), "220"),
			firstNonEmpty(paramString(params, "spacing_y"), "160"),
		)}
	case "format":
		format := firstNonEmpty(paramString(params, "value"), paramString(params, "format"))
		if format == "" {
			return []string{raw}
		}
		return []string{fmt.Sprintf("format:%s", strings.ToLower(format))}
	default:
		return []string{raw}
	}
}

func parseFilterParams(raw any) map[string]any {
	switch value := raw.(type) {
	case nil:
		return map[string]any{}
	case string:
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return map[string]any{}
		}
		var parsed map[string]any
		if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
			return parsed
		}
		return map[string]any{"value": trimmed}
	case map[string]any:
		return value
	default:
		return map[string]any{}
	}
}

func paramString(params map[string]any, key string) string {
	value, ok := params[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		if typed {
			return "true"
		}
		return "false"
	default:
		return strings.TrimSpace(fmt.Sprintf("%v", typed))
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
