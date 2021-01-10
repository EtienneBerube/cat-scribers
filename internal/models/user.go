package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	Subscriptions     []string `json:"subscriptions"`
	Photos            []string `json:"photos"`
	Balance           int64    `json:"balance"`
	SubscriptionPrice int64    `json:"subscription_price"`
}

// IsSubscribedTo checks is a user is subscribed to another user by its ID
func (u User) IsSubscribedTo(id string) bool {
	for _, sub := range u.Subscriptions {
		if id == sub {
			return true
		}
	}
	return false
}

// ToDAO transfers the data from a User struct to a UserDAO struct
func (u *User) ToDAO(user *UserDAO) {
	user.Name = u.Name
	user.Email = u.Email
	user.ID = u.ID
	user.Balance = u.Balance
	user.SubscriptionPrice = u.SubscriptionPrice

	user.Subscriptions = []primitive.ObjectID{}
	for _, id := range u.Subscriptions {
		primId, _ := primitive.ObjectIDFromHex(id)
		user.Subscriptions = append(user.Subscriptions, primId)
	}

	user.Photos = []primitive.ObjectID{}
	for _, id := range u.Photos {
		primId, _ := primitive.ObjectIDFromHex(id)
		user.Photos = append(user.Photos, primId)
	}
}

// UserDAO is a version of User that is used by MongoDB to interact with the data
type UserDAO struct {
	ID                string               `bson:"_id,omitempty"`
	Name              string               `bson:"name"`
	Email             string               `bson:"email"`
	Subscriptions     []primitive.ObjectID `json:"subscriptions"`
	Photos            []primitive.ObjectID `bson:"photos"`
	Balance           int64                `bson:"balance"`
	SubscriptionPrice int64                `bson:"subscription_price"`
}

// ToModel transfers the data from a UserDAO struct to a User struct
func (u *UserDAO) ToModel(user *User) {
	user.Name = u.Name
	user.Email = u.Email
	user.ID = u.ID
	user.Balance = u.Balance
	user.SubscriptionPrice = u.SubscriptionPrice

	user.Subscriptions = []string{}
	for _, id := range u.Subscriptions {
		user.Subscriptions = append(user.Subscriptions, id.Hex())
	}

	user.Photos = []string{}
	for _, id := range u.Photos {
		user.Photos = append(user.Photos, id.Hex())
	}

}
