package services

import (
	"errors"

	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/mlog"
)

type ApiKeyService struct {
}

func (s *ApiKeyService) GetList(resource_type, resource_id_field, resource_id string) (keys []*models.ApiToken, err error) {
	err = dbengine.Instance().DB.Model(&models.ApiToken{}).Where("type = ?", resource_type).Where(resource_id_field+"=?", resource_id).Find(&keys).Error
	return
}

func (s *ApiKeyService) KeyCount(resource_type, resource_id_field, resource_id string) int64 {
	var count int64
	err := dbengine.Instance().DB.Model(&models.ApiToken{}).Where("type = ?", resource_type).Where(resource_id_field+"=?", resource_id).Count(&count).Error
	if err != nil {
		mlog.Errorf("count api key failed:%v", err)
	}
	return count
}

func (s *ApiKeyService) Create(val *models.ApiToken) (err error) {
	if val == nil {
		return errors.New("invalid args")
	}
	err = dbengine.Instance().DB.Create(val).Error
	return err
}

func (s *ApiKeyService) Get(id, resource_type, resource_id_field, resource_id string) (*models.ApiToken, error) {
	key := new(models.ApiToken)
	err := dbengine.Instance().DB.Model(&models.ApiToken{}).Where("id = ? and type = ?", id, resource_type).Where(resource_id_field+"=?", resource_id).First(key).Error
	if err != nil {
		mlog.Errorf("get api key failed:%v", err)
		return nil, err
	}
	return key, nil
}

func (s *ApiKeyService) Delete(id string) error {
	if err := dbengine.Instance().DB.Model(&models.ApiToken{}).Delete("id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
