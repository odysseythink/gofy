package services

import (
	"gorm.io/gorm"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/models/request"
)

type DifySetupService struct {
}

// Create 创建DifySetup记录
func (s *DifySetupService) Create(val *models.DifySetup) (err error) {
	err = dbengine.Instance().DB.Create(val).Error
	return err
}

// Delete 删除DifySetup记录
func (s *DifySetupService) Delete(val models.DifySetup) (err error) {
	err = dbengine.Instance().DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.DifySetup{}).Where("version = ?", val.Version).Error; err != nil {
			return err
		}
		if err = tx.Delete(&val).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByIds 批量删除DifySetup记录
func (s *DifySetupService) DeleteByIds(ids []string) (err error) {
	if err := dbengine.Instance().DB.Model(&models.DifySetup{}).Delete("id in ?", ids).Error; err != nil {
		return err
	}

	return nil
}

// Update 更新DifySetup记录
func (s *DifySetupService) Update(val *models.DifySetup) (err error) {
	err = dbengine.Instance().DB.Save(val).Error
	return err
}

func (s *DifySetupService) Get() (*models.DifySetup, error) {
	val := new(models.DifySetup)
	err := dbengine.Instance().DB.First(val).Error
	if err != nil {
		return nil, err
	}
	return val, nil
}

// GetList 分页获取DifySetup记录
func (s *DifySetupService) GetList(info request.PageInfo) (list []models.DifySetup, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := dbengine.Instance().DB.Model(&models.DifySetup{})
	var endpoints []models.DifySetup

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&endpoints).Error
	return endpoints, total, err
}
