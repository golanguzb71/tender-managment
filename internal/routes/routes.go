package routes

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"tender-managment/internal/controller"
	"tender-managment/internal/utils"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	client := r.Group("/api/client")
	{
		client.POST("/tenders", utils.AuthMiddleware([]string{"client"}), controller.CreateTender)
		client.GET("/tenders", utils.AuthMiddleware([]string{"client", "contractor"}), controller.ListTenders)
		client.PUT("/tenders/:id", utils.AuthMiddleware([]string{"client"}), controller.UpdateTenderStatus)
		client.DELETE("/tenders/:id", utils.AuthMiddleware([]string{"client"}), controller.DeleteTender)
	}
	
}
