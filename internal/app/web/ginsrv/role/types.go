package role

import (
	gono_web "github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/model"
)

type ListRoleReq struct {
	gono_web.Pagination
	Name string `json:"name" form:"name" binding:""`
}

type ListRoleResp struct {
	List  []model.TRole `json:"list"`
	Total int64         `json:"total"`
}

type AddRoleReq struct {
	Name  string `json:"name" binding:"gt=0"`
	Intro string `json:"intro" binding:"gt=0"`
}

type EditRoleReq struct {
	Name  string `json:"name" binding:"gt=0"`
	Intro string `json:"intro" binding:"gt=0"`
}
