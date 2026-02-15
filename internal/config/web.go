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
	// 4. 设置适配性参数
	if len(path) > 1 {
		AppConfig.Gin.Casbin.Model = path[1] // casbin model file path
	}

	return nil
}

type appConfig struct {
	Name   string       `toml:"name"`
	Mode   appMode      `toml:"mode"`
	Port   int          `toml:"port"`
	Root   rootConfig   `toml:"root"`
	Crypto cryptoConfig `toml:"crypto"`
	Gin    ginConfig    `toml:"gin"`
	Gorm   gormConfig   `toml:"gorm"`
}

type appMode string

const (
	AppMode_Dev  appMode = "dev"
	AppMode_Test appMode = "test"
	AppMode_Prod appMode = "prod"
)

type rootConfig struct {
	Email    string `toml:"email"`
	Password string `toml:"password"`
}

type cryptoConfig struct {
	PasswordSalt    string `toml:"passwordSalt"`
	DefaultPassword string `toml:"defaultPassword"`
}

type ginConfig struct {
	Swagdoc bool `toml:"swagdoc"`
	Jwt     struct {
		SigningKey string        `toml:"signingKey"`
		Expired    time.Duration `toml:"expired"` // seconds
	} `toml:"jwt"`
	Casbin struct {
		Log              bool          `toml:"log"`
		Model            string        `toml:"model"`
		Enforce          bool          `toml:"enforce"`
		AutoSave         bool          `toml:"autoSave"`
		AutoLoad         bool          `toml:"autoLoad"`
		AutoLoadInterval time.Duration `toml:"autoLoadInterval"`
	} `toml:"casbin"`
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
