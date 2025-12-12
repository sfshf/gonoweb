package ginsrv

import (
	"context"
	"errors"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	ginmw "github.com/sfshf/gonoweb/internal/app/web/ginsrv/middlewares"
	"github.com/sfshf/gonoweb/internal/config"
)

func ginMode(appMode string) string {
	var mode string
	switch appMode {
	case "dev":
		mode = "debug"
	case "test":
		mode = "test"
	case "prod":
		mode = "release"
	}
	return mode
}

func InitGin(ctx context.Context) (*gin.Engine, error) {
	// gin mode
	gin.SetMode(ginMode(string(config.AppConfig.Mode)))
	// gin engine instance
	var r *gin.Engine
	switch gin.Mode() {
	case gin.DebugMode:
		r = gin.Default(func(e *gin.Engine) {
			e.HandleMethodNotAllowed = true
		})
	case gin.TestMode:
		r = gin.Default(func(e *gin.Engine) {
			e.HandleMethodNotAllowed = true
		})
	case gin.ReleaseMode:
		r = gin.New(func(e *gin.Engine) {
			e.HandleMethodNotAllowed = true
		})
		r.Use(gin.Logger())   // TODO：自定义日志器
		r.Use(gin.Recovery()) // TODO：自定义宕机重启
	}
	if r == nil {
		return nil, errors.New("fail to init gin engine")
	}
	// set gin's default validator to global
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		Validator = v
	}
	// 装配全局型中间件 -- CORS、TraceID、GZIP、NoMethod、NoRoute等
	r.Use(
		cors.Default(),
		ginmw.TraceID(),
		gzip.Gzip(gzip.DefaultCompression),
	)
	r.NoMethod(ginmw.NoMethod())
	r.NoRoute(ginmw.NoRoute())
	// v1 api router
	v1 := r.Group("/api/v1")
	LoadRoutes_V1(v1)
	return r, nil
}
