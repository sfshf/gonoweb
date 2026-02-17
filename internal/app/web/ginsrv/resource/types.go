package resource

import (
	gono_web "github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/model"
)

type ListResourceReq struct {
	gono_web.Pagination
	ID   string `json:"id" form:"id" binding:""`
	Name string `json:"name" form:"name" binding:""`
}

type ListResourceResp struct {
	List  []model.TResource `json:"list"`
	Total int64             `json:"total"`
}

type AddResourceReq struct {
	Type  int32  `json:"type" binding:"oneof=1 2 3"`
	ID    string `json:"id" binding:"gt=0"`
	Name  string `json:"name" binding:"gt=0"`
	Intro string `json:"intro" binding:"gt=0"`
	Icon  string `json:"icon" binding:""`
}

type EditResourceReq struct {
	Name  string `json:"name" binding:"gt=0"`
	Intro string `json:"intro" binding:"gt=0"`
	Icon  string `json:"icon" binding:""`
}
