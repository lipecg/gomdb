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

func InsertOne(entity *domain.Entity, collection string) (*mongo.InsertOneResult, error) {

	col := getMongoCollection(collection)

	result, err := col.InsertOne(context.TODO(), entity)

	return result, err
}

func InsertMany(entities []interface{}, collection string) (*mongo.InsertManyResult, error) {

	col := getMongoCollection(collection)

	result, err := col.InsertMany(context.TODO(), entities)

	return result, err
}

func Upsert(entity *domain.Entity, collection string) (*mongo.UpdateResult, error) {

	entity.Updated = time.Now()

	col := getMongoCollection(collection)

	filter := bson.M{"id": entity.ID}
	options := options.Replace().SetUpsert(true)

	result, err := col.ReplaceOne(context.TODO(), filter, entity, options)

	return result, err
}

func Get(entity *domain.Entity, collection string) error {

	col := getMongoCollection(collection)

	filter := bson.M{"id": entity.ID}

	err := col.FindOne(context.Background(), filter).Decode(entity)

	return err
}

func List(field string, value any, collection string, list *[]domain.Entity) error {

	col := getMongoCollection(collection)

	filter := bson.M{}
	if field != "" && value != nil {
		filter = bson.M{field: value}
	}

	cursor, err := col.Find(context.Background(), filter)
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var entity domain.Entity
		err := cursor.Decode(&entity)
		if err != nil {
			return err
		}
		*list = append(*list, entity)
	}

	return nil
}
