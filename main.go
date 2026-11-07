package main

import (
	"log"
	"os"
	"time"

	"week_06_web/auth"
	"week_06_web/repository"
	"week_06_web/router"
	"week_06_web/service"
)

const (
	tokenTTL  = 2 * time.Hour
	devSecret = "week_06_web_dev_secret_change_me_1234567890"
)

func main() {
	UserRepo, err := repository.NewFileUserRepository("data/user.json")
	if err != nil {
		log.Fatalf("create repository.NewFileUserRepository fail: %v", err)
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = devSecret
	}

	TokenManager := auth.NewTokenManager(secret, tokenTTL)

	UserService := service.NewUserService(UserRepo)
	engine := router.New(UserService, TokenManager)

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
