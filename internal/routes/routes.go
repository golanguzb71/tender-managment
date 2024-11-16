package routes

import (
	"github.com/gin-gonic/gin"
	"tender-managment/internal/controller"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
}
