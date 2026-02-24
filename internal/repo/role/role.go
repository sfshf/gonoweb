package role

import (
	"time"

	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"gorm.io/gorm"
)

func FirstByXid(xid string) (*TRole, error) {
	var record TRole
	if err := repo.GormDB.
		Table(TableNameTRole).
		Where("xid=?", xid).
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

func FindAll() ([]TRole, error) {
	var list []TRole
	if err := repo.GormDB.
		Table(TableNameTRole).
		Where(`deleted_at=0`).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func FirstUnscopedByName(name string) (*TRole, error) {
	var record TRole
	if err := repo.GormDB.
		Table(TableNameTRole).
		Unscoped().
		Where("name=?", name).
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

func ReliveByXid(xid string, m *TRole) error {
	db := repo.GormDB
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Table(TableNameTRole).
			Where("xid=?", xid).
			Updates(m).Error; err != nil {
			return err
		}
		// 激活登录
		return tx.Unscoped().
			Table(TableNameTRole).
			Where("xid=?", xid).
			Update("deleted_at", 0).Error
	})
}

func UpdateByXid(xid string, m *TRole) error {
	m.UpdatedAt = time.Now()
	return repo.GormDB.
		Table(TableNameTRole).
		Where("xid=?", xid).
		Updates(m).Error
}

func DeleteByXid(xid string) error {
	return repo.GormDB.
		Table(TableNameTRole).
		Where("xid=?", xid).
		Update("deleted_at", time.Now().Unix()).
		Error
}
