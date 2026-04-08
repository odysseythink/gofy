package services

import (
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"time"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/storage"

	uuid "github.com/satori/go.uuid"
)

type FileService struct{}

func (s *FileService) UploadFile(tenantID, userID, createdByRole string, file io.Reader, filename string, fileSize int64, mimeType string) (*models.UploadFile, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	key := fmt.Sprintf("uploads/%s/%s/%s", tenantID, time.Now().Format("20060102"), filename)
	if err := storage.Save(key, data); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}
	ext := filepath.Ext(filename)
	if mimeType == "" {
		mimeType = mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
	}
	record := &models.UploadFile{
		ID:            uuid.NewV4().String(),
		TenantID:      tenantID,
		Name:          filename,
		Size:          int(fileSize),
		Extension:     ext,
		MimeType:      mimeType,
		Key:           key,
		CreatedBy:     userID,
		CreatedByRole: createdByRole,
	}
	if err := dbengine.Instance().DB.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (s *FileService) GetFileByID(fileID string) *models.UploadFile {
	var f models.UploadFile
	if err := dbengine.Instance().DB.Where("id = ?", fileID).First(&f).Error; err != nil {
		return nil
	}
	return &f
}

func (s *FileService) GetFileContent(fileID string) ([]byte, error) {
	f := s.GetFileByID(fileID)
	if f == nil {
		return nil, fmt.Errorf("file not found")
	}
	return storage.LoadOnce(f.Key)
}

func (s *FileService) GetFileStream(fileID string) (io.ReadCloser, error) {
	f := s.GetFileByID(fileID)
	if f == nil {
		return nil, fmt.Errorf("file not found")
	}
	return storage.LoadStream(f.Key)
}

func (s *FileService) DeleteFile(fileID string) error {
	f := s.GetFileByID(fileID)
	if f == nil {
		return fmt.Errorf("file not found")
	}
	storage.Delete(f.Key)
	return dbengine.Instance().DB.Delete(f).Error
}

func (s *FileService) GetUploadFilesByIDs(tenantID string, fileIDs []string) []*models.UploadFile {
	var files []*models.UploadFile
	dbengine.Instance().DB.Where("tenant_id = ? AND id IN ?", tenantID, fileIDs).Find(&files)
	return files
}

func (s *FileService) IsFileSizeWithinLimit(extension string, fileSize int64) bool {
	maxSize := int64(15 * 1024 * 1024) // 15MB default
	switch extension {
	case ".pdf", ".docx", ".xlsx":
		maxSize = 50 * 1024 * 1024
	case ".png", ".jpg", ".jpeg", ".gif", ".webp":
		maxSize = 10 * 1024 * 1024
	}
	return fileSize <= maxSize
}
