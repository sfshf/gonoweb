package repo

import (
	"context"
	"errors"

	"github.com/sfshf/gonoweb/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	GormDB *gorm.DB
)

func InitGorm(ctx context.Context) error {
	dsn := config.AppConfig.Gorm.Mysql.DSN
	if dsn == "" {
		return errors.New("mysql dsn in config is empty")
	}
	gcfg := &gorm.Config{
		DryRun: config.AppConfig.Gorm.DryRun,
	}
	var err error
	GormDB, err = gorm.Open(mysql.Open(dsn), gcfg)
	if err != nil {
		return err
	}
	return nil
}

func Create[M any](m M) error {
	return GormDB.Create(m).Error
}
