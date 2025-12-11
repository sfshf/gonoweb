package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	AppConfig appConfig
)

func InitAppConfig(ctx context.Context, path ...string) error {
	// 1. 开发配置为默认路径（优先级最低）
	p := `config/dev/gono.toml`
	if len(path) > 0 {
		p = path[0]
	}
	// 2. 读取环境变量 APP_CONFIG_PATH
	if e := os.Getenv("APP_CONFIG_PATH"); e != "" {
		p = e
	}
	// TODO: 3. 读取命令行参数（优先级最高）
	_, err := toml.DecodeFile(p, &AppConfig)
	if err != nil {
		return err
	}
	if AppConfig.Mode == AppMode_Dev {
		b, err := json.MarshalIndent(AppConfig, "", "\t")
		if err != nil {
			return err
		}
		log.Printf("%s\n", b)
	}
	return nil
}

type appConfig struct {
	Name string     `toml:"name"`
	Mode appMode    `toml:"mode"`
	Gin  ginConfig  `toml:"gin"`
	Gorm gormConfig `toml:"gorm"`
}

type appMode string

const (
	AppMode_Dev  appMode = "dev"
	AppMode_Test appMode = "test"
	AppMode_Prod appMode = "prod"
)

type ginConfig struct {
	Jwt struct {
		SigningKey string        `toml:"signingKey"`
		Expired    time.Duration `toml:"expired"` // seconds
	} `toml:"jwt"`
}

type gormConfig struct {
	DryRun bool `toml:"dryRun"`
	Mysql  struct {
		DSN             string `toml:"dsn"`
		ConnMaxLifetime int    `toml:"connMaxLifetime"`
		MaxOpenConns    int    `toml:"maxOpenConns"`
		MaxIdleConns    int    `toml:"maxIdleConns"`
	} `toml:"mysql"`
}
