package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tender-managment/internal/model"
	"tender-managment/internal/service"
)

var (
	bidService *service.BidService
)

func SetBidService(bidSer *service.BidService) {
	bidService = bidSer
}

func CreateBidHandler(c *gin.Context) {
	tenderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	var bid model.CreateBid
	if err := c.ShouldBindJSON(&bid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if bid.Price <= 0 || bid.DeliveryTime <= 0 || bid.Comments == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid bid data"})
		return
	}
	contractorId := c.GetInt("user_id")

	createdBid, status, err := bidService.CreateBid(contractorId, tenderId, bid)
	if err != nil {
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(status, createdBid)
}

func GetBidsByTenderID(c *gin.Context) {
	tenderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	bids, err := bidService.GetBidsByTenderID(tenderId, c.GetInt("user_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}

func GetBidByIDHandler(c *gin.Context) {
	bidId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	contractorId := c.GetInt("user_id")

	bid, err := bidService.GetBidByID(contractorId, bidId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bid)
}

func UpdateBidStatusHandler(c *gin.Context) {
	contractorId := c.GetInt("user_id")
	bidId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	var updateData struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	bid, err := bidService.GetBidByID(contractorId, bidId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = bidService.UpdateBidStatus(contractorId, bidId, updateData.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bid status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bid status updated successfully",
		"bid":     bid,
	})
}

func GetBidsByContractor(c *gin.Context) {
	contractorId := c.GetInt("user_id")
	bids, err := bidService.GetBidsByContractor(contractorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bids"})
		return
	}

	c.JSON(http.StatusOK, bids)
}

func DeleteBidHandler(c *gin.Context) {
	bidId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	contractorId := c.GetInt("user_id")
	err = bidService.DeleteBid(contractorId, bidId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bid deleted successfully"})
}

func AwardBidHandler(c *gin.Context) {
	tenderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	bidId, err := strconv.Atoi(c.Param("bidId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	clientId := c.GetInt("user_id")

	err = bidService.AwardBid(clientId, tenderId, bidId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bid awarded successfully"})
}
