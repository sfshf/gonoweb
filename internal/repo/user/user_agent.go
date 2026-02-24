package user

import (
	"time"

	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"gorm.io/gorm"
)

func UserAgent_FirstByToken(token string) (*TUserAgent, error) {
	var record TUserAgent
	if err := repo.GormDB.
		Table(TableNameTUserAgent).
		Where("token=?", token).
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

func UserAgent_FirstByIP(ip string) (*TUserAgent, error) {
	var record TUserAgent
	if err := repo.GormDB.
		Table(TableNameTUserAgent).
		Where("ip=?", ip).
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

func UserAgent_FirstUnscopedByXidAndIP(userXid, ip string) (*TUserAgent, error) {
	var record TUserAgent
	if err := repo.GormDB.Unscoped().
		Table(TableNameTUserAgent).
		Where("ip=?", ip).
		Where("user_xid=?", userXid).
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

func UserAgent_FirstByUserXidAndIP(userXid, ip string) (*TUserAgent, error) {
	var record TUserAgent
	if err := repo.GormDB.
		Table(TableNameTUserAgent).
		Where("ip=?", ip).
		Where("user_xid=?", userXid).
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

func UserAgent_FirstByUserXid(userXid string) (*TUserAgent, error) {
	var record TUserAgent
	if err := repo.GormDB.
		Table(TableNameTUserAgent).
		Where("user_xid=?", userXid).
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

func UserAgent_UpdateByIP(ip string, m *TUserAgent) error {
	return repo.GormDB.
		Table(TableNameTUserAgent).
		Where("ip=?", ip).
		Updates(m).Error
}

func UserAgent_ReliveByXidAndIP(userXid, ip string, m *TUserAgent) error {
	db := repo.GormDB
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Table(TableNameTUserAgent).
			Where("user_xid=?", userXid).
			Where("ip=?", ip).
			Updates(m).Error; err != nil {
			return err
		}
		// 激活登录
		return tx.Unscoped().
			Table(TableNameTUserAgent).
			Where("user_xid=?", userXid).
			Where("ip=?", ip).
			Update("deleted_at", 0).Error
	})
}

// UserAgent_DeleteByToken soft delete record
func UserAgent_DeleteByToken(token string) error {
	return repo.GormDB.
		Table(TableNameTUserAgent).
		Where("token=?", token).
		Update("deleted_at", time.Now().Unix()).Error
}
