package services

import (
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/models/request"
	"gorm.io/gorm"
)

type TenantAccountJoinService struct {
}

// Create 创建TenantAccountJoin记录
func (s *TenantAccountJoinService) Create(ta models.TenantAccountJoin) (err error) {
	err = dbengine.Instance().DB.Create(&ta).Error
	return err
}

// Delete 删除TenantAccountJoin记录
func (s *TenantAccountJoinService) Delete(ta models.TenantAccountJoin) (err error) {
	err = dbengine.Instance().DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.TenantAccountJoin{}).Where("id = ?", ta.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&ta).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByIds 批量删除TenantAccountJoin记录
func (s *TenantAccountJoinService) DeleteByIds(ids []string) (err error) {
	if err := dbengine.Instance().DB.Model(&models.TenantAccountJoin{}).Delete("id in ?", ids).Error; err != nil {
		return err
	}

	return nil
}

// Update 更新TenantAccountJoin记录
func (s *TenantAccountJoinService) Update(ta *models.TenantAccountJoin) (err error) {
	err = dbengine.Instance().DB.Save(ta).Error
	return err
}

// Get 根据id获取TenantAccountJoin记录
func (s *TenantAccountJoinService) Get(id string) (ta *models.TenantAccountJoin, err error) {
	ta = &models.TenantAccountJoin{}
	err = dbengine.Instance().DB.Where("id = ?", id).First(ta).Error
	if err != nil {
		ta = nil
	}
	return
}

// GetList 分页获取TenantAccountJoin记录
func (s *TenantAccountJoinService) GetList(info request.PageInfo) (list []models.TenantAccountJoin, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := dbengine.Instance().DB.Model(&models.TenantAccountJoin{})
	var endpoints []models.TenantAccountJoin

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&endpoints).Error
	return endpoints, total, err
}
