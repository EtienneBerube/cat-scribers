package repositories

import (
	"errors"
	"github.com/EtienneBerube/only-cats/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserById(id string) (*models.User, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("only-cats").Collection("users")

	userDAO := models.UserDAO{}
	user := models.User{}
	user.ToDAO(&userDAO)

	query := bson.M{"_id": primitive.ObjectIDFromHex(id)}

	err := col.FindOne(ctx, query).Decode(&userDAO)
	if err != nil {
		return nil, err
	}

	userDAO.ToModel(&user)

	return &user, nil
}

func SaveUser(user models.User) (string, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("only-cats").Collection("users")

	userDAO := models.UserDAO{}
	user.ToDAO(&userDAO)

	result, err := col.InsertOne(ctx, userDAO)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("Could not Cast InsertedID to ObjectID")
	}

	return oid.Hex(), nil
}

func UpdateUser(id string, newUser *models.User) (bool, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("only-cats").Collection("users")

	userDAO := models.UserDAO{}
	newUser.ToDAO(&userDAO)

	query := bson.M{"_id": id}

	_, err := col.ReplaceOne(ctx, query, userDAO)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteUser(id string) error {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("only-cats").Collection("users")

	_, err := col.DeleteOne(ctx, bson.M{"_id": primitive.ObjectIDFromHex(id)})

	return err
}
