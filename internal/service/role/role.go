package role

import (
	"fmt"

	"github.com/rs/xid"
	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	role "github.com/sfshf/gonoweb/internal/repo/role"
	. "github.com/sfshf/gonoweb/internal/service"
)

func ListRole(page, pageSize int, wheres map[string][]any) ([]TRole, int64, *SvcErr) {
	db := repo.GormDB.Table(TableNameTRole)
	for query, args := range wheres {
		db = db.Where(query, args...)
	}
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
	role, err := role.FirstByXid(xid)
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
	roleM, err := role.FirstUnscopedByName(name)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if roleM == nil {
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
		if roleM.DeletedAt == 0 {
			return nil, &SvcErr{Err: fmt.Errorf("角色[name=%s]已存在", name)}
		}
		// 如果是死角色则激活
		if err := role.ReliveByXid(roleM.Xid, &TRole{
			Intro: intro,
		}); err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	}
	return roleM, nil
}

func EditRole(xid, name, intro string) *SvcErr {
	if err := role.UpdateByXid(xid, &TRole{
		Name:  name,
		Intro: intro,
	}); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}

func DeleteRole(xid string) *SvcErr {
	if err := role.DeleteByXid(xid); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}
