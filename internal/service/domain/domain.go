package domain

import (
	"fmt"

	"github.com/rs/xid"
	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"github.com/sfshf/gonoweb/internal/repo/domain"
	. "github.com/sfshf/gonoweb/internal/service"
)

func ListDomain(page, pageSize int, wheres map[string][]any) ([]TDomain, int64, *SvcErr) {
	db := repo.GormDB.Table(TableNameTDomain)
	for query, args := range wheres {
		db = db.Where(query, args...)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	var list []TDomain
	rows, err := db.Where(`deleted_at=0`).Rows()
	if err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	defer rows.Close()
	for rows.Next() {
		var record TDomain
		if err := db.ScanRows(rows, &record); err != nil {
			return nil, total, &SvcErr{Internal: true, Err: err}
		}
		list = append(list, record)
	}
	return list, total, nil
}

func DomainInfo(xid string) (*TDomain, *SvcErr) {
	domain, err := domain.FirstByXid(xid)
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
	domainM, err := domain.FirstUnscopedByName(name)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if domainM == nil {
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
		if domainM.DeletedAt == 0 {
			return nil, &SvcErr{Err: fmt.Errorf("域租户[name=%s]已存在", name)}
		}
		// 如果是死域租户则激活
		if err := domain.ReliveByXid(domainM.Xid, &TDomain{
			Intro: intro,
		}); err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	}
	return domainM, nil
}

func EditDomain(xid, name, intro string) *SvcErr {
	if err := domain.UpdateByXid(xid, &TDomain{
		Name:  name,
		Intro: intro,
	}); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}

func DeleteDomain(xid string) *SvcErr {
	if err := domain.DeleteByXid(xid); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}
