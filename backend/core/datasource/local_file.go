package datasource

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalFileProvider struct{}

func NewLocalFileProvider() *LocalFileProvider {
	return &LocalFileProvider{}
}

func (p *LocalFileProvider) Type() DatasourceType {
	return DatasourceTypeLocalFile
}

func (p *LocalFileProvider) ValidateCredentials(credentials map[string]any) error {
	return nil
}

func (p *LocalFileProvider) FetchDocuments(config map[string]any, credentials map[string]any) ([]*DatasourceDocument, error) {
	filePath, ok := config["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("file_path is required")
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	doc := &DatasourceDocument{
		Title:    filepath.Base(filePath),
		MimeType: detectMimeType(filePath),
		Size:     info.Size(),
		Metadata: map[string]any{
			"file_path": filePath,
			"modified":  info.ModTime(),
		},
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	doc.Content = string(content)

	return []*DatasourceDocument{doc}, nil
}

func (p *LocalFileProvider) FetchDocumentContent(config map[string]any, credentials map[string]any, documentID string) (io.ReadCloser, error) {
	filePath, ok := config["file_path"].(string)
	if !ok {
		return nil, fmt.Errorf("file_path is required")
	}
	return os.Open(filePath)
}

func detectMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".txt":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	case ".md":
		return "text/markdown"
	case ".html", ".htm":
		return "text/html"
	case ".csv":
		return "text/csv"
	case ".json":
		return "application/json"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return "application/octet-stream"
	}
}
