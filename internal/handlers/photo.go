package handlers

import (
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"github.com/EtienneBerube/cat-scribers/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UploadPhoto Handles any request to upload a photo. The user must be authenticated for the photo to be persisted.
func UploadPhoto(c *gin.Context) {
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

// UploadMultiplePhotos handles any request to upload multiple photos in a single request
func UploadMultiplePhotos(c *gin.Context) {
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
	for _, photo := range photos {
		photo.OwnerID = currentUserID
	}

	ids, rejected, err := services.CreateMultiplePhotos(photos)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ids": ids, "rejected": rejected})
	return
}

// GetPhotoByID handles any request to get a photo by its ID
func GetPhotoByID(c *gin.Context) {
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	photoID := c.Param("id")

	photo, err := services.GetPhotoByID(currentUserID, photoID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, photo)
	return
}

/* GetAllPhotosByOwnerID handles requests where the user wants all the photos from a given user
   This function can perform text search on the name when providing a query parameter called "name"
 */
func GetPhotosByOwnerID(c *gin.Context) {
	id, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No current User found"})
		return
	}
	currentUserID := id.(string)

	ownerID := c.Param("id")
	nameSearch :=  c.Query("name")

	var photos []models.Photo
	var err error
	if nameSearch == "" {
		photos, err = services.GetAllPhotosFromOwner(currentUserID, ownerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}else{
		photos, err = services.SearchPhotosOfOwnerByName(currentUserID, ownerID, nameSearch)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, photos)
	return
}

// DeletePhoto deletes a photo from a user's repository
func DeletePhoto(c *gin.Context) {
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
