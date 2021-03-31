package utils

import (
	"log"

	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertStructToBson(model interface{}) bson.M {
	user_json, err := json.Marshal(model)
	if err != nil {
		log.Fatal("Convert struct to json failed.")
	}
	bsonMap := bson.M{}
	err = bson.UnmarshalExtJSON([]byte(user_json), false, &bsonMap)
	if err != nil {
		log.Fatal("Convert json to bson failed.")
	}

	// map[_id:6063f51a033c71164cc13694 の _idをbytesにする
	// (何もしないとstringが入り文字列長が合わんってキレられる)
	// Avoid "an ObjectID string must be exactly 12 bytes long (got 24)" error
	if _, ok := bsonMap["_id"]; ok {
		bsonMap["_id"] = primitive.NewObjectID()
	}
	return bsonMap
}
