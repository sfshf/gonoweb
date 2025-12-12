package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/sfshf/gonoweb/internal/app/web/ginsrv"
	"github.com/sfshf/gonoweb/internal/config"
	"github.com/sfshf/gonoweb/internal/repo"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// 加载应用配置
	if err := config.InitAppConfig(ctx); err != nil {
		log.Fatalln(err)
	}
	// 初始化数据库
	if err := repo.InitGorm(ctx); err != nil {
		log.Fatalln(err)
	}
	// 初始化服务路由
	router, err := ginsrv.InitGin(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// 开启API服务
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.AppConfig.Port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Printf("Gonoweb Server is listening at :%d\n", config.AppConfig.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Gonoweb Server listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown: ", err)
	}

	log.Println("Gonoweb Server exiting")
}
