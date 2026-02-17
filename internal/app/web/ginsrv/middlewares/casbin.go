package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/config"
	"github.com/sfshf/gonoweb/internal/repo/user"
	"github.com/sfshf/gonoweb/internal/service/casbin"
)

// Casbin middleware
func Casbin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 检查核心服务是否启动
		if casbin.Enforcer == nil {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  "服务组件错误：Casbin服务未启动",
			})
			c.Abort()
			return
		}
		// 2. 从JWT中间件获取用户xid
		jwtClaims := JwtClaims(c)
		if jwtClaims.Subject == "" {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  "服务组件错误：请检查JWT/Casbin中间件是否正确使用",
			})
			c.Abort()
			return
		}
		// 3. 从数据库获取用户记录
		user, err := user.User_FirstByXid(jwtClaims.Subject)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  err.Error(),
			})
			c.Abort()
			return
		}
		if user == nil {
			c.JSON(http.StatusUnauthorized, &web.Response{
				Code: web.ResponseCode_RequestError,
				Msg:  "无有效的用户信息",
			})
			c.Abort()
			return
		}
		if user.Email == config.AppConfig.Root.Email {
			c.Next()
			return
		}
		authorized, err := casbin.Enforcer.
			Enforce(jwtClaims.Subject, jwtClaims.Domain, c.FullPath(), c.Request.Method)
		if err != nil {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  err.Error(),
			})
			c.Abort()
			return
		}
		if !authorized {
			c.JSON(http.StatusUnauthorized, &web.Response{
				Code: web.ResponseCode_RequestError,
				Msg:  "用户无权访问",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
