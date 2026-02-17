package resource

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sfshf/gonoweb/internal/app/web"
	"github.com/sfshf/gonoweb/internal/service/resource"
)

// ListResource 获取菜单/控件/API列表
// @Summary      获取菜单/控件/API列表
// @Description  获取菜单/控件/API列表
// @Tags         菜单/控件/API
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request query ListResourceReq false "列表搜索条件"
// @Success      200  {object}  ListResourceResp
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /resource [GET]
func ListResource(c *gin.Context) {
	// 检查入参
	var req ListResourceReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	wheres := make(map[string][]any)
	if req.Type > 0 {
		wheres["type=?"] = []any{req.Type}
	}
	if req.Name != "" {
		wheres["name LIKE ?"] = []any{"%" + req.Name + "%"}
	}
	if req.Identifier != "" {
		wheres["identifier LIKE ?"] = []any{"%" + req.Identifier + "%"}
	}
	// 调用服务
	list, total, svcErr := resource.ListResource(req.Page, req.PageSize, wheres)
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
		Data: &ListResourceResp{
			List:  list,
			Total: total,
		},
	})
}

// ResourceInfo 获取菜单/控件/API信息
// @Summary      获取菜单/控件/API信息
// @Description  获取菜单/控件/API信息
// @Tags         菜单/控件/API
// @Accept       plain
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        id path string false "菜单/控件/API的id"
// @Success      200  {object}  github_com_sfshf_gonoweb_internal_model.TResource
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /resource/:id [GET]
func ResourceInfo(c *gin.Context) {
	// 检查入参
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "菜单/控件/API的id为空"),
		})
		return
	}
	idN, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "菜单/控件/API的id非整型数字"),
		})
		return
	}
	// 调用服务
	result, svcErr := resource.ResourceInfo(idN)
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

func validateIdentifier(typ int32, identifier string) error {
	switch typ {
	case 1: // menu
		matched, err := regexp.MatchString(`^/([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`, identifier)
		if err != nil {
			return err
		}
		if !matched {
			return errors.New("menu identifier must be `^/([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`")
		}
	case 2: // widget
		matched, err := regexp.MatchString(`^[A-Za-z0-9]+(_[A-Za-z0-9]+)*$`, identifier)
		if err != nil {
			return err
		}
		if !matched {
			return errors.New("widget identifier must be `^[A-Za-z0-9]+(_[A-Za-z0-9]+)*$`")
		}
	case 3: // api
		matched, err := regexp.MatchString(`^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD) /([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`, identifier)
		if err != nil {
			return err
		}
		if !matched {
			return errors.New("API identifier must be `^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD) /([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`")
		}
	default:
		return errors.New("unsupported identifier type")
	}
	return nil
}

// AddResource 新增菜单/控件/API
// @Summary      新增菜单/控件/API
// @Description  新增菜单/控件/API
// @Tags         菜单/控件/API
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        request body AddResourceReq false "菜单/控件/API信息"
// @Success      200  {object}  github_com_sfshf_gonoweb_internal_model.TResource
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /resource [POST]
func AddResource(c *gin.Context) {
	// 检查入参
	var req AddResourceReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	// 检查identifer的格式
	if err := validateIdentifier(req.Type, req.Identifier); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	// 调用服务
	result, svcErr := resource.AddResource(
		req.Type,
		req.Identifier,
		req.Name,
		req.Intro,
		req.Icon,
	)
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

// EditResource 编辑菜单/控件/API
// @Summary      编辑菜单/控件/API
// @Description  编辑菜单/控件/API
// @Tags         菜单/控件/API
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        id path string false "菜单/控件/API的id"
// @Param        request body EditResourceReq false "菜单/控件/API信息"
// @Success      200  {object}  web.Response
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /resource/:id [PUT]
func EditResource(c *gin.Context) {
	// 检查入参
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "菜单/控件/API的id为空"),
		})
		return
	}
	idN, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "菜单/控件/API的id非整型数字"),
		})
		return
	}
	var req EditResourceReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", err.Error()),
		})
		return
	}
	// 调用服务
	svcErr := resource.EditResource(
		idN,
		req.Name,
		req.Intro,
		req.Icon,
	)
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

// DeleteResource 删除菜单/控件/API
// @Summary      删除菜单/控件/API
// @Description  删除菜单/控件/API
// @Tags         菜单/控件/API
// @Accept       json
// @Produce      json
// @Param        Authorization header string false "登录token"
// @Param        id path string false "菜单/控件/API的id"
// @Success      200  {object}  web.Response
// @Failure      400  {object}  web.Response
// @Failure      404  {object}  web.Response
// @Failure      500  {object}  web.Response
// @Router       /resource/:id [DELETE]
func DeleteResource(c *gin.Context) {
	// 检查入参
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "菜单/控件/API的id为空"),
		})
		return
	}
	idN, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &web.Response{
			Code: web.ResponseCode_RequestError,
			Msg:  fmt.Sprintf("请求参数错误：%s", "菜单/控件/API的id非整型数字"),
		})
		return
	}
	// 调用服务
	if svcErr := resource.DeleteResource(idN); svcErr != nil {
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
