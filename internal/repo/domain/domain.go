package domain_repo

import (
	"time"

	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"gorm.io/gorm"
)

func FirstByXid(xid string) (*TDomain, error) {
	var record TDomain
	if err := repo.GormDB.
		Table(TableNameTDomain).
		Where("xid=?", xid).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}

func FindAll() ([]TDomain, error) {
	var list []TDomain
	if err := repo.GormDB.
		Table(TableNameTDomain).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func FirstUnscopedByName(name string) (*TDomain, error) {
	var record TDomain
	if err := repo.GormDB.
		Table(TableNameTDomain).
		Unscoped().
		Where("name=?", name).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}

func ReliveByXid(xid string, m *TDomain) error {
	db := repo.GormDB
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Table(TableNameTDomain).
			Where("xid=?", xid).
			Updates(m).Error; err != nil {
			return err
		}
		// 激活登录
		return tx.Unscoped().
			Table(TableNameTDomain).
			Where("xid=?", xid).
			Update("deleted_at", gorm.DeletedAt{}).Error
	})
}

func UpdateByXid(xid string, m *TDomain) error {
	m.UpdatedAt = time.Now()
	return repo.GormDB.
		Table(TableNameTDomain).
		Where("xid=?", xid).
		Updates(m).Error
}

func DeleteByXid(xid string) error {
	return repo.GormDB.Delete(&TDomain{}, "xid=?", xid).Error
}
