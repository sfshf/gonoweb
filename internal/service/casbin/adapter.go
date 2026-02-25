package casbin

import (
	"errors"
	"time"

	"github.com/casbin/casbin/v3"
	casbinModel "github.com/casbin/casbin/v3/model"
	"github.com/sfshf/gonoweb/internal/config"
	"github.com/sfshf/gonoweb/internal/repo"
	casbinRepo "github.com/sfshf/gonoweb/internal/repo/casbin"
	"gorm.io/gorm"
)

var (
	Enforcer *casbin.SyncedEnforcer
)

func Launch() (func(), error) {
	opt := config.AppConfig.Gin.Casbin
	// 1. 检查依赖项有没有加载成功
	if repo.GormDB == nil {
		return nil, errors.New("系统组件错误：初始化Casbin服务，缺少核心组件")
	}
	// 2. 加载casbin模型
	m, err := casbinModel.NewModelFromFile(opt.Model)
	if err != nil {
		return nil, err
	}
	// 3. 创建Enforcer实例
	Enforcer, err = casbin.NewSyncedEnforcer(m, casbinRepo.Adapter(func() *gorm.DB { return repo.GormDB }))
	if err != nil {
		return nil, err
	}
	// 4. 装配Enforcer选项
	Enforcer.EnableLog(opt.Log)
	if opt.AutoLoad {
		Enforcer.StartAutoLoadPolicy(time.Second * time.Duration(opt.AutoLoadInterval))
	}
	Enforcer.EnableAutoSave(opt.AutoSave)
	Enforcer.EnableEnforce(opt.Enforce)
	if err = Enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	return func() {
		Enforcer.SavePolicy()
	}, nil
}
