package database

import (
	"context"
	"gomdb/internal/pkg/domain"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbConnString = "mongodb://gomdb-root:8lURb24nnHE8Kht3@10.0.0.126:27017/?retryWrites=true&w=majority"

func getMongoCollection(ms mongoStore) *mongo.Collection {
	//TODO: Need to get proper collecton based on entity type
	return ms.Client.Database("gomdb").Collection("movies")
}

type mongoStore struct {
	Client *mongo.Client
}

func NewMongoStore() (domain.EntityDB, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbConnString))
	if err != nil {
		return nil, err
	}
	return mongoStore{Client: client}, nil
}

func (ms mongoStore) Get(id int) (*interface{}, error) {
	col := getMongoCollection(ms)
	filter := bson.M{"id": id}
	var entity interface{}
	err := col.FindOne(context.TODO(), filter).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// TODO: This won't work because of the query
func (ms mongoStore) List(query string) ([]*interface{}, error) {
	col := getMongoCollection(ms)
	filter := bson.M{"query": query}
	var entities []*interface{}
	cursor, err := col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &entities); err != nil {
		return nil, err
	}
	return entities, nil
}

func (ms mongoStore) Upsert(entity *interface{}) error {
	col := getMongoCollection(ms)

	filter := bson.M{"id": (reflect.ValueOf(entity).Elem())}
	options := options.Replace().SetUpsert(true)

	_, err := col.ReplaceOne(context.TODO(), filter, entity, options)

	return err
}
