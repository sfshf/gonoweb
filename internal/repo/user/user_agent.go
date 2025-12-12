package user_repo

import (
	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"gorm.io/gorm"
)

func UserAgent_FirstByToken(token string) (*TUserAgent, error) {
	var record TUserAgent
	if err := repo.GormDB.
		Table(TableNameTUserAgent).
		Where("token=?", token).
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
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}

func UserAgent_FirstDeletedByIP(ip string) (*TUserAgent, error) {
	var record TUserAgent
	if err := repo.GormDB.
		Table(TableNameTUserAgent).
		Where("ip=?", ip).
		Where("deleted_at IS NOT NULL").
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

func UserAgent_ReliveByIP(ip string, m *TUserAgent) error {
	m.DeletedAt = gorm.DeletedAt{} // 将deleted_at置NULL
	return repo.GormDB.
		Table(TableNameTUserAgent).
		Where("ip=?", ip).
		Where("deleted_at IS NOT NULL").
		Updates(m).Error
}
