package ginsrv

import (
	"github.com/gin-gonic/gin"
	ginmw "github.com/sfshf/gonoweb/internal/app/web/ginsrv/middlewares"
	ping_handles "github.com/sfshf/gonoweb/internal/app/web/ginsrv/ping"
	user_handles "github.com/sfshf/gonoweb/internal/app/web/ginsrv/user"
)

func LoadRoutes_V1(rg *gin.RouterGroup) {
	// 注意身份认证中间件的加载位置：
	// 不需要身份认证的路由放在前面，需要身份认证的路由放在后面
	// load ping handles
	LoadPingHandles_V1(rg)
	// load public handles
	LoadPubHandles_V1(rg)
	// load auth handles
	rg.Use(ginmw.Jwt(), ginmw.Casbin())
	LoadAuthHandles_V1(rg)
}

// ping handles
func LoadPingHandles_V1(rg *gin.RouterGroup) {
	rg.GET("/ping", ping_handles.Ping)
}

func LoadPubHandles_V1(rg *gin.RouterGroup) {
	// user group
	ug := rg.Group("/user")
	{
		ug.POST("/visit", user_handles.Visit)
		ug.POST("/signIn", user_handles.SignIn)
	}
}

func LoadAuthHandles_V1(rg *gin.RouterGroup) {
	// user group
	ug := rg.Group("/user")
	{
		ug.POST("/signOut", user_handles.SignOut)
	}
}
