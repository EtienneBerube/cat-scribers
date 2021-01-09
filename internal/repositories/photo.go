package repositories

import (
	"errors"
	"github.com/EtienneBerube/only-cats/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const SEARCH_INDEX_NAME = "photo-search-index-name"

func GetPhotoByID(id string) (*models.Photo, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("only-cats").Collection("photos")

	query := bson.M{"_id": primitive.ObjectIDFromHex(id)}

	photoDAO := models.PhotoDAO{}

	err := col.FindOne(ctx, query).Decode(&photoDAO)
	if err != nil {
		return nil, err
	}

	photo := models.Photo{}

	photoDAO.ToModel(&photo)

	return &photo, nil
}

func GetPhotoByName(name string) (*models.Photo, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("only-cats").Collection("photos")

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

func GetAllPhotosByOwnerID(ownerID string) ([]models.Photo, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("only-cats").Collection("photos")

	query := bson.M{"owner_id": primitive.ObjectIDFromHex(ownerID)}

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

func SearchPhotosByNameContaining(name string, ownerID string) ([]models.Photo, error) {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("only-cats").Collection("photos")

	searchQuery := bson.D{{
		"$search", bson.D{
			{"index", SEARCH_INDEX_NAME},
			{"text", bson.D{
				{"query", name},
				{"path", "name"},
			},
			},
		},
	}}

	ownerMatch := bson.D{{
		"$match", bson.D{
			{"owner_id", primitive.ObjectIDFromHex(ownerID)},
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

func SavePhoto(photo *models.Photo) (string, error) {
	client, ctx, cancel := getDBConnection()
	defer client.Disconnect(ctx)
	defer cancel()

	col := client.Database("only-cats").Collection("photos")

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

func DeletePhoto(id string) error {
	client, ctx, cancel := getDBConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	col := client.Database("only-cats").Collection("photos")

	_, err := col.DeleteOne(ctx, bson.M{"_id": primitive.ObjectIDFromHex(id)})

	return err
}
