package resource

import (
	"time"

	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"gorm.io/gorm"
)

type ResourceType = int32

const (
	ResourceType_Menu   = 1
	ResourceType_Widget = 2
	ResourceType_API    = 3
)

func FindAllMenuWidgets() ([]TResource, error) {
	var list []TResource
	if err := repo.GormDB.
		Table(TableNameTResource).
		Where("(t_resource.type=? OR t_resource.type=?)",
			ResourceType_Menu,
			ResourceType_Widget,
		).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func FindMenuWidgetsByDomainAndRole(domain, role string) ([]TResource, error) {
	var list []TResource
	if err := repo.GormDB.
		Table(TableNameTResource).
		Joins(`LEFT JOIN t_casbin_rule ON t_resource.identifier=t_casbin_rule.v2`).
		Where("t_casbin_rule.ptype=p").
		Where("t_casbin_rule.v0=?", role).
		Where("t_casbin_rule.v1=?", domain).
		Where("t_casbin_rule.v3=''").
		Where("(t_resource.type=? OR t_resource.type=?)",
			ResourceType_Menu,
			ResourceType_Widget,
		).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func FirstByID(id int64) (*TResource, error) {
	var record TResource
	if err := repo.GormDB.
		Table(TableNameTResource).
		Where("id=?", id).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}

func FindAll() ([]TResource, error) {
	var list []TResource
	if err := repo.GormDB.
		Table(TableNameTResource).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func FirstUnscopedByIdentifier(identifier string) (*TResource, error) {
	var record TResource
	if err := repo.GormDB.
		Table(TableNameTResource).
		Unscoped().
		Where("identifier=?", identifier).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}

func ReliveByIdentifier(identifier string, m *TResource) error {
	db := repo.GormDB
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Table(TableNameTResource).
			Where("identifier=?", identifier).
			Updates(m).Error; err != nil {
			return err
		}
		// 激活登录
		return tx.Unscoped().
			Table(TableNameTResource).
			Where("identifier=?", identifier).
			Update("deleted_at", gorm.DeletedAt{}).Error
	})
}

func UpdateByID(id int64, m *TResource) error {
	m.UpdatedAt = time.Now()
	return repo.GormDB.
		Table(TableNameTResource).
		Where("id=?", id).
		Updates(m).Error
}

func DeleteByIdentifier(id int64) error {
	return repo.GormDB.Delete(&TResource{}, "id=?", id).Error
}
