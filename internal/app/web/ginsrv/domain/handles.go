package domain

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	gono_web "github.com/sfshf/gonoweb/internal/app/web"
	domain_svc "github.com/sfshf/gonoweb/internal/service/domain"
)

// ListDomain 获取域租户列表
// @Summary      获取域租户列表
// @Description  获取域租户列表
// @Tags         域租户
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request query ListDomainReq false "列表搜索条件"
// @Success      200  {object}  ListDomainResp
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /domain [GET]
func ListDomain(c *gin.Context) {
	// 检查入参
	var req ListDomainReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
	}
	wheres := make(map[string][]any)
	if req.Name != "" {
		wheres["name=?"] = []any{req.Name}
	}
	// 调用服务
	list, total, svcErr := domain_svc.ListDomain(req.Page, req.PageSize, wheres)
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
		Data: &ListDomainResp{
			List:  list,
			Total: total,
		},
	})
}

// DomainInfo 获取域租户信息
// @Summary      获取域租户信息
// @Description  获取域租户信息
// @Tags         域租户
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        xid path string false "域租户xid"
// @Success      200  {object}  model.TDomain
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /domain/:xid [GET]
func DomainInfo(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "域租户xid为空"),
		})
	}
	// 调用服务
	result, svcErr := domain_svc.DomainInfo(xid)
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

// AddDomain 新增域租户
// @Summary      新增域租户
// @Description  新增域租户
// @Tags         域租户
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request body AddDomainReq false "域租户信息"
// @Success      200  {object}  model.TDomain
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /domain [POST]
func AddDomain(c *gin.Context) {
	// 检查入参
	var req AddDomainReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
	}
	// 调用服务
	result, svcErr := domain_svc.AddDomain(req.Name, req.Intro)
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

// EditDomain 编辑域租户
// @Summary      编辑域租户
// @Description  编辑域租户
// @Tags         域租户
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        xid path string false "域租户xid"
// @Param        request body EditDomainReq false "域租户信息"
// @Success      200  {object}  gono_web.Response
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /domain/:xid [PUT]
func EditDomain(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "域租户xid为空"),
		})
	}
	var req EditDomainReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
	}
	// 调用服务
	svcErr := domain_svc.EditDomain(xid, req.Name, req.Intro)
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

// DeleteDomain 删除域租户
// @Summary      删除域租户
// @Description  删除域租户
// @Tags         域租户
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        xid path string false "域租户xid"
// @Success      200  {object}  gono_web.Response
// @Failure      400  {object}  gono_web.Response
// @Failure      404  {object}  gono_web.Response
// @Failure      500  {object}  gono_web.Response
// @Router       /domain/:xid [DELETE]
func DeleteDomain(c *gin.Context) {
	// 检查入参
	xid := c.Param("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, &gono_web.Response{
			Code: gono_web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "域租户xid为空"),
		})
	}
	// 调用服务
	if svcErr := domain_svc.DeleteDomain(xid); svcErr != nil {
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
