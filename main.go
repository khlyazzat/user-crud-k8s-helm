package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/khlyazzat/user-crud-k8s-helm/utils"

	"github.com/gin-gonic/gin"

	"github.com/khlyazzat/user-crud-k8s-helm/internal/config"
	"github.com/khlyazzat/user-crud-k8s-helm/internal/db/postgres"
	apiRouter "github.com/khlyazzat/user-crud-k8s-helm/internal/router"
	"github.com/khlyazzat/user-crud-k8s-helm/internal/user"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("env not loaded")
	}

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	db := postgres.New(cfg.DBConfig)

	router := gin.New()

	v1 := router.Group("/v1")

	healthCClient := apiRouter.NewHealthClient()
	healthCClient.RegisterRouter(v1)

	userService := user.NewUserService(db)
	authClient := apiRouter.NewUserClient(userService)
	authClient.RegisterRouter(v1)

	srv := &http.Server{
		Addr:    cfg.HTTPConfig.Port,
		Handler: router,
	}

	go func() {
		log.Println("Starting server on", cfg.HTTPConfig.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Web server error: %s", err)
		}
	}()

	ctx, cancel := utils.GracefulShutdown(context.TODO())
	defer cancel()

	<-ctx.Done()
}
