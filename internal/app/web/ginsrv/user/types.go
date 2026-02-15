package user

import (
	gono_web "github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/model"
)

type SignInReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ListUserAgentReq struct {
	gono_web.Pagination
}

type ListUserAgentResp struct {
	List  []model.TUserAgent `json:"list"`
	Total int64              `json:"total"`
}

type ListUserReq struct {
	gono_web.Pagination
	Email    string `json:"email" form:"email" binding:""`
	Nickname string `json:"nickname" form:"nickname" binding:""`
	Realname string `json:"realname" form:"realname" binding:""`
}

type ListUserResp struct {
	List  []model.TUser `json:"list"`
	Total int64         `json:"total"`
}

type AddUserReq struct {
	Email    string `json:"email" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type EditUserReq struct {
	Email    string `json:"email" binding:"required"`
	NickName string `json:"nickName" binding:"required"`
}
