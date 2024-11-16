package main

import (
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
	userRepo := repository.NewUserRepository(database)
	authService := service.NewAuthService(userRepo)
	controller.SetAuthService(authService)

	routes.SetupRoutes(r)
	r.Run(":8888")
}
