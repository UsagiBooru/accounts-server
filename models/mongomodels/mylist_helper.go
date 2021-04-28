package mongomodels

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoMylistHelper is helper struct requires *mongo.Collection
type MongoMylistHelper struct {
	col *mongo.Collection
}

// NewMongoMylistHelper creates a helper for handle mylists endpoints
func NewMongoMylistHelper(md *mongo.Client) MongoMylistHelper {
	return MongoMylistHelper{md.Database("accounts").Collection("mylists")}
}

// CreateMylist inserts specified mylist to database
func (h *MongoMylistHelper) CreateMylist(MylistID int32, targetType string, targetID int32) (*MongoMylistStruct, error) {
	newMylistForNew := MongoMylistStruct{
		MylistID:    MylistID,
		Name:        "",
		Description: "",
		CreatedDate: time.Time{},
		UpdatedDate: time.Time{},
		Private:     true,
		Arts:        []MongoLightArtStruct{},
		Owner:       LightMongoAccountStruct{},
	}
	if _, err := h.col.InsertOne(context.Background(), newMylistForNew); err != nil {
		return nil, errors.New("insert Mylist failed")
	}
	return &newMylistForNew, nil
}

// FindMylist finds specified mylist from database
func (h *MongoMylistHelper) FindMylist(MylistID int32) (*MongoMylistStruct, error) {
	filter := bson.M{
		"MylistID": MylistID,
	}
	var Mylist MongoMylistStruct
	if err := h.col.FindOne(context.Background(), filter).Decode(&Mylist); err != nil {
		return nil, errors.New("mylist was not found")
	}
	return &Mylist, nil
}

// FindMylistUsingFilter finds specified mylist from database
func (h *MongoMylistHelper) FindMylistUsingFilter(filter bson.M) (*MongoMylistStruct, error) {
	var Mylist MongoMylistStruct
	if err := h.col.FindOne(context.Background(), filter).Decode(&Mylist); err != nil {
		return nil, errors.New("mylist was not found")
	}
	return &Mylist, nil
}

// DeleteMylist deletes specified mylist from database
func (h *MongoMylistHelper) DeleteMylist(MylistID int32) error {
	filter := bson.M{
		"MylistID": MylistID,
	}
	if _, err := h.col.DeleteOne(context.Background(), filter); err != nil {
		return errors.New("delete Mylist failed")
	}
	return nil
}
