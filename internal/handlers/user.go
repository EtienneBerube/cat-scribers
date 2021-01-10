package handlers

import (
	"fmt"
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"github.com/EtienneBerube/cat-scribers/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UpdateUser handles any request to update a user. The authenticated user must be the owner of this account
func UpdateUser(c *gin.Context) {
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	updatedUser := models.User{}
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ok, err := services.UpdateUser(currentUserID, updatedUser)
	if err != nil || !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successful"})
	return
}

// Returns the current authenticated user. This information is taken from the context (provided by the auth middleware)
func GetCurrentUser(c *gin.Context) {
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	user, err := services.GetUserById(currentUserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No user with id %s found", currentUserID)})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

// GetAllUsers handles any request to get all the users registered
func GetAllUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
	return
}

// GetUserByID handles any requests to get a user by its ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := services.GetUserById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No user with id %s found", id)})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

// SubscribeTo handles requests to subscribe the currently authenticated user to another user
func SubscribeTo(c *gin.Context) {
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	subscribedToID := c.Param("id")

	user, err := services.SubscribeTo(currentUserID, subscribedToID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

// UnsubscribeFrom handles requests to unsubscribe the currently authenticated user from another user
func UnsubscribeFrom(c *gin.Context) {
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	subscribedToID := c.Param("id")

	user, err := services.UnsubscribeFrom(currentUserID, subscribedToID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}
