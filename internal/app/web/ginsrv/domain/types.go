package domain

import (
	"github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/model"
)

type ListDomainReq struct {
	web.Pagination
	Name string `json:"name" form:"name" binding:""`
}

type ListDomainResp struct {
	List  []model.TDomain `json:"list"`
	Total int64           `json:"total"`
}

type AddDomainReq struct {
	Name  string `json:"name" binding:"gt=0"`
	Intro string `json:"intro" binding:"gt=0"`
}

type EditDomainReq struct {
	Name  string `json:"name" binding:"gt=0"`
	Intro string `json:"intro" binding:"gt=0"`
}
