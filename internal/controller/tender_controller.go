package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tender-managment/internal/db"
	"tender-managment/internal/model"
	"tender-managment/internal/service"
	"time"
)

const (
	tenderListCacheKey = "tenders:client:%d"
	cacheExpiration    = 5 * time.Minute
)

var (
	tenderService *service.TenderService
	redisClient   *db.Redis
)

func SetTenderService(tenderSer *service.TenderService, redis *db.Redis) {
	tenderService = tenderSer
	redisClient = redis
}

// CreateTenderHandler godoc
// @Summary Create a new tender
// @Description Creates a new tender with provided details
// @Tags Tender
// @Accept json
// @Produce json
// @Param tender body model.CreateTender true "Tender details"
// @Success 201 {object} model.Tender
// @Failure 400 {object} map[string]string "Invalid input or tender data"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /api/client/tenders [post]
func CreateTenderHandler(c *gin.Context) {
	var payload model.CreateTender

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	if payload.Budget < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender data"})
		return
	}

	parsedTime, err := time.Parse(time.RFC3339, payload.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid deadline format"})
		return
	}

	tender := model.Tender{
		Title:       payload.Title,
		Description: payload.Description,
		Deadline:    parsedTime,
		Budget:      payload.Budget,
		Attachment:  payload.Attachment,
		ClientID:    c.GetInt("user_id"),
		Status:      "open",
	}

	createdTender, err := tenderService.CreateTender(&tender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create tender", "error": err.Error()})
		return
	}

	cacheKey := fmt.Sprintf(tenderListCacheKey, tender.ClientID)
	_ = redisClient.Del(c.Request.Context(), cacheKey)

	c.JSON(http.StatusCreated, createdTender)
}

// ListTendersHandler godoc
// @Summary List all tenders for a client
// @Description Retrieves a list of tenders for a specific client
// @Tags Tender
// @Produce json
// @Success 200 {array} model.Tender
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /api/client/tenders [get]
func ListTendersHandler(c *gin.Context) {
	clientID := c.GetInt("user_id")
	cacheKey := fmt.Sprintf(tenderListCacheKey, clientID)

	cachedData, err := redisClient.Get(c.Request.Context(), cacheKey)
	if err == nil && cachedData != "" {
		var tenders []model.Tender
		if err := json.Unmarshal([]byte(cachedData), &tenders); err == nil {
			c.JSON(http.StatusOK, tenders)
			return
		}
	}

	tenders, err := tenderService.ListTenders(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch tenders"})
		return
	}

	if tendersJSON, err := json.Marshal(tenders); err == nil {
		_ = redisClient.Set(c.Request.Context(), cacheKey, tendersJSON, cacheExpiration)
	}

	c.JSON(http.StatusOK, tenders)
}

// UpdateTenderStatusHandler godoc
// @Summary Update the status of a tender
// @Description Updates the status of a tender by ID
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path int true "Tender ID"
// @Param status body model.UpdateTenderStatusRequest true "New tender status"
// @Success 200 {object} map[string]string "Tender status updated successfully"
// @Failure 400 {object} map[string]string "Invalid input or tender status"
// @Failure 404 {object} map[string]string "Tender not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /api/client/tenders/{id} [put]
func UpdateTenderStatusHandler(c *gin.Context) {
	tenderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender ID"})
		return
	}

	var statusUpdate model.UpdateTenderStatusRequest
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	clientID := c.GetInt("user_id")
	if err := tenderService.UpdateTenderStatus(clientID, tenderID, statusUpdate.Status); err != nil {
		switch err.Error() {
		case "tender not found":
			c.JSON(http.StatusNotFound, gin.H{"message": "Tender not found"})
		case "invalid status":
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender status"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update tender", "error": err.Error()})
		}
		return
	}

	cacheKey := fmt.Sprintf(tenderListCacheKey, clientID)
	_ = redisClient.Del(c.Request.Context(), cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Tender status updated successfully"})
}

// DeleteTenderHandler godoc
// @Summary Delete a tender
// @Description Deletes a tender by ID
// @Tags Tender
// @Param id path int true "Tender ID"
// @Success 200 {object} map[string]string "Tender deleted successfully"
// @Failure 400 {object} map[string]string "Invalid tender ID"
// @Failure 404 {object} map[string]string "Tender not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /api/client/tenders/{id} [delete]
func DeleteTenderHandler(c *gin.Context) {
	tenderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender ID"})
		return
	}

	clientID := c.GetInt("user_id")
	if err := tenderService.DeleteTender(clientID, tenderID); err != nil {
		if err.Error() == "tender not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Tender not found or access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete tender", "error": err.Error()})
		return
	}

	cacheKey := fmt.Sprintf(tenderListCacheKey, clientID)
	_ = redisClient.Del(c.Request.Context(), cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "Tender deleted successfully"})
}

// GetClientTenderHistory godoc
// @Summary Retrieve Client's Tender History
// @Description Retrieves a list of tenders posted by a specific client
// @Tags User
// @Produce json
// @Param id path int true "Client ID"
// @Success 200 {array} model.Tender "List of tenders posted by the client"
// @Failure 400 {object} map[string]string "Invalid client ID"
// @Failure 404 {object} map[string]string "No tenders found for the client"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security Bearer
// @Router /api/users/{id}/tenders [get]
func GetClientTenderHistory(ctx *gin.Context) {
	clientID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid client ID"})
		return
	}

	cacheKey := fmt.Sprintf(tenderListCacheKey, clientID)
	cachedData, err := redisClient.Get(ctx.Request.Context(), cacheKey)
	if err == nil && cachedData != "" {
		var tenders []model.Tender
		if err := json.Unmarshal([]byte(cachedData), &tenders); err == nil {
			ctx.JSON(http.StatusOK, tenders)
			return
		}
	}

	tenders, err := tenderService.ListTenders(clientID)
	if err != nil {
		if err.Error() == "no tenders found" {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "No tenders found for this client"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch tender history", "error": err.Error()})
		return
	}

	if tendersJSON, err := json.Marshal(tenders); err == nil {
		_ = redisClient.Set(ctx.Request.Context(), cacheKey, tendersJSON, cacheExpiration)
	}

	ctx.JSON(http.StatusOK, tenders)
}
