package resource

import (
	gono_web "github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/model"
)

type ListResourceReq struct {
	gono_web.Pagination
	Type       int    `json:"type" form:"type" binding:""`
	Name       string `json:"name" form:"name" binding:""`
	Identifier string `json:"identifier" form:"identifier" binding:""`
}

type ListResourceResp struct {
	List  []model.TResource `json:"list"`
	Total int64             `json:"total"`
}

type AddResourceReq struct {
	Type       int32  `json:"type" binding:"oneof=1 2 3"`
	Identifier string `json:"identifier" binding:"gt=0"`
	Name       string `json:"name" binding:"gt=0"`
	Intro      string `json:"intro" binding:"gt=0"`
	Icon       string `json:"icon" binding:""`
}

type EditResourceReq struct {
	Name  string `json:"name" binding:"gt=0"`
	Intro string `json:"intro" binding:"gt=0"`
	Icon  string `json:"icon" binding:""`
}
