package mongomodels

import (
	"context"
	"errors"

	"github.com/UsagiBooru/accounts-server/utils/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoInviteHelper is helper struct requires *mongo.Collection
type MongoInviteHelper struct {
	col *mongo.Collection
}

// NewMongoInviteHelper creates a helper for handle account endpoints
func NewMongoInviteHelper(md *mongo.Client) MongoInviteHelper {
	return MongoInviteHelper{md.Database("accounts").Collection("invites")}
}

// CreateInvite inserts specified code to database
func (h *MongoInviteHelper) CreateInvite(code string, inviter AccountID) error {
	newInviteForNew := MongoInvite{
		ID:      primitive.NewObjectID(),
		Code:    code,
		Inviter: inviter,
		Invitee: 0,
	}
	if _, err := h.col.InsertOne(context.Background(), newInviteForNew); err != nil {
		return errors.New("insert invite failed")
	}
	return nil
}

// FindInvite finds specified invite info from database
func (h *MongoInviteHelper) FindInvite(code string) (*MongoInvite, error) {
	filter := bson.M{
		"code":    code,
		"invitee": 0,
	}
	var invite MongoInvite
	if err := h.col.FindOne(context.Background(), filter).Decode(&invite); err != nil {
		return nil, server.ErrInviteNotFound
	}
	return &invite, nil
}

// UseInvite set invitee to specified mongo invite data
func (h *MongoInviteHelper) UseInvite(mongoInviteID primitive.ObjectID, consumer AccountID) error {
	filter := bson.M{
		"_id": mongoInviteID,
	}
	set := bson.M{"$set": bson.M{"invitee": consumer}}
	if _, err := h.col.UpdateOne(context.Background(), filter, set); err != nil {
		return errors.New("update invite invitee failed")
	}
	return nil
}
