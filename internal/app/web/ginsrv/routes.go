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
	usrg := rg.Group("/user")
	{
		usrg.POST("/visit", user.Visit)
		usrg.POST("/signIn", user.SignIn)
	}
}

func LoadAuthHandles_V1(rg *gin.RouterGroup) {
	// user group
	usrg := rg.Group("/user")
	{
		usrg.POST("/signOut", user.SignOut)
		usrg.GET("/agent", user.ListUserAgent)
		usrg.GET("", user.ListUser)
		usrg.POST("", user.AddUser)
		usrg.PUT("", user.SwitchRole)
		susrg := usrg.Group("/:xid")
		{
			susrg.GET("", user.UserInfo)
			susrg.PUT("", user.EditUser)
			susrg.DELETE("", user.DeleteUser)
		}
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
	resg := rg.Group("/resource")
	{
		resg.GET("", resource.ListResource)
		resg.GET("/:id", resource.ResourceInfo)
		resg.POST("", resource.AddResource)
		resg.PUT("/:id", resource.EditResource)
		resg.DELETE("/:id", resource.DeleteResource)
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
			susrg := usrg.Group("/:xid")
			{
				susrg.POST("", casbin.AllocRoleInDomain)
				domg := susrg.Group("/domain")
				{
					domg.GET("", casbin.UserDomains)
					roleg := domg.Group("/:dxid/role")
					{
						roleg.GET("", casbin.UserRolesInDomain)
					}
				}
			}
		}
	}
}
