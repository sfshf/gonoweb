package casbin

import (
	"errors"
	"fmt"
	"time"

	"github.com/casbin/casbin/v3"
	casbinModel "github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
	"github.com/sfshf/gonoweb/internal/config"
	"github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"gorm.io/gorm"
)

var (
	Enforcer *casbin.SyncedEnforcer
)

func Launch() (func(), error) {
	opt := config.AppConfig.Gin.Casbin
	// 1. 检查依赖项有没有加载成功
	if repo.GormDB == nil {
		return nil, errors.New("系统组件错误：初始化Casbin服务，缺少核心组件")
	}
	// 2. 加载casbin模型
	m, err := casbinModel.NewModelFromFile(opt.Model)
	if err != nil {
		return nil, err
	}
	// 3. 创建Enforcer实例
	Enforcer, err = casbin.NewSyncedEnforcer(m, Adapter(func() *gorm.DB { return repo.GormDB }))
	if err != nil {
		return nil, err
	}
	// 4. 装配Enforcer选项
	Enforcer.EnableLog(opt.Log)
	if opt.AutoLoad {
		Enforcer.StartAutoLoadPolicy(time.Second * time.Duration(opt.AutoLoadInterval))
	}
	Enforcer.EnableAutoSave(opt.AutoSave)
	Enforcer.EnableEnforce(opt.Enforce)
	if err = Enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	return func() {
		Enforcer.SavePolicy()
	}, nil
}

// 自定义 Casbin Adapter

var _ persist.Adapter = (Adapter)(nil)
var _ persist.BatchAdapter = (Adapter)(nil)

// A implementation of Adapter, BatchAdapter, FilteredAdapter interfaces of github.com/casbin/casbin/v3/persist package.
type Adapter func() *gorm.DB

func loadPolicyLine(line *model.TCasbinRule, m casbinModel.Model) {
	var p string
	if line.Ptype != "" && line.V0 != "" {
		p += line.Ptype + ", " + line.V0
	} else {
		return
	}
	if line.V1 != "" {
		p += ", " + line.V1
	}
	if line.V2 != "" {
		p += ", " + line.V2
	}
	if line.V3 != "" {
		p += ", " + line.V3
	}
	if line.V4 != "" {
		p += ", " + line.V4
	}
	if line.V5 != "" {
		p += ", " + line.V5
	}
	persist.LoadPolicyLine(p, m)
}

func lineToModel(ptype string, rule []string) *model.TCasbinRule {
	m := &model.TCasbinRule{
		Ptype: ptype,
	}
	if len(rule) > 0 {
		m.V0 = rule[0]
	}
	if len(rule) > 1 {
		m.V1 = rule[1]
	}
	if len(rule) > 2 {
		m.V2 = rule[2]
	}
	if len(rule) > 3 {
		m.V3 = rule[3]
	}
	if len(rule) > 4 {
		m.V4 = rule[4]
	}
	if len(rule) > 5 {
		m.V5 = rule[5]
	}
	return m
}

// LoadPolicy loads all policy rules from the storage.
func (a Adapter) LoadPolicy(m casbinModel.Model) error {
	db := a()
	// load enabled policies.
	rows, err := db.Model(&model.TCasbinRule{}).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var p model.TCasbinRule
		if err := db.ScanRows(rows, &p); err != nil {
			return err
		}
		loadPolicyLine(&p, m)
	}
	return rows.Close()
}

// SavePolicy saves all policy rules to the storage.
func (a Adapter) SavePolicy(m casbinModel.Model) error {
	var ms []*model.TCasbinRule
	for ptype, ast := range m["p"] {
		for _, rule := range ast.Policy {
			m := lineToModel(ptype, rule)
			ms = append(ms, m)
		}
	}
	for ptype, ast := range m["g"] {
		for _, rule := range ast.Policy {
			m := lineToModel(ptype, rule)
			ms = append(ms, m)
		}
	}
	db := a()
	return db.Transaction(func(tx *gorm.DB) error {
		// remove all old records
		if err := tx.Unscoped().Delete(&model.TCasbinRule{}, "1=1").Error; err != nil {
			return err
		}
		// insert all new records
		if len(ms) > 0 {
			return tx.Create(ms).Error
		}
		return nil
	})
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return a().Create(lineToModel(ptype, rule)).Error
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	m := lineToModel(ptype, rule)
	return a().Unscoped().Delete(
		&model.TCasbinRule{},
		`ptype=? AND v0=? AND v1=? AND v2=? AND v3=? AND v4=? AND v5=?`,
		m.Ptype, m.V0, m.V1, m.V2, m.V3, m.V4, m.V5).Error
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return a().Transaction(func(tx *gorm.DB) error {
		if len(fieldValues) > 0 {
			return tx.Unscoped().Delete(
				&model.TCasbinRule{},
				fmt.Sprintf("ptype=%s AND v%d=?", ptype, fieldIndex),
				fieldValues[0],
			).Error
		}
		return nil
	})

}

// AddPolicies adds policy rules to the storage.
// This is part of the Auto-Save feature.
func (a Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	var ms []*model.TCasbinRule
	for _, rule := range rules {
		m := lineToModel(ptype, rule)
		ms = append(ms, m)
	}
	return a().Create(ms).Error
}

// RemovePolicies removes policy rules from the storage.
// This is part of the Auto-Save feature.
func (a Adapter) RemovePolicies(sec string, pType string, rules [][]string) error {
	return a().Transaction(func(tx *gorm.DB) error {
		for _, rule := range rules {
			m := lineToModel(pType, rule)
			if err := tx.Unscoped().Delete(
				&model.TCasbinRule{},
				`ptype=? AND v0=? AND v1=? AND v2=? AND v3=? AND v4=? AND v5=?`,
				m.Ptype, m.V0, m.V1, m.V2, m.V3, m.V4, m.V5).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
