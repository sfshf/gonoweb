package role

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/service/casbin"
	"github.com/sfshf/gonoweb/internal/service/role"
)

// ListRole 获取角色列表
// @Summary      获取角色列表
// @Description  获取角色列表
// @Tags         角色
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request query ListRoleReq false "列表搜索条件"
// @Success      200  {object}  ListRoleResp
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /role [GET]
func ListRole(c *gin.Context) {
	// 检查入参
	var req ListRoleReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	wheres := make(map[string][]any)
	if req.Name != "" {
		wheres["name LIKE ?"] = []any{"%" + req.Name + "%"}
	}
	if req.Dxid != "" {
		// 从casbin中读取该域租户下的角色xid列表
		// TODO 解决rxids数组太长，导致SQL的IN语句失效的问题
		rxids, err := casbin.Enforcer.GetAllRolesByDomain(req.Dxid)
		if err != nil {
			c.JSON(http.StatusBadRequest, &web.Response{
				Code: web.ResponseCode_RequestError,
				Msg:  fmt.Sprintf("请求参数错误：域租户xid错误 %s", err.Error()),
			})
			return
		}
		wheres["xid IN (?)"] = []any{rxids}
	}
	// 调用服务
	list, total, svcErr := role.ListRole(req.Page, req.PageSize, wheres)
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
		Data: &ListRoleResp{
			List:  list,
			Total: total,
		},
	})
}

// RoleInfo 获取角色信息
// @Summary      获取角色信息
// @Description  获取角色信息
// @Tags         角色
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        xid path string false "角色xid"
// @Success      200  {object}  github_com_sfshf_gonoweb_internal_model.TRole
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /role/:xid [GET]
func RoleInfo(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "角色xid为空"),
		})
		return
	}
	// 调用服务
	result, svcErr := role.RoleInfo(xid)
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

// AddRole 新增角色
// @Summary      新增角色
// @Description  新增角色
// @Tags         角色
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request body AddRoleReq false "角色信息"
// @Success      200  {object}  github_com_sfshf_gonoweb_internal_model.TRole
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /role [POST]
func AddRole(c *gin.Context) {
	// 检查入参
	var req AddRoleReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	// 调用服务
	result, svcErr := role.AddRole(req.Name, req.Intro)
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

// EditRole 编辑角色
// @Summary      编辑角色
// @Description  编辑角色
// @Tags         角色
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        xid path string false "角色xid"
// @Param        request body EditRoleReq false "角色信息"
// @Success      200  {object}  web.Response
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /role/:xid [PUT]
func EditRole(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "角色xid为空"),
		})
		return
	}
	var req EditRoleReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	// 调用服务
	svcErr := role.EditRole(xid, req.Name, req.Intro)
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

// DeleteRole 删除角色
// @Summary      删除角色
// @Description  删除角色
// @Tags         角色
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        xid path string false "角色xid"
// @Success      200  {object}  web.Response
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /role/:xid [DELETE]
func DeleteRole(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "角色xid为空"),
		})
		return
	}
	// 调用服务
	if svcErr := role.DeleteRole(xid); svcErr != nil {
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
