package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tender-managment/internal/model"
	"tender-managment/internal/service"
	"time"
)

var (
	bidService *service.BidService
)

func SetBidService(bidSer *service.BidService) {
	bidService = bidSer
}

const (
	bidsByTenderKey     = "bids:tender:%d"
	bidsByContractorKey = "bids:contractor:%d"
	bidDetailKey        = "bid:%d"

	bidListCacheDuration   = 5 * time.Minute
	bidDetailCacheDuration = 5 * time.Minute
)

// CreateBidHandler godoc
// @Summary Create a bid for a tender
// @Description Create a bid for a given tender with the specified price, delivery time, and comments
// @Tags bids
// @Accept json
// @Produce json
// @Param id path int true "Tender ID"
// @Param bid body model.CreateBid true "Bid Information (e.g., { \"price\": 1000, \"deliveryTime\": 30, \"comments\": \"Delivery within a month\" })"
// @Success 201 {object} model.Bid "Details of the created bid"
// @Failure 400 {object} map[string]string "Invalid tender ID or request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /api/contractor/tenders/{id}/bid [post]
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

	tenderBidsKey := fmt.Sprintf(bidsByTenderKey, tenderId)
	contractorBidsKey := fmt.Sprintf(bidsByContractorKey, contractorId)
	_ = redisClient.Del(c.Request.Context(), tenderBidsKey)
	_ = redisClient.Del(c.Request.Context(), contractorBidsKey)

	c.JSON(status, createdBid)
}

// GetBidsByTenderID godoc
// @Summary Get all bids for a tender
// @Description Retrieve all bids for a given tender, with optional filtering and sorting
// @Tags bids
// @Accept json
// @Produce json
// @Param id path int true "Tender ID"
// @Param price query float64 false "Filter bids by price"
// @Param delivery_time query string false "Filter bids by delivery time"
// @Param sort_by query string false "Sort by 'price' or 'delivery_time'"
// @Success 200 {array} model.Bid "List of bids"
// @Failure 400 {object} map[string]string "Invalid tender ID or query parameters"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "No bids found"
// @Security Bearer
// @Router /api/client/tenders/{id}/bids [get]
func GetBidsByTenderID(c *gin.Context) {
	tenderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	userId := c.GetInt("user_id")
	cacheKey := fmt.Sprintf(bidsByTenderKey, tenderId)

	cachedData, err := redisClient.Get(c.Request.Context(), cacheKey)
	if err == nil && cachedData != "" {
		var bids []model.Bid
		if err := json.Unmarshal([]byte(cachedData), &bids); err == nil {
			c.JSON(http.StatusOK, bids)
			return
		}
	}

	priceFilter, _ := strconv.ParseFloat(c.DefaultQuery("price", "0"), 64)
	deliveryTimeFilter := c.DefaultQuery("delivery_time", "")
	sortBy := c.DefaultQuery("sort_by", "")

	bids, err := bidService.GetBidsByTenderID(tenderId, userId, priceFilter, deliveryTimeFilter, sortBy)
	if err != nil {
		if err.Error() == "Tender not found or access denied" {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if bidsJSON, err := json.Marshal(bids); err == nil {
		_ = redisClient.Set(c.Request.Context(), cacheKey, bidsJSON, bidListCacheDuration)
	}

	c.JSON(http.StatusOK, bids)
}

// GetBidByIDHandler godoc
// @Summary Get bid details by bid ID
// @Description Retrieve details of a specific bid
// @Tags bids
// @Accept json
// @Produce json
// @Param id path int true "Bid ID"
// @Success 200 {object} model.Bid "Bid details"
// @Failure 400 {object} map[string]string "Invalid bid ID"
// @Failure 404 {object} map[string]string "Bid not found"
// @Security Bearer
// @Router /api/contractor/bids/{id} [get]
func GetBidByIDHandler(c *gin.Context) {
	bidId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	contractorId := c.GetInt("user_id")
	cacheKey := fmt.Sprintf(bidDetailKey, bidId)
	cachedData, err := redisClient.Get(c.Request.Context(), cacheKey)
	if err == nil && cachedData != "" {
		var bid model.Bid
		if err := json.Unmarshal([]byte(cachedData), &bid); err == nil {
			if bid.ContractorID == contractorId {
				c.JSON(http.StatusOK, bid)
				return
			}
		}
	}

	bid, err := bidService.GetBidByID(contractorId, bidId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if bidJSON, err := json.Marshal(bid); err == nil {
		_ = redisClient.Set(c.Request.Context(), cacheKey, bidJSON, bidDetailCacheDuration)
	}

	c.JSON(http.StatusOK, bid)
}

// UpdateBidStatusHandler godoc
// @Summary Update the status of a bid
// @Description Update the status of an existing bid (e.g., accepted, rejected)
// @Tags bids
// @Accept json
// @Produce json
// @Param id path int true "Bid ID"
// @Param updateData body model.UpdateBid true "Update bid request body"
// @Success 200 {object} map[string]interface{} "Bid status updated successfully"
// @Failure 400 {object} map[string]string "Invalid bid ID or request data"
// @Failure 404 {object} map[string]string "Bid not found"
// @Failure 500 {object} map[string]string "Failed to update bid status"
// @Security Bearer
// @Router /api/contractor/bids/{id} [put]
func UpdateBidStatusHandler(c *gin.Context) {
	contractorId := c.GetInt("user_id")
	bidId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	var updateData model.UpdateBid
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
	bidDetailKey := fmt.Sprintf(bidDetailKey, bidId)
	tenderBidsKey := fmt.Sprintf(bidsByTenderKey, bid.TenderID)
	contractorBidsKey := fmt.Sprintf(bidsByContractorKey, contractorId)
	_ = redisClient.Del(c.Request.Context(), bidDetailKey)
	_ = redisClient.Del(c.Request.Context(), tenderBidsKey)
	_ = redisClient.Del(c.Request.Context(), contractorBidsKey)

	c.JSON(http.StatusOK, gin.H{
		"message": "Bid status updated successfully",
		"bid":     bid,
	})
}

// GetBidsByContractor godoc
// @Summary Get all bids by a contractor
// @Description Retrieve all bids submitted by a specific contractor
// @Tags bids
// @Accept json
// @Produce json
// @Success 200 {array} model.Bid "List of bids"
// @Failure 500 {object} map[string]string "Failed to fetch bids"
// @Security Bearer
// @Router /api/contractor/bids [get]
func GetBidsByContractor(c *gin.Context) {
	contractorId := c.GetInt("user_id")
	cacheKey := fmt.Sprintf(bidsByContractorKey, contractorId)

	cachedData, err := redisClient.Get(c.Request.Context(), cacheKey)
	if err == nil && cachedData != "" {
		var bids []model.Bid
		if err := json.Unmarshal([]byte(cachedData), &bids); err == nil {
			c.JSON(http.StatusOK, bids)
			return
		}
	}

	bids, err := bidService.GetBidsByContractor(contractorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bids"})
		return
	}

	if bidsJSON, err := json.Marshal(bids); err == nil {
		_ = redisClient.Set(c.Request.Context(), cacheKey, bidsJSON, bidListCacheDuration)
	}

	c.JSON(http.StatusOK, bids)
}

// DeleteBidHandler godoc
// @Summary Delete a bid
// @Description Delete a specific bid by ID
// @Tags bids
// @Accept json
// @Produce json
// @Param id path int true "Bid ID"
// @Success 200 {object} map[string]string "Bid deleted successfully"
// @Failure 400 {object} map[string]string "Invalid bid ID"
// @Failure 404 {object} map[string]string "Bid not found"
// @Security Bearer
// @Router /api/contractor/bids/{id} [delete]
func DeleteBidHandler(c *gin.Context) {
	bidId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	contractorId := c.GetInt("user_id")

	bid, _ := bidService.GetBidByID(contractorId, bidId)

	err = bidService.DeleteBid(contractorId, bidId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	if bid != nil {
		bidDetailKey := fmt.Sprintf(bidDetailKey, bidId)
		tenderBidsKey := fmt.Sprintf(bidsByTenderKey, bid.TenderID)
		contractorBidsKey := fmt.Sprintf(bidsByContractorKey, contractorId)
		_ = redisClient.Del(c.Request.Context(), bidDetailKey)
		_ = redisClient.Del(c.Request.Context(), tenderBidsKey)
		_ = redisClient.Del(c.Request.Context(), contractorBidsKey)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bid deleted successfully"})
}

// AwardBidHandler godoc
// @Summary Award a bid for a tender
// @Description Award a specific bid for a tender, indicating it has been selected
// @Tags bids
// @Accept json
// @Produce json
// @Param id path int true "Tender ID"
// @Param bidId path int true "Bid ID"
// @Success 200 {object} map[string]string "Bid awarded successfully"
// @Failure 400 {object} map[string]string "Invalid tender or bid ID"
// @Failure 404 {object} map[string]string "Bid or tender not found"
// @Security Bearer
// @Router /api/client/tenders/{id}/award/{bidId} [post]
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

	tenderBidsKey := fmt.Sprintf(bidsByTenderKey, tenderId)
	bidDetailKey := fmt.Sprintf(bidDetailKey, bidId)
	_ = redisClient.Del(c.Request.Context(), tenderBidsKey)
	_ = redisClient.Del(c.Request.Context(), bidDetailKey)

	c.JSON(http.StatusOK, gin.H{"message": "Bid awarded successfully"})
}

// GetContractorBidHistory godoc
// @Summary Retrieve Contractor's Bid History
// @Description Retrieves a list of bids placed by a specific contractor
// @Tags User
// @Produce json
// @Param id path int true "Contractor ID"
// @Success 200 {array} model.Bid "List of bids placed by the contractor"
// @Failure 400 {object} map[string]string "Invalid contractor ID"
// @Failure 404 {object} map[string]string "No bids found for the contractor"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /api/users/{id}/bids [get]
func GetContractorBidHistory(ctx *gin.Context) {
	contractorID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid contractor ID"})
		return
	}

	cacheKey := fmt.Sprintf(bidsByContractorKey, contractorID)

	cachedData, err := redisClient.Get(ctx.Request.Context(), cacheKey)
	if err == nil && cachedData != "" {
		var bids []model.Bid
		if err := json.Unmarshal([]byte(cachedData), &bids); err == nil {
			ctx.JSON(http.StatusOK, bids)
			return
		}
	}

	bids, err := bidService.GetBidsByContractor(contractorID)
	if err != nil {
		if err.Error() == "no bids found" {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "No bids found for this contractor"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch bid history", "error": err.Error()})
		return
	}

	if bidsJSON, err := json.Marshal(bids); err == nil {
		_ = redisClient.Set(ctx.Request.Context(), cacheKey, bidsJSON, bidListCacheDuration)
	}

	ctx.JSON(http.StatusOK, bids)
}

// ManaulWebSocketSwag godoc
// @Summary Try to using it by postman it is not work on swagger
// @Description Allows users to receive real-time notifications via WebSocket connection.
// @Accept  json
// @Produce  json
// @Tags User
// @Security Bearer
// @Success 200 {string} string "Successfully connected to WebSocket."
// @Failure 401 {string} string "Unauthorized"
// @Router /api/users/notification/ws [get]
func ManaulWebSocketSwag() {

}
