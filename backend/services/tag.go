package services

import (
	"fmt"
	"time"

	"github.com/odysseythink/gofy/backend/core/exceptions"
	dbengine "github.com/odysseythink/gofy/backend/db_engine"
	"github.com/odysseythink/gofy/backend/models"
	"github.com/odysseythink/gofy/backend/models/request"
	"github.com/odysseythink/gofy/backend/proto/pbapi"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type TagService struct {
}

// Create 创建Tag记录
func (s *TagService) Create(tag models.Tag) (err error) {
	err = dbengine.Instance().DB.Create(&tag).Error
	return err
}

// Delete 删除Tag记录
func (s *TagService) Delete(tag models.Tag) (err error) {
	err = dbengine.Instance().DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Tag{}).Where("id = ?", tag.ID).Error; err != nil {
			return err
		}
		if err = tx.Delete(&tag).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteByIds 批量删除Tag记录
func (s *TagService) DeleteByIds(ids []string) (err error) {
	if err := dbengine.Instance().DB.Model(&models.Tag{}).Delete("id in ?", ids).Error; err != nil {
		return err
	}

	return nil
}

// Get 根据id获取Tag记录
func (s *TagService) Get(id string) (tag *models.Tag, err error) {
	tag = &models.Tag{}
	err = dbengine.Instance().DB.Where("id = ?", id).First(tag).Error
	if err != nil {
		tag = nil
	}
	return
}

// GetList 分页获取Tag记录
func (s *TagService) GetList(info request.PageInfo) (list []models.Tag, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := dbengine.Instance().DB.Model(&models.Tag{})
	var endpoints []models.Tag

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&endpoints).Error
	return endpoints, total, err
}

func (s *TagService) GetTargetIDsByIDs(tag_type, current_tenant_id string, tag_ids []string) []string {
	localTagIDs := make([]string, 0)
	err := dbengine.Instance().DB.Model(&models.Tag{}).Select("id").Where("`id` in ? and tenant_id = ? and `type` = ?", tag_ids, current_tenant_id, tag_type).Find(&localTagIDs).Error
	if err != nil || len(localTagIDs) < 1 {
		mlog.Errorf("get tags by(ids=%#v, tenant_id=%s, type=%s) failed:%v", tag_ids, current_tenant_id, tag_type, err)
		return nil
	}
	tag_ids = append(tag_ids, localTagIDs...)

	tag_bindings := make([]string, 0)
	err = dbengine.Instance().DB.Model(&models.TagBinding{}).Select("target_id").Where("`tag_id` in ? and tenant_id = ?", tag_ids, current_tenant_id).Find(&tag_bindings).Error
	if err != nil || len(tag_bindings) < 1 {
		mlog.Errorf("get TagBinding by(ids=%#v, tenant_id=%s, type=%s) failed:%v", tag_ids, current_tenant_id, tag_type, err)
		return nil
	}

	return tag_bindings
}

func (s *TagService) GetTags(tag_type, current_tenant_id, keyword string) []*pbapi.TagField {
	var results []*pbapi.TagField

	sql := fmt.Sprintf("select t.id, t.name, t.`type`, count(tb.id) as binding_count from tags t left join tag_bindings tb on tb.tag_id = t.id where t.`type`='%s' and t.tenant_id ='%s'", tag_type, current_tenant_id)
	if keyword != "" {
		sql += fmt.Sprintf(" t.name like '%%%s%%'", keyword)
		// db = db.Where("t.name like ?", "%"+keyword+"%")
	}
	sql += " group by t.id order by t.created_at DESC;"
	mlog.Debugf("---sql=%s", sql)
	err := dbengine.Instance().DB.Raw(sql).Find(&results).Error
	if err != nil {
		mlog.Errorf("exec sql(%s) failed:%v", sql, err)
		return nil
	}
	return results
}
func (s *TagService) SaveTag(current_user *models.Account, name, tag_type string) *models.Tag {
	now := time.Now()
	tag := &models.Tag{
		ID:        uuid.NewV4().String(),
		TenantID:  current_user.CurrentTenantID(),
		Type:      tag_type,
		Name:      name,
		CreatedBy: current_user.ID,
		CreatedAt: &now,
	}
	dbengine.Instance().DB.Create(tag)
	return tag
}
func (s *TagService) UpdateTag(name, tag_id string) *models.Tag {
	tag := new(models.Tag)
	if err := dbengine.Instance().DB.Model(&models.Tag{}).Where("id = ?", tag_id).First(tag).Error; err != nil {
		mlog.Errorf("Tag not found")
		panic(exceptions.NewValueError("Tag not found"))
	}
	dbengine.Instance().DB.Updates(&models.Tag{ID: tag_id, Name: name})
	tag.Name = name
	return tag
}
func (s *TagService) GetTagBindingCount(tag_id string) int64 {
	var count int64
	if err := dbengine.Instance().DB.Model(&models.TagBinding{}).Where("tag_id = ?", tag_id).Count(&count).Error; err != nil {
		mlog.Errorf("count TagBinding failed:%v", err)
		return 0
	}
	return count
}
func (s *TagService) CheckTargetExists(tenant_id, tag_type, target_id string) any {
	switch tag_type {
	case "knowledge":
		dataset := new(models.Dataset)
		if err := dbengine.Instance().DB.Model(&models.Dataset{}).Where("tenant_id = ? and id = ?", tenant_id, target_id).First(dataset).Error; err != nil {
			mlog.Errorf("get Dataset failed:%v", err)
			panic(exceptions.NewValueError("Dataset not found"))
		}
		return dataset
	case "app":
		app := new(models.App)
		if err := dbengine.Instance().DB.Model(&models.App{}).Where("tenant_id = ? and id = ?", tenant_id, target_id).First(app).Error; err != nil {
			mlog.Errorf("get App failed:%v", err)
			panic(exceptions.NewValueError("App not found"))
		}
		return app
	default:
		panic(exceptions.NewValueError("Invalid binding type"))

	}
}
func (s *TagService) AddTagBinding(current_user *models.Account, tag_type, target_id string, tag_ids []string) {
	// check if target exists
	s.CheckTargetExists(current_user.CurrentTenantID(), tag_type, target_id)
	// save tag binding
	for _, tag_id := range tag_ids {
		tag_binding := new(models.TagBinding)
		if err := dbengine.Instance().DB.Model(&models.TagBinding{}).Where("tenant_id = ? and tag_id = ?", current_user.CurrentTenantID(), tag_id).First(tag_binding).Error; err != nil {
			mlog.Errorf("get TagBinding failed:%v", err)
			tag_binding = nil
		}
		if tag_binding != nil {
			continue
		}
		now := time.Now()
		new_tag_binding := &models.TagBinding{
			ID:        uuid.NewV4().String(),
			TenantID:  current_user.CurrentTenantID(),
			TagID:     tag_id,
			TargetID:  target_id,
			CreatedBy: current_user.ID,
			CreatedAt: &now,
		}
		dbengine.Instance().DB.Create(new_tag_binding)
	}
}
func (s *TagService) DeleteTagBinding(current_user *models.Account, tag_id, tag_type, target_id string) {
	// check if target exists
	s.CheckTargetExists(current_user.CurrentTenantID(), tag_type, target_id)
	// delete tag binding
	tag_binding := new(models.TagBinding)
	if err := dbengine.Instance().DB.Model(&models.TagBinding{}).Where("target_id = ? and tag_id = ?", target_id, tag_id).First(tag_binding).Error; err != nil {
		mlog.Errorf("get TagBinding failed:%v", err)
		tag_binding = nil
	}

	if tag_binding != nil {
		dbengine.Instance().DB.Model(&models.Tag{}).Delete("id in ?", []string{tag_binding.ID})
	}
}
