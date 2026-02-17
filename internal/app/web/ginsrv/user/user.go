package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/app/web/ginsrv/middlewares"
	"github.com/sfshf/gonoweb/internal/service/user"
	"github.com/sfshf/gonoweb/internal/util/jwt"
)

// SignIn 用户登录
// @Summary      用户登录
// @Description  用户登录
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param 		 request body SignInReq true "登录所需参数"
// @Success      200  {object}  user.SignInData
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /user/signIn [POST]
func SignIn(c *gin.Context) {
	// 检查请求入参
	var req SignInReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	// IP  from HTTP headers
	ip := c.ClientIP()
	// User-Agent from HTTP headers
	ua := c.GetHeader("User-Agent")
	// TracdID from HTTP headers
	tid := c.GetHeader(middlewares.HeaderKey_TraceID)
	if tid == "" {
		tid = c.GetString(middlewares.HeaderKey_TraceID)
	}
	// 密码登录
	data, svcErr := user.SignInByPassword(req.Account, req.Password, ip, ua, tid)
	if svcErr != nil {
		if svcErr.Internal {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  fmt.Sprintf("系统报错：%s", svcErr.Error()),
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, &web.Response{
				Code: web.ResponseCode_RequestError,
				Msg:  fmt.Sprintf("密码登录失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 将jwt写入头部
	c.Header("Authorization", jwt.BearerPrefix+data.Token)
	c.JSON(http.StatusOK, &web.Response{
		Code: web.ResponseCode_OK,
		Msg:  web.ResponseMsg_OK,
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
// @Success      200  {object}  web.Response
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /user/signOut [POST]
func SignOut(c *gin.Context) {
	// 从gin.Context拿取用户信息
	claims := middlewares.JwtClaims(c)
	if claims == nil {
		c.JSON(http.StatusOK, &web.Response{
			Code: web.ResponseCode_OK,
			Msg:  web.ResponseMsg_OK,
		})
	}
	// 拿取Authorization头部
	token := strings.TrimPrefix(c.GetHeader("Authorization"), jwt.BearerPrefix)
	if err := user.SignOut(token); err != nil {
		c.JSON(http.StatusInternalServerError, &web.Response{
			Code: web.ResponseCode_InternalError,
			Msg:  fmt.Sprintf("系统报错：%s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, &web.Response{
		Code: web.ResponseCode_OK,
		Msg:  web.ResponseMsg_OK,
	})
}

// ListUser 获取用户列表
// @Summary      获取用户列表
// @Description  获取用户列表
// @Tags         用户
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request query ListUserReq false "列表搜索条件"
// @Success      200  {object}  ListUserResp
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /user [GET]
func ListUser(c *gin.Context) {
	// 检查入参
	var req ListUserReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	wheres := make(map[string][]any)
	if req.Email != "" {
		wheres["email=?"] = []any{req.Email}
	}
	if req.Nickname != "" {
		wheres["nick_name LIKE ?"] = []any{"%" + req.Nickname + "%"}
	}
	if req.Realname != "" {
		wheres["real_name=?"] = []any{req.Realname}
	}
	// 调用服务
	list, total, svcErr := user.ListUser(req.Page, req.PageSize, wheres)
	if svcErr != nil {
		if svcErr.Internal {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  fmt.Sprintf("系统报错：%s", svcErr.Error()),
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, &web.Response{
				Code: web.ResponseCode_RequestError,
				Msg:  fmt.Sprintf("获取列表失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 返回结果
	c.JSON(http.StatusOK, &web.Response{
		Code: web.ResponseCode_OK,
		Msg:  web.ResponseMsg_OK,
		Data: &ListUserResp{
			List:  list,
			Total: total,
		},
	})
}

// UserInfo 获取用户信息
// @Summary      获取用户信息
// @Description  获取用户信息
// @Tags         用户
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        xid path string false "用户xid"
// @Success      200  {object}  github_com_sfshf_gonoweb_internal_model.TUser
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /user/:xid [GET]
func UserInfo(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "用户xid为空"),
		})
		return
	}
	// 调用服务
	result, svcErr := user.UserInfo(xid)
	if svcErr != nil {
		if svcErr.Internal {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  fmt.Sprintf("系统报错：%s", svcErr.Error()),
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, &web.Response{
				Code: web.ResponseCode_RequestError,
				Msg:  fmt.Sprintf("获取信息失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 返回结果
	c.JSON(http.StatusOK, &web.Response{
		Code: web.ResponseCode_OK,
		Msg:  web.ResponseMsg_OK,
		Data: result,
	})
}

// AddUser 新增用户
// @Summary      新增用户
// @Description  新增用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request body AddUserReq false "用户信息"
// @Success      200  {object}  github_com_sfshf_gonoweb_internal_model.TUser
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /user [POST]
func AddUser(c *gin.Context) {
	// 检查入参
	var req AddUserReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	// 调用服务
	result, svcErr := user.AddUser(req.Email, req.Nickname)
	if svcErr != nil {
		if svcErr.Internal {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  fmt.Sprintf("系统报错：%s", svcErr.Error()),
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, &web.Response{
				Code: web.ResponseCode_RequestError,
				Msg:  fmt.Sprintf("新增失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 返回结果
	c.JSON(http.StatusOK, &web.Response{
		Code: web.ResponseCode_OK,
		Msg:  web.ResponseMsg_OK,
		Data: result,
	})
}

// EditUser 编辑用户
// @Summary      编辑用户
// @Description  编辑用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        xid path string false "用户xid"
// @Param        request body EditUserReq false "用户信息"
// @Success      200  {object}  web.Response
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /user/:xid [PUT]
func EditUser(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "用户xid为空"),
		})
		return
	}
	var req EditUserReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	// 调用服务
	svcErr := user.EditUser(xid, req.Email, req.NickName)
	if svcErr != nil {
		if svcErr.Internal {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  fmt.Sprintf("系统报错：%s", svcErr.Error()),
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, &web.Response{
				Code: web.ResponseCode_RequestError,
				Msg:  fmt.Sprintf("更新失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 返回结果
	c.JSON(http.StatusOK, &web.Response{
		Code: web.ResponseCode_OK,
		Msg:  web.ResponseMsg_OK,
	})
}

// DeleteUser 删除用户
// @Summary      删除用户
// @Description  删除用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        xid path string false "用户xid"
// @Success      200  {object}  web.Response
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /user/:xid [DELETE]
func DeleteUser(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "用户xid为空"),
		})
		return
	}
	// 调用服务
	if svcErr := user.DeleteUser(xid); svcErr != nil {
		if svcErr.Internal {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  fmt.Sprintf("系统报错：%s", svcErr.Error()),
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, &web.Response{
				Code: web.ResponseCode_RequestError,
				Msg:  fmt.Sprintf("删除失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 返回结果
	c.JSON(http.StatusOK, &web.Response{
		Code: web.ResponseCode_OK,
		Msg:  web.ResponseMsg_OK,
	})
}
