package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignUpRequest struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	SubscriptionPrice int64 `json:"subscription_price"`
}

type SignUpResponse struct {
	ID    string `json:"user_id"`
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID    string `json:"user_id"`
	Token string `json:"token"`
}

type UserAuth struct {
	ID           string `json:"id,omitempty"`
	Email        string `json:"email"`
	UserID       string `json:"user_id"`
	PasswordHash string `json:"password_hash"`
}

func (u *UserAuth) ToDAO(dao *UserAuthDAO) {
	dao.ID, _ = primitive.ObjectIDFromHex(u.ID)
	dao.Email = u.Email
	dao.UserID, _ = primitive.ObjectIDFromHex(u.UserID)
	dao.PasswordHash = u.PasswordHash
}

type UserAuthDAO struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email"`
	UserID       primitive.ObjectID `bson:"user_id"`
	PasswordHash string             `bson:"password_hash"`
}

func (dao *UserAuthDAO) ToModel(auth *UserAuth) {
	auth.ID = dao.ID.Hex()
	auth.Email = dao.Email
	auth.UserID = dao.UserID.Hex()
	auth.PasswordHash = dao.PasswordHash
}
