package role_repo

import (
	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"gorm.io/gorm"
)

func FirstByXid(xid string) (*TRole, error) {
	var record TRole
	if err := repo.GormDB.
		Table(TableNameTRole).
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

func FindAll() ([]TRole, error) {
	var list []TRole
	if err := repo.GormDB.
		Table(TableNameTRole).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
