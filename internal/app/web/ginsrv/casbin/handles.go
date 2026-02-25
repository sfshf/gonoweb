package casbin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/config"
	"github.com/sfshf/gonoweb/internal/repo/resource"
	"github.com/sfshf/gonoweb/internal/repo/user"
	"github.com/sfshf/gonoweb/internal/service/casbin"
	"github.com/sfshf/gonoweb/internal/util/strs"
)

// DomainRoleResources 获取域租户的某个角色的资源列表
// @Summary      获取域租户的某个角色的资源列表
// @Description  获取域租户的某个角色的资源列表
// @Tags         Casbin
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        dxid path string false "域租户的xid"
// @Param        rxid path string false "角色的xid"
// @Success      200  {object}  []string
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /casbin/domain/:dxid/role/:rxid/resource [GET]
func DomainRoleResources(c *gin.Context) {
	// 检查入参
	dxid := c.Param("dxid")
	if dxid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "域租户xid为空"),
		})
		return
	}
	rxid := c.Param("rxid")
	if rxid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "角色xid为空"),
		})
		return
	}
	policies, err := casbin.Enforcer.GetFilteredPolicy(0, rxid, dxid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &web.Response{
			Code: web.ResponseCode_InternalError,
			Msg:  err.Error(),
		})
		return
	}
	var list []string
	for _, p := range policies {
		list = append(list, strings.TrimSpace(
			strings.Join([]string{
				p[3], // act
				p[2], // obj
			}, " "),
		))
	}
	// 返回结果
	c.JSON(http.StatusOK, &web.Response{
		Code: web.ResponseCode_OK,
		Msg:  web.ResponseMsg_OK,
		Data: list,
	})
}

// AllocDomainRoleResources 给域租户的角色分配资源
// @Summary      给域租户的角色分配资源
// @Description  给域租户的角色分配资源
// @Tags         Casbin
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        dxid path string false "域租户的xid"
// @Param        rxid path string false "角色的xid"
// @Param        request body AllocDomainRoleResourcesReq false "分配所需参数"
// @Success      200  {object}  nil
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /casbin/domain/:dxid/role/:rxid/resource [POST]
func AllocDomainRoleResources(c *gin.Context) {
	// 检查入参
	dxid := c.Param("dxid")
	if dxid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "域租户xid为空"),
		})
		return
	}
	rxid := c.Param("rxid")
	if rxid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "角色xid为空"),
		})
		return
	}
	var req AllocDomainRoleResourcesReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	if _, err := casbin.Enforcer.RemoveFilteredPolicy(0, rxid, dxid); err != nil {
		c.JSON(http.StatusInternalServerError, &web.Response{
			Code: web.ResponseCode_InternalError,
			Msg:  err.Error(),
		})
		return
	}
	if len(req.Identifiers) > 0 {
		var rules [][]string
		for _, id := range req.Identifiers {
			// 资源identifier的格式检查
			val, err := strs.ValidateResourceIdentifier(-1, id)
			if err != nil {
				c.JSON(http.StatusBadRequest, &web.Response{
					Code: web.ResponseCode_RequestError,
					Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
				})
				return
			}
			// 检查应用服务是否有该资源
			var identifier string
			switch val.Type {
			case 1, 2: // 菜单/控件
				identifier = val.Obj
			case 3: // API
				identifier = strings.Join([]string{val.Act, val.Obj}, " ")
			}
			result, err := resource.FirstByIdentifier(identifier)
			if err != nil {
				c.JSON(http.StatusInternalServerError, &web.Response{
					Code: web.ResponseCode_InternalError,
					Msg:  err.Error(),
				})
				return
			}
			if result == nil {
				c.JSON(http.StatusBadRequest, &web.Response{
					Code: web.ResponseCode_RequestError,
					Msg:  fmt.Sprintf("资源不存在：%s", identifier),
				})
				return
			}
			rules = append(rules, []string{rxid, dxid, val.Obj, val.Act})
		}
		if _, err := casbin.Enforcer.AddPoliciesEx(rules); err != nil {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  err.Error(),
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

// DomainRoles 获取域租户的角色列表
// @Summary      获取域租户的角色列表
// @Description  获取域租户的角色列表
// @Tags         Casbin
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        dxid path string false "域租户的xid"
// @Success      200  {object}  []string
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /casbin/domain/:dxid/role [GET]
func DomainRoles(c *gin.Context) {
	// 检查入参
	dxid := c.Param("dxid")
	if dxid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "域租户xid为空"),
		})
		return
	}
	rxids, err := casbin.Enforcer.GetAllRolesByDomain(dxid)
	if err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：域租户xid错误 %s", err.Error()),
		})
		return
	}
	// 返回结果
	c.JSON(http.StatusOK, &web.Response{
		Code: web.ResponseCode_OK,
		Msg:  web.ResponseMsg_OK,
		Data: rxids,
	})
}

// AllocRoleInDomain 给用户分配域租户内的角色
// @Summary      给用户分配域租户内的角色
// @Description  给用户分配域租户内的角色
// @Tags         Casbin
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request body AllocRolesInDomainReq false "分配所需参数"
// @Success      200  {object}  nil
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /casbin/user/:xid [POST]
func AllocRoleInDomain(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "用户xid为空"),
		})
		return
	}
	// !IMPORTANT! 不可修改root账户权限
	user, err := user.User_FirstByXid(xid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &web.Response{
			Code: web.ResponseCode_InternalError,
			Msg:  err.Error(),
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusForbidden, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  "无有效的用户信息",
		})
		return
	}
	if user.Email == config.AppConfig.Root.Email {
		c.JSON(http.StatusForbidden, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  "权限不够",
		})
		return
	}
	// TODO 当前登录用户是否有权限修改目标用户
	var req AllocRolesInDomainReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	// 更新casbin
	if _, err := casbin.Enforcer.DeleteRolesForUserInDomain(xid, req.Domain); err != nil {
		c.JSON(http.StatusInternalServerError, &web.Response{
			Code: web.ResponseCode_InternalError,
			Msg:  err.Error(),
		})
		return
	}
	if len(req.Roles) > 0 {
		if _, err := casbin.Enforcer.AddRolesForUser(xid, req.Roles, req.Domain); err != nil {
			c.JSON(http.StatusInternalServerError, &web.Response{
				Code: web.ResponseCode_InternalError,
				Msg:  err.Error(),
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
