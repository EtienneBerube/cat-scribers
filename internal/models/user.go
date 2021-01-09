package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Subscriptions []string `json:"subscriptions"`
 	Photos []string `json:"photos"`
	Balance int64 `json:"balance"`
	SubscriptionPrice int64 `json:"subscription_price"`
}

func (u User) IsSubscribedTo(id string) bool{
	for _, sub := range u.Subscriptions {
		if id == sub{
			return true
		}
	}
	return false
}

func (u *User) ToDAO(user *UserDAO){
	user.Name = u.Name
	user.Email = u.Email
	user.ID = u.ID
	user.Balance = u.Balance
	user.SubscriptionPrice = u.SubscriptionPrice

	user.Subscribers = []primitive.ObjectID{}

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


type UserDAO struct {
	ID string `bson:"id"`
	Name string `bson:"name"`
	Email string `bson:"email"`
	Subscribers []primitive.ObjectID `bson:"subscribers"`
	Subscriptions []primitive.ObjectID `json:"subscriptions"`
	Photos []primitive.ObjectID `bson:"photos"`
	Balance int64 `bson:"balance"`
	SubscriptionPrice int64 `bson:"subscription_price"`
}

func (u *UserDAO)ToModel(user *User){
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