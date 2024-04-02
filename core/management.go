package core

import (
	"encoding/json"
	"example/web-service-gin/models"
	"example/web-service-gin/utils/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AddBasketInput struct {
	Data  map[string]interface{} `json:"data" gorm:"type:jsonb" binding:"required"`
	State models.BasketState     `json:"state" binding:"required"`
}
type EditBasketInput struct {
	Data  map[string]interface{} `json:"data,omitempty"`
	State models.BasketState     `json:"state,omitempty"`
}

func AddBasket(c *gin.Context) {
	// Extract user ID from the token
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Get user
	user, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "user": user})
		return
	}

	var input AddBasketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert map to JSON-formatted string
	dataJSON, err := json.Marshal(input.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert data to JSON"})
		return
	}

	dataBytes := []byte(dataJSON)
	// Create new Basket instance
	newBasket := models.Basket{
		UserID: userID,
		Data:   dataBytes,
		State:  input.State,
	}

	// Save the Basket to the database
	basket, err := newBasket.SaveBasket()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "basket": basket})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Basket created successfully", "basket": basket})
}

func UpdateBasket(c *gin.Context) {
	// Extract user ID from the token
	userID, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-unauthorized": err.Error()})
		return
	}

	// Get user
	user, err := models.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "user": user})
		return
	}

	// Parse request body
	var input EditBasketInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the existing basket from the database
	basketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid basket ID"})
		return
	}

	existingBasket, err := models.GetBasketObjectByID(uint(basketID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check is basket for current user
	if userID != uint(existingBasket.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "error : this basket is not yours"})
		return
	}
	// Check if the state is "completed"
	if existingBasket.State == models.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This basket can NOT be updated. Its a COMPLETED basket!"})
		return
	}

	// Update the basket fields if provided
	if input.Data != nil {
		// Convert map to JSON-formatted string
		dataJSON, err := json.Marshal(input.Data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert data to JSON"})
			return
		}
		existingBasket.Data = dataJSON
	}

	if input.State != "" {
		if input.State != models.Completed && input.State != models.Pending {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
			return
		}
		existingBasket.State = input.State
	}

	// Save the updated basket to the database
	basket, err := existingBasket.SaveBasket()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Basket updated successfully", "basket": basket})
}

func GetBasketByID(c *gin.Context) {

	// Fetch the existing basket from the database
	basketID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid basket ID"})
		return
	}
	b, err := models.GetBasketByID(uint(basketID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID from the token
	user_id, err := token.ExtractTokenID(c)

	// Check is basket for current user
	if user_id != uint(b.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "error : this basket is not yours"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success : successfull fetch", "basket": b})
}
func DeleteBasketByID(c *gin.Context) {

	// Fetch the existing basket from the database
	basketID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid basket ID"})
		return
	}
	b, err := models.GetBasketByID(uint(basketID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID from the token
	user_id, err := token.ExtractTokenID(c)

	// Check is basket for current user
	if user_id != uint(b.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "error : this basket is not yours"})
		return
	}

	// Delete basket
	if err := models.DeleteBasket(uint(basketID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success : successfull delete"})
}
func GetBasketsList(c *gin.Context) {

	// Extract user ID from the token
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error-unauthorized": err.Error()})
		return
	}

	// Fetch baskets of current user
	baskets, err := models.GetBasketsByUserID(int64(user_id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success : successfull fetch", "baskets": baskets})
}
