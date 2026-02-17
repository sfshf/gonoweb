package user

import (
	"errors"
	"fmt"

	"github.com/rs/xid"
	"github.com/sfshf/gonoweb/internal/config"
	"github.com/sfshf/gonoweb/internal/model"
	. "github.com/sfshf/gonoweb/internal/model"
	"github.com/sfshf/gonoweb/internal/repo"
	"github.com/sfshf/gonoweb/internal/repo/casbin"
	"github.com/sfshf/gonoweb/internal/repo/domain"
	"github.com/sfshf/gonoweb/internal/repo/resource"
	"github.com/sfshf/gonoweb/internal/repo/role"
	"github.com/sfshf/gonoweb/internal/repo/user"
	. "github.com/sfshf/gonoweb/internal/service"
	"github.com/sfshf/gonoweb/internal/util/crypto"
	"github.com/sfshf/gonoweb/internal/util/jwt"
)

func Launch() (func(), error) {
	root := config.AppConfig.Root
	// 1. 检查依赖项有没有加载成功
	if repo.GormDB == nil {
		return nil, errors.New("系统组件错误：初始化用户服务，缺少核心组件")
	}
	// 2. 检查数据库中有没有root账号记录
	user, err := user.User_FirstByEmail(root.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// 没有，则新增
		if err := repo.Create(&model.TUser{
			Xid:      xid.New().String(),
			NickName: "root",
			Email:    root.Email,
			Password: EncryptPlainPassword(root.Password),
		}); err != nil {
			return nil, err
		}
	}
	return func() {
	}, nil
}

func EncryptPlainPassword(plain string) string {
	return crypto.Md5Hex(plain + config.AppConfig.Crypto.PasswordSalt)
}

type SignInData struct {
	Token   string            `json:"token"`
	User    *model.TUser      `json:"user"`
	Domain  *model.TDomain    `json:"domain"`
	Role    *model.TRole      `json:"role"`
	Menus   []model.TResource `json:"menus"`
	Widgets []model.TResource `json:"widgets"`
}

// SignInByPassword 登录成功，则返回用户最近所在的域、角色，以及资源（菜单、控件、API）列表
// 入参 userInfo -- 用户的登录信息；0->ip，1->ua，2->traceID
func SignInByPassword(account, password string, userInfo ...string) (*SignInData, *SvcErr) {
	// 检查账号
	userM, err := user.User_FirstByNickname(account)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if userM == nil {
		userM, err = user.User_FirstByEmail(account)
		if err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	}
	if userM == nil {
		return nil, &SvcErr{Err: fmt.Errorf("用户账号[%s]不存在", account)}
	}
	// 检查密码
	if userM.Password != password {
		return nil, &SvcErr{Err: errors.New("用户密码错误")}
	}
	// 如果该用户是超管账号，则直接返回所有资源
	if userM.Email == config.AppConfig.Root.Email {
		return signIn_Root(userM, userInfo...)
	}
	return signIn_NonRoot(userM, userInfo...)
}

// root 账号登录，加载所有资源
func signIn_Root(user *model.TUser, userInfo ...string) (*SignInData, *SvcErr) {
	// 获取资源
	menuWidgets, err := resource.FindAllMenuWidgets()
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	var menus []model.TResource
	var widgets []model.TResource
	for _, item := range menuWidgets {
		switch item.Type {
		case resource.ResourceType_Menu:
			menus = append(menus, item)
		case resource.ResourceType_Widget:
			widgets = append(widgets, item)
		}
	}
	// 生成登录token
	token, err := jwt.GenerateToken(
		jwt.DefaultSigningMethod,
		config.AppConfig.Gin.Jwt.SigningKey,
		jwt.NewJwtClaims(
			user.Xid,
			"",
			"",
			config.AppConfig.Gin.Jwt.Expired,
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
func signIn_NonRoot(userM *model.TUser, userInfo ...string) (*SignInData, *SvcErr) {
	var err error
	// 从用户最近使用的token里分析出用户最近使用的域租户和角色
	var userAgent *model.TUserAgent
	var ip string
	if len(userInfo) > 0 {
		ip = userInfo[0]
	}
	if ip != "" {
		userAgent, err = user.UserAgent_FirstByUserXidAndIP(userM.Xid, ip)
		if err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
		if userAgent == nil {
			userAgent, err = user.UserAgent_FirstByUserXid(userM.Xid)
			if err != nil {
				return nil, &SvcErr{Internal: true, Err: err}
			}
		}
	} else {
		userAgent, err = user.UserAgent_FirstByUserXid(userM.Xid)
		if err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	}
	var domainXid string
	var roleXid string
	if userAgent == nil {
		// 用户首次登录，获取用户可获得的第一个域租户的第一个角色，以及角色资源
		gRule, err := casbin.FirstGByRsub(userM.Xid)
		if err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
		if gRule == nil {
			// 用户没有配给任何域租户、角色、资源
			token, err := jwt.GenerateToken(
				jwt.DefaultSigningMethod,
				config.AppConfig.Gin.Jwt.SigningKey,
				jwt.NewJwtClaims(
					userM.Xid,
					"",
					"",
					config.AppConfig.Gin.Jwt.Expired,
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
			if err := UpsertUserAgent(ip, ua, tid, userM.Xid, token); err != nil {
				return nil, &SvcErr{Internal: true, Err: err}
			}
			return &SignInData{
				Token: token,
				User:  userM,
			}, nil
		}
		roleXid = gRule.V1
		domainXid = gRule.V2
	} else {
		// 解析用户token
		claims, err := jwt.ParseToken(
			jwt.DefaultSigningMethod,
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
	domain, err := domain.FirstByXid(domainXid)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	role, err := role.FirstByXid(roleXid)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	// 获取用户当前域角色的资源
	menuWidgets, err := resource.FindMenuWidgetsByDomainAndRole(domain.Xid, role.Xid)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	var menus []model.TResource
	var widgets []model.TResource
	for _, item := range menuWidgets {
		switch item.Type {
		case resource.ResourceType_Menu:
			menus = append(menus, item)
		case resource.ResourceType_Widget:
			widgets = append(widgets, item)
		}
	}
	// 生成登录token
	token, err := jwt.GenerateToken(
		jwt.DefaultSigningMethod,
		config.AppConfig.Gin.Jwt.SigningKey,
		jwt.NewJwtClaims(
			userM.Xid,
			domain.Xid,
			role.Xid,
			config.AppConfig.Gin.Jwt.Expired,
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
	if err := UpsertUserAgent(ip, ua, tid, userM.Xid, token); err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	return &SignInData{
		Token:   token,
		User:    userM,
		Domain:  domain,
		Role:    role,
		Menus:   menus,
		Widgets: widgets,
	}, nil
}

func SignOut(token string) error {
	// 找到user_agent记录
	record, err := user.UserAgent_FirstByToken(token)
	if err != nil {
		return err
	}
	if record == nil {
		return nil
	}
	return user.UserAgent_DeleteByToken(token)
}

func ListUser(page, pageSize int, wheres map[string][]any) ([]TUser, int64, *SvcErr) {
	db := repo.GormDB.Table(TableNameTUser)
	for query, args := range wheres {
		db = db.Where(query, args...)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	var list []TUser
	if err := db.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, total, &SvcErr{Internal: true, Err: err}
	}
	return list, total, nil
}

func UserInfo(xid string) (*TUser, *SvcErr) {
	user, err := user.User_FirstByXid(xid)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if user == nil {
		return nil, &SvcErr{Err: fmt.Errorf("用户[xid=%s]不存在", xid)}
	}
	return user, nil
}

func AddUser(email, nickname string) (*TUser, *SvcErr) {
	// 搜索有没有重复的、删除的记录
	userM, err := user.User_FirstUnscopedByEmail(email)
	if err != nil {
		return nil, &SvcErr{Internal: true, Err: err}
	}
	if userM == nil {
		// 新增
		user := &TUser{
			Xid:      xid.New().String(),
			Email:    email,
			Password: EncryptPlainPassword(config.AppConfig.Crypto.DefaultPassword),
			NickName: nickname,
		}
		if err := repo.Create(user); err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	} else {
		// 如果是活用户则报错
		if userM.DeletedAt == 0 {
			return nil, &SvcErr{Err: fmt.Errorf("用户[email=%s]已存在", email)}
		}
		// 如果是死用户则激活
		if err := user.User_ReliveByXid(userM.Xid, &TUser{
			Password: EncryptPlainPassword(config.AppConfig.Crypto.DefaultPassword),
		}); err != nil {
			return nil, &SvcErr{Internal: true, Err: err}
		}
	}
	return userM, nil
}

func EditUser(xid, email, nickName string) *SvcErr {
	if err := user.User_UpdateByXid(xid, &TUser{
		Email:    email,
		NickName: nickName,
	}); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}

func DeleteUser(xid string) *SvcErr {
	if err := user.User_DeleteByXid(xid); err != nil {
		return &SvcErr{Internal: true, Err: err}
	}
	return nil
}
