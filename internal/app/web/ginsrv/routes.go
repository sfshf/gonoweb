package ginsrv

import (
	"github.com/gin-gonic/gin"
	"github.com/sfshf/gonoweb/internal/app/web/ginsrv/casbin"
	"github.com/sfshf/gonoweb/internal/app/web/ginsrv/domain"
	"github.com/sfshf/gonoweb/internal/app/web/ginsrv/middlewares"
	"github.com/sfshf/gonoweb/internal/app/web/ginsrv/ping"
	"github.com/sfshf/gonoweb/internal/app/web/ginsrv/resource"
	"github.com/sfshf/gonoweb/internal/app/web/ginsrv/role"
	"github.com/sfshf/gonoweb/internal/app/web/ginsrv/user"
)

func LoadRoutes_V1(rg *gin.RouterGroup) {
	// 注意身份认证中间件的加载位置：
	// 不需要身份认证的路由放在前面，需要身份认证的路由放在后面
	// load ping handles
	LoadPingHandles_V1(rg)
	// load public handles
	LoadPubHandles_V1(rg)
	// load auth handles
	rg.Use(middlewares.Jwt(), middlewares.Casbin())
	LoadAuthHandles_V1(rg)
}

// ping handles
func LoadPingHandles_V1(rg *gin.RouterGroup) {
	rg.GET("/ping", ping.Ping)
}

func LoadPubHandles_V1(rg *gin.RouterGroup) {
	// user group
	userg := rg.Group("/user")
	{
		userg.POST("/visit", user.Visit)
		userg.POST("/signIn", user.SignIn)
	}
}

func LoadAuthHandles_V1(rg *gin.RouterGroup) {
	// user group
	userg := rg.Group("/user")
	{
		userg.POST("/signOut", user.SignOut)
		userg.GET("/agent", user.ListUserAgent)
		userg.GET("", user.ListUser)
		userg.GET("/:xid", user.UserInfo)
		userg.POST("", user.AddUser)
		userg.PUT("/:xid", user.EditUser)
		userg.DELETE("/:xid", user.DeleteUser)
	}

	// role group
	roleg := rg.Group("/role")
	{
		roleg.GET("", role.ListRole)
		roleg.GET("/:xid", role.RoleInfo)
		roleg.POST("", role.AddRole)
		roleg.PUT("/:xid", role.EditRole)
		roleg.DELETE("/:xid", role.DeleteRole)
	}

	// domain group
	domaing := rg.Group("/domain")
	{
		domaing.GET("", domain.ListDomain)
		domaing.GET("/:xid", domain.DomainInfo)
		domaing.POST("", domain.AddDomain)
		domaing.PUT("/:xid", domain.EditDomain)
		domaing.DELETE("/:xid", domain.DeleteDomain)
	}

	// resource(menu/widget/api) group
	resourceg := rg.Group("/resource")
	{
		resourceg.GET("", resource.ListResource)
		resourceg.GET("/:id", resource.ResourceInfo)
		resourceg.POST("", resource.AddResource)
		resourceg.PUT("/:id", resource.EditResource)
		resourceg.DELETE("/:id", resource.DeleteResource)
	}

	// casbin group
	casbing := rg.Group("/casbin")
	{
		domg := casbing.Group("/domain")
		{
			roleg := domg.Group("/:dxid/role")
			{
				roleg.GET("", casbin.DomainRoles)
				resourceg := roleg.Group("/:rxid/resource")
				{
					resourceg.GET("", casbin.DomainRoleResources)
					resourceg.POST("", casbin.AllocDomainRoleResources)
				}
			}
		}
		usrg := casbing.Group("/user")
		{
			usrg.POST("/:xid", casbin.AllocRoleInDomain)
		}
	}
}
