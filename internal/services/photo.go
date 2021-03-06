package services

import (
	"errors"
	"fmt"
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"github.com/EtienneBerube/cat-scribers/internal/repositories"
	"github.com/EtienneBerube/cat-scribers/pkg/vision"
)

// GetPhotoByID returns a saved photo by its ID
func GetPhotoByID(currentUserID string, photoID string) (*models.Photo, error) {
	currentUser, err := GetUserById(currentUserID)
	if err != nil {
		return nil, err
	}

	photo, err := repositories.GetPhotoByID(photoID)
	if err != nil {
		return nil, err
	}

	if !currentUser.IsSubscribedTo(photo.OwnerID) {
		return nil, errors.New("Unauthorized: Not subscribed to Owner")
	}

	return photo, err
}

/* GetAllPhotosFromOwner returns all the pictures from a user, given that the currently authenticated user is subscribed
to the owner of the pictures
 */
func GetAllPhotosFromOwner(currentUserID string, ownerID string) ([]models.Photo, error) {
	currentUser, err := GetUserById(currentUserID)
	if err != nil {
		return nil, err
	}

	if !currentUser.IsSubscribedTo(ownerID) {
		return nil, errors.New("Unauthorized: Not subscribed to Owner")
	}

	photos, err := repositories.GetAllPhotosByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	return photos, nil
}
/* SearchPhotosOfOwnerByName will search all photos from a given owner for photos where the name partially match the
name provided
*/
func SearchPhotosOfOwnerByName(currentUserID string, ownerID string, name string) ([]models.Photo, error) {
	currentUser, err := GetUserById(currentUserID)
	if err != nil {
		return nil, err
	}

	if !currentUser.IsSubscribedTo(ownerID) {
		return nil, errors.New("Unauthorized: Not subscribed to Owner")
	}

	photos, err := repositories.SearchPhotosByNameContaining(name, ownerID)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

// CreatePhoto creates a new photo
func CreatePhoto(photo models.Photo) (string, error) {
	ok, err := vision.HasCat(photo.Base64, photo.Name)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("This image does not contain a cat. You clearly didn't read the terms and services... liar")
	}
	id, err := repositories.SavePhoto(&photo)
	if err != nil {
		return "", err
	}
	return id, nil
}

//CreateMultiplePhotos creates multiple photos
func CreateMultiplePhotos(photos []models.Photo) ([]string, []string, error) {
	names := make([]string, len(photos))
	for i, photo := range photos {
		names[i] = photo.Name
	}
	b64s := make([]string, len(photos))
	for i, photo := range photos {
		b64s[i] = photo.Base64
	}
	oks, err := vision.HasCatMultiple(b64s, names)
	if err != nil {
		return nil, nil, err
	}
	var rejected = []string{}
	var accepted = []models.Photo{}

	for i, ok := range oks {
		if ok {
			accepted = append(accepted, photos[i])
		} else {
			rejected = append(rejected, photos[i].Name)
		}
	}

	ids, err := repositories.SaveMultiplePhotos(accepted)
	if err != nil {
		return nil, nil, err
	}
	return ids, rejected, nil
}

// DeletePhoto deletes a photo
func DeletePhoto(currentUserID string, photoID string) error {
	photo, err := repositories.GetPhotoByID(photoID)
	if err != nil {
		return err
	}

	if photo.OwnerID != currentUserID {
		return errors.New(fmt.Sprintf("Cannot delete photo. %s is not the owner", currentUserID))
	}

	return repositories.DeletePhoto(photoID)
}
