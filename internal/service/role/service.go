package role_svc

import (
	"fmt"

	"github.com/rs/xid"
	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	role_repo "github.com/sfshf/gonoweb/internal/repo/role"
	. "github.com/sfshf/gonoweb/internal/service"
)

func ListRole(page, pageSize int) ([]TRole, int64, *SvcErr) {
	db := repo.GormDB.Table(TableNameTRole)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	var list []TRole
	if err := db.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	return list, total, nil
}

func RoleInfo(xid string) (*TRole, *SvcErr) {
	role, err := role_repo.FirstByXid(xid)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if role == nil {
		return nil, &SvcErr{Err: fmt.Errorf("角色[xid=%s]不存在", xid)}
	}
	return role, nil
}

func AddRole(name, intro string) (*TRole, *SvcErr) {
	// 搜索有没有重复的、删除的记录
	role, err := role_repo.FirstUnscopedByName(name)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if role == nil {
		// 新增
		role := &TRole{
			Xid:   xid.New().String(),
			Name:  name,
			Intro: intro,
		}
		if err := repo.Create(role); err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	} else {
		// 如果是活角色则报错
		if !role.DeletedAt.Valid {
			return nil, &SvcErr{Err: fmt.Errorf("角色[name=%s]已存在", name)}
		}
		// 如果是死角色则激活
		if err := role_repo.ReliveByXid(role.Xid, &TRole{
			Intro: intro,
		}); err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	}
	return role, nil
}

func EditRole(xid, name, intro string) *SvcErr {
	if err := role_repo.UpdateByXid(xid, &TRole{
		Name:  name,
		Intro: intro,
	}); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}

func DeleteRole(xid string) *SvcErr {
	if err := role_repo.DeleteByXid(xid); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}
