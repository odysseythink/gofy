package services

import (
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/models/request"

	"gorm.io/gorm"
)

type GofySetupService struct {
}

// Create 创建GofySetup记录
func (s *GofySetupService) Create(val *models.GofySetup) (err error) {
	err = dbengine.Instance().DB.Create(val).Error
	return err
}

// Delete 删除GofySetup记录
func (s *GofySetupService) Delete(val models.GofySetup) (err error) {
	err = dbengine.Instance().DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.GofySetup{}).Where("version = ?", val.Version).Error; err != nil {
			return err
		}
		if err = tx.Delete(&val).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByIds 批量删除GofySetup记录
func (s *GofySetupService) DeleteByIds(ids []string) (err error) {
	if err := dbengine.Instance().DB.Model(&models.GofySetup{}).Delete("id in ?", ids).Error; err != nil {
		return err
	}

	return nil
}

// Update 更新GofySetup记录
func (s *GofySetupService) Update(val *models.GofySetup) (err error) {
	err = dbengine.Instance().DB.Save(val).Error
	return err
}

func (s *GofySetupService) Get() (*models.GofySetup, error) {
	val := new(models.GofySetup)
	err := dbengine.Instance().DB.First(val).Error
	if err != nil {
		return nil, err
	}
	return val, nil
}

// GetList 分页获取GofySetup记录
func (s *GofySetupService) GetList(info request.PageInfo) (list []models.GofySetup, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := dbengine.Instance().DB.Model(&models.GofySetup{})
	var endpoints []models.GofySetup

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&endpoints).Error
	return endpoints, total, err
}
