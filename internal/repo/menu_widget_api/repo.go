package mwa_repo

import (
	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
)

type MenuWidgetApiType = int32

const (
	MenuWidgetApiType_Menu   = 1
	MenuWidgetApiType_Widget = 2
	MenuWidgetApiType_API    = 3
)

func FindAllMenuWidgets() ([]TMenuWidgetAPI, error) {
	var list []TMenuWidgetAPI
	if err := repo.GormDB.
		Table(TableNameTMenuWidgetAPI).
		Where("(t_menu_widget_api.type=? OR t_menu_widget_api.type=?)",
			MenuWidgetApiType_Menu,
			MenuWidgetApiType_Widget,
		).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func FindMenuWidgetsByDomainAndRole(domain, role string) ([]TMenuWidgetAPI, error) {
	var list []TMenuWidgetAPI
	if err := repo.GormDB.
		Table(TableNameTMenuWidgetAPI).
		Joins(`LEFT JOIN t_casbin_rule ON t_menu_widget_api.identifier=t_casbin_rule.v2`).
		Where("t_casbin_rule.ptype=p").
		Where("t_casbin_rule.v0=?", role).
		Where("t_casbin_rule.v1=?", domain).
		Where("t_casbin_rule.v3=''").
		Where("(t_menu_widget_api.type=? OR t_menu_widget_api.type=?)",
			MenuWidgetApiType_Menu,
			MenuWidgetApiType_Widget,
		).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
