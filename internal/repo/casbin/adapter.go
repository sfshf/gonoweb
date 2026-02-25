package casbin

import (
	"encoding/json"
	"fmt"
	"log"

	casbinModel "github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
	"github.com/sfshf/gonoweb/internal/model"
	"gorm.io/gorm"
)

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
	rows, err := db.Model(&model.TCasbinRule{}).Where("deleted_at=0").Rows()
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
		data, _ := json.MarshalIndent(ast.Policy, "", "\t")
		log.Printf("ptype=%s\tast=%s\n", ptype, data)
		for _, rule := range ast.Policy {
			m := lineToModel(ptype, rule)
			ms = append(ms, m)
		}
	}
	for ptype, ast := range m["g"] {
		data, _ := json.MarshalIndent(ast.Policy, "", "\t")
		log.Printf("ptype=%s\tast=%s\n", ptype, data)
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
				fmt.Sprintf("ptype=? AND v%d=?", fieldIndex),
				ptype,
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
