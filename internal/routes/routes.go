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
	client.POST("/tenders", utils.AuthMiddleware([]string{"client"}), controller.CreateTenderHandler)
	client.GET("/tenders", utils.AuthMiddleware([]string{"client", "contractor"}), controller.ListTendersHandler)
	client.PUT("/tenders/:id", utils.AuthMiddleware([]string{"client"}), controller.UpdateTenderStatusHandler)
	client.DELETE("/tenders/:id", utils.AuthMiddleware([]string{"client"}), controller.DeleteTenderHandler)
	client.GET("/tenders/:id/bids", utils.AuthMiddleware([]string{"client"}), controller.GetBidsByTenderID)
	client.POST("/tenders/:id/award/:bidId", utils.AuthMiddleware([]string{"client"}), controller.AwardBidHandler)

	contractor := r.Group("/api/contractor")
	contractor.POST("/tenders/:id/bid", utils.AuthMiddleware([]string{"contractor"}), controller.CreateBidHandler)
	contractor.GET("/bids", utils.AuthMiddleware([]string{"contractor"}), controller.GetBidsByContractor)
	contractor.GET("/bids/:id", utils.AuthMiddleware([]string{"contractor"}), controller.GetBidByIDHandler)
	contractor.DELETE("/bids/:id", utils.AuthMiddleware([]string{"contractor"}), controller.DeleteBidHandler)
	contractor.PUT("/bids/:id", utils.AuthMiddleware([]string{"contractor"}), controller.UpdateBidStatusHandler)

	user := r.Group("/api/users")
	user.GET("/:id/bids", utils.AuthMiddleware([]string{"client", "contractor"}), controller.GetContractorBidHistory)
	user.GET("/:id/tenders", utils.AuthMiddleware([]string{"client", "contractor"}), controller.GetClientTenderHistory)

}
