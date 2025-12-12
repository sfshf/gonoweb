package domain_repo

import (
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
