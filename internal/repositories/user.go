package repositories

import (
	"errors"
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUserById returns a user by its ID from MongoDB
func GetUserById(id string) (*models.User, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("cat-scribers").Collection("users")

	userDAO := models.UserDAO{}
	user := models.User{}
	user.ToDAO(&userDAO)

	oid, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": oid}

	err := col.FindOne(ctx, query).Decode(&userDAO)
	if err != nil {
		return nil, err
	}

	userDAO.ToModel(&user)

	return &user, nil
}

// GetAllUsers returns all the user registered from MongoDB
func GetAllUsers() ([]models.User, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("cat-scribers").Collection("users")

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var usersDAO []models.UserDAO
	if err = cursor.All(ctx, &usersDAO); err != nil {
		return nil, err
	}

	var users []models.User

	for _, dao := range usersDAO {
		user := models.User{}
		dao.ToModel(&user)
		users = append(users, user)
	}

	return users, nil
}

// GetAllUsersSubscribedTo Gets all the users subscribed to another user from MongoDB
func GetAllUsersSubscribedTo(id string) ([]models.User, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("cat-scribers").Collection("users")

	oid, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{
		"subscriptions": bson.M{
			"$elemMatch": bson.M{
				"$eq": oid,
			},
		},
	}
	cursor, err := col.Find(ctx, query)

	if err != nil {
		return nil, err
	}
	var usersDAO []models.UserDAO
	if err = cursor.All(ctx, &usersDAO); err != nil {
		return nil, err
	}

	var users []models.User

	for _, dao := range usersDAO {
		user := models.User{}
		dao.ToModel(&user)
		users = append(users, user)
	}

	return users, nil
}
// SaveUser saves a user to MongoDB
func SaveUser(user models.User) (string, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("cat-scribers").Collection("users")

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

// UpdateUser updates a user in MongoDB
func UpdateUser(id string, newUser *models.User) (bool, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("cat-scribers").Collection("users")

	userDAO := models.UserDAO{}
	newUser.ToDAO(&userDAO)

	userDAO.ID = "" // To preserve ID

	oid, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": oid}

	_, err := col.ReplaceOne(ctx, query, userDAO)
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteUser deletes a user in mongodb
func DeleteUser(id string) error {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("users")

	oid, _ := primitive.ObjectIDFromHex(id)

	_, err := col.DeleteOne(ctx, bson.M{"_id": oid})

	return err
}
