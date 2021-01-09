package services

import (
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"github.com/EtienneBerube/cat-scribers/internal/repositories"
	"log"
)

func GetUserById(id string) (*models.User, error) {
	user, err := repositories.GetUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetAllUsers() ([]models.User, error) {
	users, err := repositories.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func PaySubscription(user *models.User) {
	for _, subscribedToID := range user.Subscriptions {
		subscribedTo, err := repositories.GetUserById(subscribedToID)
		if err != nil {
			UnsubscribeFrom(user.ID, subscribedToID)
			log.Printf("ERROR ON CRON: %s", err.Error())
			continue
		}
		if user.Balance - subscribedTo.SubscriptionPrice < 0 {
			UnsubscribeFrom(user.ID, subscribedToID)
			log.Printf("CRON: %s cannot pay for %s's subscription. Unsubscribing...", user.Name, subscribedTo.Name)
			continue
		}else{
			user.Balance -= subscribedTo.SubscriptionPrice
			subscribedTo.Balance += subscribedTo.SubscriptionPrice
			ok , err := repositories.UpdateUser(subscribedTo.ID, subscribedTo)
			if err != nil || !ok {
				log.Printf("CRON: Error while giving the money to %s's. Error: %s", subscribedTo.Name, err.Error())
			}
		}
	}

	ok, err := repositories.UpdateUser(user.ID, user)
	if err != nil || !ok {
		log.Printf("CRON: Error while updating %s's balance. Error: %s", user.Name, err.Error())
	}
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





