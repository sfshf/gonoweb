package ginsrv

import (
	"github.com/gin-gonic/gin"
	casbin_handles "github.com/sfshf/gonoweb/internal/app/web/ginsrv/casbin"
	domain_handles "github.com/sfshf/gonoweb/internal/app/web/ginsrv/domain"
	ginmw "github.com/sfshf/gonoweb/internal/app/web/ginsrv/middlewares"
	ping_handles "github.com/sfshf/gonoweb/internal/app/web/ginsrv/ping"
	role_handles "github.com/sfshf/gonoweb/internal/app/web/ginsrv/role"
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
	userg := rg.Group("/user")
	{
		userg.POST("/visit", user_handles.Visit)
		userg.POST("/signIn", user_handles.SignIn)
	}
}

func LoadAuthHandles_V1(rg *gin.RouterGroup) {
	// user group
	userg := rg.Group("/user")
	{
		userg.POST("/signOut", user_handles.SignOut)
		userg.GET("/agent", user_handles.ListUserAgent)
		userg.GET("", user_handles.ListUser)
		userg.GET("/:xid", user_handles.UserInfo)
		userg.POST("", user_handles.AddUser)
		userg.PUT("/:xid", user_handles.EditUser)
		userg.DELETE("/:xid", user_handles.DeleteUser)
	}

	// role group
	roleg := rg.Group("/role")
	{
		roleg.GET("", role_handles.ListRole)
		roleg.GET("/:xid", role_handles.RoleInfo)
		roleg.POST("", role_handles.AddRole)
		roleg.PUT("/:xid", role_handles.EditRole)
		roleg.DELETE("/:xid", role_handles.DeleteRole)
	}

	// domain group
	domaing := rg.Group("/domain")
	{
		domaing.GET("", domain_handles.ListDomain)
		domaing.GET("/:xid", domain_handles.DomainInfo)
		domaing.POST("", domain_handles.AddDomain)
		domaing.PUT("/:xid", domain_handles.EditDomain)
		domaing.DELETE("/:xid", domain_handles.DeleteDomain)
	}

	// casbin group
	casbing := rg.Group("/casbin")
	{
		casbing.PUT("/role", casbin_handles.AuthorizeRole)
		casbing.PUT("/resource", casbin_handles.AuthorizeResource)
	}
}
