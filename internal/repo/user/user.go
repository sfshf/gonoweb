package user_repo

import (
	"time"

	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"gorm.io/gorm"
)

func User_FirstByNickname(nickname string) (*TUser, error) {
	var record TUser
	if err := repo.GormDB.
		Table(TableNameTUser).
		Where("nick_name=?", nickname).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}

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

func User_FirstUnscopedByEmail(email string) (*TUser, error) {
	var record TUser
	if err := repo.GormDB.
		Table(TableNameTUser).
		Unscoped().
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

func User_FirstByXid(xid string) (*TUser, error) {
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

func User_UpdateByXid(xid string, m *TUser) error {
	m.UpdatedAt = time.Now()
	return repo.GormDB.
		Table(TableNameTUser).
		Where("xid=?", xid).
		Updates(m).Error
}

func User_ReliveByXid(xid string, m *TUser) error {
	db := repo.GormDB
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Table(TableNameTUser).
			Where("xid=?", xid).
			Updates(m).Error; err != nil {
			return err
		}
		// 激活登录
		return tx.Unscoped().
			Table(TableNameTUser).
			Where("xid=?", xid).
			Update("deleted_at", gorm.DeletedAt{}).Error
	})
}

func User_DeleteByXid(xid string) error {
	return repo.GormDB.Delete(&TUser{}, "xid=?", xid).Error
}
