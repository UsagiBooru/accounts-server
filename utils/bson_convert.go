package utils

import (
	"log"

	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func ConvertOpenApiStructToBson(openapi_model interface{}) bson.M {
	user_json, err := json.Marshal(openapi_model)
	if err != nil {
		log.Fatal("Convert struct to json failed.")
	}
	bsonMap := bson.M{}
	err = bson.UnmarshalExtJSON([]byte(user_json), false, &bsonMap)
	if err != nil {
		log.Fatal("Convert json to bson failed.")
	}
	return bsonMap
}
