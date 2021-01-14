package repositories

import (
	"errors"
	"github.com/EtienneBerube/cat-scribers/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetPhotoByID returns a Photo by its id from MongoDB
func GetPhotoByID(id string) (*models.Photo, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("photos")

	oid, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": oid}

	photoDAO := models.PhotoDAO{}

	err := col.FindOne(ctx, query).Decode(&photoDAO)
	if err != nil {
		return nil, err
	}

	photo := models.Photo{}

	photoDAO.ToModel(&photo)

	return &photo, nil
}

// GetPhotoByID returns a Photo by its name from MongoDB
func GetPhotoByName(name string) (*models.Photo, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("photos")

	query := bson.M{"name": name}

	photoDAO := models.PhotoDAO{}

	err := col.FindOne(ctx, query).Decode(&photoDAO)
	if err != nil {
		return nil, err
	}

	photo := models.Photo{}

	photoDAO.ToModel(&photo)

	return &photo, nil
}

// GetAllPhotosByOwnerID returns all the photos from a user by its ID from MongoDB
func GetAllPhotosByOwnerID(ownerID string) ([]models.Photo, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("photos")

	oid, _ := primitive.ObjectIDFromHex(ownerID)

	query := bson.M{"owner_id": oid}

	cursor, err := col.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	var results []models.PhotoDAO
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var photos []models.Photo

	for _, dao := range results {
		photo := models.Photo{}
		dao.ToModel(&photo)
		photos = append(photos, photo)
	}
	return photos, nil
}

// SearchPhotosByNameContaining returns all the photos from a user where the name partially matches the provided name
func SearchPhotosByNameContaining(name string, ownerID string) ([]models.Photo, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("photos")

	oid, _ := primitive.ObjectIDFromHex(ownerID)

	ownerMatch := bson.D{{
		"$match", bson.D{
			{"owner_id", oid},
		},
	}}

	searchQuery := bson.D{{
		"$search", bson.D{
			{"text", bson.D{
				{"query", name},
				{"path", "name"},
			},
			},
		},
	}}

	var results []models.PhotoDAO

	cursor, err := col.Aggregate(ctx, mongo.Pipeline{searchQuery, ownerMatch})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var photos []models.Photo

	for _, dao := range results {
		photo := models.Photo{}
		dao.ToModel(&photo)
		photos = append(photos, photo)
	}
	return photos, nil
}

// SavePhoto saves a photo to MongoDB
func SavePhoto(photo *models.Photo) (string, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("cat-scribers").Collection("photos")

	photoDAO := models.PhotoDAO{}
	photo.ToDAO(&photoDAO)

	result, err := col.InsertOne(ctx, photoDAO)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("Could not Cast InsertedID to ObjectID")
	}

	return oid.Hex(), nil
}

// SaveMultiplePhotos saves multiple photos to MongoDB
func SaveMultiplePhotos(photos []models.Photo) ([]string, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("cat-scribers").Collection("photos")

	var daos []interface{}

	for _, photo := range photos {
		photoDAO := models.PhotoDAO{}
		photo.ToDAO(&photoDAO)
		daos = append(daos, photoDAO)
	}

	results, err := col.InsertMany(ctx, daos)
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, oid := range results.InsertedIDs {
		oid, _ := oid.(primitive.ObjectID)
		ids = append(ids, oid.Hex())
	}

	return ids, nil
}

// DeletePhoto deletes a photo from MongoDB
func DeletePhoto(id string) error {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("cat-scribers").Collection("photos")

	oid, _ := primitive.ObjectIDFromHex(id)

	_, err := col.DeleteOne(ctx, bson.M{"_id": oid})

	return err
}
