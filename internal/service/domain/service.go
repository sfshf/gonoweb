package domain_svc

import (
	"fmt"

	"github.com/rs/xid"
	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	domain_repo "github.com/sfshf/gonoweb/internal/repo/domain"
	. "github.com/sfshf/gonoweb/internal/service"
)

func ListDomain(page, pageSize int) ([]TDomain, int64, *SvcErr) {
	db := repo.GormDB.Table(TableNameTDomain)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	var list []TDomain
	if err := db.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	return list, total, nil
}

func DomainInfo(xid string) (*TDomain, *SvcErr) {
	domain, err := domain_repo.FirstByXid(xid)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if domain == nil {
		return nil, &SvcErr{Err: fmt.Errorf("域租户[xid=%s]不存在", xid)}
	}
	return domain, nil
}

func AddDomain(name, intro string) (*TDomain, *SvcErr) {
	// 搜索有没有重复的、删除的记录
	domain, err := domain_repo.FirstUnscopedByName(name)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if domain == nil {
		// 新增
		domain := &TDomain{
			Xid:   xid.New().String(),
			Name:  name,
			Intro: intro,
		}
		if err := repo.Create(domain); err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	} else {
		// 如果是活域租户则报错
		if !domain.DeletedAt.Valid {
			return nil, &SvcErr{Err: fmt.Errorf("域租户[name=%s]已存在", name)}
		}
		// 如果是死域租户则激活
		if err := domain_repo.ReliveByXid(domain.Xid, &TDomain{
			Intro: intro,
		}); err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	}
	return domain, nil
}

func EditDomain(xid, name, intro string) *SvcErr {
	if err := domain_repo.UpdateByXid(xid, &TDomain{
		Name:  name,
		Intro: intro,
	}); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}

func DeleteDomain(xid string) *SvcErr {
	if err := domain_repo.DeleteByXid(xid); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}
