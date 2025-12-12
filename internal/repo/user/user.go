package user_repo

import (
	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"gorm.io/gorm"
)

func User_FirstByEmail(email string) (*TUser, error) {
	var record TUser
	if err := repo.GormDB.
		Table(TableNameTUser).
		Where("email=?", email).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}

func User_FirstByXID(xid string) (*TUser, error) {
	var record TUser
	if err := repo.GormDB.
		Table(TableNameTUser).
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
