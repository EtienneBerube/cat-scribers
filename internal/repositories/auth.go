package repositories

import (
	"errors"
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAuthByEmail(email string) (*models.UserAuth, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("auth")

	query := bson.M{"email": email}

	userAuthDAO := models.UserAuthDAO{}

	err := col.FindOne(ctx, query).Decode(&userAuthDAO)
	if err != nil {
		return nil, err
	}

	userAuth := models.UserAuth{}

	userAuthDAO.ToModel(&userAuth)

	return &userAuth, nil
}

func SaveAuth(auth *models.UserAuth) (string, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("auth")

	userAuthDAO := models.UserAuthDAO{}
	auth.ToDAO(&userAuthDAO)

	result, err := col.InsertOne(ctx, userAuthDAO)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("Could not Cast InsertedID to ObjectID")
	}

	return oid.Hex(), nil
}

func UpdateAuth(id string, auth *models.UserAuth) (string, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("auth")

	userAuthDAO := models.UserAuthDAO{}
	auth.ToDAO(&userAuthDAO)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	_, err = col.ReplaceOne(
		ctx,
		bson.M{"_id": oid},
		userAuthDAO,
	)

	if err != nil {
		return "", nil
	}

	return id, nil
}

func DeleteAuth(id string) error {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("auth")

	_, err := col.DeleteOne(ctx, bson.M{"_id": primitive.ObjectIDFromHex(id)})

	return err
}

func IsEmailUsed(email string) (bool, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("auth")

	count, err := col.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
