package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	gono_web "github.com/sfshf/gonoweb/internal/app/web"
	ginmw "github.com/sfshf/gonoweb/internal/app/web/ginsrv/middlewares"
	"github.com/sfshf/gonoweb/internal/model"
	user_svc "github.com/sfshf/gonoweb/internal/service/user"
	jwt_util "github.com/sfshf/gonoweb/internal/util/jwt"
)

type SignInReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignIn 用户登录
// @Summary      用户登录
// @Description  用户登录
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param 		 request body SignInReq true "登录所需参数"
// @Success      200  {object}  user_svc.SignInData
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
	data, svcErr := user_svc.SignInByPassword(req.Account, req.Password, ip, ua, tid)
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

// ListUser 获取用户列表
// @Summary      获取用户列表
// @Description  获取用户列表
// @Tags         用户
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request query ListUserReq false "列表搜索条件"
// @Success      200  {object}  ListUserResp
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /user [GET]
func ListUser(c *gin.Context) {
	// 检查入参
	var req ListUserReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
	}
	wheres := make(map[string][]interface{})
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
	list, total, svcErr := user_svc.ListUser(req.Page, req.PageSize, wheres)
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
				Msg:  fmt.Sprintf("获取列表失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 返回结果
	c.JSON(http.StatusOK, &gono_web.Response{
		Code: gono_web.ResponseCode_OK,
		Msg:  gono_web.ResponseMsg_OK,
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
// @Success      200  {object}  model.TUser
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /user/:xid [GET]
func UserInfo(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "用户xid为空"),
		})
	}
	// 调用服务
	result, svcErr := user_svc.UserInfo(xid)
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
				Msg:  fmt.Sprintf("获取信息失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 返回结果
	c.JSON(http.StatusOK, &gono_web.Response{
		Code: gono_web.ResponseCode_OK,
		Msg:  gono_web.ResponseMsg_OK,
		Data: result,
	})
}

type AddUserReq struct {
	Email    string `json:"email" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

// AddUser 新增用户
// @Summary      新增用户
// @Description  新增用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request body AddUserReq false "用户信息"
// @Success      200  {object}  model.TUser
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /user [POST]
func AddUser(c *gin.Context) {
	// 检查入参
	var req AddUserReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
	}
	// 调用服务
	result, svcErr := user_svc.AddUser(req.Email, req.Nickname)
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
				Msg:  fmt.Sprintf("新增失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 返回结果
	c.JSON(http.StatusOK, &gono_web.Response{
		Code: gono_web.ResponseCode_OK,
		Msg:  gono_web.ResponseMsg_OK,
		Data: result,
	})
}

type EditUserReq struct {
	Email    string `json:"email" binding:"required"`
	NickName string `json:"nickName" binding:"required"`
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
// @Success      200  {object}  gono_web.Response
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /user/:xid [PUT]
func EditUser(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "用户xid为空"),
		})
	}
	var req EditUserReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
	}
	// 调用服务
	svcErr := user_svc.EditUser(xid, req.Email, req.NickName)
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
				Msg:  fmt.Sprintf("更新失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 返回结果
	c.JSON(http.StatusOK, &gono_web.Response{
		Code: gono_web.ResponseCode_OK,
		Msg:  gono_web.ResponseMsg_OK,
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
// @Success      200  {object}  gono_web.Response
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /user/:xid [DELETE]
func DeleteUser(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "用户xid为空"),
		})
	}
	// 调用服务
	if svcErr := user_svc.DeleteUser(xid); svcErr != nil {
		if svcErr.Internal {
			c.JSON(http.StatusInternalServerError, &gono_web.Response{
				Code: gono_web.ResponseCode_InternalError,
				Msg:  fmt.Sprintf("系统报错：%s", svcErr.Error()),
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, &gono_web.Response{
				Code: gono_web.ResponseCode_RequestError,
				Msg:  fmt.Sprintf("删除失败：%s", svcErr.Error()),
			})
			return
		}
	}
	// 返回结果
	c.JSON(http.StatusOK, &gono_web.Response{
		Code: gono_web.ResponseCode_OK,
		Msg:  gono_web.ResponseMsg_OK,
	})
}
