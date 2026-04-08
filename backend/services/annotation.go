package services

import (
	"fmt"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"

	uuid "github.com/satori/go.uuid"
)

type AnnotationService struct{}

func (s *AnnotationService) GetAnnotationList(appID string, page, limit int, keyword string) ([]*models.MessageAnnotation, int64) {
	var annotations []*models.MessageAnnotation
	var total int64
	query := dbengine.Instance().DB.Where("app_id = ?", appID)
	if keyword != "" {
		query = query.Where("content LIKE ? OR question LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query.Model(&models.MessageAnnotation{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&annotations)
	return annotations, total
}

func (s *AnnotationService) CreateAnnotation(appID, question, content, accountID string) (*models.MessageAnnotation, error) {
	annotation := &models.MessageAnnotation{
		ID:        uuid.NewV4().String(),
		AppID:     appID,
		Content:   content,
		Question:  question,
		AccountID: accountID,
	}
	if err := dbengine.Instance().DB.Create(annotation).Error; err != nil {
		return nil, err
	}
	return annotation, nil
}

func (s *AnnotationService) UpdateAnnotation(annotationID, question, content string) (*models.MessageAnnotation, error) {
	var annotation models.MessageAnnotation
	if err := dbengine.Instance().DB.Where("id = ?", annotationID).First(&annotation).Error; err != nil {
		return nil, fmt.Errorf("annotation not found")
	}
	updates := map[string]any{"content": content, "question": question}
	dbengine.Instance().DB.Model(&annotation).Updates(updates)
	return &annotation, nil
}

func (s *AnnotationService) DeleteAnnotation(appID, annotationID string) error {
	return dbengine.Instance().DB.Where("id = ? AND app_id = ?", annotationID, appID).Delete(&models.MessageAnnotation{}).Error
}

func (s *AnnotationService) DeleteAnnotationsInBatch(appID string, annotationIDs []string) error {
	return dbengine.Instance().DB.Where("id IN ? AND app_id = ?", annotationIDs, appID).Delete(&models.MessageAnnotation{}).Error
}

func (s *AnnotationService) GetAnnotationByID(annotationID string) *models.MessageAnnotation {
	var annotation models.MessageAnnotation
	if err := dbengine.Instance().DB.Where("id = ?", annotationID).First(&annotation).Error; err != nil {
		return nil
	}
	return &annotation
}

func (s *AnnotationService) GetAnnotationHitHistories(appID, annotationID string, page, limit int) ([]*models.AppAnnotationHitHistory, int64) {
	var histories []*models.AppAnnotationHitHistory
	var total int64
	query := dbengine.Instance().DB.Where("app_id = ? AND annotation_id = ?", appID, annotationID)
	query.Model(&models.AppAnnotationHitHistory{}).Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&histories)
	return histories, total
}

func (s *AnnotationService) GetAppAnnotationSetting(appID string) *models.AppAnnotationSetting {
	var setting models.AppAnnotationSetting
	if err := dbengine.Instance().DB.Where("app_id = ?", appID).First(&setting).Error; err != nil {
		return nil
	}
	return &setting
}

func (s *AnnotationService) EnableAppAnnotation(appID string, scoreThreshold float64, embeddingModelProvider, embeddingModelName string) (*models.AppAnnotationSetting, error) {
	setting := s.GetAppAnnotationSetting(appID)
	if setting != nil {
		dbengine.Instance().DB.Model(setting).Updates(map[string]any{
			"score_threshold": scoreThreshold,
		})
		return setting, nil
	}
	setting = &models.AppAnnotationSetting{
		ID:             uuid.NewV4().String(),
		AppID:          appID,
		ScoreThreshold: scoreThreshold,
	}
	if err := dbengine.Instance().DB.Create(setting).Error; err != nil {
		return nil, err
	}
	return setting, nil
}

func (s *AnnotationService) DisableAppAnnotation(appID string) error {
	return dbengine.Instance().DB.Where("app_id = ?", appID).Delete(&models.AppAnnotationSetting{}).Error
}
