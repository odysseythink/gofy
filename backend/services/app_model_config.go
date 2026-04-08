package services

import (
	"gorm.io/gorm"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/request"
)

type AppModelConfigService struct {
}

// Create 创建App记录
func (s *AppModelConfigService) Create(data *models.AppModelConfig) (err error) {
	err = dbengine.Instance().DB.Create(data).Error
	return err
}

// Delete 删除App记录
func (s *AppModelConfigService) Delete(data models.AppModelConfig) (err error) {
	err = dbengine.Instance().DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.AppModelConfig{}).Where("id = ?", data.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&data).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByIds 批量删除App记录
func (s *AppModelConfigService) DeleteByIds(ids []string) (err error) {
	if err := dbengine.Instance().DB.Model(&models.AppModelConfig{}).Delete("id in ?", ids).Error; err != nil {
		return err
	}

	return nil
}

// Update 更新App记录
func (s *AppModelConfigService) Update(data *models.AppModelConfig) (err error) {
	err = dbengine.Instance().DB.Save(data).Error
	return err
}

// Get 根据id获取App记录
func (s *AppModelConfigService) Get(id string) (data models.AppModelConfig, err error) {
	err = dbengine.Instance().DB.Where("id = ?", id).First(&data).Error
	return
}

// GetList 分页获取App记录
func (s *AppModelConfigService) List(info request.PageInfo) (list []models.AppModelConfig, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := dbengine.Instance().DB.Model(&models.AppModelConfig{})
	var endpoints []models.AppModelConfig

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&endpoints).Error
	return endpoints, total, err
}
