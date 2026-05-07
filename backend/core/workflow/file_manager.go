package workflow

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/odysseythink/gofy/backend/storage"
	"github.com/odysseythink/mlog"
)

// FileManager handles file operations within workflow execution.
type FileManager struct {
	signSecret string
}

// NewFileManager creates a file manager.
func NewFileManager(signSecret string) *FileManager {
	if signSecret == "" {
		signSecret = "default-workflow-secret"
	}
	return &FileManager{signSecret: signSecret}
}

// UploadWorkflowFile stores a file produced during workflow execution.
func (fm *FileManager) UploadWorkflowFile(tenantID, workflowRunID, nodeID, filename string, data []byte) (string, error) {
	key := fmt.Sprintf("workflow_files/%s/%s/%s/%s", tenantID, workflowRunID, nodeID, filename)
	if err := storage.Save(key, data); err != nil {
		return "", fmt.Errorf("failed to save workflow file: %w", err)
	}
	return key, nil
}

// GetWorkflowFile retrieves a workflow file.
func (fm *FileManager) GetWorkflowFile(fileKey string) ([]byte, error) {
	return storage.LoadOnce(fileKey)
}

// DeleteWorkflowFiles removes all files for a workflow run.
func (fm *FileManager) DeleteWorkflowFiles(tenantID, workflowRunID string) error {
	// TODO: Implement directory-level deletion in storage backend
	mlog.Infof("cleaning workflow files for run %s", workflowRunID)
	return nil
}

// GenerateSignedURL creates a signed URL for accessing a workflow file.
func (fm *FileManager) GenerateSignedURL(fileKey string, expiresIn time.Duration) string {
	expiry := time.Now().Add(expiresIn).Unix()
	payload := fmt.Sprintf("%s:%d", fileKey, expiry)

	mac := hmac.New(sha256.New, []byte(fm.signSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("/files/workflow/%s?expires=%d&sign=%s", fileKey, expiry, signature)
}

// VerifySignedURL verifies a signed file URL.
func (fm *FileManager) VerifySignedURL(fileKey string, expires int64, signature string) bool {
	if time.Now().Unix() > expires {
		return false
	}
	payload := fmt.Sprintf("%s:%d", fileKey, expires)
	mac := hmac.New(sha256.New, []byte(fm.signSecret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

// ToolFileParser parses tool output files.
type ToolFileParser struct{}

// ParseToolFile extracts file information from tool output.
func (p *ToolFileParser) ParseToolFile(toolOutput map[string]any) *WorkflowFile {
	url, _ := toolOutput["url"].(string)
	filename, _ := toolOutput["filename"].(string)
	mimeType, _ := toolOutput["mime_type"].(string)
	size, _ := toolOutput["size"].(float64)

	if url == "" && filename == "" {
		return nil
	}

	return &WorkflowFile{
		URL:      url,
		Filename: filename,
		MimeType: mimeType,
		Size:     int64(size),
	}
}

// WorkflowFile represents a file in the workflow context.
type WorkflowFile struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	Key      string `json:"key,omitempty"`
}
