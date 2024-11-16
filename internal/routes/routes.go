package routes

import (
	"github.com/gin-gonic/gin"
	"tender-managment/internal/controller"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)

	r.POST("/tenders", controller.CreateTender)
	r.GET("/tenders", controller.ListTenders)
	r.GET("/tenders/:id", controller.GetTender)
	r.PUT("/tenders/:id", controller.UpdateTenderStatus)
	r.DELETE("/tenders/:id", controller.DeleteTender)

	r.POST("/tenders/:id/bids", controller.SubmitBid)
	r.GET("/tenders/:id/bids", controller.ListBids)
	r.GET("/tenders/:id/bids", controller.FilterAndSortBids)

	r.POST("/tenders/:id/award/:bid_id", controller.AwardTender)

	r.GET("/users/:id/tenders", controller.ListUserTenders)
	r.GET("/users/:id/bids", controller.ListUserBids)

	r.POST("/tenders/:id/bids", controller.SubmitBid)

}
