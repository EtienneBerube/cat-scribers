package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Photo struct {
	ID          string `json:"id"`
	Base64      string `json:"base64"`
	OwnerID     string `json:"owner_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// ToDAO Transfers the data from a Photo struct to a PhotoDAO struct.
func (p Photo) ToDAO(dao *PhotoDAO) {
	dao.ID, _ = primitive.ObjectIDFromHex(p.ID)
	dao.OwnerID, _ = primitive.ObjectIDFromHex(p.OwnerID)
	dao.Base64 = p.Base64
	dao.Description = p.Description
	dao.Name = p.Name
	dao.Type = p.Type
}

// PhotoDAO is a version of Photo that is used by MongoDB to interact with the data
type PhotoDAO struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Base64      string             `bson:"base64"`
	OwnerID     primitive.ObjectID `bson:"owner_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Type        string             `bson:"type"`
}

// ToModel Transfers the data from a PhotoDAO struct to a Photo struct.
func (dao PhotoDAO) ToModel(p *Photo) {
	p.ID = dao.ID.Hex()
	p.OwnerID = dao.OwnerID.Hex()
	p.Base64 = dao.Base64
	p.Description = dao.Description
	p.Name = dao.Name
	p.Type = dao.Type
}
