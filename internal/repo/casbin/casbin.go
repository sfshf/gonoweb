package casbin_repo

import (
	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"gorm.io/gorm"
)

func FirstGByRsub(rsub string) (*TCasbinRule, error) {
	var record TCasbinRule
	if err := repo.GormDB.
		Table(TableNameTCasbinRule).
		Where("ptype=g").
		Where("v0=?", rsub).
		First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &record, nil
}

// 获取API的policies：v3 存值（GET、POST、DELETE等）
func FindApiPByDomainAndRole(domain, role string) ([]TCasbinRule, error) {
	var list []TCasbinRule
	if err := repo.GormDB.
		Table(TableNameTCasbinRule).
		Where("ptype=p").
		Where("v0=?", role).
		Where("v1=?", domain).
		Where("v3!=''").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// 获取非API的policies：v3 存空字符串
func FindNonApiPByDomainAndRole(domain, role string) ([]TCasbinRule, error) {
	var list []TCasbinRule
	if err := repo.GormDB.
		Table(TableNameTCasbinRule).
		Where("ptype=p").
		Where("v0=?", role).
		Where("v1=?", domain).
		Where("v3=''").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
