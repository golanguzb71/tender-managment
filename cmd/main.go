package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	_ "tender-managment/docs"
	"tender-managment/internal/config"
	"tender-managment/internal/controller"
	"tender-managment/internal/db"
	repository "tender-managment/internal/db/repo"
	"tender-managment/internal/routes"
	"tender-managment/internal/service"
)

// @title Tender Managment Swagger

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	r := gin.Default()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error while loading config %v", err)
	}
	database := db.NewDatabase(&cfg.Database)
	redis := db.NewRedisClient(&cfg.Database)
	fmt.Println(redis)
	userRepo := repository.NewUserRepository(database)
	authService := service.NewAuthService(userRepo)
	tenderRepo := repository.NewTenderRepository(database)
	tenderService := service.NewTenderService(tenderRepo)
	bidRepo := repository.NewBidRepository(database)
	bidService := service.NewBidService(*bidRepo, *tenderRepo, *userRepo)
	userService := service.NewUserService(userRepo)
	controller.SetAuthService(authService)
	controller.SetTenderService(tenderService)
	controller.SetBidService(bidService)
	controller.SetUserService(userService)
	routes.SetupRoutes(r)
	r.Run(":8888")
}
