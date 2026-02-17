package resource

import (
	"fmt"

	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"github.com/sfshf/gonoweb/internal/repo/resource"
	. "github.com/sfshf/gonoweb/internal/service"
)

func ListResource(page, pageSize int, wheres map[string][]any) ([]TResource, int64, *SvcErr) {
	db := repo.GormDB.Table(TableNameTResource)
	for query, args := range wheres {
		db = db.Where(query, args...)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	var list []TResource
	if err := db.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	return list, total, nil
}

func ResourceInfo(id int64) (*TResource, *SvcErr) {
	domain, err := resource.FirstByID(id)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if domain == nil {
		return nil, &SvcErr{Err: fmt.Errorf("菜单/控件/API[id=%s]不存在", id)}
	}
	return domain, nil
}

func AddResource(typ int32, identifier, name, intro, icon string) (*TResource, *SvcErr) {
	// 搜索有没有重复的、删除的记录
	mwa, err := resource.FirstUnscopedByIdentifier(identifier)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if mwa == nil {
		// 新增
		mwa = &TResource{
			Type:       typ,
			Identifier: identifier,
			Name:       name,
			Intro:      intro,
			Icon:       icon,
		}
		if err := repo.Create(mwa); err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	} else {
		// 如果是活域租户则报错
		if mwa.DeletedAt == 0 {
			return nil, &SvcErr{Err: fmt.Errorf("域租户[name=%s]已存在", name)}
		}
		// 如果是死域租户则激活
		if err := resource.ReliveByIdentifier(mwa.Identifier, &TResource{
			Intro: intro,
		}); err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	}
	return mwa, nil
}

func EditResource(id int64, name, intro, icon string) *SvcErr {
	if err := resource.UpdateByID(id, &TResource{
		Name:  name,
		Intro: intro,
	}); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}

func DeleteResource(id int64) *SvcErr {
	if err := resource.DeleteByIdentifier(id); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}
