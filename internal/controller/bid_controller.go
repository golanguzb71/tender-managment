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

// CreateBidHandler handles bid submission
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
	contractorId := c.GetInt("user_id") // Assuming the contractor ID comes from the JWT token

	// Create the bid
	createdBid, status, err := bidService.CreateBid(contractorId, tenderId, bid)
	if err != nil {
		c.JSON(status, gin.H{"message": err.Error()})
		return
	}

	c.JSON(status, createdBid)
}

// GetBidsByTenderID retrieves all bids for a specific tender (Client view)
func GetBidsByTenderID(c *gin.Context) {
	tenderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	// Get the bids for the tender
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

	contractorId := c.GetInt("user_id") // Assuming this is a utility to extract the contractor ID from the token

	// Get the specific bid
	bid, err := bidService.GetBidByID(contractorId, bidId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// UpdateBidStatusHandler handles the request to update the status of a bid.
func UpdateBidStatusHandler(c *gin.Context) {
	contractorId := c.GetInt("user_id") // Assuming user_id is extracted from the token
	bidId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	var updateData struct {
		Status string `json:"status"`
	}

	// Bind the status from the request body
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get the bid by ID and ensure it's the contractor's bid
	bid, err := bidService.GetBidByID(contractorId, bidId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Update the bid status
	err = bidService.UpdateBidStatus(contractorId, bidId, updateData.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bid status"})
		return
	}

	// Respond with the updated bid
	c.JSON(http.StatusOK, gin.H{
		"message": "Bid status updated successfully",
		"bid":     bid, // You can choose to return the updated bid if needed
	})
}

// GetBidsByContractor retrieves all bids placed by a contractor
func GetBidsByContractor(c *gin.Context) {
	contractorId := c.GetInt("user_id") // Assuming this is a utility to extract user ID from the token
	bids, err := bidService.GetBidsByContractor(contractorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bids"})
		return
	}

	c.JSON(http.StatusOK, bids)
}

// DeleteBidHandler deletes a contractor's bid
func DeleteBidHandler(c *gin.Context) {
	bidId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	contractorId := c.GetInt("user_id")
	err = bidService.DeleteBid(contractorId, bidId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bid deleted successfully"})
}

// AwardBidHandler allows the client to award a bid
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

	clientId := c.GetInt("user_id") // Extract client ID from the token

	// Award the bid
	err = bidService.AwardBid(clientId, tenderId, bidId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"nessage": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bid awarded successfully"})
}
