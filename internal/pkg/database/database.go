package database

import (
	"context"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbConnString = "mongodb+srv://gomdb:AS6pjaMXYuCpaMkr@movies.ojtvno4.mongodb.net/?retryWrites=true&w=majority"

var client *mongo.Client
var db *mongo.Database

func init() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(dbConnString))
	if err != nil {
		logging.Panic(err.Error())
	}
	db = client.Database("movies")
}

func getMongoCollection(collection string) *mongo.Collection {
	return db.Collection(collection)
}

func Upsert(entity *domain.Entity, collection string) (*mongo.UpdateResult, error) {
	col := getMongoCollection(collection)

	filter := bson.M{"id": entity.ID}
	options := options.Replace().SetUpsert(true)

	result, err := col.ReplaceOne(context.TODO(), filter, entity.Data, options)

	return result, err
}
