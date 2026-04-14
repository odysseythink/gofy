package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"strings"

	"github.com/gin-gonic/gin"
)

// extractWebhookData extracts webhook data from an HTTP request.
// Handles JSON, form data, multipart form data, text/plain, and binary.
func extractWebhookData(c *gin.Context) (map[string]any, error) {
	result := make(map[string]any)

	// Extract headers
	headers := make(map[string]any)
	for k, v := range c.Request.Header {
		if len(v) == 1 {
			headers[k] = v[0]
		} else {
			headers[k] = v
		}
	}
	result["headers"] = headers

	// Extract query parameters
	queryParams := make(map[string]any)
	for k, v := range c.Request.URL.Query() {
		if len(v) == 1 {
			queryParams[k] = v[0]
		} else {
			queryParams[k] = v
		}
	}
	result["query_params"] = queryParams

	// Extract method
	result["method"] = c.Request.Method

	// Extract body based on content type
	contentType := c.ContentType()
	mediaType, _, _ := mime.ParseMediaType(contentType)

	body, err := extractBody(c, mediaType)
	if err != nil {
		return nil, fmt.Errorf("failed to extract body: %w", err)
	}
	result["body"] = body

	// Extract files from multipart
	if mediaType == "multipart/form-data" {
		files := extractFiles(c)
		if len(files) > 0 {
			result["files"] = files
		}
	}

	return result, nil
}

func extractBody(c *gin.Context, mediaType string) (map[string]any, error) {
	switch mediaType {
	case "application/json":
		return extractJSONBody(c)
	case "application/x-www-form-urlencoded":
		return extractFormBody(c)
	case "multipart/form-data":
		return extractMultipartBody(c)
	case "text/plain":
		return extractTextBody(c)
	case "application/octet-stream":
		return extractBinaryBody(c)
	default:
		// Try JSON first, fall back to raw
		if jsonBody, err := extractJSONBody(c); err == nil {
			return jsonBody, nil
		}
		return extractTextBody(c)
	}
}

func extractJSONBody(c *gin.Context) (map[string]any, error) {
	body := make(map[string]any)
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return body, err
	}
	if len(rawBody) == 0 {
		return body, nil
	}

	var jsonData any
	if err := json.Unmarshal(rawBody, &jsonData); err != nil {
		return body, fmt.Errorf("invalid JSON body: %w", err)
	}

	switch v := jsonData.(type) {
	case map[string]any:
		return v, nil
	default:
		body["raw"] = jsonData
		return body, nil
	}
}

func extractFormBody(c *gin.Context) (map[string]any, error) {
	body := make(map[string]any)
	if err := c.Request.ParseForm(); err != nil {
		return body, err
	}
	for k, v := range c.Request.PostForm {
		if len(v) == 1 {
			body[k] = v[0]
		} else {
			body[k] = v
		}
	}
	return body, nil
}

func extractMultipartBody(c *gin.Context) (map[string]any, error) {
	body := make(map[string]any)
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		return body, nil
	}
	if c.Request.MultipartForm != nil {
		for k, v := range c.Request.MultipartForm.Value {
			if len(v) == 1 {
				body[k] = v[0]
			} else {
				body[k] = v
			}
		}
	}
	return body, nil
}

func extractTextBody(c *gin.Context) (map[string]any, error) {
	body := make(map[string]any)
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return body, err
	}
	body["raw"] = strings.TrimSpace(string(rawBody))
	return body, nil
}

func extractBinaryBody(c *gin.Context) (map[string]any, error) {
	body := make(map[string]any)
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return body, err
	}
	body["raw"] = map[string]any{
		"size":         len(rawBody),
		"content_type": c.ContentType(),
	}
	return body, nil
}

func extractFiles(c *gin.Context) map[string]any {
	files := make(map[string]any)
	if c.Request.MultipartForm == nil || c.Request.MultipartForm.File == nil {
		return files
	}
	for name, headers := range c.Request.MultipartForm.File {
		if len(headers) > 0 {
			fh := headers[0]
			files[name] = map[string]any{
				"filename":     fh.Filename,
				"size":         fh.Size,
				"content_type": fh.Header.Get("Content-Type"),
			}
		}
	}
	return files
}
