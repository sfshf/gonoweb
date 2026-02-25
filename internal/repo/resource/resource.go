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
		Where("(type=? OR type=?)",
			ResourceType_Menu,
			ResourceType_Widget,
		).
		Where(`deleted_at=0`).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func FindMenuWidgetsByIdentifiers(identifiers []string) ([]TResource, error) {
	var list []TResource
	if err := repo.GormDB.
		Table(TableNameTResource).
		// TODO 解决identifiers数组太长，导致SQL的IN语句失效的问题
		Where("identifier IN (?)", identifiers).
		Where("(type=? OR type=?)",
			ResourceType_Menu,
			ResourceType_Widget,
		).
		Where(`deleted_at=0`).
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
		Where(`deleted_at=0`).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}

func FirstByIdentifier(identifier string) (*TResource, error) {
	var record TResource
	if err := repo.GormDB.
		Table(TableNameTResource).
		Where("identifier=?", identifier).
		Where(`deleted_at=0`).
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
		Where(`deleted_at=0`).
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
		Where(`deleted_at=0`).
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
			Update("deleted_at", 0).Error
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
	return repo.GormDB.
		Table(TableNameTResource).
		Where("id=?", id).
		Update("deleted_at", time.Now().Unix()).Error
}
