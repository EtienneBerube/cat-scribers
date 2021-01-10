package repositories

import (
	"errors"
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAuthByEmail returns a UserAuth by its email from MongoDB
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
// GetAuthByEmail saves a UserAuth to MongoDB
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

// UpdateAuth updates a UserAuth in MongoDB
func UpdateAuth(id string, auth *models.UserAuth) (bool, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("auth")

	userAuthDAO := models.UserAuthDAO{}
	auth.ToDAO(&userAuthDAO)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	_, err = col.ReplaceOne(
		ctx,
		bson.M{"_id": oid},
		userAuthDAO,
	)

	if err != nil {
		return false, nil
	}

	return true, nil
}

// DeleteAuth deletes a UserAuth in MongoDB
func DeleteAuth(id string) error {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("auth")

	oid, _ := primitive.ObjectIDFromHex(id)
	_, err := col.DeleteOne(ctx, bson.M{"_id": oid})

	return err
}

// IsEmailUsed checks if a given email is used by another user in MongoDB
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
