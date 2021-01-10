package handlers

import (
	"fmt"
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"github.com/EtienneBerube/cat-scribers/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateUser(c *gin.Context){
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

func GetCurrentUser(c *gin.Context){
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

func GetUserByID(c *gin.Context){
	id := c.Param("id")

	user, err := services.GetUserById(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No user with id %s found", id)})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

func SubscribeTo(c *gin.Context){
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

func UnsubscribeFrom(c *gin.Context){
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	subscribedToID := c.Param("id")

	user, err := services.UnsubscribeFrom(currentUserID, subscribedToID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
	return
}
