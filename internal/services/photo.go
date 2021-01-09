package services

import (
	"errors"
	"github.com/EtienneBerube/only-cats/internal/models"
	"github.com/EtienneBerube/only-cats/internal/repositories"
	"github.com/EtienneBerube/only-cats/pkg/vision"
)

func GetPhotoByID(id string) (*models.Photo, error) {
	photo, err := repositories.GetPhotoByID(id)
	if err != nil {
		return nil, err
	}
	return photo, err
}

func GetAllPhotosFromOwner(currentUserID string, ownerID string) ([]models.Photo, error){
	currentUser, err := GetUserById(currentUserID)
	if err != nil {
		return nil, err
	}

	if !currentUser.IsSubscribedTo(ownerID){
		return nil, errors.New("Unauthorized: Not subscribed to Owner")
	}

	photos, err := repositories.GetAllPhotosByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

func SearchPhotosOfOwnerByName(currentUserID string, ownerID string, name string) ([]models.Photo, error){
	currentUser, err := GetUserById(currentUserID)
	if err != nil {
		return nil, err
	}

	if !currentUser.IsSubscribedTo(ownerID){
		return nil, errors.New("Unauthorized: Not subscribed to Owner")
	}

	photos, err := repositories.SearchPhotosByNameContaining(name, ownerID)
	if err != nil {
		return nil, err
	}
	return photos, nil
}

func CreatePhoto(photo models.Photo) (string, error){
	ok ,err := vision.HasCat(photo.Base64, photo.Name)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("This image does not contain a cat. You clearly didn't read the terms and services... liar")
	}
	id ,err := repositories.SavePhoto(&photo)
	if err != nil {
		return "", err
	}
	return id, nil
}

func DeletePhoto(id string) error{
	err := repositories.DeletePhoto(id)
	return err
}