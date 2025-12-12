package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	gono_web "github.com/sfshf/gonoweb/internal/app/web"
	ginmw "github.com/sfshf/gonoweb/internal/app/web/ginsrv/middlewares"
	user_svc "github.com/sfshf/gonoweb/internal/service/user"
)

// Visit 用户首次打开页面，上报用户代理、时区、语言等信息
func Visit(c *gin.Context) {
	// IP  from HTTP headers
	ip := c.ClientIP()
	// User-Agent from HTTP headers
	ua := c.GetHeader("User-Agent")
	// TracdID from HTTP headers
	tid := c.GetHeader(ginmw.HeaderKey_TraceID)
	if tid == "" {
		tid = c.GetString(ginmw.HeaderKey_TraceID)
	}
	if err := user_svc.UpsertUserAgent(ip, ua, tid); err != nil {
		c.JSON(http.StatusInternalServerError, &gono_web.Response{
			Code: gono_web.ResponseCode_InternalError,
			Msg:  fmt.Sprintf("新增/更新用户代理信息失败：%s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, &gono_web.Response{
		Code: gono_web.ResponseCode_OK,
		Msg:  gono_web.ResponseMsg_OK,
	})
}
