package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	gono_web "github.com/sfshf/gonoweb/internal/app/web"
	ginmw "github.com/sfshf/gonoweb/internal/app/web/ginsrv/middlewares"
	user_svc "github.com/sfshf/gonoweb/internal/service/user"
	jwt_util "github.com/sfshf/gonoweb/internal/util/jwt"
)

type SignInReq struct {
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"required"`
}

// SignIn 用户登录
// @Summary      用户登录
// @Description  用户登录
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param 		 request body SignInReq true "登录所需参数"
// @Success      200  {object}  gono_web.Response
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /user/signIn [POST]
func SignIn(c *gin.Context) {
	// 检查请求入参
	var req SignInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
	}
	// IP  from HTTP headers
	ip := c.ClientIP()
	// User-Agent from HTTP headers
	ua := c.GetHeader("User-Agent")
	// TracdID from HTTP headers
	tid := c.GetHeader(ginmw.HeaderKey_TraceID)
	if tid == "" {
		tid = c.GetString(ginmw.HeaderKey_TraceID)
	}
	// 密码登录
	data, svcErr := user_svc.SignInByPassword(req.Email, req.Password, ip, ua, tid)
	if svcErr != nil {
		if svcErr.Internal {
			c.JSON(http.StatusInternalServerError, &gono_web.Response{
				Code: gono_web.ResponseCode_InternalError,
				Msg:  fmt.Sprintf("系统报错：%s", svcErr.Error()),
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, &gono_web.Response{
				Code: gono_web.ResponseCode_RequestError,
				Msg:  fmt.Sprintf("密码登录失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 将jwt写入头部
	c.Header("Authorization", jwt_util.BearerPrefix+data.Token)
	c.JSON(http.StatusOK, &gono_web.Response{
		Code: gono_web.ResponseCode_OK,
		Msg:  gono_web.ResponseMsg_OK,
		Data: data,
	})
}

// SignOut 用户登出
// @Summary      用户登出
// @Description  用户登出
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Success      200  {object}  gono_web.Response
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /user/signOut [POST]
func SignOut(c *gin.Context) {
	// 从gin.Context拿取用户信息
	claims := ginmw.JwtClaims(c)
	if claims == nil {
		c.JSON(http.StatusOK, &gono_web.Response{
			Code: gono_web.ResponseCode_OK,
			Msg:  gono_web.ResponseMsg_OK,
		})
	}
	// 拿取Authorization头部
	token := strings.TrimPrefix(c.GetHeader("Authorization"), jwt_util.BearerPrefix)
	if err := user_svc.SignOut(token); err != nil {
		c.JSON(http.StatusInternalServerError, &gono_web.Response{
			Code: gono_web.ResponseCode_InternalError,
			Msg:  fmt.Sprintf("系统报错：%s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, &gono_web.Response{
		Code: gono_web.ResponseCode_OK,
		Msg:  gono_web.ResponseMsg_OK,
	})
}
