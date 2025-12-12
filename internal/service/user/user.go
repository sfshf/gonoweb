package user_svc

import (
	"errors"
	"fmt"
	"time"

	"github.com/sfshf/gonoweb/internal/config"
	. "github.com/sfshf/gonoweb/internal/model"
	casbin_repo "github.com/sfshf/gonoweb/internal/repo/casbin"
	domain_repo "github.com/sfshf/gonoweb/internal/repo/domain"
	mwa_repo "github.com/sfshf/gonoweb/internal/repo/menu_widget_api"
	role_repo "github.com/sfshf/gonoweb/internal/repo/role"
	user_repo "github.com/sfshf/gonoweb/internal/repo/user"
	. "github.com/sfshf/gonoweb/internal/service"
	jwt_util "github.com/sfshf/gonoweb/internal/util/jwt"
)

type SignInData struct {
	Token   string           `json:"token"`
	User    *TUser           `json:"user"`
	Domain  *TDomain         `json:"domain"`
	Role    *TRole           `json:"role"`
	Menus   []TMenuWidgetAPI `json:"menus"`
	Widgets []TMenuWidgetAPI `json:"widgets"`
}

// SignInByPassword 登录成功，则返回用户最近所在的域、角色，以及资源（菜单、控件、API）列表
// 入参 userInfo -- 用户的登录信息；0->ip，1->ua，2->traceID
func SignInByPassword(email, password string, userInfo ...string) (*SignInData, *SvcErr) {
	// 检查账号
	user, err := user_repo.User_FirstByEmail(email)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if user == nil {
		return nil, &SvcErr{Err: fmt.Errorf("用户账号[%s]不存在", email)}
	}
	// 检查密码
	if user.Password != password {
		return nil, &SvcErr{Err: errors.New("用户密码错误")}
	}
	// 如果该用户是超管账号，则直接返回所有资源
	if user.Email == config.AppConfig.Root.Account {
		return signIn_Root(user, userInfo...)
	}
	return signIn_NonRoot(user, userInfo...)
}

// root 账号登录，加载所有资源
func signIn_Root(user *TUser, userInfo ...string) (*SignInData, *SvcErr) {
	// 获取资源
	menuWidgets, err := mwa_repo.FindAllMenuWidgets()
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	var menus []TMenuWidgetAPI
	var widgets []TMenuWidgetAPI
	for _, item := range menuWidgets {
		switch item.Type {
		case mwa_repo.MenuWidgetApiType_Menu:
			menus = append(menus, item)
		case mwa_repo.MenuWidgetApiType_Widget:
			widgets = append(widgets, item)
		}
	}
	// 生成登录token
	token, err := jwt_util.GenerateToken(
		jwt_util.DefaultSigningMethod,
		config.AppConfig.Gin.Jwt.SigningKey,
		jwt_util.NewJwtClaims(
			user.Xid,
			"",
			"",
			config.AppConfig.Gin.Jwt.Expired*time.Second,
		),
	)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	// 将新token更新到t_user_agent表
	var ip string
	if len(userInfo) > 0 {
		ip = userInfo[0]
	}
	var ua string
	if len(userInfo) > 1 {
		ua = userInfo[1]
	}
	var tid string
	if len(userInfo) > 2 {
		tid = userInfo[2]
	}
	if err := UpsertUserAgent(ip, ua, tid, user.Xid, token); err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	return &SignInData{
		Token:   token,
		User:    user,
		Menus:   menus,
		Widgets: widgets,
	}, nil
}

// 非 root 账号登录，检查账号最近状态信息，加载相关资源
func signIn_NonRoot(user *TUser, userInfo ...string) (*SignInData, *SvcErr) {
	var err error
	// 从用户最近使用的token里分析出用户最近使用的域租户和角色
	var userAgent *TUserAgent
	var ip string
	if len(userInfo) > 0 {
		ip = userInfo[0]
	}
	if ip != "" {
		userAgent, err = user_repo.UserAgent_FirstByUserXidAndIP(user.Xid, ip)
		if err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
		if userAgent == nil {
			userAgent, err = user_repo.UserAgent_FirstByUserXid(user.Xid)
			if err != nil {
				return nil, &SvcErr{Internal: true, Err: err}
			}
		}
	} else {
		userAgent, err = user_repo.UserAgent_FirstByUserXid(user.Xid)
		if err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	}
	var domainXid string
	var roleXid string
	if userAgent == nil {
		// 用户首次登录，获取用户可获得的第一个域租户的第一个角色，以及角色资源
		gRule, err := casbin_repo.FirstGByRsub(user.Xid)
		if err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
		if gRule == nil {
			// 用户没有配给任何域租户、角色、资源
			token, err := jwt_util.GenerateToken(
				jwt_util.DefaultSigningMethod,
				config.AppConfig.Gin.Jwt.SigningKey,
				jwt_util.NewJwtClaims(
					user.Xid,
					"",
					"",
					config.AppConfig.Gin.Jwt.Expired*time.Second,
				),
			)
			if err != nil {
				return nil, &SvcErr{Internal: true, Err: err}
			}
			// 将新token更新到t_user_agent表
			var ua string
			if len(userInfo) > 1 {
				ua = userInfo[1]
			}
			var tid string
			if len(userInfo) > 2 {
				tid = userInfo[2]
			}
			if err := UpsertUserAgent(ip, ua, tid, user.Xid, token); err != nil {
				return nil, &SvcErr{Internal: true, Err: err}
			}
			return &SignInData{
				Token: token,
				User:  user,
			}, nil
		}
		roleXid = gRule.V1
		domainXid = gRule.V2
	} else {
		// 解析用户token
		claims, err := jwt_util.ParseToken(
			jwt_util.DefaultSigningMethod,
			config.AppConfig.Gin.Jwt.SigningKey,
			userAgent.Token,
		)
		if err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
		domainXid = claims.Domain
		roleXid = claims.Role
	}
	// 获取用户域租户信息、角色信息、资源信息
	domain, err := domain_repo.FirstByXid(domainXid)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	role, err := role_repo.FirstByXid(roleXid)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	// 获取用户当前域角色的资源
	menuWidgets, err := mwa_repo.FindMenuWidgetsByDomainAndRole(domain.Xid, role.Xid)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	var menus []TMenuWidgetAPI
	var widgets []TMenuWidgetAPI
	for _, item := range menuWidgets {
		switch item.Type {
		case mwa_repo.MenuWidgetApiType_Menu:
			menus = append(menus, item)
		case mwa_repo.MenuWidgetApiType_Widget:
			widgets = append(widgets, item)
		}
	}
	// 生成登录token
	token, err := jwt_util.GenerateToken(
		jwt_util.DefaultSigningMethod,
		config.AppConfig.Gin.Jwt.SigningKey,
		jwt_util.NewJwtClaims(
			user.Xid,
			domain.Xid,
			role.Xid,
			config.AppConfig.Gin.Jwt.Expired*time.Second,
		),
	)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	// 将新token更新到t_user_agent表
	var ua string
	if len(userInfo) > 1 {
		ua = userInfo[1]
	}
	var tid string
	if len(userInfo) > 2 {
		tid = userInfo[2]
	}
	if err := UpsertUserAgent(ip, ua, tid, user.Xid, token); err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	return &SignInData{
		Token:   token,
		User:    user,
		Domain:  domain,
		Role:    role,
		Menus:   menus,
		Widgets: widgets,
	}, nil
}
