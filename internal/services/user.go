package services

import (
	"github.com/EtienneBerube/only-cats/internal/models"
	"github.com/EtienneBerube/only-cats/internal/repositories"
)

func GetUserById(id string) (*models.User, error) {
	user, err := repositories.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateNewUser(user models.User) (string, error) {
	id, err := repositories.SaveUser(user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func SubscribeTo(currentUserID string, newSubscriptionID string) (bool, error) {
	currentUser, err := GetUserById(currentUserID)
	if err != nil {
		return false, err
	}
	currentUser.Subscriptions = append(currentUser.Subscriptions, newSubscriptionID)

	ok, err := repositories.UpdateUser(currentUserID, currentUser)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func UnsubscribeFrom(currentUserID string, subscriptionIDToRemove string) (bool, error) {
	currentUser, err := GetUserById(currentUserID)
	if err != nil {
		return false, err
	}
	var removeIndex int
	for index, id := range currentUser.Subscriptions {
		if id == subscriptionIDToRemove {
			removeIndex = index
			break
		}
	}

	currentUser.Subscriptions = append(currentUser.Subscriptions[:removeIndex], currentUser.Subscriptions[removeIndex+1:]...)

	ok, err := repositories.UpdateUser(currentUserID, currentUser)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func DeleteUser(id string) error {
	err := repositories.DeleteUser(id)
	return err
}





