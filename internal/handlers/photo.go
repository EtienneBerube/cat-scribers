package handlers

import (
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"github.com/EtienneBerube/cat-scribers/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadPhoto( c *gin.Context){
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	photo := models.Photo{}

	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo.OwnerID = currentUserID
	id, err := services.CreatePhoto(photo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
	return
}

func UploadMultiplePhotos (c *gin.Context){
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	var photos []models.Photo

	if err := c.ShouldBindJSON(&photos); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, photo := range photos{
		photo.OwnerID = currentUserID
	}

	ids,rejected,err := services.CreateMulitplePhotos(photos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ids": ids, "rejected":rejected})
	return
}

func GetPhotoByID( c *gin.Context){
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	photoID := c.Param("id")

	photo , err := services.GetPhotoByID(currentUserID, photoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, photo)
	return
}

func GetPhotosByOwnerID(c *gin.Context){
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	ownerID := c.Param("id")

	photos, err := services.GetAllPhotosFromOwner(currentUserID, ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK,photos)
	return
}

func DeletePhoto(c *gin.Context){
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	photoID := c.Param("id")

	if err := services.DeletePhoto(currentUserID, photoID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful"})
	return
}
