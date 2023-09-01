package database

import (
	"context"
	"gomdb/internal/pkg/domain"
	"gomdb/internal/pkg/logging"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbConnString = "mongodb://gomdb-root:8lURb24nnHE8Kht3@10.0.0.126:27017/?retryWrites=true&w=majority"

var client *mongo.Client
var db *mongo.Database

func init() {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(dbConnString))
	if err != nil {
		logging.Panic(err.Error())
	}
	db = client.Database("gomdb")
}

func getMongoCollection(collection string) *mongo.Collection {
	return db.Collection(collection)
}

func Upsert(entity *domain.Entity, collection string) (*mongo.UpdateResult, error) {

	entity.Updated = time.Now()

	col := getMongoCollection(collection)

	filter := bson.M{"id": entity.ID}
	options := options.Replace().SetUpsert(true)

	result, err := col.ReplaceOne(context.TODO(), filter, entity, options)

	return result, err
}
